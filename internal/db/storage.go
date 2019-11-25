// Package db implements database interactions.
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

// Params is a database connection parameters.
type Params struct {
	URL        string
	DB         string
	Collection string
	Username   string
	Password   string
}

// Connect establish connection to passed database type
func Connect(dbType StorageType, params Params) (Repository, error) {
	switch dbType {
	case StorageTypeMongo:
		return mongo.NewMongoDBConnection(params.URL, params.DB, params.Collection, params.Username, params.Password)
	default:
		return nil, errors.Errorf("not supported database type [%s]", dbType)
	}
}
