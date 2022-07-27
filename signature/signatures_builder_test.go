package signature

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSignatureBuilder(t *testing.T) {
	type args struct {
		reader    io.ReadSeeker
		blockSize int
	}

	r := bytes.NewReader([]byte("Lorem ipsum"))
	tests := []struct {
		name    string
		args    args
		want    *SignatureBuilder
		isErr   bool
		textErr string
	}{
		{
			name: "success",
			args: args{
				reader:    r,
				blockSize: 12,
			},
			want: &SignatureBuilder{
				reader:    r,
				buf:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				curPos:    0,
				size:      12,
				signature: nil,
			},
		},
		{
			name: "error 0 blocksize",
			args: args{
				reader:    r,
				blockSize: 0,
			},
			isErr:   true,
			textErr: "The blocksize cannot be 0",
		},
	}
	for _, tt := range tests {
		if tt.isErr {
			assert.PanicsWithError(t, tt.textErr,
				func() {
					_ = NewSignatureBuilder(tt.args.reader, tt.args.blockSize)
				})
		} else {
			got := NewSignatureBuilder(tt.args.reader, tt.args.blockSize)
			assert.Equal(t, got, tt.want)
		}
	}
}
