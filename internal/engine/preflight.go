package engine

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Preflight checks common startup failures BEFORE launching mihomo so the user
// gets a helpful error instead of the opaque "API did not become ready in 10s".
type PortCheck struct {
	Label string // e.g. "mixed", "api", "dns"
	Port  int
	Bind  string // "0.0.0.0" or "127.0.0.1"
}

// PortOwner describes the process currently holding a TCP port.
type PortOwner struct {
	Name string // e.g. "mihomo"
	PID  int
}

// PortConflict is one entry from a preflight check.
type PortConflict struct {
	Check PortCheck
	Owner *PortOwner // nil if we couldn't identify it
	Err   error      // raw net.Listen error
}

// PortConflictError is returned by CheckPorts so callers can offer recovery
// (e.g. "occupier is stale mihomo, kill it and retry?").
type PortConflictError struct {
	Conflicts []PortConflict
}

func (e *PortConflictError) Error() string {
	var b strings.Builder
	b.WriteString("端口冲突：\n")
	for _, c := range e.Conflicts {
		fmt.Fprintf(&b, "  • %s 端口 %d 被占用", c.Check.Label, c.Check.Port)
		if c.Owner != nil {
			fmt.Fprintf(&b, " 占用者: %s (PID %d)", c.Owner.Name, c.Owner.PID)
		}
		b.WriteString("\n")
	}
	b.WriteString("\n解决：\n  1) 主菜单 → 2 分流 & 规则 → 9 高级设置 里关 DNS 或换端口\n  2) 主菜单 → 4 启动 / 重启 / 停止 → 4 清理残留 mihomo（如果是残留的 mihomo 在占端口）")
	return b.String()
}

// HasStaleMihomo returns true if any conflict is caused by a mihomo we can kill.
func (e *PortConflictError) HasStaleMihomo() bool {
	for _, c := range e.Conflicts {
		if c.Owner != nil && isMihomo(c.Owner.Name) {
			return true
		}
	}
	return false
}

// StaleMihomoPIDs returns the PIDs of mihomo occupants (deduped).
func (e *PortConflictError) StaleMihomoPIDs() []int {
	seen := map[int]bool{}
	var out []int
	for _, c := range e.Conflicts {
		if c.Owner != nil && isMihomo(c.Owner.Name) && !seen[c.Owner.PID] {
			seen[c.Owner.PID] = true
			out = append(out, c.Owner.PID)
		}
	}
	return out
}

// CheckPorts returns a *PortConflictError describing any port conflicts, or nil.
func CheckPorts(checks []PortCheck) error {
	var conflicts []PortConflict
	for _, c := range checks {
		addr := fmt.Sprintf("%s:%d", c.Bind, c.Port)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			conflicts = append(conflicts, PortConflict{
				Check: c,
				Owner: LookupPortOwner(c.Port),
				Err:   err,
			})
			continue
		}
		_ = ln.Close()
	}
	if len(conflicts) == 0 {
		return nil
	}
	return &PortConflictError{Conflicts: conflicts}
}

// LookupPortOwner names the process holding a port. Uses lsof on mac/linux,
// netstat + tasklist on Windows. Returns nil if unknown.
func LookupPortOwner(port int) *PortOwner {
	if runtime.GOOS == "windows" {
		return lookupPortOwnerWindows(port)
	}
	out, err := exec.Command("lsof", "-nP", "-iTCP:"+fmt.Sprint(port), "-sTCP:LISTEN").Output()
	if err != nil || len(out) == 0 {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return nil
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return nil
	}
	pid, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil
	}
	return &PortOwner{Name: fields[0], PID: pid}
}

// FindStaleMihomoPIDs returns PIDs of every running mihomo on this host
// (regardless of which port they're on). Used by the "cleanup" menu action
// so users can recover from prior crashes without knowing which ports.
func FindStaleMihomoPIDs() []int {
	if runtime.GOOS == "windows" {
		return findStaleMihomoPIDsWindows()
	}
	// pgrep is installed on mac + most linux distros.
	out, err := exec.Command("pgrep", "-x", "mihomo").Output()
	if err != nil {
		return nil
	}
	var pids []int
	for _, line := range strings.Fields(strings.TrimSpace(string(out))) {
		if p, err := strconv.Atoi(line); err == nil {
			pids = append(pids, p)
		}
	}
	return pids
}

// KillPID sends SIGTERM then SIGKILL as fallback on Unix. On Windows it shells
// out to `taskkill /T /F` so child processes (mihomo's TUN worker) also die —
// TerminateProcess alone leaves them dangling and the next start hits a port
// conflict. Requires privileges if the target was started as root/Administrator.
func KillPID(pid int) error {
	if runtime.GOOS == "windows" {
		return exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F").Run()
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := proc.Signal(syscall.SIGTERM); err == nil {
		// give it a moment to shut down cleanly before we escalate
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			if !pidAlive(pid) {
				return nil
			}
		}
	}
	return proc.Kill()
}

// findStaleMihomoPIDsWindows enumerates mihomo.exe via `tasklist`. The output
// comes back CSV-formatted: `"mihomo.exe","1234","Console","1","12,345 K"`.
// When no matches exist, tasklist prints a locale-dependent info line (CJK on
// Chinese Windows), so we filter by the image name rather than by error code.
func findStaleMihomoPIDsWindows() []int {
	out, err := exec.Command("tasklist", "/FI", "IMAGENAME eq mihomo.exe", "/FO", "CSV", "/NH").Output()
	if err != nil {
		return nil
	}
	output := strings.TrimSpace(string(out))
	if !strings.Contains(strings.ToLower(output), "mihomo.exe") {
		return nil
	}
	var pids []int
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Split(strings.TrimSpace(line), ",")
		if len(fields) < 2 {
			continue
		}
		pidStr := strings.Trim(fields[1], "\" \r")
		if pid, err := strconv.Atoi(pidStr); err == nil {
			pids = append(pids, pid)
		}
	}
	return pids
}

// lookupPortOwnerWindows uses `netstat -ano` to map port → PID, then
// `tasklist` to resolve PID → image name.
func lookupPortOwnerWindows(port int) *PortOwner {
	out, err := exec.Command("netstat", "-ano", "-p", "TCP").Output()
	if err != nil {
		return nil
	}
	suffix := fmt.Sprintf(":%d", port)
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, "LISTENING") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		// fields[1] is the local address, e.g. "0.0.0.0:53" or "[::]:53".
		if !strings.HasSuffix(fields[1], suffix) {
			continue
		}
		pid, err := strconv.Atoi(fields[len(fields)-1])
		if err != nil {
			continue
		}
		return &PortOwner{Name: lookupProcessNameWindows(pid), PID: pid}
	}
	return nil
}

func lookupProcessNameWindows(pid int) string {
	out, err := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV", "/NH").Output()
	if err != nil {
		return ""
	}
	output := strings.TrimSpace(string(out))
	if output == "" {
		return ""
	}
	// CSV first field: "imagename.exe"
	fields := strings.Split(output, ",")
	if len(fields) < 1 {
		return ""
	}
	name := strings.Trim(fields[0], "\" \r")
	return strings.TrimSuffix(name, ".exe")
}

func isMihomo(name string) bool {
	n := strings.ToLower(name)
	return n == "mihomo" || strings.HasPrefix(n, "mihomo")
}

// TailLog returns the last N lines of a log file, for post-mortem on startup failure.
func TailLog(path string, n int) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	start := 0
	if len(lines) > n {
		start = len(lines) - n
	}
	return strings.Join(lines[start:], "\n")
}

// WaitForFile polls a path until it exists or timeout, used for mihomo's config write-out.
func WaitForFile(path string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}
