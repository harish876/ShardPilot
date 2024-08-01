package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/exp/rand"
)

var dbMap map[int]*sql.DB

func getConnStr(shardId int) string {
	switch shardId {
	case 1:
		return fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			host1,
			port1,
			dbname1,
			user,
			password,
		)
	case 2:
		return fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			host2,
			port2,
			dbname2,
			user,
			password,
		)
	case 3:
		return fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			host3,
			port3,
			dbname3,
			user,
			password,
		)
	}
	return ""
}

func Init() error {
	rand.Seed(uint64(time.Now().UnixNano()))
	dbMap = make(map[int]*sql.DB)

	for i := 1; i <= 3; i++ {
		connStr := getConnStr(i)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			return fmt.Errorf("error connecting to database: %v", err)
		}

		db.SetMaxOpenConns(10)                  // Max open connections to the database
		db.SetMaxIdleConns(5)                   // Max idle connections in the pool
		db.SetConnMaxIdleTime(30 * time.Second) // Max time a connection can remain idle
		db.SetConnMaxLifetime(1 * time.Hour)    // Max lifetime of a connection

		dbMap[i] = db
	}
	return nil
}

func GetConnectionPoolForShard(shardId uint32) (*sql.DB, error) {
	if _, ok := dbMap[int(shardId)]; !ok {
		return nil, fmt.Errorf("shardId pool not available")
	}
	return dbMap[int(shardId)], nil
}
