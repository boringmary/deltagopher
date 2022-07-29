package delta

import (
	"fmt"
	"io"

	"deltagopher/hash"
)

// BuildDelta roll across the new file content and generate delta obj on the fly.
func (db *DeltaBuilder) BuildDelta() *Delta {
	weakChecksumsMap := map[uint32][16]byte{}

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
		db.foundHashes[db.sum] = true
	}
	for _, item := range db.signature.Checksums {
		var d *SingleDelta
		if _, ok := db.foundHashes[item.WeakCheaksum]; !ok {
			d = &SingleDelta{
				WeakCheaksum:   item.WeakCheaksum,
				StrongChecksum: &item.StrongChecksum,
				Start:          item.Start,
				End:            item.End,
			}
			if db.full && len(item.Content) != 0 {
				d.DiffBytes = append(d.DiffBytes, item.Content...)
			}
			db.delta.Deleted = append(db.delta.Deleted, d)
		}
	}

	return db.delta
}

// AppendFilteredCandidate appends only the candidate whose old strong hash is matched with new strong hash
func (db *DeltaBuilder) AppendFilteredCandidate(strongHash [16]byte) bool {
	newStrongHash := hash.MD5Checksum(db.window)
	if strongHash == newStrongHash {
		d := &SingleDelta{
			WeakCheaksum:   db.sum,
			StrongChecksum: &newStrongHash,
			Start:          db.pos - db.blockSize,
			End:            db.pos,
		}
		if db.full && len(db.window) != 0 {
			d.DiffBytes = append(d.DiffBytes, db.window...)
		}
		db.delta.Copied = append(db.delta.Copied, d)
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
		db.foundHashes[db.sum] = true
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
