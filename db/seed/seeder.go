package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/harish876/ShardPilot/lib/hash"
	_ "github.com/lib/pq"
	"golang.org/x/exp/rand"
)

var dbMap map[int]*sql.DB

const (
	user     = "shardPilot"
	password = "shardPilot@123"

	host1   = "localhost"
	dbname1 = "postgres"
	port1   = 5431

	host2   = "localhost"
	dbname2 = "postgres"
	port2   = 5432

	host3   = "localhost"
	dbname3 = "postgres"
	port3   = 5433

	NUMBER_OF_SHARDS = 3
	NUMBER_OF_USERS  = 2500
)

// todo move to config file
func getShard(shardId int) string {
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

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
	dbMap = make(map[int]*sql.DB)

	for i := 1; i <= 3; i++ {
		connStr := getShard(i)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

		db.SetMaxOpenConns(10)                  // Max open connections to the database
		db.SetMaxIdleConns(5)                   // Max idle connections in the pool
		db.SetConnMaxIdleTime(30 * time.Second) // Max time a connection can remain idle
		db.SetConnMaxLifetime(1 * time.Hour)    // Max lifetime of a connection

		dbMap[i] = db
	}
}

func main() {
	for userId := 1; userId <= NUMBER_OF_USERS; userId++ {
		shardId, _ := hash.CalculateShardId(hash.IntToBytes(userId), NUMBER_OF_SHARDS)
		shardDB, ok := dbMap[int(shardId)]
		if !ok {
			log.Fatalf("Error connecting to database")
		}
		defer shardDB.Close()

		name := GenerateRandomName()
		phoneNumber := fmt.Sprintf("+%d-%010d", rand.Intn(100), rand.Int63n(10000000000))
		emailAddr := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
		email := fmt.Sprintf("%s@example.com", emailAddr)

		_, err := shardDB.Exec(
			"INSERT INTO users (user_id, name, phone_number, email) VALUES ($1, $2, $3, $4)",
			userId,
			name,
			phoneNumber,
			email,
		)
		if err != nil {
			log.Printf("Error inserting user %d: %v", userId, err)
		}
	}
	fmt.Println("Seeder Executed Successfully")
}

var firstNames = []string{
	"Alice", "Bob", "Charlie", "David", "Eva", "Frank", "Grace", "Hannah", "Ivy", "Jack",
	"Karen", "Leo", "Mia", "Nate", "Olivia", "Paul", "Quinn", "Rita", "Sam", "Tina",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Jones", "Brown", "Davis", "Miller", "Wilson", "Moore", "Taylor",
	"Anderson", "Thomas", "Jackson", "White", "Harris", "Martin", "Thompson", "Garcia", "Martinez", "Robinson",
}

func GenerateRandomName() string {
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	return fmt.Sprintf("%s %s", firstName, lastName)
}
