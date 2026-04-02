package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tght/lan-proxy-gateway/internal/config"
	"github.com/tght/lan-proxy-gateway/internal/egress"
	"github.com/tght/lan-proxy-gateway/internal/mihomo"
)

type consoleAction int

const (
	consoleActionNone consoleAction = iota
	consoleActionExit
	consoleActionRestart
	consoleActionStop
	consoleActionOpenConfig
	consoleActionOpenChainsSetup
)

type pendingConfirm struct {
	prompt string
	action consoleAction
}

type pickerMode int

const (
	pickerModeNone pickerMode = iota
	pickerModeGroups
	pickerModeNodes
)

type snapshot struct {
	modeSummary   string
	egressSummary string
	panelURL      string
	configPath    string
	iface         string
	currentNode   string
	shareEntry    string
}

type petTickMsg time.Time

type runtimeConsoleModel struct {
	width    int
	height   int
	logFile  string
	ip       string
	iface    string
	dataDir  string
	cfg      *config.Config
	client   *mihomo.Client
	snapshot snapshot
	update   *updateNotice

	input    textinput.Model
	viewport viewport.Model
	spin     spinner.Model

	history []string
	action  consoleAction

	pending *pendingConfirm

	picker      pickerMode
	groups      []mihomo.ProxyGroup
	groupCursor int
	nodeCursor  int

	petFrame int
}

func runRuntimeConsole(logFile, ip, iface, dataDir string) consoleAction {
	model := newRuntimeConsoleModel(logFile, ip, iface, dataDir)
	program := tea.NewProgram(model)
	finalModel, err := program.Run()
	if err != nil {
		fmt.Println("runtime console error:", err)
		return consoleActionExit
	}

	if m, ok := finalModel.(runtimeConsoleModel); ok {
		return m.action
	}
	return consoleActionExit
}

func newRuntimeConsoleModel(logFile, ip, iface, dataDir string) runtimeConsoleModel {
	cfg := loadConfigOrDefault()
	update := loadUpdateNotice()
	ti := textinput.New()
	ti.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#7dd3fc")).Render("/")
	ti.Placeholder = "status, config, chains, groups, logs, help"
	focusCmd := ti.Focus()
	ti.CharLimit = 512
	ti.SetWidth(48)

	sp := spinner.New()
	sp.Spinner = spinner.Points
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b"))

	vp := viewport.New()

	m := runtimeConsoleModel{
		logFile:  logFile,
		ip:       ip,
		iface:    iface,
		dataDir:  dataDir,
		cfg:      cfg,
		client:   newConsoleClient(cfg),
		update:   update,
		input:    ti,
		viewport: vp,
		spin:     sp,
		history: []string{
			"[system] Gateway Console 已连接。输入 /help 查看命令。",
			"[tip] 按 Ctrl+P 打开策略组选择器，像 CLI 版 Clash Verge 一样切节点。",
		},
	}
	if update != nil {
		m.history = append(m.history, noteLine(fmt.Sprintf("发现新版本 %s，输入 /update 查看升级方式。", update.Latest)))
	}
	m.refreshSnapshot()
	m.refreshViewport()
	if focusCmd != nil {
		_, _ = m.Update(focusCmd())
	}
	return m
}

func (m runtimeConsoleModel) Init() tea.Cmd {
	return tea.Batch(m.spin.Tick, petTickCmd(), m.input.Focus())
}

func (m runtimeConsoleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resize()
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		return m, cmd

	case petTickMsg:
		m.petFrame = (m.petFrame + 1) % len(petFrames())
		return m, petTickCmd()

	case tea.KeyMsg:
		if m.picker != pickerModeNone {
			return m.handlePickerKey(msg)
		}
		switch msg.String() {
		case "ctrl+c":
			m.action = consoleActionExit
			return m, tea.Quit
		case "ctrl+p":
			return m.openGroupPicker()
		case "enter":
			value := strings.TrimSpace(m.input.Value())
			m.input.SetValue("")
			return m.handleCommand(value)
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m runtimeConsoleModel) View() tea.View {
	if m.width == 0 || m.height == 0 {
		v := tea.NewView("loading...")
		v.AltScreen = true
		return v
	}

	header := m.renderHeader()
	main := m.renderMain()
	input := m.renderInput()

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left, header, main, input))
	v.AltScreen = true
	return v
}

func (m *runtimeConsoleModel) handleCommand(value string) (tea.Model, tea.Cmd) {
	if value == "" {
		return *m, nil
	}

	if m.pending != nil {
		return m.handleConfirm(value)
	}

	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}

	m.pushHistory(commandLine(value))

	fields := strings.Fields(strings.TrimPrefix(value, "/"))
	if len(fields) == 0 {
		m.pushHistory(noteLine("输入 /help 查看命令。"))
		m.refreshViewport()
		return *m, nil
	}

	cmd := strings.ToLower(fields[0])
	args := fields[1:]

	switch cmd {
	case "help", "?":
		m.pushHistory(
			noteLine("/status        查看完整运行状态"),
			noteLine("/summary       查看配置摘要"),
			noteLine("/config        打开配置中心"),
			noteLine("/chains        查看链式代理 / 扩展状态"),
			noteLine("/chains setup  打开链式代理向导"),
			noteLine("/groups        打开策略组选择器"),
			noteLine("/device        查看设备接入说明"),
			noteLine("/logs          查看最近日志"),
			noteLine("/guide         查看功能导航"),
			noteLine("/update        查看升级提示"),
			noteLine("/clear         清空主屏记录"),
			noteLine("/restart       重启网关（需确认）"),
			noteLine("/stop          停止网关（需确认）"),
			noteLine("/exit          退出控制台"),
		)
	case "status":
		m.pushHistory(outputBlock(m.capture(func() { runStatus(nil, nil) }))...)
	case "summary":
		m.pushHistory(outputBlock(m.capture(func() { printConfigSummary(loadConfigOrDefault()) }))...)
	case "config":
		m.action = consoleActionOpenConfig
		return *m, tea.Quit
	case "chains":
		if len(args) > 0 && args[0] == "setup" {
			m.action = consoleActionOpenChainsSetup
			return *m, tea.Quit
		}
		cfg := loadConfigOrDefault()
		if cfg.Extension.Mode == "chains" {
			m.pushHistory(outputBlock(m.capture(func() { runChainsStatus(nil, nil) }))...)
		} else {
			m.pushHistory(outputBlock(m.capture(func() { printExtensionStatus(cfg) }))...)
		}
	case "groups":
		return m.openGroupPicker()
	case "device":
		cfg := loadConfigOrDefault()
		m.pushHistory(outputBlock(m.capture(func() { printDeviceSetupPanel(m.ip, cfg.Runtime.Ports.API) }))...)
	case "logs", "log":
		m.pushHistory(m.captureLogLines(30)...)
	case "guide":
		m.pushHistory(outputBlock(m.capture(func() { printStartGuide(loadConfigOrDefault(), m.logFile) }))...)
	case "update":
		if m.update == nil {
			m.pushHistory(noteLine("当前已经是最新版本，或本次未检测到更新。"))
		} else {
			for _, line := range renderUpdateNoticeLines(m.update) {
				m.pushHistory(noteLine(line))
			}
		}
	case "clear":
		m.history = []string{noteLine("主屏记录已清空。输入 /help 查看命令。")}
	case "restart":
		m.pending = &pendingConfirm{prompt: "确认重启网关？", action: consoleActionRestart}
		m.pushHistory(noteLine("等待确认: 输入 y / n"))
	case "stop":
		m.pending = &pendingConfirm{prompt: "确认停止网关？", action: consoleActionStop}
		m.pushHistory(noteLine("等待确认: 输入 y / n"))
	case "exit", "quit":
		m.action = consoleActionExit
		return *m, tea.Quit
	default:
		m.pushHistory(errorLine("未识别的命令。输入 /help 查看可用命令。"))
	}

	m.refreshSnapshot()
	m.refreshViewport()
	return *m, nil
}

func (m *runtimeConsoleModel) handleConfirm(value string) (tea.Model, tea.Cmd) {
	answer := strings.ToLower(strings.TrimSpace(value))
	m.pushHistory(commandLine(answer))

	if answer != "y" && answer != "yes" {
		m.pushHistory(noteLine("已取消。"))
		m.pending = nil
		m.refreshViewport()
		return *m, nil
	}

	action := m.pending.action
	m.pending = nil
	switch action {
	case consoleActionRestart:
		m.pushHistory(successLine("准备重启网关..."))
		m.action = consoleActionRestart
	case consoleActionStop:
		m.pushHistory(successLine("准备停止网关..."))
		m.action = consoleActionStop
	}
	m.refreshViewport()
	return *m, tea.Quit
}

func (m runtimeConsoleModel) handlePickerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		if m.picker == pickerModeNodes {
			m.picker = pickerModeGroups
			return m, nil
		}
		m.picker = pickerModeNone
		return m, nil
	case "up", "k":
		if m.picker == pickerModeGroups && m.groupCursor > 0 {
			m.groupCursor--
		}
		if m.picker == pickerModeNodes && m.nodeCursor > 0 {
			m.nodeCursor--
		}
	case "down", "j":
		if m.picker == pickerModeGroups && m.groupCursor < len(m.groups)-1 {
			m.groupCursor++
		}
		if m.picker == pickerModeNodes && m.groupCursor < len(m.groups) {
			current := m.groups[m.groupCursor]
			if m.nodeCursor < len(current.All)-1 {
				m.nodeCursor++
			}
		}
	case "enter":
		if m.picker == pickerModeGroups {
			if len(m.groups) == 0 {
				return m, nil
			}
			m.picker = pickerModeNodes
			m.nodeCursor = 0
			return m, nil
		}
		if m.picker == pickerModeNodes {
			if len(m.groups) == 0 {
				return m, nil
			}
			group := m.groups[m.groupCursor]
			if len(group.All) == 0 {
				return m, nil
			}
			target := group.All[m.nodeCursor]
			if err := m.client.SelectProxy(group.Name, target); err != nil {
				m.pushHistory(errorLine("切换失败: " + err.Error()))
			} else {
				m.pushHistory(successLine(fmt.Sprintf("已切换策略组 %s -> %s", group.Name, target)))
			}
			m.picker = pickerModeNone
			m.refreshSnapshot()
			m.refreshViewport()
			return m, nil
		}
	}

	return m, nil
}

func (m *runtimeConsoleModel) openGroupPicker() (tea.Model, tea.Cmd) {
	groups, err := m.client.ListProxyGroups()
	if err != nil {
		m.pushHistory(errorLine("无法读取策略组: " + err.Error()))
		m.refreshViewport()
		return *m, nil
	}
	if len(groups) == 0 {
		m.pushHistory(noteLine("当前没有可切换的策略组。"))
		m.refreshViewport()
		return *m, nil
	}

	m.groups = groups
	m.groupCursor = 0
	m.nodeCursor = 0
	m.picker = pickerModeGroups
	return *m, nil
}

func (m *runtimeConsoleModel) refreshSnapshot() {
	m.cfg = loadConfigOrDefault()
	m.client = newConsoleClient(m.cfg)

	report := egress.Collect(m.cfg, m.dataDir, m.client)
	m.snapshot = snapshot{
		modeSummary:   compactModeSummary(m.cfg),
		egressSummary: compactEgressSummary(m.cfg, report),
		panelURL:      fmt.Sprintf("http://%s:%d/ui", m.ip, m.cfg.Runtime.Ports.API),
		configPath:    displayConfigPath(),
		iface:         m.iface,
		shareEntry:    m.ip,
		currentNode:   "未知",
	}
	if pg, err := m.client.GetProxyGroup("Proxy"); err == nil && strings.TrimSpace(pg.Now) != "" {
		m.snapshot.currentNode = pg.Now
	}
}

func (m *runtimeConsoleModel) refreshViewport() {
	m.viewport.SetContent(strings.Join(m.history, "\n"))
	m.viewport.GotoBottom()
}

func (m *runtimeConsoleModel) resize() {
	headerHeight := 7
	inputHeight := 3
	mainHeight := m.height - headerHeight - inputHeight
	if mainHeight < 8 {
		mainHeight = 8
	}

	leftWidth := m.width - 34
	if leftWidth < 40 {
		leftWidth = m.width
	}
	m.viewport.SetWidth(leftWidth - 4)
	m.viewport.SetHeight(mainHeight - 2)
	m.input.SetWidth(leftWidth - 8)
}

func (m runtimeConsoleModel) renderHeader() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f8fafc")).
		Background(lipgloss.Color("#0f172a")).
		Padding(0, 1)
	subStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#94a3b8"))

	line1 := titleStyle.Render(fmt.Sprintf("%s Gateway Console", m.spin.View()))
	line2 := subStyle.Render("CLI 版局域网共享 + 链式代理工作台")
	line3 := subStyle.Render("Slash 命令、策略组切换、节点选择、确认流全部在这里完成")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#38bdf8")).
		Padding(0, 1).
		Width(max(36, m.width-2))

	return box.Render(lipgloss.JoinVertical(lipgloss.Left, line1, line2, line3))
}

func (m runtimeConsoleModel) renderMain() string {
	leftWidth := m.width - 34
	if leftWidth < 40 {
		leftWidth = m.width - 2
	}
	rightWidth := m.width - leftWidth - 2
	if rightWidth < 0 {
		rightWidth = 0
	}

	left := m.renderTranscript(leftWidth)
	if rightWidth == 0 {
		return left
	}
	right := m.renderSidebar(rightWidth)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m runtimeConsoleModel) renderTranscript(width int) string {
	title := "操作记录"
	content := m.viewport.View()
	if m.picker != pickerModeNone {
		title = "策略组选择器"
		content = m.renderPicker(width - 4)
	}
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#334155")).
		Padding(0, 1).
		Width(width)

	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#e2e8f0")).Render(title)
	return box.Render(lipgloss.JoinVertical(lipgloss.Left, header, "", content))
}

func (m runtimeConsoleModel) renderSidebar(width int) string {
	if width <= 0 {
		return ""
	}

	statusCard := m.renderCard(width, "系统摘要", []string{
		"共享入口: " + m.snapshot.shareEntry,
		"当前节点: " + m.snapshot.currentNode,
		"运行模式: " + plainText(m.snapshot.modeSummary),
		"出口摘要: " + plainText(m.snapshot.egressSummary),
	})

	pet := petFrames()[m.petFrame]
	petCard := m.renderCard(width, "Gateway Pet", []string{
		pet,
		"mood: " + petMood(m.cfg),
		"tip: Ctrl+P 切节点",
	})

	commandCard := m.renderCard(width, "快捷命令", []string{
		"/status",
		"/config",
		"/chains",
		"/groups",
		"/update",
		"/device",
		"/logs",
	})

	return lipgloss.JoinVertical(lipgloss.Left, statusCard, petCard, commandCard)
}

func (m runtimeConsoleModel) renderCard(width int, title string, lines []string) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7dd3fc"))
	bodyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#e5e7eb"))

	renderedLines := make([]string, 0, len(lines))
	for _, line := range lines {
		renderedLines = append(renderedLines, bodyStyle.Render(line))
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#1e293b")).
		Padding(0, 1).
		Width(width)

	return box.Render(lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render(title), "", strings.Join(renderedLines, "\n")))
}

func (m runtimeConsoleModel) renderInput() string {
	title := "输入命令"
	placeholder := lipgloss.NewStyle().Foreground(lipgloss.Color("#94a3b8")).Render("支持 /命令。示例: /status /groups /chains setup")
	if m.pending != nil {
		title = "确认操作"
		placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("#fbbf24")).Render(m.pending.prompt + "  输入 y / n")
	}
	if m.picker != pickerModeNone {
		title = "选择器"
		placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("#fbbf24")).Render("↑/↓ 选择，Enter 确认，Esc 返回")
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#475569")).
		Padding(0, 1).
		Width(max(36, m.width-2))

	head := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#e2e8f0")).Render(title)
	return box.Render(lipgloss.JoinVertical(lipgloss.Left, head, "", placeholder, m.input.View()))
}

func (m runtimeConsoleModel) renderPicker(width int) string {
	if len(m.groups) == 0 {
		return "暂无可用策略组"
	}

	if m.picker == pickerModeGroups {
		lines := []string{"选择一个策略组，然后回车进入节点列表。", ""}
		for i, group := range m.groups {
			cursor := "  "
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("#cbd5e1"))
			if i == m.groupCursor {
				cursor = m.spin.View() + " "
				style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7dd3fc"))
			}
			lines = append(lines, style.Render(fmt.Sprintf("%s%s  [%s]  当前: %s", cursor, group.Name, group.Type, group.Now)))
		}
		return strings.Join(lines, "\n")
	}

	group := m.groups[m.groupCursor]
	lines := []string{
		fmt.Sprintf("策略组: %s", group.Name),
		fmt.Sprintf("当前节点: %s", group.Now),
		"",
	}
	for i, node := range group.All {
		cursor := "  "
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#cbd5e1"))
		if i == m.nodeCursor {
			cursor = m.spin.View() + " "
			style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f59e0b"))
		}
		if node == group.Now {
			node += "  (current)"
		}
		lines = append(lines, style.Render(cursor+node))
	}
	if width > 0 {
		return lipgloss.NewStyle().Width(width).Render(strings.Join(lines, "\n"))
	}
	return strings.Join(lines, "\n")
}

func (m *runtimeConsoleModel) pushHistory(lines ...string) {
	m.history = append(m.history, lines...)
	if len(m.history) > 240 {
		m.history = m.history[len(m.history)-240:]
	}
}

func (m runtimeConsoleModel) capture(fn func()) string {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	r, w, err := os.Pipe()
	if err != nil {
		return "无法捕获输出"
	}

	os.Stdout = w
	os.Stderr = w

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()

	_ = w.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	_ = r.Close()
	return <-done
}

func (m runtimeConsoleModel) captureLogLines(n int) []string {
	data, err := os.ReadFile(m.logFile)
	if err != nil {
		return []string{errorLine("无法读取日志文件。")}
	}
	lines := splitLines(string(data))
	start := len(lines) - n
	if start < 0 {
		start = 0
	}

	out := []string{noteLine("最近日志:")}
	for _, line := range lines[start:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		out = append(out, "  "+line)
	}
	out = append(out, noteLine("实时查看: tail -f "+m.logFile))
	return out
}

func newConsoleClient(cfg *config.Config) *mihomo.Client {
	apiURL := mihomo.FormatAPIURL("127.0.0.1", cfg.Runtime.Ports.API)
	return mihomo.NewClient(apiURL, cfg.Runtime.APISecret)
}

func petTickCmd() tea.Cmd {
	return tea.Tick(320*time.Millisecond, func(t time.Time) tea.Msg {
		return petTickMsg(t)
	})
}

func petFrames() []string {
	return []string{
		" /\\_/\\\\\n( ^.^ )\n / >🌐",
		" /\\_/\\\\\n( o.o )\n / >✨",
		" /\\_/\\\\\n( -.- )\n / >🛰",
	}
}

func petMood(cfg *config.Config) string {
	if cfg.Extension.Mode == "chains" {
		return "guarding ai traffic"
	}
	if cfg.Runtime.Tun.Enabled {
		return "sharing lan gateway"
	}
	return "waiting for commands"
}

func commandLine(text string) string {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7dd3fc")).Render("› " + text)
}
func noteLine(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#94a3b8")).Render("[note] " + text)
}
func successLine(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#4ade80")).Render("[ok] " + text)
}
func errorLine(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#f87171")).Render("[error] " + text)
}

func outputBlock(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	raw := splitLines(text)
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		if strings.TrimSpace(line) == "" {
			continue
		}
		out = append(out, "  "+line)
	}
	return out
}

func plainText(s string) string {
	replacer := strings.NewReplacer(
		"\x1b[0m", "",
		"\x1b[1m", "",
		"\x1b[2m", "",
		"\x1b[31m", "",
		"\x1b[32m", "",
		"\x1b[33m", "",
		"\x1b[34m", "",
		"\x1b[35m", "",
		"\x1b[36m", "",
		"\x1b[37m", "",
	)
	return replacer.Replace(s)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
