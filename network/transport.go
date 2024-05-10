package network

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/tonny-zhang/goutils/env"
	"github.com/tonny-zhang/goutils/logger"
	"golang.org/x/net/proxy"
)

var globalTransport *http.Transport
var proxyForWebsocket func(*http.Request) (*url.URL, error)
var loggerNetwork = logger.PrefixLogger("[network]")
var transportSync sync.Once

func initData() {
	transportSync.Do(func() {
		env.AutoLoad()
		globalTransport = &http.Transport{
			// IdleConnTimeout:       90 * time.Second,
			// TLSHandshakeTimeout:   10 * time.Second,
			// ExpectContinueTimeout: 1 * time.Second,
			// ResponseHeaderTimeout: 30 * time.Second,
		}
		if v := os.Getenv("PROXY_URL"); v != "" {
			// 使用代理
			if u, e := url.Parse(v); e == nil {
				if u.Scheme == "socks5h" {
					dialer, err := proxy.SOCKS5("tcp", u.Host, nil, proxy.Direct)
					if err != nil {
						fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
						os.Exit(1)
					}
					globalTransport.DialContext = func(ctx context.Context, network, addr string) (c net.Conn, e error) {
						c, e = dialer.Dial(network, addr)
						return
					}
				} else {
					globalTransport.Proxy = http.ProxyURL(u)
				}

				loggerNetwork.Info("使用代理: %v", v)
			}
		}

		if v := os.Getenv("PROXY_URL_WS"); v != "" {
			if u, e := url.Parse(v); e == nil {
				proxyForWebsocket = http.ProxyURL(u)
			}

		}
	})
}

// GetTransport 得到共用的transport
func GetTransport() *http.Transport {
	initData()
	return globalTransport
}

// GetProxyForWebsocket 得到websocket使用的代理
func GetProxyForWebsocket() func(*http.Request) (*url.URL, error) {
	return proxyForWebsocket
}
