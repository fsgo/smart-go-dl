// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/21

package internal

import (
	"io/ioutil"
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

func (c *Config) getTarUrL(fp string) string {
	var b strings.Builder
	p := c.getTarURLPrefix()
	b.WriteString(p)
	if !strings.HasSuffix(p, "/") {
		b.WriteString("/")
	}
	b.WriteString(fp)
	return b.String()
}

func (c *Config) trySetProxyEnv() {
	if len(c.Proxy) == 0 {
		return
	}
	os.Setenv("HTTP_PROXY", c.Proxy)
	os.Setenv("HTTPS_PROXY", c.Proxy)
}

const tarURLPrefixDefault = "https://dl.google.com/go/"

func (c *Config) getTarURLPrefix() string {
	if len(c.TarURLPrefix) > 0 {
		return c.TarURLPrefix
	}
	return tarURLPrefixDefault
}

func (c *Config) isDefaultRarURLPrefix() bool {
	return strings.Contains(c.getTarURLPrefix(), "dl.google.com")
}

var defaultConfig = &Config{}

func loadConfig() {
	fp := filepath.Join(DataDir(), "app.toml")
	logPrint("config", fp)
	content, err := os.ReadFile(fp)
	if err != nil && os.IsNotExist(err) {
		_ = ioutil.WriteFile(fp, []byte(cfgTpl), 0644)
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
	req, _ := http.NewRequest(http.MethodGet, tarURLPrefixDefault, nil)
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
# 默认值是 "https://dl.google.com/go/"
#TarURLPrefix="https://dl.google.com/go/"

`
