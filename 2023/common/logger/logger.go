// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// logHandler is a simplistic handler that ignores groups and attributes. Its
// sole purpose is to format log lines in a compact way.
type logHandler struct {
	minLevel slog.Leveler
	logTime  bool
	mu       *sync.Mutex
	out      io.Writer
}

func newLogHandler(out io.Writer, minLevel slog.Leveler) *logHandler {
	return &logHandler{
		minLevel: minLevel,
		logTime:  false,
		out:      out,
		mu:       &sync.Mutex{},
	}
}

func (h *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.minLevel.Level()
}

func levelToChar(level slog.Level) byte {
	switch {
	case level < slog.LevelInfo:
		return 'D'
	case level == slog.LevelInfo:
		return 'I'
	case level == slog.LevelWarn:
		return 'W'
	case level == slog.LevelError:
		return 'E'
	case level >= 10:
		return '*'
	default:
		return byte('0' + int(level))
	}
}

func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	buf := make([]byte, 0, 1024)
	buf = append(buf, levelToChar(r.Level))

	if h.logTime && !r.Time.IsZero() {
		buf = fmt.Appendf(buf, "%s", r.Time.Format(time.RFC3339))
	}
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf = fmt.Appendf(buf, fmt.Sprintf(" %s:%d", filepath.Base(f.File), f.Line))
	}
	buf = append(buf, ' ')
	buf = append(buf, []byte(r.Message)...)
	buf = append(buf, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}

func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *logHandler) WithGroup(name string) slog.Handler {
	return h
}

func Init(log bool) {
	level := slog.Level(999)
	if log {
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(newLogHandler(os.Stderr, level)))
}

func Infof(format string, args ...any) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(format, args...), pcs[0])
	_ = slog.Default().Handler().Handle(context.Background(), r)
}

func Errorf(format string, args ...any) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Errorf]
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(format, args...), pcs[0])
	_ = slog.Default().Handler().Handle(context.Background(), r)
}

func Fatalf(format string, args ...any) {
	Errorf(format, args...)
	os.Exit(1)
}
