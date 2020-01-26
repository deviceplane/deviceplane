package device

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/deviceplane/deviceplane/cmd/deviceplane/cliutils"
	"github.com/deviceplane/deviceplane/pkg/models"
	"golang.org/x/sync/errgroup"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func deviceListAction(c *kingpin.ParseContext) error {
	var filters []models.Filter
	for _, textFilter := range *deviceFilterListFlag {
		filter, err := parseTextFilter(textFilter)
		if err != nil {
			return err
		}

		filters = append(filters, filter)
	}

	devices, err := config.APIClient.ListDevices(context.TODO(), filters, *config.Flags.Project)
	if err != nil {
		return err
	}

	if *deviceOutputFlag == cliutils.FormatTable {
		table := cliutils.DefaultTable()
		table.SetHeader([]string{"Name", "Status", "IP", "OS", "Labels", "Last Seen", "Created"})
		for _, d := range devices {
			createdStr := cliutils.DurafmtSince(d.CreatedAt).String() + " ago"
			lastSeenStr := cliutils.DurafmtSince(d.LastSeenAt).String() + " ago"

			labelsArr := make([]string, len(d.Labels))
			i := 0
			for k, v := range d.Labels {
				labelsArr[i] = fmt.Sprintf("%s:%s", k, v)
				i += 1
			}
			labelsStr := strings.Join(labelsArr, "\n")

			table.Append([]string{
				d.Name,
				string(d.Status),
				d.Info.IPAddress,
				d.Info.OSRelease.Name,
				labelsStr,
				lastSeenStr,
				createdStr,
			})
		}
		table.Render()
		return nil
	}

	return cliutils.PrintWithFormat(devices, *deviceOutputFlag)
}

func deviceRebootAction(c *kingpin.ParseContext) error {
	err := config.APIClient.RebootDevice(context.TODO(), *config.Flags.Project, *deviceArg)
	if err != nil {
		return err
	}

	fmt.Println("Successfully initiated reboot")
	return nil
}

func deviceInspectAction(c *kingpin.ParseContext) error {
	device, err := config.APIClient.GetDevice(context.TODO(), *config.Flags.Project, *deviceArg)
	if err != nil {
		return err
	}

	return cliutils.PrintWithFormat(device, *deviceOutputFlag)
}

func deviceSSHAction(c *kingpin.ParseContext) error {
	conn, err := config.APIClient.InitiateSSH(context.TODO(), *config.Flags.Project, *deviceArg)
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
			"-o",
			"NoHostAuthenticationForLocalhost yes",
			"127.0.0.1",
			"-o",
			fmt.Sprintf("ConnectTimeout=%d", *sshTimeoutFlag),
		}, *sshCommandsArg...)

		cmd := exec.CommandContext(
			ctx,
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
}
