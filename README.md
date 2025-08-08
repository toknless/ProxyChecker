# ProxyChecker

ProxyChecker is a **fast, concurrent proxy testing tool** written in Go.  
It reads proxies from a `proxies.txt` file, determines their working protocol,  
saves HTTP/HTTPS ones to `working_proxies.txt`, removes dead proxies from the list,  
and can optionally save dead proxies to a separate file.  

It supports multiple protocols (`http`, `https`, `socks4`, `socks5`) and can  
automatically guess the protocol based on the port number for faster checking.  

---

## 📜 Description
This repository contains **ProxyChecker**, a high-performance tool for validating and sorting proxy servers.  
It’s optimized for speed using Go’s concurrency features, allowing you to check thousands of proxies in seconds.  
Perfect for developers, testers, or anyone managing large proxy lists.  

With ProxyChecker you can:
- Automatically detect proxy protocols by port number
- Test all protocols or only one specific protocol
- Run checks with as many threads as you want
- Save dead proxies for later inspection (optional)
- Keep a clean `proxies.txt` with only working entries

---

## 📦 Prerequisites
Before using ProxyChecker, make sure you have:
- **Go 1.18+** installed  
  [Download Go here](https://go.dev/dl/)
- A `proxies.txt` file with one proxy per line in the format:
  ```
  ip:port
  ```

---

## 🚀 Installation

### Option 1: Download from GitHub
You can download the precompiled binary from this repository:
```
https://github.com/toknless/ProxyChecker/releases
```
Choose the binary for your operating system, make it executable, and run it.

### Option 2: Build from source
```bash
# Clone the repository
git clone https://github.com/toknless/ProxyChecker.git
cd ProxyChecker

# Build the program
go build proxy_checker.go
```

---

## ⚙️ Usage
```bash
./proxy_checker [flags]
```

### Command-line Flags
| Flag            | Description                                                        | Default                |
|-----------------|--------------------------------------------------------------------|------------------------|
| `-threads`      | Number of concurrent proxy checks                                  | CPU cores              |
| `-protocol`     | Test only one protocol (`http`, `https`, `socks4`, `socks5`)       | Test all               |
| `-save-dead`    | Save non-working proxies to `dead_proxies.txt`                     | Disabled               |
| `-auto-detect`  | Guess protocol from port before testing (80→HTTP, 443→HTTPS, etc.) | Disabled               |

---

## 📂 Output Files
- **`proxies.txt`** → Updated with only working proxies (any protocol)
- **`working_proxies.txt`** → Only working HTTP(S) proxies in `protocol:ip:port` format
- **`dead_proxies.txt`** → (Optional) Saved when `-save-dead` is used, contains non-working proxies

---

## 💡 Examples
```bash
# Default: test all protocols with threads = CPU cores
./proxy_checker

# Test all protocols with 50 threads
./proxy_checker -threads 50

# Test only HTTP proxies with 20 threads
./proxy_checker -threads 20 -protocol http

# Test only SOCKS5 proxies, save dead ones
./proxy_checker -protocol socks5 -save-dead

# Enable auto-detect mode for faster results
./proxy_checker -threads 100 -auto-detect
```

---

## 📜 License
MIT License — feel free to use, modify, and share.
