package utils

import (
	"encoding/hex"
	"strings"

	"github.com/OneOfOne/xxhash"
)

func GetHash(text string) (string, error) {
	hasher := xxhash.New64()
	_, err := hasher.Write([]byte(text))
	if err != nil {
		return "", nil
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// Contains returns true if target string is present in the strings slice.
// Comparison is case-insensitive.
func Contains(slice []string, lookup string) bool {
	for _, val := range slice {
		if strings.EqualFold(val, lookup) {
			return true
		}
	}
	return false
}
