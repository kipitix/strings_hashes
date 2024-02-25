package hashlog

import (
	formatter "github.com/fabienm/go-logrus-formatters"
	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	log "github.com/sirupsen/logrus"
)

type LogCfg struct {
	LogLevel    int    `arg:"--log-level,env:LOG_LEVEL" default:"4" help:"0-panic, 1-fatal, 2-error, 3-warn, 4-info, 5-debug, 6-trace"`
	LogGELF     bool   `arg:"--log-gelf,env:LOG_GELF" default:"false" help:"Enable of disable GELF format of logs"`
	LogURL      string `arg:"--log-url,env:LOG_URL" default:"localhost:12201" help:"Host and port of log server, keep it empty to disable sending logs to server"`
	LogHostname string `arg:"--log-hostname,env:LOG_HOSTNAME" default:"localhost" help:"Name of instance in logs"`
}

func InitLog(cfg LogCfg) {
	if cfg.LogGELF {
		gelfFmt := formatter.NewGelf(cfg.LogHostname)
		log.SetFormatter(gelfFmt)
	}

	log.SetLevel(log.Level(cfg.LogLevel))

	if cfg.LogGELF && cfg.LogURL != "" {
		hook := graylog.NewGraylogHook(cfg.LogURL, map[string]interface{}{})
		hook.Level = log.Level(cfg.LogLevel)
		log.AddHook(hook)
	}
}
