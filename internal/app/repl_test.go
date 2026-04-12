package app

import (
	"strings"
	"testing"
)

func TestExecuteCommand_empty(t *testing.T) {
	if got := ExecuteCommand(""); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
}

func TestExecuteCommand_echo(t *testing.T) {
	if got := ExecuteCommand("echo a b"); got != "a b" {
		t.Fatalf("got %q", got)
	}
}

func TestExecuteCommand_unknown(t *testing.T) {
	got := ExecuteCommand("notacommand")
	if !strings.Contains(got, "unknown command") {
		t.Fatalf("got %q", got)
	}
}

func TestExecuteCommand_helpSorted(t *testing.T) {
	out := ExecuteCommand("help")
	idxEcho := strings.Index(out, "\n  echo")
	idxHelp := strings.Index(out, "\n  help")
	if idxEcho < 0 || idxHelp < 0 {
		t.Fatalf("unexpected help output:\n%s", out)
	}
	if idxEcho > idxHelp {
		t.Fatalf("commands should be sorted (echo before help):\n%s", out)
	}
}
