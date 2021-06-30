package cmd

import (
	"strconv"

	"github.com/oidc-proxy-ecosystem/proxy-server/config"
	"github.com/urfave/cli/v2"
)

var (
	write = &cli.Command{
		Name:        "write",
		Description: "write or generate command",
		Aliases:     []string{"w", "g"},
		Subcommands: []*cli.Command{
			caGenerate,
		},
	}
	run = &cli.Command{
		Name:        "run",
		Description: "run proxy server",
		Aliases:     []string{"r", "s"},
		Action: func(c *cli.Context) error {
			return runProxy()
		},
	}
	sshCmd = &cli.Command{
		Name:        "ssh",
		Description: "ssh client",
		Action: func(c *cli.Context) error {
			if v, err := newClientCertificate(config.File.Certificate); err != nil {
				return err
			} else {
				// if err := v.SaveFile(); err != nil {
				// 	return err
				// }
				return newSSH(c.String("host"), strconv.Itoa(c.Int("port")), c.String("user"), v.GetRsaPrivateKey(), v.GetCertificate())
			}
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"H"},
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   22,
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
			},
		},
	}
)

func Command() *cli.App {
	app := cli.NewApp()
	app.Name = "proxy server"
	app.Commands = []*cli.Command{
		write,
		run,
		sshCmd,
	}
	return app
}
