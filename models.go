package helpers

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
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
	CreatedAt      time.Time `bson:"created_at" json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at" db:"updated_at"`
}

func (u User) String() string {
	return fmt.Sprintf("%4d|%s|%s|%s", u.Id, u.AccountCode, u.Email, u.Login)
}

type Account struct {
	Id                      FKInt
	Title                   string
	Nickname                string
	HomepageImageUrl        string    `bson:"homepage_image_url" json:"homepage_image_url" db:"homepage_image_url"`
	DefaultFromEmailAddress string    `bson:"default_email_from_address" json:"default_email_from_address" db:"default_email_from_address"`
	CreatedAt               time.Time `bson:"created_at" json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `bson:"updated_at" json:"updated_at" db:"updated_at"`
}

type Recipient struct {
	Id          FKInt
	Type        string
	DisplayName string `bson:"display_name"`
	Email       string
}

func (r *Recipient) UnmarshalJSON(data []byte) error {
	var aux struct {
		Id          string
		Type        string
		DisplayName string
		Email       string
	}
	fmt.Println("UNMARSHALLING RECIPIENT", string(data))
	err := json.Unmarshal(data, aux)
	id, err := strconv.Atoi(aux.Id)
	r.Id = FKInt(id)
	r.Type = aux.Type
	r.DisplayName = aux.DisplayName
	r.Email = r.Email
	return err
}

type Message struct {
	Id                bson.ObjectId `bson:"_id"`
	Body              string
	AuthorId          FKInt       `bson:"author_id"`
	AuthorDisplayName string      `bson:"author_display_name"`
	Recipients        []Recipient `bson:"recipients"`
	Subject           string
	CreatedAt         time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at" json:"updated_at"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var aux struct {
		Id                string
		Body              string
		AuthorId          string `json:"author_id"`
		AuthorDisplayName string `json:"author_display_name"`
		Recipients        []Recipient
		Subject           string
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
	}
	fmt.Println("UNMARSHALLING MESSAGE", string(data))
	err := json.Unmarshal(data, aux)
	m.Id = bson.ObjectIdHex(aux.Id)
	m.Body = aux.Body
	author_id, err := strconv.Atoi(aux.AuthorId)
	m.AuthorId = FKInt(author_id)
	m.AuthorDisplayName = aux.AuthorDisplayName
	m.Recipients = aux.Recipients
	m.Subject = aux.Subject
	m.CreatedAt = aux.CreatedAt
	m.UpdatedAt = aux.UpdatedAt
	return err
}

type CaseDocument struct {
	Id            bson.ObjectId `bson:"_id"`
	AccountId     FKInt         `bson:"account_id" json:"account_id"`
	AccountCode   string
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
