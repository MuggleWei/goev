package demo

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLog(logLevel, logFile string, enableConsole bool) error {
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		return err
	}

	var logger *lumberjack.Logger = nil
	if logFile != "" {
		logger = &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    128, // megabytes
			MaxBackups: 8,
			MaxAge:     28, //days
			Compress:   false,
		}
	}

	if !enableConsole {
		if logger == nil {
			// don't output
		} else {
			log.SetOutput(logger)
		}
	} else {
		if logger == nil {
			log.SetOutput(os.Stdout)
		} else {
			log.SetOutput(io.MultiWriter(os.Stdout, logger))
		}
	}
	log.SetReportCaller(true)
	log.SetFormatter(&LogFormat{})
	log.SetLevel(lvl)

	return nil
}

type LogFormat struct{}

func (this *LogFormat) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format(time.RFC3339)
	goid := this.Goid()

	var newLog string
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("%s|%s|%d|%s:%d %s - %s\n",
			entry.Level, timestamp, goid, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		newLog = fmt.Sprintf("%s|%s - %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func (this *LogFormat) Goid() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		//panic(fmt.Sprintf("cannot get goroutine id: %v", err))
		return 0
	}
	return id
}
