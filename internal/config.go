// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/21

package internal

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config 当前程序的配置
type Config struct {
	// Proxy 代理服务器地址，可选
	// 若为空，会使用环境变量中的 Proxy 配置
	Proxy string

	// GoProxy 可选
	// 若为空 会读取 go env GOPROXY 的值
	GoProxy string

	// TarURLPrefix 下载 go 打包文件的 url 地址前缀，可选
	// 为空时使用默认值 "https://dl.google.com/go/"
	TarURLPrefix string

	// InsecureSkipVerify 是否跳过证书校验
	InsecureSkipVerify bool
}

func (c *Config) getProxy() func(*http.Request) (*url.URL, error) {
	if len(c.Proxy) == 0 {
		return http.ProxyFromEnvironment
	}
	return func(request *http.Request) (*url.URL, error) {
		return url.Parse(c.Proxy)
	}
}

func (c *Config) getTarURLs(fp string) []string {
	ps := c.getTarURLPrefix()
	var result []string
	for _, p := range ps {
		p = strings.TrimSpace(p)
		if len(p) == 0 {
			continue
		}
		var b strings.Builder
		b.WriteString(p)
		if !strings.HasSuffix(p, "/") {
			b.WriteString("/")
		}
		b.WriteString(fp)
		result = append(result, b.String())
	}
	return result
}

func (c *Config) trySetProxyEnv() {
	if len(c.Proxy) == 0 {
		return
	}
	os.Setenv("HTTP_PROXY", c.Proxy)
	os.Setenv("http_proxy", c.Proxy)
	os.Setenv("HTTPS_PROXY", c.Proxy)
	os.Setenv("https_proxy", c.Proxy)
}

var tarURLPrefixDefault = []string{
	"https://dl.google.com/go/",
	"https://dl-ssl.google.com/go/", // 部分不能使用 tls 的尝试这个
}

func (c *Config) getTarURLPrefix() []string {
	if len(c.TarURLPrefix) > 0 {
		return strings.Split(c.TarURLPrefix, ",")
	}
	return tarURLPrefixDefault
}

var defaultConfig = &Config{}

func loadConfig() {
	fp := filepath.Join(DataDir(), "app.toml")
	logPrint("config", fp)
	content, err := os.ReadFile(fp)
	if err != nil && os.IsNotExist(err) {
		_ = os.WriteFile(fp, []byte(cfgTpl), 0644)
		return
	}
	var cfg *Config
	if err = toml.Unmarshal(content, &cfg); err != nil {
		logPrint("config", "ignored,parser", fp, "failed,", err)
		return
	}
	cfg.Proxy = strings.TrimSpace(cfg.Proxy)
	cfg.TarURLPrefix = strings.TrimSpace(cfg.TarURLPrefix)
	defaultConfig = cfg
	cfg.trySetProxyEnv()
}

func printProxy() {
	req, _ := http.NewRequest(http.MethodGet, tarURLPrefixDefault[0], nil)
	proxyFn := defaultConfig.getProxy()
	pu, err := proxyFn(req)
	if err != nil {
		logPrint("proxy", "parser proxy failed:", err)
	} else if pu != nil {
		logPrint("proxy", pu.String())
	}
}

var cfgTpl = `
# smart-go-dl
# https://github.com/fsgo/smart-go-dl

# 下载时使用的 Proxy，可选
# 不配置或者为空时，会使用环境变量的代理配置
# Proxy="http://127.0.0.1:8128"

# 下载文件时，是否跳过证书校验，可选，默认 false
# InsecureSkipVerify = true

# 下载 Go tar 文件的地址前缀，可选
# 默认值是 "https://dl-ssl.google.com/go/"
#TarURLPrefix="https://dl-ssl.google.com/go/"

`
