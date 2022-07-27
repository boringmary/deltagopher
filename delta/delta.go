package delta

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Delta represents the difference betweeb 2 byte arrays
type Delta struct {
	Inserted []*SingleDelta `yaml:"insert,omitempty"`
	Deleted  []*SingleDelta `yaml:"delete,omitempty"`
	Copied   []*SingleDelta `yaml:"copy,omitempty"`
}

func NewDelta() *Delta {
	return &Delta{
		Inserted: nil,
		Deleted:  nil,
		Copied:   nil,
	}
}

func (s *Delta) MarshalYAML() ([]byte, error) {
	y, err := yaml.Marshal(s)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	return y, nil
}

func UnmarshalYAML(bytes []byte) (*Delta, error) {
	s := &Delta{}
	err := yaml.Unmarshal(bytes, s)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	return s, nil
}
