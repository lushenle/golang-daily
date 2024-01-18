package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// randomInt generates a random integer between min and max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// randomString generates a random string of length n
func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// writeRedis set some random key-value to rds
func writeRedis(ctx context.Context, rdb *redis.Client) error {
	times := randomInt(50, 100)

	for i := 0; i < times; i++ {
		key := randomString(5)
		value := randomString(randomInt(10, 50))
		if err := rdb.Set(ctx, key, value, 10*time.Second).Err(); err != nil {
			return err
		}
	}

	return nil
}

// cron custom scheduling
func cron(ctx context.Context, startTime time.Time, delay time.Duration) <-chan time.Time {
	// Create the channel which we will return
	stream := make(chan time.Time, 1)

	// Calculating the first start time in the future
	// Need to check if the time is zero (e.g. if time.Time{} was used)
	if !startTime.IsZero() {
		diff := time.Until(startTime)
		if diff < 0 {
			total := diff - delay
			times := total / delay * -1

			startTime = startTime.Add(times * delay)
		}
	}

	// Run this in a goroutine, or our function will block until the first event
	go func() {

		// Run the first event after it gets to the start time
		t := <-time.After(time.Until(startTime))
		stream <- t

		// Open a new ticker
		ticker := time.NewTicker(delay)
		// Make sure to stop the ticker when we're done
		defer ticker.Stop()

		// Listen on both the ticker and the context done channel to know when to stop
		for {
			select {
			case t2 := <-ticker.C:
				stream <- t2
			case <-ctx.Done():
				close(stream)
				return
			}
		}
	}()

	return stream
}

func main() {
	// Setup rdb
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       "redis-master",
		SentinelAddrs:    []string{"redis-sentinel1:26379", "redis-sentinel2:26379", "redis-sentinel3:26379"},
		SentinelPassword: "sentinel_pass",
		Password:         "redis_pass",
	})

	ctx := context.Background()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("ping Error: %v", err)
	}

	fmt.Println("ping response:", pong)

	startTime := time.Now().Add(time.Millisecond * 10)
	delay := time.Second * 5 // 5 seconds

	for t := range cron(ctx, startTime, delay) {
		log.Println(t.Format("2006-01-02 15:04:05"))
		err := writeRedis(ctx, rdb)
		if err != nil {
			log.Printf("set rdb err: %v", err)
		}
	}

	//signalChan := make(chan os.Signal, 1)
	//signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM) //os.Interrupt, os.Kill
	//<-signalChan
}
