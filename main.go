package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "stork"
	app.Usage = "get connected with our vault installation"

	app.Commands = []cli.Command{
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "login to vault via EC2 authentication",
			Action: func(c *cli.Context) error {
				return login_to_vault(c)
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "verbose",
					Usage: "Verbose mode",
				},
				cli.StringFlag{
					Name:  "nonce",
					Usage: "A file for storing the Vault nonce",
				},
				cli.StringFlag{
					Name:  "token",
					Usage: "A file for storing the Vault token",
				},
				cli.StringFlag{
					Name:  "role",
					Usage: "Role to authenticate to Vault as. If unset, we will try to guess from IAM data.",
				},
				cli.StringFlag{
					Name:  "pkcs7",
					Usage: "You probably don't need to set this, we will get it from IAM data.",
				},
				cli.StringFlag{
					Name:  "server",
					Usage: "URL to Vault server",
				},
			},
		},
	}

	app.Run(os.Args)
}
