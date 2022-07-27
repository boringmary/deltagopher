package signature

import (
	"bytes"
	"crypto/md5"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"gordiff/hash"
)

const text1 = "Lorem Ipsum is simply dummy text of the printing and typesetting industry."
const text2 = "Lorem Ipsum 123456789qwemmy text of the printing and typeswwwwwwwwwwwwwwwwwwww."
const blockSize = 12

func MakeTestSig() *Signature {
	return &Signature{
		Checksums: []*Checksum{},
		BlockSize: blockSize,
		Hashing:   "adler32",
	}
}

func MakeChecksums() []*Checksum {
	cs := []*Checksum{}
	r := bytes.NewReader([]byte(text1))
	buf := make([]byte, blockSize)
	cur := 0
	for {
		n, err := r.Read(buf)
		cur += n
		if err == io.EOF {
			break
		}
		_, _, sum := hash.Adler32Checksums(buf)
		cs = append(cs, &Checksum{
			WeakCheaksum:   sum,
			StrongChecksum: md5.Sum(buf),
			Start:          cur - blockSize,
			End:            cur,
		})
	}

	return cs
}

func TestSignatureBuilder_BuildSignature(t *testing.T) {
	type fields struct {
		reader    io.ReadSeeker
		buf       []byte
		curPos    int
		size      int
		signature *Signature
	}

	tests := []struct {
		name    string
		fields  fields
		want    *Signature
		isErr   bool
		textErr string
	}{
		{
			name: "success build",
			fields: fields{
				reader:    bytes.NewReader([]byte(text1)),
				buf:       make([]byte, blockSize),
				curPos:    0,
				size:      blockSize,
				signature: MakeTestSig(),
			},
			want: &Signature{
				Checksums: MakeChecksums(),
				BlockSize: blockSize,
				Hashing:   "adler32",
			},
		},
		{
			name: "empty reader",
			fields: fields{
				reader:    bytes.NewReader([]byte("")),
				buf:       make([]byte, blockSize),
				curPos:    0,
				size:      blockSize,
				signature: MakeTestSig(),
			},
			isErr:   true,
			textErr: "The file is empty",
		},
		{
			name: "0 buff",
			fields: fields{
				reader:    bytes.NewReader([]byte(text1)),
				buf:       make([]byte, 0),
				curPos:    0,
				size:      blockSize,
				signature: MakeTestSig(),
			},
			isErr:   true,
			textErr: "Adler32 hash generated invalid hashsum",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &SignatureBuilder{
				reader:    tt.fields.reader,
				buf:       tt.fields.buf,
				curPos:    tt.fields.curPos,
				size:      tt.fields.size,
				signature: tt.fields.signature,
			}

			if tt.isErr {
				assert.PanicsWithError(t, tt.textErr,
					func() {
						sb.BuildSignature()
					})
			} else {
				got := sb.BuildSignature()

				for i, ch := range got.Checksums {
					assert.Equal(t, ch.StrongChecksum, tt.want.Checksums[i].StrongChecksum)
					assert.Equal(t, ch.WeakCheaksum, tt.want.Checksums[i].WeakCheaksum)
					assert.Equal(t, ch.Start, tt.want.Checksums[i].Start)
					assert.Equal(t, ch.End, tt.want.Checksums[i].End)
				}
			}
		})
	}
}
