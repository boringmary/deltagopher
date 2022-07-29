package delta

import (
	"errors"
	"io"

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

	visited     []byte
	foundHashes map[uint32]bool

	// should deleted/copy content be included to the delta
	full bool
}

func NewDeltaBuilder(signature *signature.Signature, reader io.ReadSeeker, full bool) *DeltaBuilder {
	if signature.BlockSize == 0 {
		panic(errors.New("The size of the block is invalid: 0"))
	}
	return &DeltaBuilder{
		signature:   signature,
		reader:      reader,
		pos:         0,
		blockSize:   signature.BlockSize,
		window:      make([]byte, signature.BlockSize),
		visited:     make([]byte, 0),
		foundHashes: map[uint32]bool{},
		full:        full,
		delta:       NewDelta(),
	}
}
