package shell

import "testing"

func TestInit(t *testing.T) {
	shells := []string{"bash", "zsh", "fish", "powershell", "pwsh", "nushell", "nu"}
	for _, shell := range shells {
		_, err := Init(shell)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
	init, err := Init("unknown")
	if init != "" {
		t.Errorf("expected empty string, got %s", init)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
