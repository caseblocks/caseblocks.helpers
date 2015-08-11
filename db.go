package helpers

import (
	"time"

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

func NowTS() time.Time {
	return time.Now().UTC()
}

type IArray []interface{}
type IMap map[string]interface{}

func BSONSafeElem(param interface{}) interface{} {
	// fmt.Printf("Element: %s\n", reflect.TypeOf(param))
	switch w := param.(type) {
	case []interface{}:
		return BSONSafeArray(w)
	case bson.M:
		return BSONSafeMap(w)
	case bson.ObjectId:
		return w.Hex()
	case bson.Symbol:
		return string(w)
	}
	return param
}

func BSONSafeMap(params bson.M) IMap {
	results := make(IMap)
	for k, v := range params {
		results[k] = BSONSafeElem(v)
	}
	return results
}

func BSONSafeArray(params IArray) IArray {
	l := len(params)
	results := make(IArray, l)
	for i, elem := range params {
		results[i] = BSONSafeElem(elem)
	}
	return results
}

func BSONSafe(param interface{}) interface{} {
	return BSONSafeElem(param)
}
