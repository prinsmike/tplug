package tplug

import (
	"os"
	"path/filepath"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("tplug", caddy.Plugin{
		ServerType: "http",
		Action:     Setup,
	})
}

func Setup(c *caddy.Controller) (err error) {
	cfg := httpserver.GetConfig(c)
	var config *Config

	config, err = ParseTPlugConfig(c, cfg)
	if err != nil {
		return err
	}

	tplug := &TPlug{
		Config: config,
	}

	cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		tplug.Next = next
		return tplug
	})

	return
}

type Config struct {
	HostName string
	Endpoint string
	SiteRoot string
	Template *template.Template
}

func ParseTPlugConfig(c *caddy.Controller, cnf *httpserver.SiteConfig) (*Config, error) {
	conf := &Config{
		HostName: cnf.Host(),
		Endpoint: `/tplug`,
		SiteRoot: cnf.Root,
		Template: nil,
	}

	_, err := os.Stat(conf.SiteRoot)
	if err != nil {
		return nil, c.Err("[tplug]: `invalid root directory`")
	}

	for c.Next() {
		args := c.RemainingArgs()

		if len(args) == 1 {
			conf.Endpoint = args[0]
		}

		for c.NextBlock() {
			switch c.Val() {
			case "template":
				var err error
				if c.NextArg() {
					conf.Template, err = template.ParseFiles(filepath.Join(conf.SiteRoot, c.Val()))
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	if conf.Template == nil {
		var err error
		conf.Template, err = template.New("tplug").Parse(defaultTemplate)
		if err != nil {
			return nil, err
		}
	}

	return conf, nil
}

const defaultTemplate = `<!DOCTYPE html>
<html>
	<head>
		<title>TPlug</title>
		<meta charset="utf-8">
		<style>
		body {
			padding: 1% 2%;
			font: 16px Arial;
		}
		</style>
	</head>
	<body>
		<h1>TPlug</h1>
		<p>
			Hello from TPlug, {{$.Req.Host}}.
		</p>
	</body>
</html>`
