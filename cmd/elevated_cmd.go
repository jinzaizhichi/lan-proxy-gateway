package cmd

import "runtime"

// elevatedCmd returns a shell-ready invocation string for `gateway <sub>` with
// the right prefix for the OS:
//   - Windows: bare `gateway …` (the shell should already be an admin
//     PowerShell; there is no `sudo`)
//   - Unix: `sudo gateway …`
//
// Pass sub="" to get the bare invocation (no subcommand).
func elevatedCmd(sub string) string {
	name := "gateway"
	if sub != "" {
		name = name + " " + sub
	}
	if runtime.GOOS == "windows" {
		return name
	}
	return "sudo " + name
}
