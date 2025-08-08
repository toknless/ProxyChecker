package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type ProxyResult struct {
	Proxy    string
	Protocol string
	Working  bool
}

var (
	testURL         = "https://nextjs-ai-chatbot-mu-ten-84ftfptu6e.vercel.app/"
	uploadSize      = 8 * 1024              // 8 KB upload
	downloadSize    = int64(8 * 1024)       // 8 KB download
	uploadBody      = bytes.Repeat([]byte("A"), uploadSize)
)

func checkProxy(proxy, protocol string) ProxyResult {
	proxyURL, _ := url.Parse(fmt.Sprintf("%s://%s", protocol, proxy))
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	req, err := http.NewRequest("POST", testURL, bytes.NewBuffer(uploadBody))
	if err != nil {
		return ProxyResult{Proxy: proxy, Protocol: protocol, Working: false}
	}
	req.ContentLength = int64(len(uploadBody))

	resp, err := client.Do(req)
	if err != nil {
		return ProxyResult{Proxy: proxy, Protocol: protocol, Working: false}
	}
	defer resp.Body.Close()

	_, err = io.CopyN(io.Discard, resp.Body, downloadSize)
	if err != nil && err != io.EOF {
		return ProxyResult{Proxy: proxy, Protocol: protocol, Working: false}
	}

	return ProxyResult{Proxy: proxy, Protocol: protocol, Working: true}
}

func main() {
	threads := flag.Int("threads", 10, "Number of concurrent threads")
	protocolFlag := flag.String("protocol", "", "Force protocol: http or socks5")
	flag.Parse()

	file, err := os.Open("proxies.txt")
	if err != nil {
		fmt.Println("Error opening proxies.txt:", err)
		return
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			proxies = append(proxies, line)
		}
	}

	if len(proxies) == 0 {
		fmt.Println("No proxies found in proxies.txt")
		return
	}

	results := []ProxyResult{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, *threads)

	for _, proxy := range proxies {
		wg.Add(1)
		go func(proxy string) {
			defer wg.Done()
			sem <- struct{}{}

			protocolsToTry := []string{}
			if *protocolFlag != "" {
				protocolsToTry = []string{*protocolFlag}
			} else {
				protocolsToTry = []string{"http", "socks5"}
			}

			var workingResult ProxyResult
			for _, proto := range protocolsToTry {
				res := checkProxy(proxy, proto)
				if res.Working {
					workingResult = res
					break
				}
			}

			if workingResult.Working {
				mu.Lock()
				results = append(results, workingResult)
				mu.Unlock()
			}

			<-sem
		}(proxy)
	}

	wg.Wait()

	// Sort results before saving
	sort.Slice(results, func(i, j int) bool {
		return results[i].Proxy < results[j].Proxy
	})

	outFile, err := os.Create("working_proxies.txt")
	if err != nil {
		fmt.Println("Error creating working_proxies.txt:", err)
		return
	}
	defer outFile.Close()

	for _, r := range results {
		outFile.WriteString(fmt.Sprintf("%s:%s\n", r.Protocol, r.Proxy))
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("Total proxies tested: %d\n", len(proxies))
	fmt.Printf("Working proxies found: %d\n", len(results))
	fmt.Println("Results saved to working_proxies.txt")
}
