package helpers

import (
	"time"
)

type User struct {
	Id             int
	AccountId      int  `db:"account_id"`
	IsAccountAdmin bool `db:"is_account_admin"`
}

type CaseType struct {
	Id             int
	Name           string `db:"name"`
	SystemCategory string `db:"system_category"`
}

type Permission struct {
	Name           string `db:"name"`
	PermissionFlag string `db:"permission_flag"`
}

type TeamMembership struct {
	TeamId int `db:"team_id"`
	UserId int `db:"user_id"`
}

type Bucket struct {
	Id                      int
	Name                    string `db:"name"`
	CaseTypeId              int    `db:"case_type_id`
	Kpi                     string
	LastCheckedMembershipAt time.Time `db:"last_checked_membership_at"`
	LastCheckedTrippingAt   time.Time `db:"last_checked_membership_at"`
}
