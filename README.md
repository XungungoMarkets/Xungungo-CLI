# 📈 Xungungo-CLI (xgg)

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Xungungo-CLI** es una herramienta de línea de comandos escrita en Go para acceder a datos financieros en tiempo real. Cotizaciones de acciones, datos históricos y más, directamente desde tu terminal.

---

## 📋 Índice

- [Características](#características)
- [Instalación](#instalación)
- [Uso](#uso)
- [Comandos](#comandos)
- [Requisitos](#requisitos)
- [Licencia](#licencia)

---

## ✨ Características

- 📊 **Cotizaciones en tiempo real**: Obtén precios actuales de cualquier acción
- 📈 **Datos históricos**: Accede a datos OHLCV (Open, High, Low, Close, Volume)
- 🎨 **Interfaz visual**: Colores y formateo elegante en la terminal
- ⚡ **Rápido**: Múltiples símbolos en una sola solicitud
- 📅 **Periodos configurables**: 5 días, 1 mes, 3 meses, 6 meses, 1 año, 5 años

---

  ## 🔽 Instalación

### Desde GitHub Releases (Recomendado) ⭐

Descarga el binario pre-compilado para tu sistema operativo desde la página de [Releases](https://github.com/XungungoMarkets/xgg/releases).

**Linux/macOS:**
```bash
# Descargar el archivo (reemplaza VERSION con la versión deseada)
wget https://github.com/XungungoMarkets/xgg/releases/download/vVERSION/xgg-linux-amd64.tar.gz

# Extraer
tar xzf xgg-linux-amd64.tar.gz

# Mover a PATH
sudo mv xgg /usr/local/bin/
```

**macOS (Apple Silicon):**
```bash
# Descargar
wget https://github.com/XungungoMarkets/xgg/releases/download/vVERSION/xgg-darwin-arm64.tar.gz

# Extraer
tar xzf xgg-darwin-arm64.tar.gz

# Mover a PATH
sudo mv xgg /usr/local/bin/
```

**Windows:**
```powershell
# Descargar el archivo (reemplaza VERSION con la versión deseada)
# https://github.com/XungungoMarkets/xgg/releases/download/vVERSION/xgg-windows-amd64.zip

# Extraer y mover a una carpeta en tu PATH
```

**Verificar checksums (opcional):**
```bash
sha256sum -c xgg-linux-amd64.tar.gz.sha256
```

### Desde Go

```bash
go install github.com/XungungoMarkets/xgg@latest
```

Esto instalará el binario `xgg` en tu `$GOPATH/bin`. Asegúrate de que esté en tu `$PATH`.

### Desde Source

```bash
git clone https://github.com/XungungoMarkets/xgg.git
cd xgg
go build -o xgg
```

Mueve el binario a tu PATH:

```bash
# Linux/macOS
sudo mv xgg /usr/local/bin/

# Windows
# Mueve xgg.exe a una carpeta en tu PATH
```

---

## 🚀 Uso

### Ver Ayuda

```bash
# Ayuda general
xgg --help

# Ayuda de comando específico
xgg stock --help
xgg history --help
```

---

## 📖 Comandos

### Cotizaciones en Vivo

Obtén el precio actual de una acción:

```bash
xgg stock NVDA
```

**Salida:**
```
┌─────────────────────────────────────────┐
│  NVDA - NVIDIA Corporation             │
│  $875.28  ▲ +12.34 (+1.43%)             │
│  Vol: 45.2M  │  Mkt Cap: 2.2T          │
└─────────────────────────────────────────┘
```

Obtén cotizaciones de múltiples acciones a la vez:

```bash
xgg stock NVDA AAPL TSLA MSFT GOOGL
```

### Datos Históricos

Obtén datos históricos del último mes (por defecto):

```bash
xgg history NVDA
```

**Salida:**
```
  NVDA Historical Data
  Date          Open      High       Low     Close       Volume
  ──────────────────────────────────────────────────────────────────
  2026-03-10  $850.00  $880.00  $845.00  $875.28     45,234,567 ▲
  2026-03-09  $835.00  $852.00  $830.00  $848.50     52,123,456 ▲
  2026-03-08  $820.00  $840.00  $815.00  $834.20     48,567,890 ▼
  ...
```

Obtén datos históricos con un periodo específico:

```bash
# Últimos 5 días
xgg history NVDA --period 5d

# Últimos 6 meses
xgg history AAPL --period 6m

# Último año
xgg history TSLA --period 1y

# Últimos 5 años
xgg history MSFT --period 5y
```

**Periodos disponibles:**
- `5d` - Últimos 5 días
- `1m` - Último mes (default)
- `3m` - Últimos 3 meses
- `6m` - Últimos 6 meses
- `1y` - Último año
- `5y` - Últimos 5 años

### Formato Rápido

Para obtener información sin los bordes decorativos, puedes redirigir la salida o usarla en scripts:

```bash
xgg stock NVDA | grep -E "NVDA|Price"
```

---

## 🛠️ Requisitos

- **Go**: 1.26 o superior (solo para compilar desde source)
- **Sistema Operativo**: Linux, macOS, Windows
- **Conexión a Internet**: Requerida para obtener datos financieros

---

## 📝 Licencia

Este proyecto está licenciado bajo la Licencia MIT.

---

## 🔗 Links

- [GitHub Repository](https://github.com/XungungoMarkets/xgg)
- [Issues](https://github.com/XungungoMarkets/xgg/issues)

---

**Hecho con ❤️ por Xungungo Markets**

*Accede a los mercados financieros desde tu terminal, fácil y rápido.*