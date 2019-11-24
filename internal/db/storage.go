package db

import (
	"github.com/pkg/errors"

	"github.com/oleg-balunenko/logs-converter/internal/db/mongo"
	"github.com/oleg-balunenko/logs-converter/internal/models"
)

// StorageType is a storage type.
//go:generate stringer -type=StorageType -trimprefix=StorageType
type StorageType uint

const ( // database types
	storageTypeUnknown StorageType = iota

	// StorageTypeMongo - mongo db type
	StorageTypeMongo

	storageTypeSentinel // should be always last
)

// Valid checks if storage type is in a valid value range.
func (i StorageType) Valid() bool {
	return i > storageTypeUnknown && i < storageTypeSentinel
}

// Repository is a contract for databases
type Repository interface {
	Store(logModel *models.LogModel) (string, error)
	Update(id string, logModel models.LogModel) error
	Delete(id string) error
	Drop(bool) error
	Close()
}

// Connect establish connection to passed database type
func Connect(dbType StorageType, url string, dbName string, colletionName string, username string,
	password string) (Repository, error) {
	switch dbType {
	case StorageTypeMongo:
		return mongo.NewMongoDBConnection(url, dbName, colletionName, username, password)
	default:
		return nil, errors.Errorf("not supported database type [%s]", dbType)
	}
}
