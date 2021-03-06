package main

import (
	"github.com/BurntSushi/toml"
	"github.com/nynicg/cake/lib/third/uuid"
)

const pf = "./configure_server.toml"

var globalConfig *ServConfig

type ServConfig struct {
	ProxyConfig ProxyConfig `toml:"proxy"`
	ApiConfig   ApiConfig   `toml:"api"`
}

type ProxyConfig struct {
	LocalAddr string   `toml:"localAddr"`
	MaxConn   int      `toml:"maxConn"`
	LogLevel  string   `toml:"logLevel"`
	Key       string   `toml:"key"`
	Uids      []string `toml:"uids"`
}

type ApiConfig struct {
	EnableApi         bool   `toml:"enableApi"`
	LocalApiAddr      string `toml:"localApiAddr"`
	BasicAuthUser     string `toml:"basicAuthUser"`
	BasicAuthPassword string `toml:"basicAuthPassword"`
}

func init() {
	globalConfig = &ServConfig{}
	if _, e := toml.DecodeFile(pf, globalConfig); e != nil {
		panic(e)
	}
}

func override(dst, src *ServConfig) {
	sp := src.ProxyConfig
	if sp.LogLevel != "" {
		dst.ProxyConfig.LogLevel = sp.LogLevel
	}
	if sp.Key != "" {
		dst.ProxyConfig.Key = sp.Key
	}
	if sp.LocalAddr != "" {
		dst.ProxyConfig.LocalAddr = sp.LocalAddr
	}
	if sp.MaxConn != 0 {
		dst.ProxyConfig.MaxConn = sp.MaxConn
	}

	sa := src.ApiConfig
	if sa.EnableApi {
		dst.ApiConfig.EnableApi = true
	}
	if sa.BasicAuthPassword != "" {
		dst.ApiConfig.BasicAuthPassword = sa.BasicAuthPassword
	} else if dst.ApiConfig.BasicAuthPassword == "" {
		uid ,e := uuid.NewV4()
		if e != nil{
			panic(e)
		}
		dst.ApiConfig.BasicAuthPassword = uid.String()
	}

	if sa.BasicAuthUser != "" {
		dst.ApiConfig.BasicAuthUser = sa.BasicAuthUser
	} else if dst.ApiConfig.BasicAuthUser == "" {
		uid ,e := uuid.NewV4()
		if e != nil{
			panic(e)
		}
		dst.ApiConfig.BasicAuthUser = uid.String()
	}

	if sa.LocalApiAddr != "" {
		dst.ApiConfig.LocalApiAddr = sa.LocalApiAddr
	}

}
