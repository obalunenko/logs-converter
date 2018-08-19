package mongo

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/oleg.balunenko/logs-converter/config"

	mgo "gopkg.in/mgo.v2"
)

// Connect establish connection with mongoDB and return mgo.Collection
func Connect(cfg *config.Config) *mgo.Collection {

	session, errDial := mgo.Dial(cfg.MongoURL)
	if errDial != nil {
		log.Fatalf("Could not leave without Mongo: %v\nExiting...", errDial)

	}

	// Collection
	collection := session.DB(cfg.MongoDB).C(cfg.MongoCollection)

	return collection

}

// CloseConnection closes mongo connection
func CloseConnection(collection *mgo.Collection) {
	log.Infof("Closing connection...")
	collection.Database.Session.Close()

}

// StoreModel insert model to collection
func StoreModel(model interface{}, collection *mgo.Collection) error {
	log.Debugf("Storing model [%+v] to collection [%+v]", model, collection)
	errInsert := collection.Insert(model)
	if errInsert != nil {
		return fmt.Errorf("failed to store model %+v at [%v.%v]: %v", model, collection.Database, collection.Name, errInsert)
	}
	log.Debugf("Successfully stored model [%+v]", model)
	return nil
}

// DropDBCollection drops database collection
func DropDBCollection(collection *mgo.Collection) {

	if errDrop := collection.DropCollection(); errDrop != nil {
		log.Fatalf("Failed to drop the collection [%+v.%+v]:%v", collection, collection.Database, errDrop)
	}
}
