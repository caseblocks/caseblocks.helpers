package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CryptRequest struct {
	Plaintext        string `json:"plaintext"`
	Base64Ciphertext string `json:"base64ciphertext"`
}

func Decrypt(base64ciphertext string, user User) (string, error) {
	var byteArr []byte
	var err error
	var resp *http.Response
	var cryptRequest CryptRequest

	if byteArr, err = json.Marshal(CryptRequest{"", base64ciphertext}); err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/case_blocks/crypt/decrypt?auth_token=%s", FindCryptAPIEndpoint(), user.AuthToken)

	if resp, err = http.Post(url, "application/json", bytes.NewBuffer(byteArr)); err != nil {
		return "", err
	}

	if byteArr, err = ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	}

	if err = json.Unmarshal(byteArr, &cryptRequest); err != nil {
		return "", err
	}

	return cryptRequest.Plaintext, nil
}
