package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var block = "abc"

func TestAdler32Checksums(t *testing.T) {
	type args struct {
		block []byte
	}
	tests := []struct {
		name  string
		args  args
		want  uint32
		want1 uint32
		want2 uint32
	}{
		{
			name:  "success",
			args:  args{block: []byte(block)},
			want:  295,
			want1: 589,
			want2: 38592164,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := Adler32Checksums(tt.args.block)
			assert.Equal(t, got, tt.want)
			assert.Equal(t, got1, tt.want1)
			assert.Equal(t, got2, tt.want2)
		})
	}
}

func TestAdler32Slide(t *testing.T) {
	type args struct {
		a     uint32
		b     uint32
		left  byte
		right byte
		size  int
	}
	tests := []struct {
		name  string
		args  args
		want  uint32
		want1 uint32
		want2 uint32
	}{
		{
			name: "success",
			args: args{
				a:     uint32(295),
				b:     uint32(589),
				left:  []byte("a")[0],
				right: []byte("d")[0],
				size:  3,
			},
			want:  298,
			want1: 595,
			want2: 38985293,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := Adler32Slide(tt.args.a, tt.args.b, tt.args.left, tt.args.right, tt.args.size)
			assert.Equal(t, got, tt.want)
			assert.Equal(t, got1, tt.want1)
			assert.Equal(t, got2, tt.want2)
		})
	}
}
