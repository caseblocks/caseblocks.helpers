package helpers

import (
	"errors"
	"net/http"

	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	salt         = "encrypted cookie" // default value for Rails 4 app
	key_iter_num = 1000
	key_size     = 64
	dummy        = 1
)

func generateSecret(base, salt string) []byte {
	return pbkdf2.Key([]byte(base), []byte(salt), key_iter_num, key_size, sha1.New)
}

func decodeCookieData(cookie []byte) (data, iv []byte, err error) {
	vectors := strings.SplitN(string(cookie), "--", 2)

	data, err = base64.StdEncoding.DecodeString(vectors[0])
	if err != nil {
		return
	}

	iv, err = base64.StdEncoding.DecodeString(vectors[1])
	if err != nil {
		return
	}

	return
}

func decryptCookie(cookie []byte, secret []byte) (dd []byte, err error) {
	data, iv, err := decodeCookieData(cookie)

	c, err := aes.NewCipher(secret[:32])
	if err != nil {
		return
	}

	cfb := cipher.NewCBCDecrypter(c, iv)
	dd = make([]byte, len(data))
	cfb.CryptBlocks(dd, data)

	return
}

func DecryptSignedCookie(signed_cookie, secret_key_base, salt string) (session []byte, err error) {
	cookie, err := url.QueryUnescape(signed_cookie)
	if err != nil {
		return
	}

	vectors := strings.SplitN(cookie, "--", 2)
	data, err := base64.StdEncoding.DecodeString(vectors[0])
	if err != nil {
		return
	}

	session, err = decryptCookie(data, generateSecret(secret_key_base, salt))
	if err != nil {
		return
	}

	return
}

func getAuthUserId(decrypted_session_data []byte) (user_id int64, err error) {
	userKeyRegex := regexp.MustCompile(`warden.user.user.key":\[\[(\d+)`)
	session_data := string(decrypted_session_data)
	matches := userKeyRegex.FindStringSubmatch(session_data)
	userIdInt, err := strconv.Atoi(matches[1])
	user_id = int64(userIdInt)
	return user_id, err
}

func findUserIdInRequest(req *http.Request) (int64, error) {
	var user_id int64
	found := false
	cookie, err := req.Cookie("_caseblocks_session")
	// Dump("user cookie", cookie)
	if err != nil {
		return 0, errors.New("User cookie not found")
	}

	secret_key_base := os.Getenv("SECRET_KEY_BASE")
	// Dump("SECRET_KEY_BASE", secret_key_base)
	if len(secret_key_base) == 0 {
		return 0, errors.New("Please specify SECRET_KEY_BASE envvar to share Rails sessions.")
	}

	decrypted_session_data, err := DecryptSignedCookie(cookie.Value, secret_key_base, salt)
	// Dump("decrypted_session_data222", string(decrypted_session_data))
	if err != nil {
		return 0, errors.New("Unable to decrypt session data.")
	}
	user_id, err = getAuthUserId(decrypted_session_data)
	// Dump("user_id", user_id)
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
