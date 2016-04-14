package helpers

import (
	"github.com/emergeadapt/caseblocks.helpers/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

type UserRepository interface {
	FindById(id FKInt) (User, error)
	FindByAuthToken(token string) (User, error)
	FindByDisplayName(accountCode, displayName string) (User, error)
	FindByLogin(accountCode, login string) (User, error)
	FindByEmail(accountCode, email string) (User, error)
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
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login, u.created_at, u.updated_at from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.id=$1"
	} else {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login, u.created_at, u.updated_at from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.id=?"
	}

	user := User{}
	err := r.DB.Get(&user, sql, id)
	return user, err
}

func (r *relationalUserRepository) FindByAuthToken(token string) (User, error) {
	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login, u.created_at, u.updated_at from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.authentication_token=$1"
	} else {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login, u.created_at, u.updated_at from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.authentication_token=?"
	}

	user := User{}
	err := r.DB.Get(&user, sql, token)
	return user, err
}

func (r *relationalUserRepository) FindByDisplayName(accountCode, name string) (User, error) {
	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login, u.created_at, u.updated_at from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.display_name=$1 and a.nickname=$2"
	} else {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login, u.created_at, u.updated_at from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.display_name=? and a.nickname=?"
	}

	user := User{}
	err := r.DB.Get(&user, sql, name, accountCode)
	return user, err
}

func (r *relationalUserRepository) FindByLogin(accountCode, login string) (User, error) {
	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login, u.created_at, u.updated_at from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.login=$1 and a.nickname=$2"
	} else {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login, u.created_at, u.updated_at from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.login=? and a.nickname=?"
	}

	user := User{}
	err := r.DB.Get(&user, sql, login, accountCode)
	return user, err
}

func (r *relationalUserRepository) FindByEmail(accountCode, email string) (User, error) {
	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login, u.created_at, u.updated_at from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.email=$1 and a.nickname=$2"
	} else {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login, u.created_at, u.updated_at from case_blocks_users u, case_blocks_accounts a where u.account_id=a.id and u.email=? and a.nickname=?"
	}

	user := User{}
	err := r.DB.Get(&user, sql, email, accountCode)
	return user, err
}
