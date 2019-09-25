package logfmt_test

import (
	"bytes"
	"errors"
	"github.com/cep21/log/logfmt"
	"io/ioutil"
	"math"
	"testing"

	"github.com/cep21/log"
	logfmtbase "github.com/go-logfmt/logfmt"
)

func benchmarkRunner(b *testing.B, logger log.Logger, f func(log.Logger)) {
	lc := log.With(logger, "common_key", "common_value")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(lc)
	}
}

var (
	baseMessage = func(logger log.Logger) { logger.Log("foo_key", "foo_value") }
	withMessage = func(logger log.Logger) { log.With(logger, "a", "b").Log("c", "d") }
)
// These test are designed to be run with the race detector.

func testConcurrency(t *testing.T, logger log.Logger, total int) {
	n := int(math.Sqrt(float64(total)))
	share := total / n

	errC := make(chan error, n)

	for i := 0; i < n; i++ {
		go func() {
			errC <- spam(logger, share)
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errC
		if err != nil {
			t.Fatalf("concurrent logging error: %v", err)
		}
	}
}

func spam(logger log.Logger, count int) error {
	for i := 0; i < count; i++ {
		err := logger.Log("key", i)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestLogfmtLogger(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	logger := logfmt.NewLogfmtLogger(buf)

	if err := logger.Log("hello", "world"); err != nil {
		t.Fatal(err)
	}
	if want, have := "hello=world\n", buf.String(); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	buf.Reset()
	if err := logger.Log("a", 1, "err", errors.New("error")); err != nil {
		t.Fatal(err)
	}
	if want, have := "a=1 err=error\n", buf.String(); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	buf.Reset()
	if err := logger.Log("std_map", map[int]int{1: 2}, "my_map", mymap{0: 0}); err != nil {
		t.Fatal(err)
	}
	if want, have := "std_map=\""+logfmtbase.ErrUnsupportedValueType.Error()+"\" my_map=special_behavior\n", buf.String(); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
}

func BenchmarkLogfmtLoggerSimple(b *testing.B) {
	benchmarkRunner(b, logfmt.NewLogfmtLogger(ioutil.Discard), baseMessage)
}

func BenchmarkLogfmtLoggerContextual(b *testing.B) {
	benchmarkRunner(b, logfmt.NewLogfmtLogger(ioutil.Discard), withMessage)
}

func TestLogfmtLoggerConcurrency(t *testing.T) {
	t.Parallel()
	testConcurrency(t, logfmt.NewLogfmtLogger(ioutil.Discard), 10000)
}

type mymap map[int]int

func (m mymap) String() string { return "special_behavior" }
