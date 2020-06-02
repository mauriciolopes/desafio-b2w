package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB represents a database connection.
type DB struct {
	db *mongo.Database
}

// New creates a DB.
func New() *DB {
	return &DB{}
}

// Open tests mongo database connection and opens it.
func (d *DB) Open(uri, dbName string) error {
	if d.db != nil {
		return nil
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	log.Println("Trying to connect to the database...")
	var dbOK bool
	for !dbOK {
		ok := func() bool {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = client.Ping(ctx, nil)
			if err != nil {
				log.Println(err)
				time.Sleep(1 * time.Second)
				return false
			}
			return true
		}()
		if ok {
			dbOK = true
			log.Println("Database OK")
		}
	}
	d.db = client.Database(dbName)
	return nil
}

// Get gets database.
// Call this function after Open. Otherwise, nil is returned
// and a log warning message is printed.
func (d *DB) Get() *mongo.Database {
	if d.db == nil {
		log.Println("No open connection. Open function was called?")
	}
	return d.db
}

// Close closes mongo database connection.
func (d *DB) Close() (err error) {
	err = d.db.Client().Disconnect(context.TODO())
	if err != nil {
		return
	}
	d.db = nil
	return
}
