package helpers

import (
	"errors"
	"fmt"

	"github.com/kr/pretty"
)

type MockedUserRepository struct {
	Users map[FKInt]User
}

func (r *MockedUserRepository) FindById(id FKInt) (User, error) {
	return r.Users[id], nil
}
func (r *MockedUserRepository) FindByAuthToken(token string) (User, error) {
	return User{}, fmt.Errorf("Not implemented.")
}
func (r *MockedUserRepository) FindByDisplayName(accountCode, displayName string) (User, error) {
	return User{}, fmt.Errorf("Not implemented.")
}
func (r *MockedUserRepository) FindByLogin(accountCode, login string) (User, error) {
	return User{}, fmt.Errorf("Not implemented.")
}
func (r *MockedUserRepository) FindByEmail(accountCode, email string) (User, error) {
	return User{}, fmt.Errorf("Not implemented.")
}
func (r *MockedUserRepository) AddUser(user User) error {
	if user.Id > 0 {
		r.Users[user.Id] = user
		return nil
	} else {
		return fmt.Errorf("Id > 0")
	}
}

func NewMockedUserRepository() *MockedUserRepository {
	return &MockedUserRepository{Users: make(map[FKInt]User)}
}

func Dump(msg string, obj interface{}) {
	fmt.Printf("%s\n%# v\n", msg, pretty.Formatter(obj))

	// fmt.Printf("%# v")
}

// doesn't handle multiple accounts yet
type MockedTeamMemberRepository struct {
	TeamUsers map[FKInt][]User
}

func NewMockedTeamMemberRepository() *MockedTeamMemberRepository {
	teamUsers := make(map[FKInt][]User)
	return &MockedTeamMemberRepository{TeamUsers: teamUsers}
}

func (r *MockedTeamMemberRepository) FindByTeamId(accountCode string, id FKInt) ([]User, error) {
	users, ok := r.TeamUsers[id]
	if !ok {
		return []User{}, errors.New("Invalid team id")
	}
	return users, nil
}

func (r *MockedTeamMemberRepository) AddTeamMember(accountCode string, id FKInt, user User) {
	users, ok := r.TeamUsers[id]
	if !ok {
		users = make([]User, 0)
	}
	users = append(users, user)
	r.TeamUsers[id] = users
}

func NewMockedCaseTypeRepository() *MockedCaseTypeRepository {
	return &MockedCaseTypeRepository{CaseTypeMap: make(map[FKInt]CaseType)}
}

type MockedCaseTypeRepository struct {
	CaseTypeMap map[FKInt]CaseType
}

func (r *MockedCaseTypeRepository) FindCaseTypeById(accountCode string, caseTypeId FKInt) (CaseType, error) {
	caseType, ok := r.CaseTypeMap[caseTypeId]
	if !ok {
		return caseType, errors.New("Unable to find case type")
	}
	return caseType, nil
}

func (r *MockedCaseTypeRepository) AddCaseType(caseType CaseType) {
	r.CaseTypeMap[caseType.Id] = caseType
}

type MockedWorkerbotRepository struct {
	Workerbot User
}

func (r *MockedWorkerbotRepository) FindWorkerbotByAccount(accountCode string) (User, error) {
	return r.Workerbot, nil
}
