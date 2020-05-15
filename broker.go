package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type Broker struct {
	postgresDBSetup pgSetup
	postgresDB      *gorm.DB
}

func NewBroker() Broker {
	b := Broker{}
	return b
}

type pgSetup struct {
	username string
	password string
	dbName   string
	dbHost   string
}

type Post struct {
	Id        uuid.UUID `gorm:"type:uuid"`
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Broker) InitializeBroker() error {
	if b.postgresDBSetup.username == "" {
		err := errors.New("Please set up username")
		return err
	}
	if b.postgresDBSetup.password == "" {
		err := errors.New("Please set up password")
		return err
	}
	if b.postgresDBSetup.dbName == "" {
		err := errors.New("Please set up dbName")
		return err
	}
	if b.postgresDBSetup.dbHost == "" {
		err := errors.New("Please set up dbHost")
		return err
	}

	err := b.setPostgres()
	if err != nil {
		return err
	}

	return nil
}

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
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", b.postgresDBSetup.dbHost, b.postgresDBSetup.username, b.postgresDBSetup.dbName, b.postgresDBSetup.password)
	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		return err
	}
	b.postgresDB = conn
	b.postgresDB.LogMode(true)
	b.postgresDB.Debug().AutoMigrate(
		&Post{},
	)
	return nil
}
