package cmd

import "testing"

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
