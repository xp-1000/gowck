/*
Copyright © 2020 XP-1000 <xp-1000@hotmail.fr>

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package httpcheck

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/creasty/defaults"
)

// Config of the monitor
type Config struct {
	Body            string            `yaml:"body"`
	FollowRedirects bool              `yaml:"followRedirects" default:"true"`
	Headers         map[string]string `yaml:"headers"`
	Method          string            `yaml:"method" default:"GET"`
	Timeout         int               `yaml:"timeout" default:"5"`
	URLs            []string          `yaml:"urls"`
	URL             string            `yaml:"url" default:"https://cloud.manfroi.fr"`
	Regex           string            `yaml:"regex"`
}

func monitor() {
	config := &Config{}
	if err := defaults.Set(config); err != nil {
		panic(fmt.Sprintf("Config defaults are wrong types: %s", err))
	}
	parsedURL, err := url.Parse(config.URL)
	if err != nil {
		panic(fmt.Sprintf("Could not parse url: %s", err))
	}
	tlsValid := true
	host := parsedURL.Hostname()
	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	timeoutDuration := time.Duration(config.Timeout) * time.Second
	dialer := &net.Dialer{Timeout: timeoutDuration}
	client := &http.Client{
		Transport: &http.Transport{
			DialContext:       dialer.DialContext,
			DisableKeepAlives: true,
			TLSClientConfig:   tlsCfg,
		},
		Timeout: timeoutDuration,
	}
	if !config.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	var body io.Reader
	if config.Body != "" {
		body = strings.NewReader(config.Body)
	}
	req, err := http.NewRequest(config.Method, config.URL, body)
	if err != nil {
		panic(fmt.Sprintf("Could not create new request: %s", err))
	}
	for key, val := range config.Headers {
		req.Header.Add(key, val)
		if key == "Host" {
			req.Host = val
		}
	}
	now := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Could not do request: %s", err))
	}
	responseTime := time.Since(now).Seconds()
	if err == nil && parsedURL.Scheme == "https" {
		port := parsedURL.Port()
		if port == "" {
			port = "443"
		}
		conn, err := tls.DialWithDialer(dialer, "tcp", host+":"+port, tlsCfg)
		if err != nil {
			panic(fmt.Sprintf("Could not dial server: %s", err))
		}
		defer conn.Close()
		err = conn.Handshake()
		if err != nil {
			panic(fmt.Sprintf("Could not handshake: %s", err))
		}
		certs := conn.ConnectionState().PeerCertificates
		for i, cert := range certs {
			opts := x509.VerifyOptions{
				Intermediates: x509.NewCertPool(),
			}
			if i == 0 {
				opts.DNSName = host
				for j, cert := range certs {
					if j != 0 {
						opts.Intermediates.AddCert(cert)
					}
				}
				tlsLeft := cert.NotAfter.Sub(now).Seconds()
				fmt.Println(tlsLeft)
			}
			_, err := cert.Verify(opts)
			if err != nil {
				tlsValid = false
				fmt.Println(err)
			}
		}

	}
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Could not parse body: %s", err))
	}
	contentLength := len(bodyBytes)
	match := false
	if config.Regex != "" {
		regex, err := regexp.Compile(config.Regex)
		if err != nil {
			panic(fmt.Sprintf("Failed to compile regular expression %s", err))
		}
		if regex.Match(bodyBytes) {
			match = true
		}
	} else {
		match = true
	}
	fmt.Println("status code:")
	fmt.Println(statusCode)
	fmt.Println("response time:")
	fmt.Println(responseTime)
	fmt.Println("body size:")
	fmt.Println(contentLength)
	fmt.Println("regex match:")
	fmt.Println(match)
	fmt.Println("tls valid:")
	fmt.Println(tlsValid)
}
