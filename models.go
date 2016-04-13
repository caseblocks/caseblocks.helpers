package helpers

import (
	"math"
	"time"

	"github.com/emergeadapt/caseblocks.helpers/Godeps/_workspace/src/labix.org/v2/mgo/bson"
)

type TS struct {
	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
}

type FKInt int

type CaseTypeCode string

type User struct {
	Id             FKInt
	AccountId      FKInt  `db:"account_id"`
	IsAccountAdmin bool   `db:"is_account_admin"`
	AccountCode    string `db:"nickname"`
	Email          string
	DisplayName    string `db:"display_name"`
	Login          string
}

type Recipient struct {
	Id          FKInt
	Type        string
	DisplayName string `bson:"display_name"`
	Email       string
}

type Message struct {
	Id                bson.ObjectId `bson:"_id"`
	Body              string
	AuthorId          FKInt       `bson:"author_id"`
	AuthorDisplayName string      `bson:"author_display_name"`
	RecipientUsers    []Recipient `bson:"recipient_users"`
	RecipientTeams    []Recipient `bson:"recipient_teams"`
	Subject           string
	CreatedAt         time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at" json:"updated_at"`
}

type CaseDocument struct {
	Id            bson.ObjectId `bson:"_id"`
	Title         string
	Conversations []Conversation
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

type Conversation struct {
	Id        bson.ObjectId `bson:"_id"`
	Subject   string
	Messages  []Message
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type CaseType struct {
	Id             FKInt
	AccountId      FKInt  `db:"account_id"`
	Name           string `db:"name"`
	SystemCategory string `db:"system_category"`
	Schemas        []string
}

func (ct *CaseType) CurrentSchemaVersion() int {
	return len(ct.Schemas)
}

type Permission struct {
	Name           string `db:"name"`
	PermissionFlag string `db:"permission_flag"`
}

type TeamMembership struct {
	TeamId FKInt `db:"team_id"`
	UserId FKInt `db:"user_id"`
}

type Bucket struct {
	Id                      FKInt
	Name                    string `db:"name"`
	CaseTypeId              FKInt  `db:"case_type_id`
	Kpi                     string
	LastCheckedMembershipAt time.Time `db:"last_checked_membership_at"`
	LastCheckedTrippingAt   time.Time `db:"last_checked_membership_at"`
}

type Int64List []int64

type StringHistogram map[string]int64

func (sh StringHistogram) Occurences() Int64List {
	occs := make(Int64List, 0)
	for _, v := range sh {
		occs = append(occs, v)
	}
	return occs
}

func (i Int64List) Sum() int64 {
	var sum int64
	for _, ii := range i {
		sum += ii
	}
	return sum
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func RoundFloat64(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift
}

func GenerateStringHistogram(array []string) StringHistogram {
	hist := make(StringHistogram)
	for _, val := range array {
		hist[val] += 1
	}
	return hist
}

func UniqueStringArray(array []string) []string {
	results := make([]string, 0)
	for k, _ := range GenerateStringHistogram(array) {
		results = append(results, k)
	}
	return results
}

func ZipStringArrays(keys, vals []string) map[string]string {
	results := make(map[string]string)
	for i, _ := range keys {
		results[keys[i]] = vals[i]
	}
	return results
}
