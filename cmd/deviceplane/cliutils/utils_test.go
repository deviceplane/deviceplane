package cliutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSSHParsing(t *testing.T) {
	preSSH, postSSH := GetSSHArgs([]string{
		"deviceplane",
		"device",
		"ssh",
		"elegant-lamarr",
		"echo",
		"-L",
		"3000:localhost:3000",
	})
	require.Equal(t, preSSH, []string{
		"deviceplane",
		"device",
		"ssh",
		"elegant-lamarr",
	})
	require.Equal(t, postSSH, []string{
		"echo",
		"-L",
		"3000:localhost:3000",
	})
}

func TestSSHParsingWithoutPostSSH(t *testing.T) {
	preSSH, postSSH := GetSSHArgs([]string{
		"deviceplane",
		"device",
		"ssh",
		"elegant-lamarr",
	})
	require.Equal(t, preSSH, []string{
		"deviceplane",
		"device",
		"ssh",
		"elegant-lamarr",
	})
	require.Len(t, postSSH, 0)
}

func TestSSHParsingWithSinglePostSSH(t *testing.T) {
	preSSH, postSSH := GetSSHArgs([]string{
		"deviceplane",
		"device",
		"ssh",
		"elegant-lamarr",
		"ls",
	})
	require.Equal(t, preSSH, []string{
		"deviceplane",
		"device",
		"ssh",
		"elegant-lamarr",
	})
	require.Equal(t, postSSH, []string{
		"ls",
	})
}

func TestSSHParsingWithoutDevice(t *testing.T) {
	preSSH, postSSH := GetSSHArgs([]string{
		"deviceplane",
		"device",
		"ssh",
	})
	require.Equal(t, preSSH, []string{
		"deviceplane",
		"device",
		"ssh",
	})
	require.Len(t, postSSH, 0)
}
func TestSSHParsingWithoutSSH(t *testing.T) {
	preSSH, postSSH := GetSSHArgs([]string{
		"deviceplane",
		"device",
		"list",
	})
	require.Equal(t, preSSH, []string{
		"deviceplane",
		"device",
		"list",
	})
	require.Len(t, postSSH, 0)
}
