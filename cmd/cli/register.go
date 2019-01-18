package main

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

var register = cli.Command{
	Name: "register",
	Action: func(c *cli.Context) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter username: ")
		username, _ := reader.ReadString('\n')

		fmt.Print("Enter password: ")
		password, _ := terminal.ReadPassword(int(syscall.Stdin))

		fmt.Println(string(username))
		fmt.Println(string(password))

		return nil
	},
}
