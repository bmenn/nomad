package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAgentMembersCommand_Implements(t *testing.T) {
	var _ cli.Command = &AgentMembersCommand{}
}

func TestAgentMembersCommand_Run(t *testing.T) {
	agent, http, client, url := testAgent(t)
	defer agent.Shutdown()
	defer http.Shutdown()

	ui := new(cli.MockUi)
	cmd := &AgentMembersCommand{Ui: ui}

	// Get our own node name
	name, err := client.Agent().NodeName()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Query the members
	if code := cmd.Run([]string{"-http-addr=" + url}); code != 0 {
		t.Fatalf("expected exit 0, got: %d", code)
	}
	if out := ui.OutputWriter.String(); !strings.Contains(out, name) {
		t.Fatalf("expected %q in output, got: %s", name, out)
	}
	ui.OutputWriter.Reset()

	// Query members with detailed output
	if code := cmd.Run([]string{"-http-addr=" + url, "-detailed"}); code != 0 {
		t.Fatalf("expected exit 0, got: %d", code)
	}
	if out := ui.OutputWriter.String(); !strings.Contains(out, "Tags") {
		t.Fatalf("expected tags in output, got: %s", out)
	}
}

func TestMembersCommand_Fails(t *testing.T) {
	ui := new(cli.MockUi)
	cmd := &AgentMembersCommand{Ui: ui}

	// Fails on misuse
	if code := cmd.Run([]string{"some", "bad", "args"}); code != 1 {
		t.Fatalf("expected exit code 1, got: %d", code)
	}
	if out := ui.ErrorWriter.String(); !strings.Contains(out, cmd.Help()) {
		t.Fatalf("expected help output, got: %s", out)
	}

	// Fails on connection failure
	if code := cmd.Run([]string{"-http-addr=nope"}); code != 1 {
		t.Fatalf("expected exit code 1, got: %d", code)
	}
	if out := ui.ErrorWriter.String(); !strings.Contains(out, "Failed querying members") {
		t.Fatalf("expected failed query error, got: %s", out)
	}
}
