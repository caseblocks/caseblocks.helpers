package helpers

import "github.com/jmoiron/sqlx"

type CaseTypeRepository interface {
	FindCaseTypeById(accountCode string, caseTypeId FKInt) (CaseType, error)
}

func NewCaseTypeRepository(db *sqlx.DB) CaseTypeRepository {
	return &relationalCaseTypeRepository{db}
}

type relationalCaseTypeRepository struct {
	DB *sqlx.DB
}

func (r *relationalCaseTypeRepository) FindCaseTypeById(accountCode string, caseTypeId FKInt) (CaseType, error) {
	var result CaseType

	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select ct.id, ct.account_id, ct.name, ct.system_category from case_blocks_case_types ct, case_blocks_accounts a  where ct.account_id == a.id and ct.id = $1 and a.nickname = $2"
	} else {
		sql = "select ct.id, ct.account_id, ct.name, ct.system_category from case_blocks_case_types ct, case_blocks_accounts a  where ct.account_id == a.id and ct.id = ? and a.nickname = ?"
	}

	err := r.DB.Get(&result, sql, caseTypeId, accountCode)
	return result, err
}
