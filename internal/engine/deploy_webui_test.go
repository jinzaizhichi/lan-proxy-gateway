package engine

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDeployWebUIProducesReachableIndex verifies that the webui is extracted
// to <workdir>/ui with its files intact. Regression test for the Windows
// filepath.Rel bug where files ended up outside the workdir and `/ui/`
// served nothing.
func TestDeployWebUIProducesReachableIndex(t *testing.T) {
	workdir := t.TempDir()
	if err := deployWebUI(workdir); err != nil {
		t.Fatalf("deployWebUI: %v", err)
	}

	uiDir := filepath.Join(workdir, "ui")
	info, err := os.Stat(uiDir)
	if err != nil {
		t.Fatalf("ui dir missing: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("ui path is not a directory")
	}

	// Must have the SPA entry point at the root of /ui.
	for _, name := range []string{"index.html", "favicon.svg"} {
		p := filepath.Join(uiDir, name)
		fi, err := os.Stat(p)
		if err != nil {
			t.Errorf("%s missing after deploy: %v", name, err)
			continue
		}
		if fi.Size() == 0 {
			t.Errorf("%s is empty after deploy", name)
		}
	}

	// Must have the _nuxt asset subdirectory populated (index.html
	// references files under it; empty _nuxt = blank page in browser).
	nuxt := filepath.Join(uiDir, "_nuxt")
	if info, err := os.Stat(nuxt); err != nil || !info.IsDir() {
		t.Fatalf("_nuxt dir missing or not a directory: %v", err)
	}
	entries, err := os.ReadDir(nuxt)
	if err != nil {
		t.Fatalf("read _nuxt: %v", err)
	}
	if len(entries) == 0 {
		t.Error("_nuxt was deployed but empty — files leaked outside workdir?")
	}

	// Nothing should have escaped the workdir. If filepath.Rel had
	// misbehaved, sibling dirs might have appeared next to workdir.
	parent := filepath.Dir(workdir)
	siblings, err := os.ReadDir(parent)
	if err != nil {
		t.Fatalf("read workdir parent: %v", err)
	}
	for _, e := range siblings {
		if e.Name() == "webui" || e.Name() == "_nuxt" {
			t.Errorf("deploy leaked %q next to workdir — filepath handling broke", e.Name())
		}
	}
}
