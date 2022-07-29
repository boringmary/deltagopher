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

	// set true to include block content to the files
	full bool
}

func NewSignatureBuilder(reader io.ReadSeeker, blockSize int, full bool) *SignatureBuilder {
	if blockSize == 0 {
		panic(errors.New("The blocksize cannot be 0"))
	}
	return &SignatureBuilder{
		reader: reader,
		buf:    make([]byte, blockSize),
		curPos: 0,
		size:   blockSize,
		full:   full,
	}
}
