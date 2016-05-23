package helpers

import (
	"github.com/jmoiron/sqlx"
)

type TeamRepository interface {
	FindById(accountCode string, id FKInt) (Team, error)
	FindByDisplayName(accountCode, displayName string) (Team, error)
}

func NewTeamRepository(db *sqlx.DB) TeamRepository {
	return &relationalTeamRepository{db}
}

type relationalTeamRepository struct {
	DB *sqlx.DB
}

func (r *relationalTeamRepository) FindById(accountCode string, id FKInt) (Team, error) {

	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select t.id, t.account_id, a.nickname, t.display_name, t.created_at, t.updated_at, t.exclude_from_distribution, t.include_in_distribution, t.team_screen_enabled from case_blocks_teams t, case_blocks_accounts a where t.account_id=a.id and t.id=$1 and a.nickname=$2"
	} else {
		sql = "select t.id, t.account_id, a.nickname, t.display_name, t.created_at, t.updated_at, t.exclude_from_distribution, t.include_in_distribution, t.team_screen_enabled from case_blocks_teams t, case_blocks_accounts a where t.account_id=a.id and t.id=? and a.nickname=?"
	}

	team := Team{}
	err := r.DB.Get(&team, sql, id, accountCode)
	return team, err
}

func (r *relationalTeamRepository) FindByDisplayName(accountCode, displayName string) (Team, error) {

	var sql string
	if r.DB.DriverName() == "postgres" {
		sql = "select t.id, t.account_id, a.nickname, t.display_name, t.created_at, t.updated_at, t.exclude_from_distribution, t.include_in_distribution, t.team_screen_enabled from case_blocks_teams t, case_blocks_accounts a where t.account_id=a.id and t.display_name=$1 and a.nickname=$2"
	} else {
		sql = "select t.id, t.account_id, a.nickname, t.display_name, t.created_at, t.updated_at, t.exclude_from_distribution, t.include_in_distribution, t.team_screen_enabled from case_blocks_teams t, case_blocks_accounts a where t.account_id=a.id and t.display_name=? and a.nickname=?"
	}

	team := Team{}
	err := r.DB.Get(&team, sql, displayName, accountCode)
	return team, err
}
