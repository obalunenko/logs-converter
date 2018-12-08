package mongo

import (
	"fmt"
	"time"

	"github.com/oleg-balunenko/logs-converter/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DB stores mongo db connection details
type DB struct {
	Session    *mgo.Session
	database   *mgo.Database
	collection *mgo.Collection
}

// NewMongoDBConnection establishes connection with mongoDB and return MongoDB object
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
	res, errInsert := db.collection.Upsert(bson.M{"id": id}, model)
	if errInsert != nil {
		return "", fmt.Errorf("StoreModel: failed to store model %+v at [%v.%v]: %v", model, db.database.Name, db.collection.Name, errInsert)
	}

	log.Debugf("Successfully stored model [%+v]", model)

	return res.UpsertedId.(string), nil
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
		log.Errorf("DropDatabase: Failed to get list of all databases: %v", err)
		return errors.Wrap(err, "Failed to get list of all databases")
	}
	for _, c := range databases {
		if c == db.database.Name {
			if errDrop := db.database.DropDatabase(); errDrop != nil {
				return errors.Wrap(errDrop, "Failed to drop the database [%s]")
			}
			break
		}

	}
	return nil
}
