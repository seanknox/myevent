package dblayer

import (
	"fmt"

	"github.com/seanknox/myevent/lib/persistence"
	"github.com/seanknox/myevent/lib/persistence/mongolayer"
)

type DBTYPE string

const (
	MONGODB  DBTYPE = "mongodb"
	DYNAMODB DBTYPE = "dynamodb"
)

func NewPersistenceLayer(options DBTYPE, connection string) (persistence.DatabaseHandler, error) {

	switch options {
	case MONGODB:
		fmt.Println("creating new mongo connection")
		return mongolayer.NewMongoDBLayer(connection)
	}
	return nil, nil
}
