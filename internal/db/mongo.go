// Package db implements database interactions.
package db

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/oleg-balunenko/logs-converter/internal/models"
)

// mongoDB stores mongo mongoDB connection details
type mongoDB struct {
	session    *mgo.Session
	database   *mgo.Database
	collection *mgo.Collection
}

// newMongoDBConnection establishes connection with mongoDB and return DBName object
func newMongoDBConnection(url, dbName, collectionName, username, password string) (*mongoDB, error) {
	var timeout = 60 * time.Second

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{url},
		Timeout:  timeout,
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

	return &mongoDB{
		session:    session,
		database:   database,
		collection: collection,
	}, nil
}

// Store stores model in database with unique id
// return id and error
func (db *mongoDB) Store(model *models.LogModel) (string, error) {
	log.Debugf("Storing model [%+v] to collection [%+v]", model, db.collection)

	id := bson.NewObjectId().Hex()
	model.ID = id

	if err := db.collection.Insert(model); err != nil {
		return "", errors.Wrap(err, "failed to insert model")
	}

	log.Debugf("Successfully stored model [%+v]", model)

	return model.ID, nil
}

// Delete deletes model from mongoDB by id
func (db *mongoDB) Delete(id string) error {
	return db.collection.RemoveId(id)
}

// Update updates existed model by id
func (db *mongoDB) Update(id string, logModel models.LogModel) error {
	return db.collection.UpdateId(id, logModel)
}

// Close closes mongo connection
func (db *mongoDB) Close() {
	log.Infof("Closing connection...")
	db.session.Close()
}

// Drop drops database collection
func (db *mongoDB) Drop() error {
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
