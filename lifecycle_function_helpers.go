package helpers

import (
	// "fmt"
	// "net/http"
  //
	// "strings"

	// "github.com/go-martini/martini"
	"github.com/jmoiron/sqlx"
)

// LifecycleFunction is a model to describe a javascript exectued function
type LifecycleFunction struct {
	ID                   int64         `db:"id"`
	Name                 string        `db:"name"`
	Description					 string				 `db:"description"`
	DraftCode					 string				 `db:"draft_code"`
	PublishedCode			 string				 `db:"published_code"`
	DraftDocument			 string        `db:"draft_document"`
	CaseTypeID				 int64				 `db:"case_type_id"`
}

// GetFunction retrieves a function from the database and loads it into struct LifecycleFunction
func GetFunction(functionName string, caseTypeName string, db *sqlx.DB) (function LifecycleFunction, err error) {
	// db := cb.NewSqlConnection(cb.FindDBConnString(), nil)
	var emptyFunction LifecycleFunction

  var sql = ""

  if db.DriverName() == "postgres" {
		sql = "select case_blocks_lifecycle_functions.id, case_blocks_lifecycle_functions.name, case_blocks_lifecycle_functions.description, case_blocks_lifecycle_functions.draft_code, published_code, draft_document, case_blocks_lifecycle_functions.case_type_id from case_blocks_lifecycle_functions join case_blocks_case_types on case_blocks_case_types.id = case_blocks_lifecycle_functions.case_type_id where case_blocks_case_types.code = $1 and case_blocks_lifecycle_functions.name=$1"
	} else {
		sql = "select case_blocks_lifecycle_functions.id, case_blocks_lifecycle_functions.name, case_blocks_lifecycle_functions.description, case_blocks_lifecycle_functions.draft_code, published_code, draft_document, case_blocks_lifecycle_functions.case_type_id from case_blocks_lifecycle_functions join case_blocks_case_types on case_blocks_case_types.id = case_blocks_lifecycle_functions.case_type_id where case_blocks_case_types.code = ? and case_blocks_lifecycle_functions.name=?"
	}
	dbErr := db.Get(&function, sql, caseTypeName, functionName)
	if dbErr != nil {
		return emptyFunction, dbErr
	}
	return function, nil
}
