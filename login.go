package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/protectwise/stork/helper"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const NONCE_LEN = 64
const EC2_ROLE_URL = "http://169.254.169.254/2016-09-02/meta-data/iam/info"
const EC2_PKCS7_URL = "http://169.254.169.254/2016-09-02/dynamic/instance-identity/pkcs7"

func login_to_vault(c *cli.Context) error {

	nonce_file := c.String("nonce")
	token_file := c.String("token")
	role := c.String("role")
	pkcs7 := c.String("pkcs7")
	vault_server := c.String("server")
	verbose := c.BoolT("verbose")

	if nonce_file == "" {
		fmt.Println("You must specify a nonce file.")
		os.Exit(100)
	}

	if token_file == "" {
		fmt.Println("You must specify a token file.")
		os.Exit(100)
	}

	if vault_server == "" {
		fmt.Println("You must specify a Vault server.")
		os.Exit(100)
	}

	nonce := strings.TrimSpace(manage_nonce(nonce_file))

	if role == "" {
		role = helper.EC2Role(EC2_ROLE_URL)
	}

	if pkcs7 == "" {
		pkcs7 = string(helper.HttpGet(EC2_PKCS7_URL))

		// remove new lines
		re := regexp.MustCompile(`\r?\n`)
		pkcs7 = re.ReplaceAllString(pkcs7, "")
	}

	if verbose {
		fmt.Printf("Vault Server: %s\n", vault_server)
		fmt.Printf("Nonce: %s\n", nonce)
		fmt.Printf("Role: %s\n", role)
		fmt.Printf("PKCS7: %s...\n", pkcs7[0:32])
	}

	login_response := helper.Get_token_from_vault(fmt.Sprintf("%s/v1/auth/aws-ec2/login", vault_server), pkcs7, role, nonce)

	token, err := jsonparser.GetString(login_response, "auth", "client_token")
	helper.Fatal_error(fmt.Sprintf("Could not parse the json from Vault server: %s", string(login_response)), err)
	write_file(token_file, token)

	if verbose {
		fmt.Printf("Vault Token: %s\n", token)
	}

	return nil
}

func manage_nonce(nonce_file string) string {

	fileInfo, err := os.Stat(nonce_file)

	// if the file doesn't exist, no sense in checking the size
	if err == nil {
		if fileInfo.Size() > 0 {
			nonce, err := ioutil.ReadFile(nonce_file)
			helper.Fatal_error(fmt.Sprintf("Couldn't open nonce file %s", nonce_file), err)
			return string(nonce)
		}
	}

	nonce := helper.GenerateRandomString(NONCE_LEN)

	write_file(nonce_file, nonce)

	return nonce
}

func write_file(file string, content string) {
	err := ioutil.WriteFile(file, []byte(content), 0660)
	helper.Fatal_error(fmt.Sprintf("Could not write file %s", file), err)
}
