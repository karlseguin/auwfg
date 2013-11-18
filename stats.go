package auwfg

import (
  "os"
  "time"
  "bytes"
  "strconv"
  "runtime"
  "sync/atomic"
)

type stats struct {
  errors uint64
  fatals uint64
  requests uint64
}

var Stats = &stats{}

func startStats(c *Configuration) {
  go snapshot(c)
}

func (s *stats) Request() { atomic.AddUint64(&s.requests, 1) }
func (s *stats) Error() { atomic.AddUint64(&s.errors, 1) }
func (s *stats) Fatal() { atomic.AddUint64(&s.fatals, 1) }

func snapshot(c *Configuration) {
  var last = new(stats)
  buffer := new(bytes.Buffer)
  for {
    time.Sleep(c.statsSleep)
    requests := atomic.LoadUint64(&Stats.requests)
    errors := atomic.LoadUint64(&Stats.errors)
    fatals := atomic.LoadUint64(&Stats.fatals)

    buffer.Reset()
    buffer.WriteString(`{"requests":` + strconv.FormatUint(requests - last.requests , 10))
    buffer.WriteString(`,"errors":` + strconv.FormatUint(errors - last.errors, 10))
    buffer.WriteString(`,"fatals":` + strconv.FormatUint(fatals - last.fatals, 10))
    buffer.WriteString(`,"goroutines":` + strconv.Itoa(runtime.NumGoroutine()))
    buffer.WriteString("}")

    last.requests = requests
    last.errors = errors
    last.fatals = fatals

    file, _ := os.Create(c.statsFile)
    file.Write(buffer.Bytes())
    file.Close()
  }
}
