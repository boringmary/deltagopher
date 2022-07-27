package delta

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Delta represents the difference betweeb 2 byte arrays
type Delta struct {
	// blocks to be inserted
	Inserted []*SingleDelta `yaml:"insert,omitempty" json:"insert,omitempty"`

	// blocks to be deleted
	Deleted []*SingleDelta `yaml:"delete,omitempty" json:"delete,omitempty"`

	// blocks to be copied
	Copied []*SingleDelta `yaml:"copy,omitempty" json:"copy,omitempty"`
}

func NewDelta() *Delta {
	return &Delta{
		Inserted: nil,
		Deleted:  nil,
		Copied:   nil,
	}
}

func (d *Delta) MarshalYAML() ([]byte, error) {
	y, err := yaml.Marshal(d)
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

func (d *Delta) MarshalJSON() ([]byte, error) {
	y, err := json.Marshal(*d)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	return y, nil
}

func UnmarshalJSON(bytes []byte) (*Delta, error) {
	s := &Delta{}
	err := json.Unmarshal(bytes, s)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	return s, nil
}
