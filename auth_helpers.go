package helpers

import (
	"fmt"
	"net/http"

	"strings"

	"github.com/go-martini/martini"
	"github.com/jmoiron/sqlx"
)

func FindUserInAccount(userId, accountId int, db *sqlx.DB) (User, error) {
	var sql string
	if db.DriverName() == "postgres" {
		sql = "select id, account_id, is_account_admin from case_blocks_users where account_id=$1 and id=$2"
	} else {
		sql = "select id, account_id, is_account_admin from case_blocks_users where account_id=? and id=?"
	}
	user := User{}
	getUserErr := db.Get(&user, sql, accountId, userId)
	return user, getUserErr
}

func FindUserFromId(req *http.Request, res http.ResponseWriter, db *sqlx.DB) (User, error) {

	var token string
	for headerKey, headerVals := range req.Header {
		if strings.ToLower(headerKey) == "auth_token" {
			token = headerVals[0]
		}
	}
	req.ParseForm()
	for formKey, formVals := range req.Form {
		if strings.ToLower(formKey) == "auth_token" {
			token = formVals[0]
		}
	}
	var user User
	var getUserErr error
	if token != "" {

		user, getUserErr = NewUserRepository(db).FindByAuthToken(token)
		if getUserErr != nil {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
			return user, getUserErr
		}
	} else {
		userId, err := findUserIdInRequest(req)
		if err != nil {
			fmt.Println("Authentication token not found, no session cookie in request.")
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
			return user, err
		}

		user, getUserErr = NewUserRepository(db).FindById(FKInt(userId))
		if getUserErr != nil {
			fmt.Printf("User not found: %d - %s\n", userId, getUserErr.Error())
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

func AllowAccessControl(res http.ResponseWriter) {
	// headers := res.Header()
	// headers.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Token")
	// headers.Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	// headers.Set("Access-Control-Allow-Origin", "*")
	// headers.Set("Access-Control-Max-Age", "1728000")
}

func RequireAdmin(req *http.Request, res http.ResponseWriter, db *sqlx.DB, c martini.Context) {
	if user, err := FindUserFromId(req, res, db); err == nil && user.IsAccountAdmin {
		c.Map(user)
	} else {
		http.Redirect(res, req, "/users/sign_in", http.StatusFound)
	}
}
