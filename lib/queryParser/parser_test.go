package queryparser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetShardID1(t *testing.T) {
	input := []byte(`select * from users where shardId = 83310`)
	shardId, err := GetShardID(input)

	if err != nil {
		t.Fatalf("Error creating query plan %v", err)
	}

	assert.Equal(t, shardId, []string{"83310"})
	fmt.Println(shardId)
}

func TestGetShardID2(t *testing.T) {
	input := []byte(`select * from users where shardId = 83310 and userId = 123`)
	shardId, err := GetShardID(input)

	if err != nil {
		t.Fatalf("Error creating query plan %v", err)
	}

	assert.Equal(t, shardId, []string{"83310"})
	fmt.Println(shardId)
}
