package helpers

import (
	"github.com/emergeadapt/caseblocks.helpers/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

type UserRepository interface {
	FindById(id FKInt) (User, error)
	FindByAuthToken(token string) (User, error)
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &relationalUserRepository{db}
}

type relationalUserRepository struct {
	DB *sqlx.DB
}

func (r *relationalUserRepository) FindById(id FKInt) (User, error) {

	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.id=$1"
	} else {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.id=?"
	}

	user := User{}
	err := r.DB.Get(&user, sql, id)
	return user, err
}

func (r *relationalUserRepository) FindByAuthToken(token string) (User, error) {
	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.authentication_token=$1"
	} else {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.authentication_token=?"
	}

	user := User{}
	err := r.DB.Get(&user, sql, token)
	return user, err
}
