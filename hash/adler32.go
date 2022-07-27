package hash

const MOD_ADLER uint32 = 65521

// Adler32Checksums generates 3 integers - a, b and their sum
// https://en.wikipedia.org/wiki/Adler-32
func Adler32Checksums(block []byte) (uint32, uint32, uint32) {
	var a, b uint32 = 1, 0
	for i := 0; i < len(block); i++ {
		a = (a + uint32(block[i])) % MOD_ADLER
		b = (b + a) % MOD_ADLER
	}
	return a, b, a + b*MOD_ADLER
}

// Slide slides rolling window (aabb and w=3 => aab to abb)
func Adler32Slide(a, b uint32, left, right byte, size int) (uint32, uint32, uint32) {
	l, r := uint32(left), uint32(right)
	a = (a - l + r) % MOD_ADLER
	b = ((b - uint32(size)*l) + a - 1) % MOD_ADLER
	return a, b, a + b*MOD_ADLER
}
