package mongo

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/oleg-balunenko/logs-converter/model"
)

const prefix = "db/mongo"

// DB stores mongo db connection details
type DB struct {
	Session    *mgo.Session
	database   *mgo.Database
	collection *mgo.Collection
}

// NewMongoDBConnection establishes connection with mongoDB and return DBName object
func NewMongoDBConnection(url, dbName, collectionName, username, password string) *DB {

	return newConnection(url, dbName, collectionName, username, password)

}

func newConnection(url, dbName, collectionName, username, password string) *DB {
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
		Session:    session,
		database:   database,
		collection: collection,
	}

}

// Store stores model in database with unique id
// return id and error
func (db *DB) Store(model *model.LogModel) (string, error) {
	log.Debugf("Storing model [%+v] to collection [%+v]", model, db.collection)

	id := bson.NewObjectId().Hex()
	model.ID = id
	errInsert := db.collection.Insert(model)
	if errInsert != nil {
		err := fmt.Errorf("failed to store model %+v at [%v.%v]: %v",
			model, db.database.Name, db.collection.Name, errInsert)
		return "", errors.Wrap(err, prefix+": Store")
	}

	log.Debugf("Successfully stored model [%+v]", model)

	return model.ID, nil
}

// Delete deletes model from db by id
func (db *DB) Delete(id string) error {
	panic("implement me")
}

// Update updates existed model by id
func (db *DB) Update(id string, logModel model.LogModel) error {
	panic("implement me")
}

// Close closes mongo connection
func (db *DB) Close() {
	log.Infof("Closing connection...")
	db.Session.Close()

}

// Drop drops database collection
func (db *DB) Drop(shouldDrop bool) error {
	databases, err := db.Session.DatabaseNames()
	if err != nil {
		err = fmt.Errorf("failed to get list of all databases: %v", err)
		return errors.Wrap(err, prefix+": Drop")
	}
	for _, c := range databases {
		if c == db.database.Name {
			if err = db.database.DropDatabase(); err != nil {
				err = fmt.Errorf("failed to drop the database [%s]: %v", db.database.Name, err)
				return errors.Wrap(err, prefix+": Drop")
			}
			break
		}

	}
	return nil
}
