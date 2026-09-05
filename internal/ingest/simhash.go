package ingest

import (
	"hash/fnv"
	"math/bits"
	"strings"
	"unicode"
)

// SimHash64 computes a 64-bit SimHash fingerprint for the provided text.
func SimHash64(text string) uint64 {
	words := tokenize(text)
	if len(words) == 0 {
		return 0
	}

	var v [64]int

	for _, w := range words {
		h := fnvHash64(w)
		for i := 0; i < 64; i++ {
			bit := (h >> i) & 1
			if bit == 1 {
				v[i]++
			} else {
				v[i]--
			}
		}
	}

	var fingerprint uint64
	for i := 0; i < 64; i++ {
		if v[i] > 0 {
			fingerprint |= (1 << i)
		}
	}

	return fingerprint
}

// HammingDistance calculates the number of differing bits between two 64-bit hashes.
func HammingDistance(a, b uint64) int {
	return bits.OnesCount64(a ^ b)
}

// IsNearDuplicate returns true if the Hamming distance between two texts is <= threshold (default 3).
func IsNearDuplicate(textA, textB string, threshold int) bool {
	hashA := SimHash64(textA)
	hashB := SimHash64(textB)
	return HammingDistance(hashA, hashB) <= threshold
}

func tokenize(text string) []string {
	var words []string
	var current strings.Builder

	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
		} else {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		}
	}
	if current.Len() > 0 {
		words = append(words, current.String())
	}
	return words
}

func fnvHash64(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}
