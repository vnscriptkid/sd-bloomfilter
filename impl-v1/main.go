package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
)

type BloomFilterBehaviors interface {
	Add(data []byte)
	Contains(data []byte) bool
}

type BloomFilter struct {
	bitset    []bool
	size      uint
	hashFuncs []func(data []byte) uint
}

var _ BloomFilterBehaviors = (*BloomFilter)(nil)

// NewBloomFilter creates a new Bloom filter with a specified size and number of hash functions.
func NewBloomFilter(size uint, hashCount int) *BloomFilter {
	bf := &BloomFilter{
		bitset: make([]bool, size),
		size:   size,
	}

	// Add hash functions
	for i := 0; i < hashCount; i++ {
		bf.hashFuncs = append(bf.hashFuncs, bf.generateHashFunc(i))
	}

	return bf
}

// Add adds an element to the Bloom filter.
func (bf *BloomFilter) Add(data []byte) {
	for _, hashFunc := range bf.hashFuncs {
		index := hashFunc(data) % bf.size
		bf.bitset[index] = true
	}
}

// Contains checks if an element is in the Bloom filter.
func (bf *BloomFilter) Contains(data []byte) bool {
	for _, hashFunc := range bf.hashFuncs {
		index := hashFunc(data) % bf.size
		if !bf.bitset[index] {
			return false
		}
	}
	return true
}

// generateHashFunc generates a hash function based on an index.
func (bf *BloomFilter) generateHashFunc(i int) func(data []byte) uint {
	return func(data []byte) uint {
		hash1 := fnv.New64a()
		hash1.Write(data)
		hash1Sum := hash1.Sum64()

		hash2 := sha256.New()
		hash2.Write(data)
		hash2Sum := binary.BigEndian.Uint64(hash2.Sum(nil))

		return uint((hash1Sum + uint64(i)*hash2Sum) % uint64(bf.size))
	}
}

// EstimateParameters estimates the size of the bitset and the number of hash functions needed for the given parameters.
func EstimateParameters(n uint, p float64) (uint, int) {
	m := uint(math.Ceil(-float64(n) * math.Log(p) / (math.Ln2 * math.Ln2)))
	k := int(math.Round(math.Ln2 * float64(m) / float64(n)))
	return m, k
}

func main() {
	n := uint(1000 * 1000 * 1000) // Number of elements expected to be inserted
	p := 0.01                     // False positive probability

	size, hashCount := EstimateParameters(n, p)
	bf := NewBloomFilter(size, hashCount)

	fmt.Printf("Size: %d, Hash functions: %d\n", size, hashCount)
	fmt.Printf("Size in KB: %f\n", float64(size)/8/1024)
	fmt.Printf("Size in MB: %f\n", float64(size)/8/1024/1024)

	data := []byte("example")
	bf.Add(data)

	fmt.Println("Contains 'example':", bf.Contains(data))
	fmt.Println("Contains 'another':", bf.Contains([]byte("another")))
}
