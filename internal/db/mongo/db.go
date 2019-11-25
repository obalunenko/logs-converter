// Package mongo implements mongo db interactions functionality.
package mongo

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/oleg-balunenko/logs-converter/internal/models"
)

// db stores mongo db connection details
type db struct {
	session    *mgo.Session
	database   *mgo.Database
	collection *mgo.Collection
}

// NewMongoDBConnection establishes connection with mongoDB and return DBName object
func NewMongoDBConnection(url, dbName, collectionName, username, password string) (*db, error) {
	return newConnection(url, dbName, collectionName, username, password)
}

func newConnection(url, dbName, collectionName, username, password string) (*db, error) {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{url},
		Timeout:  60 * time.Second,
		Database: dbName,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return nil, err
	}

	database := session.DB(dbName)
	collection := database.C(collectionName)

	return &db{
		session:    session,
		database:   database,
		collection: collection,
	}, nil
}

// Store stores model in database with unique id
// return id and error
func (db *db) Store(model *models.LogModel) (string, error) {
	log.Debugf("Storing model [%+v] to collection [%+v]", model, db.collection)

	id := bson.NewObjectId().Hex()
	model.ID = id

	if err := db.collection.Insert(model); err != nil {
		return "", errors.Wrap(err, "failed to insert model")
	}

	log.Debugf("Successfully stored model [%+v]", model)

	return model.ID, nil
}

// Delete deletes model from db by id
func (db *db) Delete(id string) error {
	return db.collection.RemoveId(id)
}

// Update updates existed model by id
func (db *db) Update(id string, logModel models.LogModel) error {
	return db.collection.UpdateId(id, logModel)
}

// Close closes mongo connection
func (db *db) Close() {
	log.Infof("Closing connection...")
	db.session.Close()
}

// Drop drops database collection
func (db *db) Drop(shouldDrop bool) error {
	if !shouldDrop {
		return nil
	}

	databases, err := db.session.DatabaseNames()
	if err != nil {
		return errors.Wrap(err, "failed to list databases")
	}

	for _, c := range databases {
		if c == db.database.Name {
			if err = db.database.DropDatabase(); err != nil {
				return errors.Wrapf(err, "failed to drop database [%s]", db.database.Name)
			}

			break
		}
	}

	return nil
}
