package delta

import (
	"errors"
	"fmt"
	"io"

	"deltagopher/hash"
	"deltagopher/signature"
)

const MOD_ADLER uint32 = 65521

type DeltaBuilder struct {
	reader io.ReadSeeker

	signature *signature.Signature
	delta     *Delta

	blockSize int
	window    []byte
	pos       int

	a   uint32
	b   uint32
	sum uint32

	visited []byte
}

func NewDeltaBuilder(signature *signature.Signature, reader io.ReadSeeker) *DeltaBuilder {
	if signature.BlockSize == 0 {
		panic(errors.New("The size of the block is invalid: 0"))
	}
	return &DeltaBuilder{
		signature: signature,
		reader:    reader,
		pos:       0,
		blockSize: signature.BlockSize,
		window:    make([]byte, signature.BlockSize),
		visited:   make([]byte, 0),
		delta:     NewDelta(),
	}
}

// BuildDelta roll across the new file content and generate delta obj on the fly.
func (db *DeltaBuilder) BuildDelta() *Delta {
	weakChecksumsMap := map[uint32][16]byte{}
	foundHashes := map[uint32]bool{}

	for _, chs := range db.signature.Checksums {
		weakChecksumsMap[chs.WeakCheaksum] = chs.StrongChecksum
	}
	for {
		err := db.Roll(weakChecksumsMap)
		if err == io.EOF {
			fmt.Println("End of the file")
			break
		} else if err != nil {
			panic(err)
		}
		foundHashes[db.sum] = true
	}
	for _, item := range db.signature.Checksums {
		if _, ok := foundHashes[item.WeakCheaksum]; !ok {
			db.delta.Deleted = append(db.delta.Deleted, &SingleDelta{
				WeakCheaksum:   item.WeakCheaksum,
				StrongChecksum: &item.StrongChecksum,
				Start:          item.Start,
				End:            item.End,
				DiffBytes:      nil,
			})
		}
	}

	return db.delta
}

func (db *DeltaBuilder) AppendFilteredCandidate(strongHash [16]byte) bool {
	newStrongHash := hash.MD5Checksum(db.window)
	if strongHash == newStrongHash {
		db.delta.Copied = append(db.delta.Copied, &SingleDelta{
			WeakCheaksum:   db.sum,
			StrongChecksum: &newStrongHash,
			Start:          db.pos - db.blockSize,
			End:            db.pos,
			DiffBytes:      nil,
		})
		return true
	} else {
		return false
	}
}

// PullRemainingBytes check if we have some already traversed bytes in the buf that has to be inserted
func (db *DeltaBuilder) PullRemainingBytes() {
	if len(db.visited) != 0 {
		db.delta.Inserted = append(db.delta.Inserted, &SingleDelta{
			Start:     db.pos - len(db.visited),
			End:       db.pos,
			DiffBytes: db.visited,
		})
	}
}

// Roll rolls the sliding window and search for already existing blocks
func (db *DeltaBuilder) Roll(checksums map[uint32][16]byte) error {
	n, err := db.reader.Read(db.window)
	db.pos += n
	if err != nil {
		if err == io.EOF {
			db.PullRemainingBytes()
		}
		return err
	}

	db.a, db.b, db.sum = hash.Adler32Checksums(db.window)
	if strongHash, ok := checksums[db.sum]; ok {
		// Filter by strong md5 hash
		if ok := db.AppendFilteredCandidate(strongHash); ok {
			return nil
		}
	}

	db.visited = append(db.visited, db.window...)

	for {
		err = db.Next()
		if err != nil {
			if err == io.EOF {
				db.PullRemainingBytes()
				return err
			}
			panic(err)
		}
		if strongHash, ok := checksums[db.sum]; ok {
			if db.visited != nil {
				db.delta.Inserted = append(db.delta.Inserted, &SingleDelta{
					Start:     db.pos - len(db.visited),
					End:       db.pos,
					DiffBytes: db.visited,
				})
				db.visited = nil
			}

			db.AppendFilteredCandidate(strongHash)
			return nil
		}
		db.visited = append(db.visited, db.window[len(db.window)-1])
	}

	return nil
}

// Next populate buffer and genrate hash for the new rolling iteration
func (db *DeltaBuilder) Next() error {
	prev := db.window[0]
	next := make([]byte, 1)

	n, err := db.reader.Read(next)
	db.pos += n
	if err != nil {
		return err
	}

	db.a, db.b, db.sum = hash.Adler32Slide(db.a, db.b, prev, next[0], db.blockSize)
	db.window = append(db.window[1:], next[0])
	return nil
}
