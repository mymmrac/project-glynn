package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/mymmrac/project-glynn/pkg/client"
)

var cli struct {
	Host string `kong:"required,help='Server host'"`

	Join struct {
		RoomID string `kong:"arg,required,help='Room ID to connect'"`
	} `kong:"cmd,help='Connect ro room'"`

	CreateUser struct {
		Username string `kong:"arg,required,help='Username of your user'"`
	} `kong:"cmd,help='Create new user'"`
}

func main() {
	ctx := kong.Parse(&cli)

	switch ctx.Command() {
	case "join <room-id>":
		fmt.Println("Connecting...")

		c := client.NewClient(cli.Host)
		c.StartChat(cli.Join.RoomID)
	case "create-user <username>":
		fmt.Println("Creating user...")

		c := client.NewClient(cli.Host)
		c.CreateUser(cli.CreateUser.Username)
	default:
		fmt.Printf("Unknown command: %q\n", ctx.Command())
		return
	}
}
