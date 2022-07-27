package signature

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Checksum struct {
	WeakCheaksum   uint32   `yaml:"weak"`
	StrongChecksum [16]byte `yaml:"strong"`
	Start          int      `yaml:"start"`
	End            int      `yaml:"end"`
}

type Signature struct {
	Checksums []*Checksum `yaml:"checksums"`

	BlockSize int    `yaml:"size"`
	Hashing   string `yaml:"algo"`
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
