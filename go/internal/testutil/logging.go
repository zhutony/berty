package testutil

import (
	"flag"
	"os"
	"strconv"
	"testing"

	ipfs_log "github.com/ipfs/go-log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var debug = flag.Bool("debug", false, "is more verbose logging")

func Logger(t *testing.T) *zap.Logger {
	t.Helper()

	bertyDebug := parseBoolFromEnv("BERTY_DEBUG") || *debug
	libp2pDebug := parseBoolFromEnv("LIBP2P_DEBUG")
	// @NOTE(gfanton): since orbitdb use `zap.L()`, this will only
	// replace zap global logger with our logger
	orbitdbDebug := parseBoolFromEnv("ORBITDB_DEBUG")
	bertylogfile := os.Getenv("LOGFILE")

	isDebugEnabled := bertyDebug || orbitdbDebug || libp2pDebug
	if !isDebugEnabled {
		return zap.NewNop()
	}

	// setup zap config
	config := zap.NewDevelopmentConfig()
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level.SetLevel(zap.DebugLevel)
	if bertylogfile != "" {
		config.OutputPaths = []string{bertylogfile}
	}

	// build logger
	logger, err := config.Build()
	if err != nil {
		t.Errorf("setup debug logger error: `%v`", err)
		return zap.NewNop()
	}

	if libp2pDebug {
		ipfs_log.SetDebugLogging()
	}

	if orbitdbDebug {
		zap.ReplaceGlobals(logger)
	}

	if bertyDebug {
		return logger
	}

	return zap.NewNop()
}

func parseBoolFromEnv(key string) (b bool) {
	b, _ = strconv.ParseBool(os.Getenv(key))
	return
}
