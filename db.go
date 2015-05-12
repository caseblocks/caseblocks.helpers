package helpers

import (
	"github.com/jmoiron/sqlx"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func NewSqlConnection(connectionString string, logger Log) *sqlx.DB {

	if logger == nil {
		logger = NewConsoleLogger()
	}

	db, err := sqlx.Connect("mysql", connectionString)
	PanicToLogIf(err, logger)
	return db
}

func NewMongoSession(connectionString string, logger Log) *mgo.Session {
	if logger == nil {
		logger = NewConsoleLogger()
	}
	session, err := mgo.Dial(connectionString)
	PanicToLogIf(err, logger)
	return session
}

func ConvertStringToBSON(id string) bson.ObjectId {
	return bson.ObjectIdHex(id)
}

func ConvertStringsToBSONs(ids []string) []bson.ObjectId {
	results := make([]bson.ObjectId, len(ids))
	for i, id := range ids {
		results[i] = ConvertStringToBSON(id)
	}
	return results
}
