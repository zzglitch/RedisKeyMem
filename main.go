package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/go-redis/redis"
)

func main() {
	hostPtr := flag.String("h", "127.0.0.1", "Redis Host")
	portPtr := flag.Int("p", 6379, "Redis Port")
	passwordPtr := flag.String("a", "", "Redis password")
	flag.Parse()

	keyFilter := "*"
	if len(flag.Args()) > 0 {
		keyFilter = flag.Args()[0]
	}

	client := redis.NewClient(
		&redis.Options{
			Addr:     *hostPtr + ":" + strconv.Itoa(*portPtr),
			Password: *passwordPtr,
			DB:       0,
		})
	results, err := client.Keys(keyFilter).Result()
	if err != nil {
		fmt.Printf("Error fetching keys: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Number of keys: %v\n", len(results))
	var totalBytes int
	var maxBytes int
	var maxKey string
	regex, _ := regexp.Compile(`serializedlength:(\d*)`)
	for _, key := range results {
		val, err := client.DebugObject(key).Result()
		if err != nil {
			fmt.Printf("Error getting debug infomration about %s: %v\n", key, err)
			continue
		}
		matches := regex.FindStringSubmatch(val)
		var bytes int
		bytes, _ = strconv.Atoi(matches[1])
		totalBytes += bytes
		if bytes > maxBytes {
			maxBytes = bytes
			maxKey = key
		}
	}
	if len(results) > 0 {
		avgBytes := totalBytes / len(results)
		fmt.Printf("Average serialized size: %d\n", avgBytes)
	}
	fmt.Printf("Max key \"%s\" serialized size: %d\n", maxKey, maxBytes)
}
