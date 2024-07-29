package hash

import (
	"bytes"
	"encoding/binary"

	"github.com/spaolacci/murmur3"
)

func CalculateShardId(key []byte, numberOfShards int) (uint32, error) {
	hasher := murmur3.New32()
	_, err := hasher.Write(key)

	if err != nil {
		return 0, err
	}
	output := hasher.Sum32()
	shardId := output % uint32(numberOfShards)
	return shardId + 1, nil
}

func IntToBytes(n int) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int32(n))
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
