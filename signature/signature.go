package signature

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// The single chunk checksum
type Checksum struct {
	// weak hash
	WeakCheaksum uint32 `yaml:"weak"`

	// strong hash
	StrongChecksum [16]byte `yaml:"strong"`

	// start position of the chunk
	Start int `yaml:"start"`

	// end position of the chunk
	End int `yaml:"end"`
}

type Signature struct {
	Checksums []*Checksum `yaml:"checksums"`

	// size of the chunk
	BlockSize int `yaml:"size"`

	// hashing algorithm
	Hashing string `yaml:"algo"`
}

func (s *Signature) MarshalYAML() ([]byte, error) {
	y, err := yaml.Marshal(s)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	return y, nil
}

func UnmarshalYAML(bytes []byte) (*Signature, error) {
	s := &Signature{}
	err := yaml.Unmarshal(bytes, s)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	return s, nil
}
