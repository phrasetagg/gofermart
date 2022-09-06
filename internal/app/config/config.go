package config

import (
	"flag"
	"github.com/phrasetagg/gofermart/internal/app/helpers"
)

type cfg struct {
	ServerAddr  string
	DBDsn       string
	AccrualAddr string
}

func PrepareCfg() *cfg {
	cfg := new(cfg)

	flag.StringVar(&cfg.ServerAddr, "a", helpers.GetEnv("RUN_ADDRESS", "localhost:8080"), "host:port of the server")
	flag.StringVar(&cfg.DBDsn, "d", helpers.GetEnv("DATABASE_URI", ""), "Database DSN")
	flag.StringVar(&cfg.AccrualAddr, "r", helpers.GetEnv("ACCRUAL_SYSTEM_ADDRESS", ""), "accrual system address")

	flag.Parse()

	return cfg
}
