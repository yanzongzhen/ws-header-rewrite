package ws_header_rewrite

import (
	"context"
	"net/http"
	"strings"
)

type Config struct {
	Enabled bool `json:"enabled"`
}

type plugin struct {
	name   string
	next   http.Handler
	config *Config
}

func CreateConfig() *Config {
	return &Config{
		Enabled: false,
	}
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &plugin{
		name:   name,
		next:   next,
		config: config,
	}, nil
}

func (p *plugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if p.config.Enabled {
		for name, values := range rw.Header() {
			if strings.HasPrefix(strings.ToLower(name), "sec-websocket-") {
				rw.Header().Del(name)
				fixedName := "Sec-WebSocket-" + name[14:]
				rw.Header()[fixedName] = values
			}
		}
	}
	p.next.ServeHTTP(rw, req)
}
