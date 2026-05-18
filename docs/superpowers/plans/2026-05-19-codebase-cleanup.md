# Codebase cleanup — implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Clear visible clutter from `lan-proxy-gateway` (local working-tree cruft, stale root scripts, undocumented `legacy/`, lint warnings in `cmd/update.go`) without changing any user-facing behavior.

**Architecture:** All work happens on a single branch `chore/housekeeping` as four atomic commits. Each commit is independently revertible. No source code semantics change — Task 7 (the only Go change) is a pure refactor protected by existing `cmd/update_test.go`. Compatibility wrappers preserve every external script entry point.

**Tech Stack:** Go 1.21+, Cobra CLI, GNU Make, Bash shell, git. Existing test runner is `go test`. No new dependencies.

**Spec:** `docs/superpowers/specs/2026-05-19-codebase-cleanup-design.md`

**Hard invariant (from spec):** No functional behavior change. `go build ./...` and `go test ./...` must be green after every commit; existing CLI entry points (`gateway start`, `gateway update`, `make install`, root `install.sh`, root `install.ps1`, root `dev.sh`, root `dev.ps1`) must continue to work bit-for-bit.

---

## File structure summary

After the branch lands, the working tree changes are:

- **Modified:** `Makefile` (expanded `clean:` target), `cmd/update.go` (constants + 2 extracted helpers).
- **New (committed):** `legacy/README.md`, `scripts/install-mihomo.sh`, `scripts/download-mihomo.sh`, `examples/extension.js`, root `install-mihomo.sh` (now a 4-line wrapper), root `download-mihomo.sh` (now a 4-line wrapper).
- **Removed (committed):** root `install-mihomo.sh` original content (replaced by wrapper), root `download-mihomo.sh` original content (replaced by wrapper), root `script-demo.js` (moved to `examples/`).
- **Touched (working-tree only, not in any commit):** gitignored cruft (`gateway` binary, `dist/`, `.DS_Store`, `.tmp/`, `.cache/`, `.try/`, `logs/*.log`).

Every Go file modified has tests in the same package; no new test files needed because all Go changes are pure refactor.

---

## Task 1: Create cleanup branch

**Files:** none yet — just git state.

- [ ] **Step 1.1: Confirm working tree clean and on main**

  Run: `git status && git rev-parse --abbrev-ref HEAD`
  Expected: `working tree clean` and `main`. If not clean, stop and ask the user.

- [ ] **Step 1.2: Create and check out the branch**

  Run: `git checkout -b chore/housekeeping`
  Expected: `Switched to a new branch 'chore/housekeeping'`.

- [ ] **Step 1.3: Verify branch state**

  Run: `git status`
  Expected: `On branch chore/housekeeping`, clean tree.

  (No commit in this task — it sets up the branch.)

---

## Task 2: Expand Makefile `clean` target (Commit 1)

**Files:**
- Modify: `Makefile` (the existing `clean:` target near the bottom)

**Spec section:** Commit 1 in the design.

- [ ] **Step 2.1: Read the current Makefile bottom**

  Run: `tail -15 Makefile`
  Expected: ends with a `clean:` rule whose body is `rm -rf $(DIST_DIR)/ $(BINARY)`.

- [ ] **Step 2.2: Replace the `clean` target**

  Edit `Makefile`. Replace the existing two lines:
  ```make
  clean:
  	rm -rf $(DIST_DIR)/ $(BINARY)
  ```
  with:
  ```make
  clean:
  	rm -rf $(DIST_DIR)/ $(BINARY)
  	rm -rf .tmp/ .cache/ .try/
  	rm -f logs/*.log
  	find . -name '.DS_Store' -delete
  ```

  Indentation must be a real TAB on every body line (Make requirement). Do **not** include any `rm -f handoff_unfinished_tasks.txt` line — that file is intentionally preserved.

- [ ] **Step 2.3: Sanity-check the Makefile parses**

  Run: `make -n clean`
  Expected: prints the 4 commands above (one per line), no `*** missing separator` error.

- [ ] **Step 2.4: Actually clean the working tree**

  Run: `make clean`
  Expected: commands execute without error. `ls` should now show no `gateway` binary, no `dist/`, no `.tmp/`, no `.cache/`, no `.try/`, no `.DS_Store` files anywhere, and `logs/` contains no `.log` files (the directory itself may stay).

- [ ] **Step 2.5: Verify only Makefile is in the staging area**

  Run: `git status`
  Expected: only `modified: Makefile` is listed; everything else removed was gitignored, so git doesn't see anything else changed.

- [ ] **Step 2.6: Build and test, confirming nothing broke**

  Run: `go build ./... && go test ./...`
  Expected: both exit 0.

- [ ] **Step 2.7: Commit**

  ```bash
  git add Makefile
  git commit -m "$(cat <<'EOF'
  chore(make): expand clean target, wipe local cruft

  Existing clean: only removed dist/ and the gateway binary. Add the rest
  of the gitignored cruft that accumulates locally (.tmp/, .cache/, .try/,
  log files, .DS_Store throughout the tree) so a single 'make clean' fully
  resets the working directory.

  handoff_unfinished_tasks.txt is intentionally not touched — it is a
  human-meaningful hand-off document, not a build artifact.

  Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
  EOF
  )"
  ```

  Expected: commit succeeds, no hooks fail.

---

## Task 3: Add `legacy/README.md` (Commit 2)

**Files:**
- Create: `legacy/README.md`

**Spec section:** Commit 2 in the design.

- [ ] **Step 3.1: Confirm `legacy/` contents to describe accurately**

  Run: `find legacy -maxdepth 2 -type f | sort`
  Expected: shell scripts (`install.sh`, `start.sh`, `stop.sh`, `status.sh`, `switch.sh`, `setup-autostart.sh`), a launchd plist, `legacy/config/template.yaml`, `legacy/lib/*.sh`, and a `legacy/v1/` subdirectory containing its own Go module (`go.mod`, `main.go`, `cmd/`, `internal/`, `embed/`).

- [ ] **Step 3.2: Write the README**

  Create `legacy/README.md` with this exact content:

  ````markdown
  # legacy/ — archive of historical implementations

  This directory exists for git-history reference only. **Nothing here is
  maintained, tested, or supported.** Do not import from `legacy/v1/`, do
  not invoke the shell scripts in `legacy/`, and do not link to anything
  in this directory from public documentation.

  ## What's in here

  - **`legacy/*.sh`, `legacy/com.lan-proxy-gateway.plist`, `legacy/lib/`,
    `legacy/config/`** — the original shell-script-based implementation
    that predated the Go binary. macOS launchd plist + bash scripts
    (`install.sh`, `start.sh`, `stop.sh`, `status.sh`, `switch.sh`,
    `setup-autostart.sh`) that turned a Mac into a gateway by editing
    system configuration directly.
  - **`legacy/v1/`** — the first Go rewrite. It has its own `go.mod`
    (`module gateway`) and is independent from the current top-level
    module (`github.com/tght/lan-proxy-gateway`). Replaced by the
    current `cmd/` + `internal/` tree.

  ## Why kept

  `git log` is the source of truth for project history, but having the
  files reachable from the current `HEAD` makes pre-refactor behavior
  comparison fast (e.g. when investigating regressions a user reports
  against an older release). Tagging and dropping the directory is
  acceptable future work; until then, this README marks it as inert.
  ````

- [ ] **Step 3.3: Verify the file**

  Run: `head -3 legacy/README.md`
  Expected: prints the `# legacy/ — archive of historical implementations` heading followed by the first paragraph.

- [ ] **Step 3.4: Build and test, confirming nothing broke**

  Run: `go build ./... && go test ./...`
  Expected: both exit 0. (No code changed; this is a sanity check that the docs-only commit doesn't accidentally break something.)

- [ ] **Step 3.5: Commit**

  ```bash
  git add legacy/README.md
  git commit -m "$(cat <<'EOF'
  chore(legacy): document archive status

  legacy/ has been undocumented since the Go rewrite. Add a README that
  makes the inertness explicit: nothing inside is maintained, the v1/
  subdir is a separate Go module, and the bash scripts predate the
  current binary. Removing the directory entirely is acceptable future
  work; this commit only marks it.

  Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
  EOF
  )"
  ```

---

## Task 4: Relocate `install-mihomo.sh` with wrapper (Commit 3, part 1)

**Files:**
- Move: root `install-mihomo.sh` → `scripts/install-mihomo.sh`
- Create: root `install-mihomo.sh` (new wrapper)

This is the first of three moves that together form Commit 3. We'll stage all three before committing.

- [ ] **Step 4.1: Verify the target directory exists and the source is tracked**

  Run: `ls scripts/ && git ls-files --error-unmatch install-mihomo.sh`
  Expected: `scripts/` already contains `rebuild-tag-assets.sh` and `sync-release-notes.sh`; the `git ls-files` exits 0, confirming `install-mihomo.sh` is tracked.

- [ ] **Step 4.2: Record the original file's first lines (for parity check later)**

  Run: `head -5 install-mihomo.sh`
  Expected: starts with `#!/usr/bin/env bash`, `set -euo pipefail`, then `MIRROR="${MIRROR:-https://mirror.ghproxy.com/}"`, then `VERSION="${VERSION:-v1.19.8}"`. Keep that output mentally — Step 4.5 must show the same thing for the new location.

- [ ] **Step 4.3: Move the file**

  Run: `git mv install-mihomo.sh scripts/install-mihomo.sh`
  Expected: no output, exit 0. `git status` now shows `renamed: install-mihomo.sh -> scripts/install-mihomo.sh`.

- [ ] **Step 4.4: Create the root wrapper**

  Write a new file at `install-mihomo.sh` with this exact content:

  ```bash
  #!/usr/bin/env bash
  # Compatibility wrapper. The real script lives at scripts/install-mihomo.sh.
  # Kept at repo root so any documentation or muscle-memory pointing here keeps working.
  exec "$(dirname "$0")/scripts/$(basename "$0")" "$@"
  ```

  Then make it executable: `chmod +x install-mihomo.sh`.

- [ ] **Step 4.5: Verify parity at the file level**

  Run: `head -5 scripts/install-mihomo.sh`
  Expected: identical to what Step 4.2 showed (the real script moved unchanged).

  Run: `cat install-mihomo.sh`
  Expected: the 4-line wrapper above.

- [ ] **Step 4.6: Syntax-check both files**

  Run: `bash -n install-mihomo.sh && bash -n scripts/install-mihomo.sh`
  Expected: both exit 0 (no output means OK).

  (No commit yet — we batch all three moves into one commit in Task 6.)

---

## Task 5: Relocate `download-mihomo.sh` with wrapper (Commit 3, part 2)

**Files:**
- Move: root `download-mihomo.sh` → `scripts/download-mihomo.sh`
- Create: root `download-mihomo.sh` (new wrapper)

- [ ] **Step 5.1: Record original first lines**

  Run: `head -5 download-mihomo.sh`
  Expected: starts with `#!/usr/bin/env bash`, `set -euo pipefail`, `VERSION="v1.19.8"`, `ARCH="linux-amd64"  # 默认 amd64，其他架构请自行修改`.

- [ ] **Step 5.2: Move**

  Run: `git mv download-mihomo.sh scripts/download-mihomo.sh`
  Expected: exit 0.

- [ ] **Step 5.3: Create root wrapper**

  Write `download-mihomo.sh` at repo root with this exact content:

  ```bash
  #!/usr/bin/env bash
  # Compatibility wrapper. The real script lives at scripts/download-mihomo.sh.
  # Kept at repo root so any documentation or muscle-memory pointing here keeps working.
  exec "$(dirname "$0")/scripts/$(basename "$0")" "$@"
  ```

  Then `chmod +x download-mihomo.sh`.

- [ ] **Step 5.4: Verify**

  Run: `head -5 scripts/download-mihomo.sh && echo --- && cat download-mihomo.sh`
  Expected: first block matches Step 5.1 output; second block is the 4-line wrapper.

- [ ] **Step 5.5: Syntax-check**

  Run: `bash -n download-mihomo.sh && bash -n scripts/download-mihomo.sh`
  Expected: both exit 0.

---

## Task 6: Relocate `script-demo.js` to `examples/extension.js`, commit (Commit 3, part 3 + commit)

**Files:**
- Create directory: `examples/`
- Move + rename: root `script-demo.js` → `examples/extension.js`

This task ends with the Commit 3 commit covering Tasks 4, 5, 6.

- [ ] **Step 6.1: Create the examples directory**

  Run: `mkdir -p examples`
  Expected: exits 0; directory present.

- [ ] **Step 6.2: Move + rename**

  Run: `git mv script-demo.js examples/extension.js`
  Expected: exit 0. `git status` shows `renamed: script-demo.js -> examples/extension.js`.

- [ ] **Step 6.3: Verify zero callers remain**

  Run: `git grep -nE 'install-mihomo|download-mihomo|script-demo'`
  Expected: matches **only** in `scripts/install-mihomo.sh`, `scripts/download-mihomo.sh`, the two root wrappers, and possibly `examples/extension.js` (if its content references itself, which it shouldn't). No matches in `README.md`, `README_EN.md`, `docs/`, `cmd/`, `internal/`, `Makefile`, or `.github/`.

  If matches appear in any other file, stop. Either update that file in this task or revert the moves.

- [ ] **Step 6.4: Build and test**

  Run: `go build ./... && go test ./...`
  Expected: both exit 0.

- [ ] **Step 6.5: Commit all three moves together**

  ```bash
  git add install-mihomo.sh scripts/install-mihomo.sh \
          download-mihomo.sh scripts/download-mihomo.sh \
          script-demo.js examples/extension.js
  git commit -m "$(cat <<'EOF'
  refactor: relocate stand-alone helper scripts

  Three files at the repo root that nothing in the codebase, docs, or
  release notes references move into more conventional locations:

    install-mihomo.sh   -> scripts/install-mihomo.sh   (wrapper kept)
    download-mihomo.sh  -> scripts/download-mihomo.sh  (wrapper kept)
    script-demo.js      -> examples/extension.js       (rename also)

  Both shell scripts get a 4-line root wrapper that execs the real
  target, so any documentation or muscle-memory pointing at the old
  paths keeps working. The JS file is purely demo material with no
  callers; the rename clarifies its role as an example, not runtime.

  install.sh, install.ps1, dev.sh, dev.ps1 stay at root unchanged —
  README and release notes link to install.sh / install.ps1 via the
  raw.githubusercontent URL, and README_EN documents ./dev.sh as the
  developer entry point.

  Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
  EOF
  )"
  ```

---

## Task 7: Refactor `cmd/update.go` — add constants for repeated literals (Commit 4, part 1)

**Files:**
- Modify: `cmd/update.go` (add const block near top; replace 6 literals)

**Spec section:** Commit 4 in the design — repeated literal fix (SonarQube S1192).

**Verified line locations from spec:**
- `"%s: %v"` at lines 247, 276, 412.
- `"User-Agent"` and `"lan-proxy-gateway"` together at lines 288, 334, 402.

Tests covering these code paths already exist: `TestFetchLatestReleaseTagFromRedirect` exercises the User-Agent path; the `"%s: %v"` formatting is reached by any failure path test. We rely on the existing test suite.

- [ ] **Step 7.1: Read the current constant block at the top of `cmd/update.go`**

  Run: `sed -n '22,35p' cmd/update.go`
  Expected: the `const ( updateRepo ... updateAPITimeout = 20 * time.Second )` block.

- [ ] **Step 7.2: Extend the existing const block**

  Edit `cmd/update.go`. Find the existing const block:

  ```go
  const (
      updateRepo       = "Tght1211/lan-proxy-gateway"
      updateAPIBase    = "https://api.github.com/repos/" + updateRepo
      updateLatestPage = "https://github.com/" + updateRepo + "/releases/latest"
      updateAPITimeout = 20 * time.Second
  )
  ```

  Replace it with (note: same prefix retained, new constants added):

  ```go
  const (
      updateRepo       = "Tght1211/lan-proxy-gateway"
      updateAPIBase    = "https://api.github.com/repos/" + updateRepo
      updateLatestPage = "https://github.com/" + updateRepo + "/releases/latest"
      updateAPITimeout = 20 * time.Second

      updateUserAgentHeader = "User-Agent"
      updateUserAgentValue  = "lan-proxy-gateway"
      updateErrCandidateFmt = "%s: %v"
  )
  ```

- [ ] **Step 7.3: Replace the three `"%s: %v"` occurrences**

  Use the Edit tool with `replace_all: true` on this old string:

  ```go
  failures = append(failures, fmt.Sprintf("%s: %v", candidate, err))
  ```

  Replace all 3 occurrences with:

  ```go
  failures = append(failures, fmt.Sprintf(updateErrCandidateFmt, candidate, err))
  ```

  (All three call sites have identical text per the spec's verified-sites list — `replace_all: true` handles them in one operation.)

- [ ] **Step 7.4: Replace the three `req.Header.Set("User-Agent", "lan-proxy-gateway")` occurrences**

  Use the Edit tool with `replace_all: true` on this old string:

  ```go
  req.Header.Set("User-Agent", "lan-proxy-gateway")
  ```

  Replace all 3 occurrences with:

  ```go
  req.Header.Set(updateUserAgentHeader, updateUserAgentValue)
  ```

- [ ] **Step 7.5: Verify the literals are gone outside the const block**

  Run: `grep -nE '"%s: %v"|"User-Agent"|"lan-proxy-gateway"' cmd/update.go`
  Expected: matches **only** inside the new const block declarations (3 lines, one literal per line). Zero matches elsewhere in the file.

- [ ] **Step 7.6: Build and test**

  Run: `go build ./... && go test ./cmd/...`
  Expected: build exit 0; tests show `ok  github.com/tght/lan-proxy-gateway/cmd <time>`.

- [ ] **Step 7.7: Commit (this is a logical mid-point of Commit 4, but we commit it as its own step to keep diffs reviewable)**

  We will combine Tasks 7, 8, and 9 into a single commit at the end of Task 9, so do **not** commit here. Leave the changes staged-in-tree.

---

## Task 8: Refactor `cmd/update.go` — extract `prepareUpdateBinary` helper (Commit 4, part 2)

**Files:**
- Modify: `cmd/update.go` (extract sequence inside `runUpdate`)

**Spec section:** Commit 4 in the design — cognitive complexity in `runUpdate` (currently 18, target ≤15).

- [ ] **Step 8.1: Identify the extraction region**

  Run: `sed -n '98,130p' cmd/update.go`
  Expected: the block that calls `resolveUpdateTag`, prints `当前版本` / `目标版本`, calls `gatewayReleaseAsset`, computes the download URL, calls `downloadUpdateAsset`, chmod's, and prints `下载完成`. This is the block being extracted.

- [ ] **Step 8.2: Add the new helper function**

  In `cmd/update.go`, add this helper function **immediately after** `runUpdate` (before `installUpdateBinary`):

  ```go
  // prepareUpdateBinary resolves the requested tag, downloads the matching
  // release asset to a temp path under the current user's identity, and
  // returns the resolved tag + temp path. Caller is responsible for
  // removing tmpPath on the failure path before the elevated install step
  // takes ownership of it.
  func prepareUpdateBinary(ctx context.Context, requested string) (string, string, error) {
      tag, err := resolveUpdateTag(ctx, requested)
      if err != nil {
          return "", "", err
      }
      current := Version
      color.Cyan("当前版本: %s", current)
      color.Cyan("目标版本: %s", tag)
      if current == tag {
          color.Green("已是目标版本，无需更新")
          return tag, "", nil
      }

      asset, err := gatewayReleaseAsset(runtime.GOOS, runtime.GOARCH)
      if err != nil {
          return "", "", err
      }
      url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", updateRepo, tag, asset)

      color.Cyan("下载 %s ...", asset)
      tmpPath, err := downloadUpdateAsset(ctx, url)
      if err != nil {
          return "", "", err
      }
      if runtime.GOOS != "windows" {
          _ = os.Chmod(tmpPath, 0o755)
      }
      if out, err := exec.Command(tmpPath, "--version").Output(); err == nil {
          if text := strings.TrimSpace(string(out)); text != "" {
              color.Green("下载完成: %s", text)
          }
      } else {
          color.Green("下载完成")
      }
      return tag, tmpPath, nil
  }
  ```

- [ ] **Step 8.3: Replace the body of `runUpdate` wholesale**

  The cleanest mechanical step: replace the entire body of `runUpdate` (between the opening `{` and closing `}`) with the version below. This avoids needing to find-and-rename individual `target` references — the whole post-prefetched-check body changes.

  Use the Edit tool. The old string is the current body of `runUpdate`:

  ```go
  	// elevate 后的 root 子进程：跳过下载，直接接管安装阶段。
  	if updatePrefetchedAsset != "" && updatePrefetchedTag != "" {
  		return installUpdateBinary(ctx, updatePrefetchedTag, updatePrefetchedAsset)
  	}

  	admin, _ := platform.Current().IsAdmin()
  	if !admin && runtime.GOOS == "windows" {
  		color.Red("此操作需要管理员权限。")
  		color.Yellow("请关闭当前窗口，右键 PowerShell → 以管理员身份运行，再执行：")
  		fmt.Printf("  gateway update %s\n", requested)
  		return errors.New("admin required")
  	}
  	if admin && runtime.GOOS != "windows" && os.Getenv("SUDO_USER") != "" && proxyEnvMissing() {
  		color.Yellow("提示：通过 sudo 启动会清除 HTTPS_PROXY 等代理变量。")
  		color.Yellow("如下载失败，请改用：gateway update %s（不要预先 sudo，程序会按需切换 sudo 并保留代理）。", requested)
  	}

  	target, err := resolveUpdateTag(ctx, requested)
  	if err != nil {
  		return err
  	}
  	current := Version
  	color.Cyan("当前版本: %s", current)
  	color.Cyan("目标版本: %s", target)
  	if current == target {
  		color.Green("已是目标版本，无需更新")
  		return nil
  	}

  	asset, err := gatewayReleaseAsset(runtime.GOOS, runtime.GOARCH)
  	if err != nil {
  		return err
  	}
  	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", updateRepo, target, asset)

  	color.Cyan("下载 %s ...", asset)
  	tmpPath, err := downloadUpdateAsset(ctx, url)
  	if err != nil {
  		return err
  	}
  	if runtime.GOOS != "windows" {
  		_ = os.Chmod(tmpPath, 0o755)
  	}
  	if out, err := exec.Command(tmpPath, "--version").Output(); err == nil {
  		if text := strings.TrimSpace(string(out)); text != "" {
  			color.Green("下载完成: %s", text)
  		}
  	} else {
  		color.Green("下载完成")
  	}

  	if !admin {
  		color.Cyan("切换到 sudo 进行替换 ...")
  		if err := reexecForUpdateInstall(target, tmpPath); err != nil {
  			_ = os.Remove(tmpPath)
  			return err
  		}
  		// syscall.Exec 成功时不会返回到这里。
  		return nil
  	}
  	return installUpdateBinary(ctx, target, tmpPath)
  ```

  The new string is:

  ```go
  	// elevate 后的 root 子进程：跳过下载，直接接管安装阶段。
  	if updatePrefetchedAsset != "" && updatePrefetchedTag != "" {
  		return installUpdateBinary(ctx, updatePrefetchedTag, updatePrefetchedAsset)
  	}

  	admin, _ := platform.Current().IsAdmin()
  	if !admin && runtime.GOOS == "windows" {
  		color.Red("此操作需要管理员权限。")
  		color.Yellow("请关闭当前窗口，右键 PowerShell → 以管理员身份运行，再执行：")
  		fmt.Printf("  gateway update %s\n", requested)
  		return errors.New("admin required")
  	}
  	if admin && runtime.GOOS != "windows" && os.Getenv("SUDO_USER") != "" && proxyEnvMissing() {
  		color.Yellow("提示：通过 sudo 启动会清除 HTTPS_PROXY 等代理变量。")
  		color.Yellow("如下载失败，请改用：gateway update %s（不要预先 sudo，程序会按需切换 sudo 并保留代理）。", requested)
  	}

  	tag, tmpPath, err := prepareUpdateBinary(ctx, requested)
  	if err != nil {
  		return err
  	}
  	if tmpPath == "" {
  		// prepareUpdateBinary 已打印 "已是目标版本"
  		return nil
  	}

  	if !admin {
  		color.Cyan("切换到 sudo 进行替换 ...")
  		if err := reexecForUpdateInstall(tag, tmpPath); err != nil {
  			_ = os.Remove(tmpPath)
  			return err
  		}
  		// syscall.Exec 成功时不会返回到这里。
  		return nil
  	}
  	return installUpdateBinary(ctx, tag, tmpPath)
  ```

  Verify after the edit: `grep -n '\btarget\b' cmd/update.go` should return zero matches inside `runUpdate` (the variable is gone; only `tag` remains).

- [ ] **Step 8.4: Verify compile**

  Run: `go vet ./cmd/... && go build ./...`
  Expected: both exit 0. Compile errors here mean the rename `target → tag` missed a spot — re-check.

- [ ] **Step 8.5: Run tests**

  Run: `go test ./cmd/...`
  Expected: `ok  github.com/tght/lan-proxy-gateway/cmd`.

- [ ] **Step 8.6: Do not commit yet — proceed to Task 9.**

---

## Task 9: Refactor `cmd/update.go` — extract `restartGatewayAfterUpdate` helper + commit (Commit 4, part 3 + commit)

**Files:**
- Modify: `cmd/update.go` (extract tail of `installUpdateBinary`)

**Spec section:** Commit 4 in the design — cognitive complexity in `installUpdateBinary` (currently 24, target ≤15).

- [ ] **Step 9.1: Locate the tail block of `installUpdateBinary`**

  Run: `grep -n 'replaceExecutable\|gateway 已重新启动\|SetLocalDNSToLoopback' cmd/update.go`
  Expected: line numbers for `replaceExecutable(tmpPath, self)`, `Plat.SetLocalDNSToLoopback()`, and `color.Green("gateway 已重新启动")` — all inside `installUpdateBinary`.

- [ ] **Step 9.2: Determine the `a *app.App` import path**

  Run: `grep -n 'a, err := app.New()' cmd/update.go`
  Expected: a single match inside `installUpdateBinary`. The variable is of type `*app.App` (the `app.New` constructor returns one). This matters for the helper signature.

- [ ] **Step 9.3: Confirm `app.App` exported method shape**

  Run: `grep -n 'func (a \*App)' internal/app/app.go | head -10`
  Expected: methods `Start`, `Stop`, `Status` exist on `*app.App`. The field `a.Plat` is also used by callers. The helper will take `a *app.App` directly.

- [ ] **Step 9.4: Add the helper function**

  Add this function in `cmd/update.go` **immediately after** `installUpdateBinary`:

  ```go
  // restartGatewayAfterUpdate brings gateway back up after the binary swap
  // succeeded, restoring the loopback DNS pinning if it was active before
  // the stop. Caller is expected to have already replaced the binary.
  func restartGatewayAfterUpdate(ctx context.Context, a *app.App, localDNSWasLoopback bool) error {
      color.Cyan("重新启动 gateway ...")
      if err := a.Start(ctx); err != nil {
          return err
      }
      if localDNSWasLoopback && a.Plat != nil {
          if err := a.Plat.SetLocalDNSToLoopback(); err != nil {
              color.Yellow("gateway 已启动，但恢复本机 DNS 到 127.0.0.1 失败: %v", err)
          }
      }
      color.Green("gateway 已重新启动")
      return nil
  }
  ```

- [ ] **Step 9.5: Replace the tail block in `installUpdateBinary`**

  Find this existing block at the bottom of `installUpdateBinary`:

  ```go
      if wasRunning {
          color.Cyan("重新启动 gateway ...")
          if err := a.Start(ctx); err != nil {
              return err
          }
          if localDNSWasLoopback && a.Plat != nil {
              if err := a.Plat.SetLocalDNSToLoopback(); err != nil {
                  color.Yellow("gateway 已启动，但恢复本机 DNS 到 127.0.0.1 失败: %v", err)
              }
          }
          color.Green("gateway 已重新启动")
      }
      return nil
  }
  ```

  Replace with:

  ```go
      if wasRunning {
          if err := restartGatewayAfterUpdate(ctx, a, localDNSWasLoopback); err != nil {
              return err
          }
      }
      return nil
  }
  ```

- [ ] **Step 9.6: Build and test**

  Run: `go build ./... && go test ./...`
  Expected: build exits 0; full test run exits 0 (all packages, not just cmd).

- [ ] **Step 9.7: Verify the cobra command still parses**

  Run: `go run . update --help 2>&1 | head -10`
  Expected: shows `升级到最新版本，或升级/回退到指定版本`, the `Usage:` line, the `Examples:` block, and the `Flags:` line. No panic, no missing flag.

- [ ] **Step 9.8: Commit Task 7+8+9 together as Commit 4**

  ```bash
  git add cmd/update.go
  git commit -m "$(cat <<'EOF'
  refactor(update): extract helpers and constants for lint

  Pure refactor of cmd/update.go to satisfy SonarQube findings inherited
  from the recent update fixes:

  - S1192 (repeated literals): "User-Agent", "lan-proxy-gateway", and
    "%s: %v" each appeared 3x. Pulled into module-level constants
    updateUserAgentHeader / updateUserAgentValue / updateErrCandidateFmt
    and substituted at all 6 call sites.
  - S3776 (cognitive complexity):
    * runUpdate (18 -> well below 15) — the resolve/download/chmod/
      version-print sequence moved into prepareUpdateBinary.
    * installUpdateBinary (24 -> well below 15) — the wasRunning tail
      block (start + DNS pin + log line) moved into
      restartGatewayAfterUpdate.

  No behavior change. Existing cmd/update_test.go covers all extracted
  call paths.

  Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
  EOF
  )"
  ```

---

## Task 10: Whole-branch verification

**Files:** none.

- [ ] **Step 10.1: Verify the 4-commit history is clean**

  Run: `git log --oneline main..HEAD`
  Expected: exactly 4 lines:
  1. `<sha> refactor(update): extract helpers and constants for lint`
  2. `<sha> refactor: relocate stand-alone helper scripts`
  3. `<sha> chore(legacy): document archive status`
  4. `<sha> chore(make): expand clean target, wipe local cruft`

- [ ] **Step 10.2: Full build + test**

  Run: `go build ./... && go test ./...`
  Expected: build exit 0; all packages report `ok`.

- [ ] **Step 10.3: CLI sanity**

  Run: `go build -o /tmp/gateway-cleanup . && /tmp/gateway-cleanup --version && /tmp/gateway-cleanup update --help | head -8`
  Expected: version line prints (e.g. `gateway version dev`); update help shows the `Usage:` and `Examples:` sections.

- [ ] **Step 10.4: Wrapper sanity for the moved shell scripts**

  Run: `bash -n install-mihomo.sh && bash -n scripts/install-mihomo.sh && bash -n download-mihomo.sh && bash -n scripts/download-mihomo.sh`
  Expected: all four exit 0 (no output).

- [ ] **Step 10.5: Confirm install.sh and install.ps1 are untouched**

  Run: `git diff main..HEAD -- install.sh install.ps1 dev.sh dev.ps1`
  Expected: empty output (these files were not modified by the branch).

- [ ] **Step 10.6: Hand off**

  Report back to the user: 4 commits on `chore/housekeeping`, all verifications green. Ask whether to merge into main, open a PR for review, or hold for further review.

---

## Self-review checklist (run before declaring plan complete)

- [x] Every spec section has a corresponding task: C1→Task 2, C2→Task 3, C3→Tasks 4+5+6, C4→Tasks 7+8+9, hard invariant→Task 10.
- [x] No `TBD` / `TODO` / `implement later` text anywhere in the plan.
- [x] No "similar to Task N" — code is repeated in each task that uses it.
- [x] Type and name consistency: `prepareUpdateBinary`, `restartGatewayAfterUpdate`, `updateUserAgentHeader`, `updateUserAgentValue`, `updateErrCandidateFmt` are referenced identically every time.
- [x] Exact line numbers for SonarQube literals are baked in.
- [x] Every Go-touching step ends with `go build ./...` and/or `go test ./...`.
- [x] The hard invariant ("no functional behavior change") is enforced by Task 10.5 (no diff on documented entry points) and the per-task tests.
