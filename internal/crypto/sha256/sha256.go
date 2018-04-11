package sha256

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// Sum256Hex calculates the sha256 digest of date in hexadecimal
func Sum256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	encoded := hex.EncodeToString(sum[:])
	return encoded
}

// Sum256HexFile calculates the sha256 sum of a file, in hexadecimal
func Sum256HexFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	sum := hash.Sum(nil)

	encoded := hex.EncodeToString(sum[:])

	return encoded, nil
}
