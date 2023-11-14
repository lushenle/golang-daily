// A HTTP client that dials through a SOCKS5 proxy using username/password
// authentication. The auth is configured manually through a transport, without
// relying on parsing http_proxy env vars.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

func main() {
	target := flag.String("target", "http://example.org", "URL to get")
	proxyAddr := flag.String("proxy", "localhost:1080", "SOCKS5 proxy address to use")
	username := flag.String("user", "", "username for SOCKS5 proxy")
	password := flag.String("pass", "", "password for SOCKS5 proxy")
	flag.Parse()

	auth := proxy.Auth{
		User:     *username,
		Password: *password,
	}
	dialer, err := proxy.SOCKS5("tcp", *proxyAddr, &auth, nil)
	if err != nil {
		log.Fatal(err)
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.(proxy.ContextDialer).DialContext(ctx, network, addr)
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	r, err := client.Get(*target)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
