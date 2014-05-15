package helpers

import (
	"errors"
	"fmt"
	"github.com/adeven/gorails/marshal"
	"github.com/adeven/gorails/session"
	"net/http"
)

const (
	secret_key_base = "04ca87cb0ed96b33364b4a4aeace525a1a5d98b4e5852e32a2b3848b2191e334a0d776f19fd66ccbdab91f4aec45d037084d65e7afe9cf05e1a578e31a91d76f" // can be found in config/initializers/secret_token.rb
	salt            = "encrypted cookie"                                                                                                                 // default value for Rails 4 app
)

func getRailsSessionData(session_cookie string) (decrypted_cookie_data []byte, err error) {
	decrypted_cookie_data, err = session.DecryptSignedCookie(session_cookie, secret_key_base, salt)

	return
}

func getAuthUserId(decrypted_session_data []byte) (user_id int64, err error) {
	unauthorized_user := errors.New("Unauthorized user")
	invalid_auth_data := errors.New("Invalid auth data")

	session_data, err := marshal.CreateMarshalledObject(decrypted_session_data).GetAsMap()
	if err != nil {
		return
	}

	warden_data, ok := session_data["warden.user.user.key"]
	if !ok {
		return 0, unauthorized_user
	}

	warden_user_key, err := warden_data.GetAsArray()
	if err != nil {
		return
	}
	if len(warden_user_key) < 1 {
		return 0, invalid_auth_data
	}

	user_data, err := warden_user_key[0].GetAsArray()
	if err != nil {
		return
	}
	if len(user_data) < 1 {
		return 0, invalid_auth_data
	}

	user_id, err = user_data[0].GetAsInteger()

	return
}

func findUserIdInRequest(req *http.Request) (int64, error) {
	var user_id int64
	found := false
	cookie, err := req.Cookie("_caseblocks_session")
	for _, cookie := range req.Cookies() {
		fmt.Println(cookie.Name)
	}
	if err != nil {
		return 0, errors.New("User cookie not found")
	}
	decrypted_session_data, err := getRailsSessionData(cookie.Value)
	if err != nil {
		return 0, errors.New("Unable to decrypt session data.")
	}
	user_id, err = getAuthUserId(decrypted_session_data)
	if err != nil {
		return 0, errors.New("Unable to unmarshall data.")
	}
	found = true

	if found {
		return user_id, nil
	} else {
		return 0, errors.New("User Id not found.")
	}
}
