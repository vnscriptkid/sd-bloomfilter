package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/gomodule/redigo/redis"
)

func main() {
	// Connect to Redis
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer conn.Close()

	n := 10000 // Number of email addresses to insert
	data := generateEmailAddresses(n)

	// Redis Set
	setKey := "newsletterSet"
	initialSetSize, err := getRedisMemoryUsage(conn)
	if err != nil {
		log.Fatalf("Failed to get initial memory usage for set: %v", err)
	}

	for _, d := range data {
		_, err = conn.Do("SADD", setKey, d)
		if err != nil {
			log.Fatalf("Failed to add element to set: %v", err)
		}
	}

	finalSetSize, err := getRedisMemoryUsage(conn)
	if err != nil {
		log.Fatalf("Failed to get final memory usage for set: %v", err)
	}
	setMemoryUsage := finalSetSize - initialSetSize

	// Redis Bloom Filter
	bloomKey := "newsletterBloom"
	err = createBloomFilter(conn, bloomKey, n, 0.01)
	if err != nil {
		log.Fatalf("Failed to create Bloom filter: %v", err)
	}

	initialBloomSize, err := getRedisMemoryUsage(conn)
	if err != nil {
		log.Fatalf("Failed to get initial memory usage for Bloom filter: %v", err)
	}

	for _, d := range data {
		_, err = conn.Do("BF.ADD", bloomKey, d)
		if err != nil {
			log.Fatalf("Failed to add element to Bloom filter: %v", err)
		}
	}

	finalBloomSize, err := getRedisMemoryUsage(conn)
	if err != nil {
		log.Fatalf("Failed to get final memory usage for Bloom filter: %v", err)
	}
	bloomMemoryUsage := finalBloomSize - initialBloomSize

	fmt.Printf("Memory usage for Redis set: %d bytes\n", setMemoryUsage)
	fmt.Printf("Memory usage for Redis Bloom filter: %d bytes\n", bloomMemoryUsage)
}

func generateEmailAddresses(n int) []string {
	domains := []string{"example.com", "test.com", "demo.com"}
	data := make([]string, n)
	for i := 0; i < n; i++ {
		email := fmt.Sprintf("user%d@%s", i, domains[rand.Intn(len(domains))])
		data[i] = email
	}
	return data
}

func getRedisMemoryUsage(conn redis.Conn) (int64, error) {
	reply, err := redis.String(conn.Do("INFO", "memory"))
	if err != nil {
		return 0, err
	}

	var usedMemory int64
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, "used_memory:") {
			fmt.Sscanf(line, "used_memory:%d", &usedMemory)
			break
		}
	}

	return usedMemory, nil
}

func createBloomFilter(conn redis.Conn, key string, capacity int, errorRate float64) error {
	_, err := conn.Do("BF.RESERVE", key, errorRate, capacity)
	return err
}
