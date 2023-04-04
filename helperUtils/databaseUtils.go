package helperUtils

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // add this
	"log"
	"os"
	"time"
)

type DatabaseHandler struct {
	conn *sql.DB
}

func newDatabaseHandler(conn *sql.DB) *DatabaseHandler {
	return &DatabaseHandler{
		conn: conn,
	}
}

func InitDB(config *AppConfig) *sql.DB {

	hostFromEnv := os.Getenv("POSTGRESQL_SERVICE_PORT_5432_TCP_ADDR")

	dbConnectionInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.Postgres.Username, config.Postgres.Password, hostFromEnv, config.Postgres.Port, config.Postgres.Database)

	var conn *sql.DB
	var err error

	if conn, err = sql.Open("postgres", dbConnectionInfo); err != nil {
		log.Print("Error: while connecting to database. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		if conn, err = sql.Open("postgres", dbConnectionInfo); err != nil {
			log.Panicf("Error: while connecting to database. Error: %v", err)
		}
	}

	if err = conn.Ping(); err != nil {
		log.Panicf("There was an error verifying the database connection: %v", err)
	}

	statement := `CREATE TABLE IF NOT EXISTS image_metadata (id serial PRIMARY KEY UNIQUE NOT NULL, image_id VARCHAR(255) UNIQUE NOT NULL, image_uri VARCHAR(255) NOT NULL, created_date TIMESTAMP NOT NULL, vote_date TIMESTAMP, votes INT);`

	if _, err = conn.Exec(statement); err != nil {
		log.Panic(err)
	}
	log.Printf("Database connected\n")

	return conn
}

func WriteMetadata(config *AppConfig, conn *sql.DB, filename string) {
	imageId := uuid.New()
	imageURI := fmt.Sprintf("%v%v", config.DataSources.Buckets.Shoes, filename)
	createdDateTime := time.Now().Format(time.RFC3339)

	statement := `INSERT INTO image_metadata (image_id, image_uri, created_date) VALUES ($1, $2, $3 ) RETURNING id`

	id := 0
	err := conn.QueryRow(statement, imageId, imageURI, createdDateTime).Scan(&id)
	if err != nil {
		log.Panicf("Error writing to the database: %v", err)
	}
	log.Println(fmt.Sprintf("New record ID is: %v and url is %v", id, imageURI))

}
