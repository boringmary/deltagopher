package signature

import (
	"errors"
	"io"
)

//Old file signature, consisting from hashed chunks
type SignatureBuilder struct {
	reader io.ReadSeeker
	buf    []byte
	curPos int
	size   int

	signature *Signature
}

func NewSignatureBuilder(reader io.ReadSeeker, blockSize int) *SignatureBuilder {
	if blockSize == 0 {
		panic(errors.New("The blocksize cannot be 0"))
	}
	return &SignatureBuilder{
		reader: reader,
		buf:    make([]byte, blockSize),
		curPos: 0,
		size:   blockSize,
	}
}
