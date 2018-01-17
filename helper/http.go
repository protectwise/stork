package helper

import (
	"bytes"
	"encoding/json"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Token struct {
	Token string
}

type AuthRequest struct {
	Role  string `json:"role"`
	Pkcs7 string `json:"pkcs7"`
	Nonce string `json:"nonce"`
}

// get something from somewhere
func HttpGet(url string) []byte {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// convert to a string and return the body
	return body
}

// get the role from the iam profile
// url should generally be: http://169.254.169.254/2016-09-02/meta-data/iam/info
func EC2Role(url string) string {

	role_arn, err := jsonparser.GetString(HttpGet(url), "InstanceProfileArn")
	if err != nil {
		log.Fatal("Could not get the role from EC2 metadata service. " + err.Error())
	}

	return strings.Split(role_arn, "/")[1]
}

// gets something from the vault api and returns the string
// url: vault-server/v1/auth/token/lookup-self
func (t *Token) VaultAPIGet(url string) string {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("X-VAULT-TOKEN", t.Token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

// POSTS the url vault-server/auth/aws-ec2/login
// with identity_doc
func Get_token_from_vault(vault_url string, pkcs7 string, role string, nonce string) []byte {

	request := AuthRequest{
		Role:  role,
		Pkcs7: pkcs7,
		Nonce: nonce,
	}

	post, err := json.Marshal(request)
	Fatal_error("Could not form request to Vault server", err)

	resp, err := http.Post(vault_url, "application/json", bytes.NewBuffer(post))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// return the body
	return body
}
