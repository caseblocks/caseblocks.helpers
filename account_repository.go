package helpers

import (
	"database/sql"
	"time"

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

type accountRecord struct {
	ID                      FKInt          `db:"id"`
	Title                   string         `db:"title"`
	Nickname                string         `db:"nickname"`
	HomepageImageURL        sql.NullString `db:"homepage_image_url"`
	DefaultFromEmailAddress string         `db:"default_email_from_address"`
	CreatedAt               time.Time      `db:"created_at"`
	UpdatedAt               time.Time      `db:"updated_at"`
	CipherKey               sql.NullString `db:"cipher_key"`
}

func (r *relationalAccountRepository) FindById(id FKInt) (Account, error) {

	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select id, title, nickname, created_at, updated_at, default_email_from_address, homepage_image_url, cipher_key from case_blocks_accounts where id=$1"
	} else {
		sql = "select id, title, nickname, created_at, updated_at, default_email_from_address, homepage_image_url, cipher_key from case_blocks_accounts where id=?"
	}

	var a accountRecord
	err := r.DB.Get(&a, sql, id)
	return Account{a.ID, a.Title, a.Nickname, a.HomepageImageURL.String, a.DefaultFromEmailAddress, a.CreatedAt, a.UpdatedAt, a.CipherKey.String}, err
}

func (r *relationalAccountRepository) FindByAccountCode(code string) (Account, error) {
	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select id, title, nickname, created_at, updated_at, default_email_from_address, homepage_image_url, cipher_key from case_blocks_accounts where nickname=$1"
	} else {
		sql = "select id, title, nickname, created_at, updated_at, default_email_from_address, homepage_image_url, cipher_key from case_blocks_accounts where nickname=?"
	}

	var a accountRecord
	err := r.DB.Get(&a, sql, code)
	return Account{a.ID, a.Title, a.Nickname, a.HomepageImageURL.String, a.DefaultFromEmailAddress, a.CreatedAt, a.UpdatedAt, a.CipherKey.String}, err
}
