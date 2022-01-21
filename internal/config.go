// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/21

package internal

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	// Proxy 代理服务器地址，可选
	// 若为空，会使用环境变量中的 Proxy 配置
	Proxy string

	// TarURLPrefix 下载 go 打包文件的 url 地址前缀，可选
	// 为空时使用默认值 "https://dl.google.com/go/"
	TarURLPrefix string
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

var defaultConfig = &Config{}

func loadConfig() {
	fp := filepath.Join(TmpDir(), "app.toml")
	content, err := os.ReadFile(fp)
	if err != nil && os.IsNotExist(err) {
		return
	}
	var cfg *Config
	if err = toml.Unmarshal(content, &cfg); err != nil {
		log.Println("[ignore] parser", fp, "failed,", err)
		return
	}
	defaultConfig = cfg
	cfg.trySetProxyEnv()
}

func printProxy() {
	req, _ := http.NewRequest(http.MethodGet, tarURLPrefixDefault, nil)
	proxyFn := defaultConfig.getProxy()
	pu, err := proxyFn(req)
	if err != nil {
		log.Println("[proxy] error:", err)
	} else if pu != nil {
		log.Println("[proxy]", pu.String())
	}
}
