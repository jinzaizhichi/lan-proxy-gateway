package cmd

import (
	"strings"
	"testing"
)

func TestReleaseAssetName(t *testing.T) {
	cases := []struct {
		goos   string
		goarch string
		want   string
	}{
		{goos: "darwin", goarch: "arm64", want: "gateway-darwin-arm64"},
		{goos: "linux", goarch: "amd64", want: "gateway-linux-amd64"},
		{goos: "windows", goarch: "amd64", want: "gateway-windows-amd64.exe"},
	}

	for _, tc := range cases {
		if got := releaseAssetName(tc.goos, tc.goarch); got != tc.want {
			t.Fatalf("releaseAssetName(%q, %q) = %q, want %q", tc.goos, tc.goarch, got, tc.want)
		}
	}
}

func TestUpdateTempPattern(t *testing.T) {
	cases := []struct {
		goos string
		want string
	}{
		{goos: "darwin", want: "gateway-update-*"},
		{goos: "linux", want: "gateway-update-*"},
		{goos: "windows", want: "gateway-update-*.exe"},
	}

	for _, tc := range cases {
		if got := updateTempPattern(tc.goos); got != tc.want {
			t.Fatalf("updateTempPattern(%q) = %q, want %q", tc.goos, got, tc.want)
		}
	}
}

func TestBuildWindowsUpdateScript(t *testing.T) {
	script := buildWindowsUpdateScript(
		`C:\Program Files\gateway\gateway.exe`,
		`C:\Temp\gateway-update.exe`,
		`C:\Users\demo\gateway.yaml`,
		`C:\Users\demo\data`,
	)

	wants := []string{
		`set "TARGET=C:\Program Files\gateway\gateway.exe"`,
		`set "SOURCE=C:\Temp\gateway-update.exe"`,
		`"%TARGET%" start --config "%CONFIG%" --data-dir "%DATA%" >nul 2>&1 <nul`,
		`move /Y "%TARGET%" "%BACKUP%"`,
	}

	for _, want := range wants {
		if !strings.Contains(script, want) {
			t.Fatalf("expected script to contain %q, got:\n%s", want, script)
		}
	}
}

func TestEscapeWindowsBatchValue(t *testing.T) {
	got := escapeWindowsBatchValue(`C:\Users\%USERNAME%\gateway.exe`)
	want := `C:\Users\%%USERNAME%%\gateway.exe`
	if got != want {
		t.Fatalf("escapeWindowsBatchValue() = %q, want %q", got, want)
	}
}
