package scanner

import (
	"testing"
)

func TestResolveAndValidateRejectsPrivateIPs(t *testing.T) {
	t.Parallel()

	cases := []string{"127.0.0.1", "10.0.0.5", "192.168.1.10", "::1"}
	for _, input := range cases {
		input := input
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if _, err := resolveAndValidate(input); err == nil {
				t.Fatalf("expected %s to be rejected", input)
			}
		})
	}
}

func TestPrepareTargetRejectsUnsupportedScheme(t *testing.T) {
	t.Parallel()

	if _, err := PrepareTarget("ftp://example.com"); err == nil {
		t.Fatal("expected unsupported scheme to fail")
	}
}
