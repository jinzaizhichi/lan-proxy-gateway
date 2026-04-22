//go:build windows

package engine

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func configureProcAttrs(cmd *exec.Cmd) {
	// Hide child window on Windows.
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

// terminateProcess kills the target PID and its entire child tree on Windows.
// mihomo spawns a TUN helper; plain TerminateProcess (os.Process.Kill) leaves
// the child alive, which keeps holding ports like 17890/53 so the next Start
// sees a spurious "port in use". `taskkill /T /F` walks the tree.
func terminateProcess(p *os.Process) {
	if p == nil {
		return
	}
	if err := exec.Command("taskkill", "/PID", strconv.Itoa(p.Pid), "/T", "/F").Run(); err != nil {
		// taskkill missing is extremely unlikely on modern Windows, but fall
		// back to the single-PID kill just in case — better something than
		// nothing.
		_ = p.Kill()
	}
}

// pidAlive probes via `tasklist` because Process.Signal is not implemented on
// Windows (it always returns "unsupported"), so the unix-style Signal(0) trick
// would report every PID as dead — which made the menu always show "未启动"
// even when mihomo was running.
func pidAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	out, err := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/NH", "/FO", "CSV").Output()
	if err != nil {
		return false
	}
	// tasklist prints a locale-dependent "no tasks" line on miss. On hit the
	// CSV row always quotes the PID: `"imagename.exe","1234",...`. Matching the
	// quoted PID avoids false positives from the info message.
	return strings.Contains(string(out), fmt.Sprintf("\"%d\"", pid))
}
