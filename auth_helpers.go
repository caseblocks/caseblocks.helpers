package helpers

import (
	"fmt"
	"net/http"

	"strings"

	"github.com/go-martini/martini"
	"github.com/jmoiron/sqlx"
)

func FindUserFromId(req *http.Request, res http.ResponseWriter, db *sqlx.DB) (User, error) {

	var token string
	for headerKey, headerVals := range req.Header {
		if strings.ToLower(headerKey) == "auth_token" {
			token = headerVals[0]
		}
	}
	req.ParseForm()
	for formKey, formVals := range req.Form {
		fmt.Println(formKey)
		if strings.ToLower(formKey) == "auth_token" {
			token = formVals[0]
		}
	}
	user := User{}
	if token != "" {
		getUserErr := db.Get(&user, "select id, account_id, is_account_admin from case_blocks_users where authentication_token=?", token)
		if getUserErr != nil {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
			return user, getUserErr
		}
	} else {
		userId, err := findUserIdInRequest(req)
		if err != nil {
			fmt.Println("user id not found")
			fmt.Println(userId)
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
			return user, err
		}
		getUserErr := db.Get(&user, "select id, account_id, is_account_admin from case_blocks_users where id=?", userId)
		if getUserErr != nil {
			fmt.Println("user not found")
			fmt.Println(userId)
			fmt.Println(getUserErr)
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
			return user, err
		}
	}
	return user, nil
}

func RequireLogin(req *http.Request, res http.ResponseWriter, db *sqlx.DB, c martini.Context) {
	if user, err := FindUserFromId(req, res, db); err == nil {
		c.Map(user)
	} else {
		http.Redirect(res, req, "/users/sign_in", http.StatusFound)
	}
}

func RequireAdmin(req *http.Request, res http.ResponseWriter, db *sqlx.DB, c martini.Context) {
	if user, err := FindUserFromId(req, res, db); err == nil && user.IsAccountAdmin {
		c.Map(user)
	} else {
		http.Redirect(res, req, "/users/sign_in", http.StatusFound)
	}
}
