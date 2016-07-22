package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TS struct {
	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
}

type Context struct {
  CurrentUser   User
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
	AuthToken      string    `json:"-" db:"authentication_token"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at" db:"updated_at"`
}

func (u User) String() string {
	return fmt.Sprintf("%4d|%s|%s|%s", u.Id, u.AccountCode, u.Email, u.Login)
}

type Team struct {
	Id                      FKInt     `db:"id" json:"id"`
	AccountId               FKInt     `db:"account_id" json:"account_id"`
	AccountCode             string    `db:"nickname" json:"account_code"`
	DisplayName             string    `db:"display_name" json:"display_name"`
	CreatedAt               time.Time `bson:"created_at" json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `bson:"updated_at" json:"updated_at" db:"updated_at"`
	ExcludeFromDistribution bool      `db:"exclude_from_distribution" json:"exclude_from_distribution"`
	IncludeInDistribution   bool      `db:"include_in_distribution"`
	TeamScreenEnabled       bool      `db:"team_screen_enabled"`
}

type Account struct {
	Id                      FKInt
	Title                   string
	Nickname                string
	HomepageImageUrl        sql.NullString `bson:"homepage_image_url" json:"homepage_image_url" db:"homepage_image_url"`
	DefaultFromEmailAddress string         `bson:"default_email_from_address" json:"default_email_from_address" db:"default_email_from_address"`
	CreatedAt               time.Time      `bson:"created_at" json:"created_at" db:"created_at"`
	UpdatedAt               time.Time      `bson:"updated_at" json:"updated_at" db:"updated_at"`
	CipherKey               string         `json:"-" db:"cipher_key"`
}

type Recipient struct {
	Id          FKInt  `json:"id"`
	Type        string `json:"type"`
	DisplayName string `bson:"display_name" json:"display_name"`
	Email       string `json:"email"`
}

func (r *Recipient) UnmarshalJSON(data []byte) error {
	var aux struct {
		Id          interface{} `json:"id"`
		Type        string      `json:"type"`
		DisplayName string      `json:"name"`
		Email       string      `json:"email"`
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	switch v := aux.Id.(type) {
	case string:
		if id, err := strconv.Atoi(v); err != nil {
			return err
		} else {
			r.Id = FKInt(id)
		}
	case float64:
		r.Id = FKInt(v)
	case int64:
		r.Id = FKInt(v)
	}
	r.Type = aux.Type
	r.DisplayName = aux.DisplayName
	r.Email = aux.Email
	return err
}

type Message struct {
	Id                bson.ObjectId `bson:"_id" json:"_id"`
	Body              string        `json:"body"`
	CaseId            bson.ObjectId `bson:"case_id" json:"case_id"`
	ConversationId    bson.ObjectId `bson:"conversation_id" json:"conversation_id"`
	AuthorId          FKInt         `bson:"author_id" json:"author_id"`
	AuthorDisplayName string        `bson:"author_display_name" json:"author_display_name"`
	Recipients        []Recipient   `bson:"recipients" json:"recipients"`
	Subject           string        `json:"subject"`
	Attachments       []string      `json:"attachments"`
	CreatedAt         time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time     `bson:"updated_at" json:"updated_at"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var aux struct {
		Id                string
		Body              string        `json:"body"`
		CaseId            bson.ObjectId `json:"case_id"`
		ConversationId    bson.ObjectId `json:"conversation_id"`
		AuthorId          string        `json:"author_id"`
		AuthorDisplayName string        `json:"author_display_name"`
		Recipients        []Recipient   `json:"recipients"`
		Subject           string        `json:"subject"`
		Attachments       []string      `json:"attachments"`
		CreatedAt         time.Time     `json:"created_at"`
		UpdatedAt         time.Time     `json:"updated_at"`
	}
	err := json.Unmarshal(data, &aux)
	if len(aux.Id) > 0 {
		m.Id = bson.ObjectIdHex(aux.Id)
	}
	m.Body = aux.Body
	author_id, err := strconv.Atoi(aux.AuthorId)
	m.CaseId = aux.CaseId
	m.ConversationId = aux.ConversationId
	m.AuthorId = FKInt(author_id)
	m.AuthorDisplayName = aux.AuthorDisplayName
	m.Recipients = aux.Recipients
	m.Subject = aux.Subject
	m.CreatedAt = aux.CreatedAt
	m.UpdatedAt = aux.UpdatedAt
	m.Attachments = aux.Attachments
	return err
}

type CaseDocument struct {
	Id            bson.ObjectId `bson:"_id"`
	AccountId     FKInt         `bson:"account_id" json:"account_id"`
	CaseTypeId    FKInt         `bson:"case_type_id" json:"case_type_id"`
	AccountCode   string
	Title         string
	Conversations []Conversation
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

type Conversation struct {
	Id        bson.ObjectId `bson:"_id" json:"_id"`
	Subject   string        `json:"subject" db:"subject"`
	CaseId    bson.ObjectId `bson:"case_id" json:"case_id"`
	Messages  []Message     `json:"messages"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"  db:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"  db:"updated_at`
}

type CaseType struct {
	Id             FKInt
	AccountId      FKInt  `db:"account_id"`
	Name           string `db:"name"`
	SystemCategory string `db:"system_category"`
	Schemas        []string
}

func (ct CaseType) UriForCaseId(caseId string) string {

	var systemCategory string
	switch ct.SystemCategory {
	case "P":
		systemCategory = "people_types"
	case "O":
		systemCategory = "organization_types"
	default:
		systemCategory = "case_types"
	}

	return fmt.Sprintf("/#/%s/%d/case/%s/detail", systemCategory, ct.Id, caseId)
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
