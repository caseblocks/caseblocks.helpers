package helpers

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func FindUserFromId(req *http.Request, res http.ResponseWriter, db *sqlx.DB) (User, error) {
	token := req.Header.Get("AUTH_TOKEN")
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
