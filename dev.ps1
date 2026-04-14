$ErrorActionPreference = "Stop"

$RootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $RootDir

$GoBin = if ($env:GO_BIN) { $env:GO_BIN } else { "go" }

if ($env:VERSION) {
    $Version = $env:VERSION
} else {
    try {
        $Version = (git describe --tags --always --dirty 2>$null).Trim()
    } catch {
        $Version = ""
    }
    if (-not $Version) {
        $Version = "dev-local"
    }
}

$Ldflags = if ($env:LDFLAGS) { $env:LDFLAGS } else { "-X main.version=$Version" }

$CacheDir = if ($env:CACHE_DIR) { $env:CACHE_DIR } else { Join-Path $RootDir ".cache" }
$BuildDir = if ($env:BUILD_DIR) { $env:BUILD_DIR } else { Join-Path $RootDir ".tmp" }
$BinaryName = if ($env:BINARY_NAME) { $env:BINARY_NAME } else { "gateway-dev.exe" }
$BinaryPath = if ($env:BINARY_PATH) { $env:BINARY_PATH } else { Join-Path $BuildDir $BinaryName }

if (-not $env:GOCACHE) {
    $env:GOCACHE = Join-Path $CacheDir "go-build"
}

if ($env:USE_LOCAL_GOMODCACHE -eq "1" -and -not $env:GOMODCACHE) {
    $env:GOMODCACHE = Join-Path $CacheDir "go-mod"
}

function Write-Info {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Yellow
}

function Fail {
    param([string]$Message)
    throw $Message
}

function Show-Usage {
    @"
用法:
  ./dev.ps1 build
  ./dev.ps1 test
  ./dev.ps1 test-core
  ./dev.ps1 run -- <gateway 参数>
  ./dev.ps1 start [gateway start 参数]
  ./dev.ps1 console
  ./dev.ps1 stop
  ./dev.ps1 restart
  ./dev.ps1 status
  ./dev.ps1 clean

常用例子:
  ./dev.ps1 build
  ./dev.ps1 test
  ./dev.ps1 test-core
  ./dev.ps1 run -- --version
  ./dev.ps1 run -- config show
  ./dev.ps1 start
  ./dev.ps1 console

环境变量:
  GO_BIN=go                                 指定 go 可执行文件
  VERSION=dev-local                         覆盖构建版本
  BINARY_PATH=.tmp/gateway-dev.exe          覆盖输出二进制路径
  USE_LOCAL_GOMODCACHE=1                    让模块缓存也落到仓库内 .cache/
"@ | Write-Host
}

function Ensure-Go {
    if (-not (Get-Command $GoBin -ErrorAction SilentlyContinue)) {
        Fail "未找到 go，请先安装 Go 1.25+"
    }
}

function Prepare-Dirs {
    New-Item -ItemType Directory -Path $BuildDir -Force | Out-Null
    New-Item -ItemType Directory -Path $env:GOCACHE -Force | Out-Null

    if ($env:GOMODCACHE) {
        New-Item -ItemType Directory -Path $env:GOMODCACHE -Force | Out-Null
    }
}

function Build-Binary {
    Ensure-Go
    Prepare-Dirs
    Write-Info "编译本地开发二进制..."
    Write-Info "输出路径: $BinaryPath"
    & $GoBin build -ldflags $Ldflags -o $BinaryPath .
}

function Run-Tests {
    Ensure-Go
    Prepare-Dirs
    Write-Info "运行全量测试..."
    & $GoBin test ./...
}

function Run-CoreTests {
    Ensure-Go
    Prepare-Dirs
    Write-Info "运行核心包测试..."
    & $GoBin test ./cmd ./internal/config ./internal/egress ./internal/platform ./internal/rules ./internal/template
}

function Run-Binary {
    param([string[]]$BinaryArgs)

    Build-Binary
    Write-Info "运行: $BinaryPath $($BinaryArgs -join ' ')"
    & $BinaryPath @BinaryArgs
}

function Test-IsAdmin {
    try {
        $currentIdentity = [Security.Principal.WindowsIdentity]::GetCurrent()
        $principal = New-Object Security.Principal.WindowsPrincipal($currentIdentity)
        return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    } catch {
        return $false
    }
}

function Run-With-AdminIfNeeded {
    param(
        [string]$Subcommand,
        [string[]]$SubArgs
    )

    Build-Binary

    if (-not (Test-IsAdmin)) {
        Write-Warn "命令 $Subcommand 需要管理员权限，将弹出 UAC 确认。"
        $process = Start-Process -FilePath $BinaryPath -ArgumentList (@($Subcommand) + $SubArgs) -Verb RunAs -Wait -PassThru
        exit $process.ExitCode
    }

    & $BinaryPath $Subcommand @SubArgs
}

function Clean-Artifacts {
    Write-Info "清理本地开发产物..."

    if (Test-Path $BuildDir) {
        Remove-Item -LiteralPath $BuildDir -Recurse -Force
    }

    if (Test-Path $CacheDir) {
        Remove-Item -LiteralPath $CacheDir -Recurse -Force
    }
}

$Command = if ($args.Count -gt 0) { $args[0] } else { "help" }
$Rest = if ($args.Count -gt 1) { $args[1..($args.Count - 1)] } else { @() }

switch ($Command) {
    "build" {
        Build-Binary
    }
    "test" {
        Run-Tests
    }
    "test-core" {
        Run-CoreTests
    }
    "check" {
        Build-Binary
        Run-CoreTests
    }
    "run" {
        if ($Rest.Count -gt 0 -and $Rest[0] -eq "--") {
            if ($Rest.Count -gt 1) {
                $Rest = $Rest[1..($Rest.Count - 1)]
            } else {
                $Rest = @()
            }
        }
        Run-Binary -BinaryArgs $Rest
    }
    "start" {
        Run-With-AdminIfNeeded -Subcommand "start" -SubArgs $Rest
    }
    "console" {
        Run-With-AdminIfNeeded -Subcommand "console" -SubArgs $Rest
    }
    "stop" {
        Run-With-AdminIfNeeded -Subcommand "stop" -SubArgs $Rest
    }
    "restart" {
        Run-With-AdminIfNeeded -Subcommand "restart" -SubArgs $Rest
    }
    "status" {
        Run-Binary -BinaryArgs (@("status") + $Rest)
    }
    "clean" {
        Clean-Artifacts
    }
    "help" {
        Show-Usage
    }
    "-h" {
        Show-Usage
    }
    "--help" {
        Show-Usage
    }
    default {
        Show-Usage
        Fail "未知命令: $Command"
    }
}
