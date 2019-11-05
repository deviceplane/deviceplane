package main

import (
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"

	"github.com/deviceplane/deviceplane/pkg/client"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
)

var ssh = cli.Command{
	Name: "ssh",
	Flags: []cli.Flag{
		projectFlag,
		deviceFlag,
	},
	Action: func(c *cli.Context) error {
		return withClient(c, func(client *client.Client) error {
			project := c.String("project")
			device := c.String("device")

			conn, err := client.InitiateSSH(context.TODO(), project, device)
			if err != nil {
				return err
			}

			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				return err
			}
			defer listener.Close()

			g, ctx := errgroup.WithContext(context.TODO())

			g.Go(func() error {
				localConn, err := listener.Accept()
				if err != nil {
					return err
				}

				go io.Copy(conn, localConn)
				io.Copy(localConn, conn)

				return nil
			})

			g.Go(func() error {
				defer conn.Close()

				port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

				sshArguments := append([]string{
					"-p", port,
					"127.0.0.1",
				}, c.Args()...)

				cmd := exec.CommandContext(ctx,
					"ssh",
					sshArguments...,
				)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					if exitError, ok := err.(*exec.ExitError); ok {
						os.Exit(exitError.ExitCode())
						return nil
					}
					return err
				}

				return nil
			})

			return g.Wait()
		})
	},
}
