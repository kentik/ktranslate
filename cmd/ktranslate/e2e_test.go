package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// buildStaticBinary compiles the real ktranslate entrypoint the same way every
// shipped build path does post-furious-removal (Makefile's all/windows/arm
// targets, Dockerfile): CGO_ENABLED=0. It exists to catch "the binary doesn't
// even start" regressions, not to exercise features.
func buildStaticBinary(t *testing.T) string {
	t.Helper()

	if out, err := exec.Command("go", "generate", "github.com/kentik/ktranslate/pkg/version").CombinedOutput(); err != nil {
		t.Fatalf("go generate ./pkg/version failed: %v\n%s", err, out)
	}

	bin := filepath.Join(t.TempDir(), "ktranslate")
	cmd := exec.Command("go", "build", "-o", bin, "github.com/kentik/ktranslate/cmd/ktranslate")
	cmd.Env = append(os.Environ(),
		"CGO_ENABLED=0",
		"GOOS="+runtime.GOOS,
		"GOARCH="+runtime.GOARCH,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("static build failed: %v\n%s", err, out)
	}

	return bin
}

// TestStaticBuildRunsAndPrintsUsage proves the statically-built binary
// actually executes, rather than just compiling. -h is handled entirely by
// the standard flag package (main.go defines no -h flag of its own), so a
// successful run here also confirms flag registration didn't break.
func TestStaticBuildRunsAndPrintsUsage(t *testing.T) {
	bin := buildStaticBinary(t)

	cmd := exec.Command(bin, "-h")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s -h failed: %v\noutput:\n%s", bin, err, out)
	}

	got := string(out)
	if !strings.Contains(got, "Usage of") {
		t.Errorf("expected usage banner, got:\n%s", got)
	}
	for _, flagName := range []string{"-listen", "-mapping", "-snmp", "-sinks"} {
		if !strings.Contains(got, flagName) {
			t.Errorf("expected %q in usage output, got:\n%s", flagName, got)
		}
	}
}

// TestStaticBuildHasNoDynamicLinkage confirms CGO_ENABLED=0 actually produced
// a statically linked binary. That's the entire point of removing the
// furious -> gopacket/pcap cgo dependency: a Go binary with no cgo is fully
// static on Linux (ldd/file report no shared library dependencies), which is
// only meaningful to check on Linux -- Darwin and Windows binaries always
// link against OS-provided libraries regardless of CGO_ENABLED.
func TestStaticBuildHasNoDynamicLinkage(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("static linkage is only meaningful/verifiable on linux")
	}

	bin := buildStaticBinary(t)

	out, err := exec.Command("file", bin).CombinedOutput()
	if err != nil {
		t.Fatalf("file failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "statically linked") {
		t.Errorf("expected a statically linked binary, got: %s", out)
	}
}
