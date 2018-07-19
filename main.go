package main

import (
	"github.com/adwpc/logmetrics/conf"
	"github.com/adwpc/logmetrics/parser"
	"github.com/adwpc/logmetrics/zlog"
)

var (
	log = zlog.Log
	cfg = conf.GetConfig()
)

func main() {
	log.Info().Msg("main")
	if !cfg.Parse() {
		log.Panic().Msg("parse conf.toml error!")
		return
	}
	parser.Monitor(cfg)

}
