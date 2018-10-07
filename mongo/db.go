package mongo

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	mgo "gopkg.in/mgo.v2"
)

// DB stores db connection details
type DB struct {
	session    *mgo.Session
	database   *mgo.Database
	collection *mgo.Collection
}

// NewConnection establishes connection with mongoDB and return MongoDB object
func NewConnection(url, dbName, collectionName, username, password string) *DB {

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{url},
		Timeout:  60 * time.Second,
		Database: dbName,
		Username: username,
		Password: password,
	}

	session, errDial := mgo.DialWithInfo(mongoDBDialInfo)
	if errDial != nil {
		log.Fatalf("Could not leave without Mongo: %v\nExiting...", errDial)

	}
	database := session.DB(dbName)
	collection := database.C(collectionName)

	return &DB{
		session:    session,
		database:   database,
		collection: collection,
	}

}

// CloseConnection closes mongo connection
func (db *DB) CloseConnection() {
	log.Infof("Closing connection...")
	db.session.Close()

}

// StoreModel insert model to collection
func (db *DB) StoreModel(model interface{}) error {
	log.Debugf("Storing model [%+v] to collection [%+v]", model, db.collection)

	errInsert := db.collection.Insert(model)
	if errInsert != nil {
		return fmt.Errorf("StoreModel: failed to store model %+v at [%v.%v]: %v", model, db.database.Name, db.collection.Name, errInsert)
	}
	log.Debugf("Successfully stored model [%+v]", model)
	return nil
}

// DropCollection drops database collection
func (db *DB) DropCollection() {
	collections, err := db.database.CollectionNames()
	if err != nil {
		log.Errorf("DropCollection: Failed to get list of all collections: %v", err)
		return
	}

	for _, c := range collections {
		if c == db.collection.Name {
			if errDrop := db.collection.DropCollection(); errDrop != nil {
				log.Fatalf("DropCollection: Failed to drop the collection [%s.%s]: %v", db.database.Name, db.collection.Name, errDrop)
			}
			return
		}

	}

	log.Warnf("DropCollection: Collection does not yet exist")

}

// DropDatabase drops database collection
func (db *DB) DropDatabase() {
	databases, err := db.session.DatabaseNames()
	if err != nil {
		log.Errorf("DropDatabase: Failed to get list of all databases: %v", err)
		return
	}
	for _, c := range databases {
		if c == db.database.Name {
			if errDrop := db.database.DropDatabase(); errDrop != nil {
				log.Fatalf("DropDatabase: Failed to drop the database [%s]: %v", db.database.Name, errDrop)
			}
			return
		}

	}

	log.Warnf("DropDatabase: database does not yet exist")

}
