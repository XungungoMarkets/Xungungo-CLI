# 📈 Xungungo-CLI (xgg)

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Xungungo-CLI** is a command-line tool written in Go to access real-time financial data. Stock quotes, historical data, and more, directly from your terminal.

---

## 📋 Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [Requirements](#requirements)
- [License](#license)

---

## ✨ Features

- 📊 **Real-time quotes**: Get current prices for any stock
- 📈 **Historical data**: Access OHLCV data (Open, High, Low, Close, Volume)
- 🎨 **Visual interface**: Colors and elegant formatting in the terminal
- ⚡ **Fast**: Multiple symbols in a single request
- 📅 **Configurable periods**: 5 days, 1 month, 3 months, 6 months, 1 year, 5 years

---

  ## 🔽 Installation

### From GitHub Releases (Recommended) ⭐

Download the pre-compiled binary for your operating system from the [v0.1.0 Release](https://github.com/XungungoMarkets/Xungungo-CLI/releases/tag/v0.1.0) page.

**Linux/macOS:**
```bash
# Download the file
wget https://github.com/XungungoMarkets/Xungungo-CLI/releases/download/v0.1.0/xgg-linux-amd64.tar.gz

# Extract
tar xzf xgg-linux-amd64.tar.gz

# Move to PATH
sudo mv xgg /usr/local/bin/
```

**macOS (Apple Silicon):**
```bash
# Download
wget https://github.com/XungungoMarkets/Xungungo-CLI/releases/download/v0.1.0/xgg-darwin-arm64.tar.gz

# Extract
tar xzf xgg-darwin-arm64.tar.gz

# Move to PATH
sudo mv xgg /usr/local/bin/
```

**Windows:**
```powershell
# Download the file
# https://github.com/XungungoMarkets/Xungungo-CLI/releases/download/v0.1.0/xgg-windows-amd64.zip

# Extract and move to a folder in your PATH
```

**Verify checksums (optional):**
```bash
sha256sum -c xgg-linux-amd64.tar.gz.sha256
```

### From Go

```bash
go install github.com/XungungoMarkets/xgg@latest
```

This will install the `xgg` binary in your `$GOPATH/bin`. Make sure it's in your `$PATH`.

### From Source

```bash
git clone https://github.com/XungungoMarkets/xgg.git
cd xgg
go build -o xgg
```

Move the binary to your PATH:

```bash
# Linux/macOS
sudo mv xgg /usr/local/bin/

# Windows
# Move xgg.exe to a folder in your PATH
```

---

## 🚀 Usage

### View Help

```bash
# General help
xgg --help

# Specific command help
xgg stock --help
xgg history --help
```

---

## 📖 Commands

### Live Quotes

Get the current price of a stock:

```bash
xgg stock NVDA
```

**Output:**
```
┌─────────────────────────────────────────┐
│  NVDA - NVIDIA Corporation             │
│  $875.28  ▲ +12.34 (+1.43%)             │
│  Vol: 45.2M  │  Mkt Cap: 2.2T          │
└─────────────────────────────────────────┘
```

Get quotes for multiple stocks at once:

```bash
xgg stock NVDA AAPL TSLA MSFT GOOGL
```

### Historical Data

Get historical data for the last month (default):

```bash
xgg history NVDA
```

**Output:**
```
  NVDA Historical Data
  Date          Open      High       Low     Close       Volume
  ──────────────────────────────────────────────────────────────────
  2026-03-10  $850.00  $880.00  $845.00  $875.28     45,234,567 ▲
  2026-03-09  $835.00  $852.00  $830.00  $848.50     52,123,456 ▲
  2026-03-08  $820.00  $840.00  $815.00  $834.20     48,567,890 ▼
  ...
```

Get historical data with a specific period:

```bash
# Last 5 days
xgg history NVDA --period 5d

# Last 6 months
xgg history AAPL --period 6m

# Last year
xgg history TSLA --period 1y

# Last 5 years
xgg history MSFT --period 5y
```

**Available periods:**
- `5d` - Last 5 days
- `1m` - Last month (default)
- `3m` - Last 3 months
- `6m` - Last 6 months
- `1y` - Last year
- `5y` - Last 5 years

### Quick Format

To get information without decorative borders, you can redirect the output or use it in scripts:

```bash
xgg stock NVDA | grep -E "NVDA|Price"
```

---

## 🛠️ Requirements

- **Go**: 1.26 or higher (only for compiling from source)
- **Operating System**: Linux, macOS, Windows
- **Internet Connection**: Required to fetch financial data

---

## 📝 License

This project is licensed under the MIT License.

---

## 🔗 Links

- [GitHub Repository](https://github.com/XungungoMarkets/Xungungo-CLI)
- [Issues](https://github.com/XungungoMarkets/Xungungo-CLI/issues)
- [v0.1.0 Release](https://github.com/XungungoMarkets/Xungungo-CLI/releases/tag/v0.1.0)

---

**Made with ❤️ by Xungungo Markets**

*Access financial markets from your terminal, easy and fast.*