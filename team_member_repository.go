package helpers

import "github.com/jmoiron/sqlx"

type TeamMemberRepository interface {
	FindByTeamId(accountCode string, id FKInt) ([]User, error)
}

func NewTeamMemberRepository(db *sqlx.DB) TeamMemberRepository {
	return &relationalTeamMemberRepository{db}
}

type relationalTeamMemberRepository struct {
	DB *sqlx.DB
}

func (r *relationalTeamMemberRepository) FindByTeamId(accountCode string, id FKInt) ([]User, error) {

	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login, u.created_at, u.updated_at, u.authentication_token from case_blocks_users u, case_blocks_accounts a, case_blocks_team_memberships tm where u.account_id=a.id and a.nickname=$1 and u.id=tm.user_id and tm.team_id=$2"
	} else {
		sql = "select u.id, u.account_id, u.is_account_admin, a.nickname, email, display_name, login, u.created_at, u.updated_at, u.authentication_token from case_blocks_users u, case_blocks_accounts a, case_blocks_team_memberships tm where u.account_id=a.id and a.nickname=? and u.id=tm.user_id and tm.team_id=?"
	}

	users := []User{}
	err := r.DB.Select(&users, sql, accountCode, id)
	return users, err
}
