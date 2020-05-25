package main

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"

	// This is called blank import.
	// We need this for our app to understand that it needs to connect to PostgreSQL.
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Custom struct that holds the logic for the broker.
// Notice that it has only unexported fields - if we need something from here,
// we'd create a small helper function.
type Broker struct {
	postgresDBSetup pgSetup
	postgresDB      *gorm.DB
}

// Create and return our broker.
// This is ready for some extra things - I do a lot more in here in big, production-grade APIs
func NewBroker() Broker {
	b := Broker{}
	return b
}

// Neat way of storing all needed variables for connecting to psql
type pgSetup struct {
	username string
	password string
	dbName   string
	dbHost   string
}

// And our cherry on top - this struct will define our schema
// and will be our aid when unmarshaling JSON payloads.
type Post struct {
	// `gorm:"type:uuid"` helps psql digest the type of this field
	Id uuid.UUID `gorm:"type:uuid"`
	// `json:"title"` helps Go JSON engine to understand what's coming in and out
	// This way, if the JSON payload name changes, we just need to change it here.
	// No need to do a massive changes in the application.
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// InitializeBroker sets up the connection for database.
// Again, not much here, but it's better to have the logic ready before
// rather than wondering how to make the app more flexible in prod.
// (Please, don't figure out stuff in prod...)
func (b *Broker) InitializeBroker() error {
	err := b.setPostgres()
	if err != nil {
		return err
	}

	return nil
}

// Getter for our db connection.
// You'll see a lot of it in handlers.go
func (b *Broker) GetPostgres() *gorm.DB {
	return b.postgresDB
}

func (b *Broker) SetPostgresConfig(username, password, dbName, dbHost string) {
	pgs := pgSetup{
		username: username,
		password: password,
		dbName:   dbName,
		dbHost:   dbHost,
	}
	b.postgresDBSetup = pgs
}

func (b *Broker) setPostgres() error {
	// Create a connection string - I don't like concatenation that much
	// so when I don't have to use it, I don't.
	// That's why I'm using fmt.Sprintf to print the string for me and store it in dbUri
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", b.postgresDBSetup.dbHost, b.postgresDBSetup.username, b.postgresDBSetup.dbName, b.postgresDBSetup.password)
	// Connect to the database.
	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		return err
	}
	b.postgresDB = conn
	// See all the logs from database.
	// Might be nice to switch it off in production.
	b.postgresDB.LogMode(true)
	// My favourite part - this bit will automatically generate the schemas in our database
	// and if there's any changes - will do the changes as well!
	b.postgresDB.Debug().AutoMigrate(
		&Post{},
	)
	return nil
}
