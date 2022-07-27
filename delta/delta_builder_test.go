package delta

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"gordiff/hash"
	"gordiff/signature"
)

var _, _, a1 = hash.Adler32Checksums([]byte("abc"))
var _, _, a2 = hash.Adler32Checksums([]byte("def"))
var _, _, a3 = hash.Adler32Checksums([]byte("qwe"))

// abcdef
func makeSig() *signature.Signature {
	return &signature.Signature{
		Checksums: []*signature.Checksum{
			{
				WeakCheaksum:   a1,
				StrongChecksum: hash.MD5Checksum([]byte("abc")),
				Start:          0,
				End:            3,
			},
			{
				WeakCheaksum:   a2,
				StrongChecksum: hash.MD5Checksum([]byte("def")),
				Start:          3,
				End:            6,
			},
		},
		BlockSize: 3,
		Hashing:   "adler32",
	}
}
func TestDeltaBuilder_BuildDelta(t *testing.T) {
	type fields struct {
		reader    io.ReadSeeker
		signature *signature.Signature
		delta     *Delta
		blockSize int
		pos       int
		a         uint32
		b         uint32
		sum       uint32
	}

	md52 := hash.MD5Checksum([]byte("def"))
	md51 := hash.MD5Checksum([]byte("abc"))
	tests := []struct {
		name    string
		fields  fields
		want    *Delta
		wantErr bool
		textErr string
	}{
		{
			name: "success 1 inserted 1 copied 1 deleted",
			fields: fields{
				reader:    bytes.NewReader([]byte("abcqwe")),
				signature: makeSig(),
			},
			want: &Delta{
				Inserted: []*SingleDelta{
					{
						Start:     3,
						End:       6,
						DiffBytes: []byte("qwe"),
					},
				},
				Deleted: []*SingleDelta{
					{
						Start:          3,
						End:            6,
						WeakCheaksum:   a2,
						StrongChecksum: &md52,
					},
				},
				Copied: []*SingleDelta{
					{
						WeakCheaksum:   a1,
						StrongChecksum: &md51,
						Start:          0,
						End:            3,
					}},
			},
		},
		{
			name: "success 2 copied",
			fields: fields{
				reader:    bytes.NewReader([]byte("abcdef")),
				signature: makeSig(),
			},
			want: &Delta{
				Inserted: nil,
				Deleted:  nil,
				Copied: []*SingleDelta{
					{
						WeakCheaksum:   a1,
						StrongChecksum: &md51,
						Start:          0,
						End:            3,
					},
					{
						WeakCheaksum:   a2,
						StrongChecksum: &md52,
						Start:          3,
						End:            6,
					}},
			},
		},
		{
			name: "success 2 deleted",
			fields: fields{
				reader:    bytes.NewReader([]byte("")),
				signature: makeSig(),
			},
			want: &Delta{
				Inserted: nil,
				Copied:   nil,
				Deleted: []*SingleDelta{
					{
						Start:          0,
						End:            3,
						WeakCheaksum:   a1,
						StrongChecksum: &md51,
					},
					{
						Start:          3,
						End:            6,
						WeakCheaksum:   a2,
						StrongChecksum: &md52,
					},
				},
			},
		},
		{
			name: "insert all",
			fields: fields{
				reader: bytes.NewReader([]byte("abcdef")),
				signature: &signature.Signature{
					Checksums: nil,
					BlockSize: 3,
					Hashing:   "adler32",
				},
			},
			want: &Delta{
				Inserted: []*SingleDelta{
					{
						Start:     0,
						End:       6,
						DiffBytes: []byte("abcdef"),
					},
				},
				Copied:  nil,
				Deleted: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewDeltaBuilder(tt.fields.signature, tt.fields.reader)
			d := db.BuildDelta()
			assert.Equal(t, tt.want.Inserted, d.Inserted)
			assert.Equal(t, tt.want.Copied, d.Copied)
			assert.Equal(t, tt.want.Deleted, d.Deleted)
		})
	}
}
