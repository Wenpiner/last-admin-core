# install.ps1 - Windows 安装脚本

function Write-Success {
    param([string]$Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "⚠ $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "✗ $Message" -ForegroundColor Red
}

# 检查 Python 是否已安装
function Check-Python {
    try {
        $version = python --version 2>&1
        Write-Success "Python 已安装: $version"
        return $true
    } catch {
        Write-Warning "Python 未安装"
        return $false
    }
}

# 安装 Python（Windows）
function Install-Python {
    Write-Warning "正在为 Windows 安装 Python 3..."
    
    if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
        Write-Warning "正在安装 Chocolatey..."
        Set-ExecutionPolicy Bypass -Scope Process -Force
        [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
        iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
    }
    
    choco install python -y
    Write-Success "Python 3 安装完成"
}

# 安装 Python 依赖
function Install-Dependencies {
    Write-Warning "正在安装 Python 依赖..."
    
    if (-not (Test-Path "venv")) {
        python -m venv venv
        Write-Success "虚拟环境已创建"
    }
    
    & ".\venv\Scripts\Activate.ps1"
    
    python -m pip install --upgrade pip
    pip install -r requirements.txt
    
    Write-Success "Python 依赖安装完成"
}

# 运行主安装脚本
function Run-Installer {
    Write-Warning "正在启动安装向导..."
    
    & ".\venv\Scripts\Activate.ps1"
    python install.py
}

# 主函数
function Main {
    Write-Host "================================" -ForegroundColor Green
    Write-Host "Last Admin 安装向导" -ForegroundColor Green
    Write-Host "================================" -ForegroundColor Green
    Write-Host ""
    
    if (-not (Check-Python)) {
        Install-Python
    }
    
    Install-Dependencies
    Run-Installer
}

$ErrorActionPreference = "Stop"
trap {
    Write-Error "安装失败: $_"
    exit 1
}

Main

