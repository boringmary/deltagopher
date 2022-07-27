package delta

// Delta of a single chunk
type SingleDelta struct {
	WeakCheaksum   uint32    `yaml:"weak,omitempty" json:"weak,omitempty"`
	StrongChecksum *[16]byte `yaml:"strong,omitempty" json:"-"`
	Start          int       `yaml:"start" json:"start,omitempty"`
	End            int       `yaml:"end" json:"end,omitempty"`

	// For Inserted
	DiffBytes []byte `yaml:"diff,omitempty" json:"diff,omitempty"`
}
