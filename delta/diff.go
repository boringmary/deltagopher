package delta

// Delta of a single chunk
type SingleDelta struct {
	WeakCheaksum   uint32    `yaml:"weak,omitempty"`
	StrongChecksum *[16]byte `yaml:"strong,omitempty"`
	Start          int       `yaml:"start"`
	End            int       `yaml:"end"`

	// For Inserted
	DiffBytes []byte `yaml:"diff,omitempty"`
}
