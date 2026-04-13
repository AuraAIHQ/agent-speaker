package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jason/agent-speaker/internal/identity"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "agent-speaker",
		Usage: "A nostr-based agent communication CLI",
		Commands: []*cli.Command{
			{
				Name:  "identity",
				Usage: "Manage identities",
				Commands: []*cli.Command{
					{
						Name:  "create",
						Usage: "Create a new identity",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "nickname", Required: true},
							&cli.BoolFlag{Name: "default"},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							mgr, err := identity.NewManager()
							if err != nil {
								return err
							}
							defer mgr.Close()
							
							id, err := mgr.Create(c.String("nickname"))
							if err != nil {
								return err
							}
							
							if c.Bool("default") {
								mgr.SetDefault(c.String("nickname"))
							}
							
							fmt.Printf("Created identity: %s\n", id.Nickname)
							fmt.Printf("Npub: %s\n", id.Npub)
							return nil
						},
					},
					{
						Name:  "list",
						Usage: "List identities",
						Action: func(ctx context.Context, c *cli.Command) error {
							mgr, err := identity.NewManager()
							if err != nil {
								return err
							}
							defer mgr.Close()
							
							identities, err := mgr.List()
							if err != nil {
								return err
							}
							
							for _, id := range identities {
								fmt.Printf("- %s (%s...)\n", id.Nickname, id.Npub[:20])
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "msg",
				Usage: "Send a message",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "from", Required: true},
					&cli.StringFlag{Name: "to", Required: true},
					&cli.StringFlag{Name: "content", Required: true},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					// 简化版本，仅演示结构
					fmt.Printf("Message from %s to %s: %s\n", 
						c.String("from"), 
						c.String("to"),
						c.String("content"))
					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
