package helpers

import (
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
