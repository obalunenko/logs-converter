package db

import (
	"fmt"

	"github.com/oleg-balunenko/logs-converter/db/mongo"
	"github.com/oleg-balunenko/logs-converter/model"
)

const ( // database types

	// Mongo - mongo db type
	Mongo = "mongo"
)

// Repository is a contract for databases
type Repository interface {
	Store(logModel *model.LogModel) (string, error)
	Update(id string, logModel model.LogModel) error
	Delete(id string) error
	Drop(bool) error
	Close()
}

// Connect establish connection to passed database type
func Connect(dbType string, url string, dbName string, colletionName string, username string, password string) (Repository, error) {
	switch dbType {
	case Mongo:
		return mongo.NewMongoDBConnection(url, dbName, colletionName, username, password), nil
	default:
		return Repository(nil), fmt.Errorf("not supported database type [%s]", dbType)

	}

}
