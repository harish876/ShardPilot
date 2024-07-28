package hash

import (
	"fmt"
	"testing"
)

func TestCalculateShardId(t *testing.T) {
	numberOfUsers := 250000
	numberOfShards := 3

	shards := make(map[uint32]int)
	for userId := 1; userId <= numberOfUsers; userId++ {
		shardId, err := CalculateShardId(IntToBytes(userId), numberOfShards)
		if err != nil {
			t.Fatal(err)
		}
		shards[shardId]++
	}

	fmt.Println("Shard Distribution", shards)
}
