package hash

import "crypto/md5"

func MD5Checksum(data []byte) [16]byte {
	return md5.Sum(data)
}
