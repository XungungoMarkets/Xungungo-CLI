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
- 📉 **Technical analysis**: RSI, MACD, SMA, EMA, Bollinger Bands
- 🎨 **Visual interface**: Colors and elegant formatting in terminal
- ⚡ **Fast**: Multiple symbols in a single request
- 📅 **Configurable periods**: 5 days, 1 month, 3 months, 6 months, 1 year, 5 years
- 🔄 **Auto-updater**: Built-in update mechanism to keep xgg current

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
xgg technical --help
xgg search --help
```

### JSON Output

All commands support JSON output format using the `--json` flag. This is useful for scripting, automation, or integrating with other tools.

```bash
# Stock quote in JSON
xgg stock NVDA --json

# Multiple stocks in JSON
xgg stock NVDA AAPL TSLA --json

# Historical data in JSON
xgg history NVDA --json --period 1y

# Technical analysis in JSON
xgg technical NVDA --json --indicator rsi

# Multiple indicators in JSON
xgg technical NVDA --json --indicator rsi,macd,sma

# Symbol search in JSON
xgg search NVDA --json
```

**JSON Output Example (Stock):**
```json
{
  "symbol": "NVDA",
  "name": "NVIDIA Corporation",
  "price": 186.03,
  "change": 1.26,
  "change_percent": 0.68,
  "volume": 138663908,
  "market_cap": 0
}
```

**JSON Output Example (Technical):**
```json
[
  {
    "symbol": "NVDA",
    "indicator": "RSI(14)",
    "value": 51.47,
    "signal": "neutral"
  },
  {
    "symbol": "NVDA",
    "indicator": "MACD(12,26,9)",
    "macd": -0.70,
    "signal": -0.50,
    "histogram": -0.20,
    "signal_type": "bearish"
  }
]
```

### Data Provider Configuration

`xgg` now supports provider selection and Nasdaq runtime tuning with precedence:

`flags > environment variables > defaults`

```bash
# Force Nasdaq provider
xgg --provider nasdaq stock NVDA

# Auto mode (Nasdaq primary + legacy fallback)
xgg --provider auto history AAPL --period 1m

# Force legacy provider
xgg --provider legacy stock MSFT
```

Available global flags:

- `--provider` (`auto`, `nasdaq`, `legacy`)
- `--rate-limit` (requests per second)
- `--max-retries`
- `--retry-delay` (example: `2s`)
- `--timeout` (example: `30s`)
- `--watchlist-type` (example: `Rv`)

Equivalent environment variables:

- `XGG_PROVIDER`
- `XGG_RATE_LIMIT`
- `XGG_MAX_RETRIES`
- `XGG_RETRY_DELAY`
- `XGG_TIMEOUT`
- `XGG_WATCHLIST_TYPE`

---

## � Commands

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

### Symbol Search

Search symbols (stocks, ETFs, indices) using Nasdaq autosuggest:

```bash
xgg search NVDA
xgg search Apple --limit 20
xgg search semiconductors --json
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

### Technical Analysis

Calculate and display technical analysis indicators for a stock:

```bash
xgg technical NVDA
```

**Output (RSI):**
```
📊 NVDA - RSI (14)
┌─────────────────────────────┐
│ Current RSI:   65.23       │
└─────────────────────────────┘

⚪ Neutral zone
```

**Output (MACD):**
```
📊 NVDA - MACD (12, 26, 9)
┌─────────────────────────────┐
│ MACD Line:    12.34       │
│ Signal Line:  10.56       │
│ Histogram:     1.78       │
└─────────────────────────────┘

🟢 Bullish - Buy signal
```

Get a specific indicator:

```bash
# RSI (Relative Strength Index)
xgg technical NVDA --indicator rsi

# MACD (Moving Average Convergence Divergence)
xgg technical NVDA --indicator macd

# SMA (Simple Moving Average)
xgg technical NVDA --indicator sma

# EMA (Exponential Moving Average)
xgg technical NVDA --indicator ema

# Bollinger Bands
xgg technical NVDA --indicator bb
```

**Multiple indicators at once:**

```bash
# Two indicators
xgg technical NVDA --indicator rsi,macd

# Three indicators
xgg technical NVDA --indicator rsi,sma,ema

# All indicators
xgg technical NVDA --indicator all
```

Use different time periods:

```bash
xgg technical NVDA --indicator rsi --period 1y
xgg technical NVDA --indicator rsi,macd --period 3m
```

**Available indicators:**
- `rsi` - Relative Strength Index (14 period)
- `macd` - Moving Average Convergence Divergence (12, 26, 9)
- `sma` - Simple Moving Averages (20 and 50 period)
- `ema` - Exponential Moving Averages (12 and 26 period)
- `bb` - Bollinger Bands (20 period, 2 standard deviations)
- `all` - Display all indicators

**Note:** You can specify multiple indicators separated by commas (e.g., `rsi,macd,sma`).

**Technical analysis includes:**
- RSI: Overbought (>70) and oversold (<30) signals
- MACD: Trend direction and momentum
- SMA/EMA: Support, resistance, and trend identification
- Bollinger Bands: Volatility and potential breakout points

### Quick Format

To get information without decorative borders, you can redirect the output or use it in scripts:

```bash
xgg stock NVDA | grep -E "NVDA|Price"
```

### Version & Updates

Check the current version:

```bash
xgg version
```

**Output:**
```
📈 Xungungo CLI
Version: 0.1.0
GitHub: https://github.com/XungungoMarkets/Xungungo-CLI
```

Check if an update is available without installing:

```bash
xgg check-update
```

**Output:**
```
✓ You are using the latest version: 0.1.0
```

Or if an update is available:
```
⚠ A new version is available!
  Current: 0.1.0
  Latest:  0.2.0
  Release:  https://github.com/XungungoMarkets/Xungungo-CLI/releases/download/v0.2.0/xgg-linux-amd64.tar.gz

Run 'xgg update' to update to the latest version.
```

Update xgg to the latest version:

```bash
xgg update
```

The update command will:
- Check for the latest release on GitHub
- Show you the current and latest version
- Ask for confirmation before downloading
- Download and replace the binary automatically
- Verify the download using checksums

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
