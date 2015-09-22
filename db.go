package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func FindDBConnString() string {
	if os.Getenv("MYSQL_CONN") != "" {
		return os.Getenv("MYSQL_CONN")
	} else if os.Getenv("MYSQL_PORT_3306_TCP_ADDR") != "" && os.Getenv("MYSQL_PORT_3306_TCP_PORT") != "" {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_PORT_3306_TCP_ADDR"),
			os.Getenv("MYSQL_PORT_3306_TCP_PORT"),
			os.Getenv("MYSQL_DATABASE"))
	}
	return ""
}

func FindRedisConnString() string {
	if os.Getenv("REDIS_HOST") != "" {
		return os.Getenv("REDIS_HOST")
	} else if os.Getenv("REDIS_PORT_6379_TCP_ADDR") != "" && os.Getenv("REDIS_PORT_6379_TCP_PORT") != "" {
		return fmt.Sprintf("%s:%s",
			os.Getenv("REDIS_PORT_6379_TCP_ADDR"),
			os.Getenv("REDIS_PORT_6379_TCP_PORT"))
	}
	return ""
}

func FindMongoConnString() string {
	if os.Getenv("MONGO_HOST") != "" {
		return os.Getenv("MONGO_HOST")
	} else if os.Getenv("MONGODB_PORT_27017_TCP_ADDR") != "" && os.Getenv("MONGODB_PORT_27017_TCP_PORT") != "" {
		return fmt.Sprintf("%s:%s",
			os.Getenv("MONGODB_PORT_27017_TCP_ADDR"),
			os.Getenv("MONGODB_PORT_27017_TCP_PORT"))
	}
	return ""
}

func NewSqlConnection(connectionString string, logger Log) *sqlx.DB {

	if logger == nil {
		logger = NewConsoleLogger()
	}
	logger.Debug(connectionString)
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
