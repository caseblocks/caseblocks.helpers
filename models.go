package helpers

import (
	"math"
	"time"
)

type TS struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
