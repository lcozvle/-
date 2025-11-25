package main

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	val, err := client.Get("Spider:TargetRid").Result()
	if err != nil {
		fmt.Println("Redis Error:", err)
	} else {
		fmt.Println("Current Rid in Redis:", val)
	}
}
