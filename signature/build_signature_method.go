package signature

import (
	"errors"
	"fmt"
	"io"

	"deltagopher/hash"
)

func allZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

// NextChunkHashes gets the weak/strong hashes of the new obtained chunk
func (sb *SignatureBuilder) NextChunkHashes() (uint32, [16]byte, error) {
	n, err := sb.reader.Read(sb.buf)
	sb.curPos += n
	if err != nil {
		fmt.Println(err)
		return 0, [16]byte{}, err
	}
	_, _, weakHash := hash.Adler32Checksums(sb.buf)
	if weakHash == 1 || weakHash == 0 {
		panic(errors.New("Adler32 hash generated invalid hashsum"))
	}
	strongHash := hash.MD5Checksum(sb.buf)

	return weakHash, strongHash, nil
}

// BuildSignature create a Signature object representing the chunked and hashed file content
func (sb *SignatureBuilder) BuildSignature() *Signature {
	s := &Signature{
		Checksums: nil,
		BlockSize: sb.size,
		Hashing:   "adler32",
	}
	var ch *Checksum
	for {
		weakHash, strongHash, err := sb.NextChunkHashes()
		if err != nil {
			if err == io.EOF {
				if allZero(sb.buf) {
					panic(errors.New("The file is empty"))
				}
				fmt.Println("End of file")
				return s
			}
			panic(err)
		}

		ch = &Checksum{
			WeakCheaksum:   weakHash,
			StrongChecksum: strongHash,
			Start:          sb.curPos - sb.size,
			End:            sb.curPos,
		}
		if sb.full {
			ch.Content = append(ch.Content, sb.buf...)
		}
		s.Checksums = append(s.Checksums, ch)
	}
	return s
}
