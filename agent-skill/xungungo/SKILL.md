---
name: xungungo
description: Fetches real-time stock quotes, historical OHLCV data, technical analysis indicators (RSI, MACD, SMA, EMA, Bollinger Bands), generates PNG price charts, shows market breadth by sector, industry, or country, and queries the full NASDAQ stock screener using the xgg CLI. Use when the user asks for stock prices, market data, trading indicators, historical prices, ETF symbol searches, wants to visualize a chart, asks how sectors/countries are performing, or wants to browse/filter all listed stocks — even without explicitly mentioning xgg or Xungungo. Triggers — "what's NVDA trading at", "show me Apple stock", "RSI for MSFT", "MACD for Amazon last 3 months", "compare AAPL and TSLA", "search for semiconductor ETFs", "generate a candlestick chart for TSLA", "chart NVDA with Bollinger Bands", "plot AAPL over 1 year", "how is tech sector doing today", "which sectors are up", "show me stocks by country", "show me all NASDAQ stocks", "list stocks in the energy sector", "screener for US tech stocks".
---

# Xungungo (xgg) — Financial Data CLI

## Commands

```bash
xgg stock SYMBOL [SYMBOL2...]                                         # Real-time quotes
xgg history SYMBOL --period PERIOD [--interval INTERVAL]              # Historical OHLCV
xgg technical SYMBOL --indicator INDICATOR --period PERIOD            # Technical analysis
xgg search "query" [--limit N] [--market-data]                       # Symbol/ETF search
xgg chart SYMBOL [options]                                            # Generate PNG chart
xgg sectors [--by-industry] [--by-stock]                              # % change by sector
xgg country [country name...] [--by-stock]                            # % change by country
xgg screener [SYMBOL...] [--sector S] [--country C] [--industry I]   # Full screener table
xgg version                                                           # Show version
xgg update                                                            # Update xgg to latest
xgg check-update                                                      # Check for updates
```

**Global flags (all commands):** `--json`, `--provider auto|nasdaq|legacy`, `--rate-limit`, `--max-retries`, `--retry-delay`, `--timeout`

---

## chart command

Generates a PNG price chart saved to disk. Prints the output file path on success.

```bash
xgg chart SYMBOL [flags]
```

| Flag | Short | Default | Options |
|------|-------|---------|---------|
| `--type` | `-t` | `line` | `line`, `candlestick` |
| `--period` | `-p` | `1m` | see periods table |
| `--interval` | `-i` | `day` | `day`, `week`, `month` |
| `--indicator` | | _(none)_ | `sma20`, `sma50`, `sma200`, `ema12`, `ema26`, `ema50`, `bb`, `linear`, `cubic` (comma-separated) |
| `--theme` | | `dark` | `light`, `dark`, `vivid-light`, `vivid-dark`, `ant`, `grafana` |
| `--output` | `-o` | `<symbol>_chart.png` | any file path |
| `--width` | | `900` | pixels |
| `--height` | | `500` | pixels |

**Examples:**
```bash
xgg chart AAPL
xgg chart NVDA --type candlestick --period 3m
xgg chart TSLA --period 1y --interval week --indicator sma20,sma50
xgg chart NVDA --type candlestick --period 6m --indicator bb
xgg chart MSFT --period 5y --interval month --type candlestick --indicator sma20,ema50
xgg chart AAPL --output /tmp/aapl.png --theme light
```

When the user asks to "plot", "graph", or "visualize" a stock, use `chart`.

---

## search command

```bash
xgg search "query" [--limit N] [--market-data] [--json]
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--limit` | `-l` | `10` | Max number of results |
| `--market-data` | | `false` | Include real-time market data in results when available |

---

## history command

```bash
xgg history SYMBOL --period PERIOD [--interval day|week|month] [--json]
```

Use `--interval week` or `--interval month` to aggregate bars for longer periods.

---

## technical command

```bash
xgg technical SYMBOL --indicator INDICATOR --period PERIOD [--interval day|week|month] [--json]
```

**Indicators:** `rsi`, `macd`, `sma`, `ema`, `bb`, `all` (comma-separated: `rsi,macd`)

Minimum data requirements — the period is **auto-extended** if the requested period has fewer bars than needed:
- RSI: 14 bars · MACD: 35 bars (26 + 9 signal) · SMA/EMA: 200 bars · BB: 20 bars

When auto-extension happens, a note is printed to stderr: `Note: period extended to Xm to compute SMA on Y bars`.

---

## Periods

| User says | Flag |
|-----------|------|
| last week / 5 days | `5d` |
| 1–2 weeks | `1w` / `2w` |
| last month | `1m` |
| 2 months | `2m` |
| 3 months / quarter | `3m` |
| 6 months | `6m` |
| 9 months | `9m` |
| last year | `1y` |
| 2–3 years | `2y` / `3y` |
| 5 years / long term | `5y` |
| 10 years | `10y` |
| all time | `max` |

---

## Provider modes

| Mode | Behavior |
|------|----------|
| `auto` (default) | NASDAQ first, falls back to Yahoo Finance on errors |
| `nasdaq` | NASDAQ only |
| `legacy` | Yahoo Finance only (search not available) |

Fallback warnings print to stderr (suppressed with `--json`).

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `XGG_PROVIDER` | `auto` | Provider mode |
| `XGG_RATE_LIMIT` | `2` | Max requests/second |
| `XGG_MAX_RETRIES` | `3` | Retry attempts |
| `XGG_RETRY_DELAY` | `2s` | Delay between retries |
| `XGG_TIMEOUT` | `30s` | HTTP timeout |
| `XGG_WATCHLIST_TYPE` | `Rv` | NASDAQ watchlist type |

CLI flags take precedence over env vars.

---

## Installation / PATH issues

If `xgg` is not in PATH, follow this decision tree:

**1. Check the current working directory**

```bash
# Linux/macOS
ls ./xgg 2>/dev/null && echo "found"

# Windows (bash)
ls ./xgg.exe 2>/dev/null && echo "found"
```

If found → use `./xgg` (or `./xgg.exe`) as the command for the rest of the conversation.

**2. If not found — offer to download**

Tell the user:
> `xgg` was not found in PATH or in the current folder. Download the latest release here?
> → https://github.com/XungungoMarkets/Xungungo-CLI/releases

**3. After user confirms — install to current directory**

Linux/macOS:
```bash
curl -L https://github.com/XungungoMarkets/Xungungo-CLI/releases/latest/download/xgg-linux-amd64.tar.gz | tar xz
[ ! -f xgg ] && mv xgg-* xgg 2>/dev/null || true
chmod +x xgg && ./xgg version
```

Windows (bash):
```bash
curl -L -o xgg-windows.zip https://github.com/XungungoMarkets/Xungungo-CLI/releases/latest/download/xgg-windows-amd64.zip
unzip -o xgg-windows.zip
[ ! -f xgg.exe ] && mv xgg-*.exe xgg.exe 2>/dev/null || true
./xgg.exe version
```

Use `./xgg` or `./xgg.exe` for all commands after install.

---

## sectors command

Average daily % change for NASDAQ-listed stocks grouped by sector (or sector + industry).

```bash
xgg sectors [--by-industry] [--json]
```

| Flag | Description |
|------|-------------|
| `--by-industry` | Group by sector **and** industry |
| `--by-stock` | Show individual stocks within each sector (with per-stock % change) |

Optional sector filter argument (substring match):
```bash
xgg sectors --by-stock Energy
```

**Examples:**
```bash
xgg sectors
xgg sectors --by-industry
xgg sectors --by-stock
xgg sectors --by-stock Technology
```

Use when the user asks "how are sectors doing", "which sectors are up/down", "tech vs energy today", "show me stocks in the energy sector", or wants market breadth data.

---

## country command

Average daily % change for NASDAQ-listed stocks grouped by country.

```bash
xgg country [country name...] [--by-stock] [--json]
```

| Flag | Description |
|------|-------------|
| `--by-stock` | Show individual stocks within each country |

**Examples:**
```bash
xgg country
xgg country --by-stock
xgg country --by-stock uruguay
xgg country --by-stock hong kong
```

Country name is case-insensitive and can be multiple words. Use when the user asks about stocks by country or region.

---

## screener command

Full raw NASDAQ screener table with all fields for every listed stock.

```bash
xgg screener [SYMBOL...] [--sector S] [--country C] [--industry I] [--json]
```

| Flag | Description |
|------|-------------|
| `--sector` | Filter by sector (substring match, case-insensitive) |
| `--country` | Filter by country (substring match, case-insensitive) |
| `--industry` | Filter by industry (substring match, case-insensitive) |

**Columns returned:** `SYM`, `Name`, `Price`, `Change`, `Chg%`, `Volume`, `MktCap`, `Country`, `IPO`, `Sector`, `Industry`

**Examples:**
```bash
xgg screener                                   # all ~5000+ NASDAQ stocks
xgg screener AAPL MSFT NVDA                    # specific symbols
xgg screener --sector Technology               # all tech stocks
xgg screener --country USA --sector Energy     # US energy stocks
xgg screener --industry "Software" --json      # software stocks as JSON
```

Use when the user wants to browse all listed stocks, filter stocks by sector/country/industry, or get raw market data for a large set of stocks.

---

## ETF search tip

Add "ETF" to the query — `xgg search` returns no symbols for broad terms alone:
- ✅ `xgg search "semiconductor ETF"` → SOXX, SOXL, etc.
- ❌ `xgg search "semiconductors"` → only index solutions (no tickers)

If search still returns no symbols, use known ETFs directly:
- Semiconductors: `SOXX`, `SOXL`, `SMH`, `SOXQ`
- Tech: `QQQ`, `XLK`, `VGT`
- S&P 500: `SPY`, `VOO`, `IVV`
- Energy: `XLE`, `VDE`
- Financials: `XLF`, `VFH`
