package main

import (
	"fmt"
	"math/rand"

	redisbloom "github.com/RedisBloom/redisbloom-go"
)

func main() {
	// Connect to localhost with no password
	var rb = redisbloom.NewClient("localhost:6379", "nohelp", nil)

	bloomFilterName := fmt.Sprintf("myBloomFilter-%s", randomNumber())

	// Create a Bloom filter with a given error rate and initial capacity
	// 1 -> 100% error rate
	// 0 -> 0% error rate
	// 0.1 -> 10% error rate
	// 0.01 -> 1% error rate
	// 0.001 -> 0.1% error rate
	err := rb.Reserve(bloomFilterName, 0.01 /*error rate*/, 2000 /*capacity*/)
	if err != nil {
		fmt.Println("Failed to create Bloom filter:", err)
		// return
	}

	// Generate and add unique strings to the Bloom filter
	data1 := generateUniqueStrings(1000)
	for _, item := range data1 {
		_, err := rb.Add(bloomFilterName, item)
		if err != nil {
			fmt.Println("Failed to add item:", err)
			return
		}
	}

	data2 := generateUniqueStrings(1000)
	// Test for false positive by checking items in the original data
	for _, item := range data2 {
		exists, err := rb.Exists(bloomFilterName, item)
		if err != nil {
			fmt.Println("Failed to check item:", err)
			return
		}
		if exists {
			fmt.Println("False positive detected for:", item)
		}
	}
}

// Generate a set of unique strings
func generateUniqueStrings(count int) []string {
	var data []string

	randomNumber := randomNumber()

	for i := 0; i < count; i++ {
		data = append(data, fmt.Sprintf("%s-%d", randomNumber, i))
	}

	return data
}

func randomNumber() string {
	return fmt.Sprintf("%d", rand.Intn(1000))
}
