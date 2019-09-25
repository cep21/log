package level_test

import (
	"errors"
	"github.com/cep21/log/logfmt"
	"os"

	"github.com/cep21/log"
	"github.com/cep21/log/level"
)

func Example_basic() {
	logger := logfmt.NewLogfmtLogger(os.Stdout)
	level.Debug(logger).Log("msg", "this message is at thte debug level")
	level.Info(logger).Log("msg", "this message is at the info level")
	level.Warn(logger).Log("msg", "this message is at the warn level")
	level.Error(logger).Log("msg", "this message is at the error level")

	// Output:
	// level=debug msg="this message is at thte debug level"
	// level=info msg="this message is at the info level"
	// level=warn msg="this message is at the warn level"
	// level=error msg="this message is at the error level"
}

func Example_filtered() {
	// Set up logger with level filter.
	logger := logfmt.NewLogfmtLogger(os.Stdout)
	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger, "caller", log.DefaultCaller)

	// Use level helpers to log at different levels.
	level.Error(logger).Log("err", errors.New("bad data"))
	level.Info(logger).Log("event", "data saved")
	level.Debug(logger).Log("next item", 17) // filtered

	// Output:
	// level=error caller=example_test.go:33 err="bad data"
	// level=info caller=example_test.go:34 event="data saved"
}
