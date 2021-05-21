package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/mymmrac/project-glynn/pkg/client"
)

var cli struct {
	Join struct {
		Host   string `kong:"arg,required,help='Server host'"`
		RoomID string `kong:"arg,required,help='Room ID to connect'"`
	} `kong:"cmd,help='Connect ro room'"`
}

func main() {
	ctx := kong.Parse(&cli)

	switch ctx.Command() {
	case "join <host> <room-id>":
		fmt.Println("Connecting...")

		c := client.NewClient(cli.Join.Host)
		c.StartChat(cli.Join.RoomID)
	default:
		fmt.Printf("Unknown command: %q\n", ctx.Command())
		return
	}
}
