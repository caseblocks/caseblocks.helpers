package helpers

import (
	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	FindById(id FKInt) (Account, error)
	FindByAccountCode(code string) (Account, error)
}

func NewAccountRepository(db *sqlx.DB) AccountRepository {
	return &relationalAccountRepository{db}
}

type relationalAccountRepository struct {
	DB *sqlx.DB
}

func (r *relationalAccountRepository) FindById(id FKInt) (Account, error) {

	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select * from case_blocks_accounts where id=$1"
	} else {
		sql = "select * from case_blocks_accounts where id=?"
	}

	account := Account{}
	err := r.DB.Get(&account, sql, id)
	return account, err
}

func (r *relationalAccountRepository) FindByAccountCode(code string) (Account, error) {
	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select * from case_blocks_accounts where nickname=$1"
	} else {
		sql = "select * from case_blocks_accounts where nickname=?"
	}

	account := Account{}
	err := r.DB.Get(&account, sql, code)
	return account, err
}
