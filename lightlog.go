package lightlog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

//LogLevel 日志输出等级
type LogLevel int

const (
	LevelAll LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
	LevelNone
)

const (
	_ = iota
	//KB 1024 Byte
	KB uint64 = 1 << (iota * 10)
	//MB 1024 KByte
	MB
	//GB 1024 MByte
	GB
	//TB 1024 GByte
	TB
)

//颜色代码
var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

var (
	defaultConsoleOut = os.Stdout
	defaultFileOut    io.Writer
	defaultTimeFormat = "2006-01-02 15:04:05"
	defaultPrefix     = "lightlog"
	defaultLevel      = LevelInfo
)

//LogMsg 日志详细信息
type LogMsg struct {
	timestamp time.Time
	Level     LogLevel
	e         error
}

//Logger 日志记录
type Logger struct {
	Level      LogLevel
	ConsoleOut io.Writer
	FileOut    io.Writer
	queuen     chan *LogMsg
	Prefix     string
	TimeFormat string
}

func color(t LogLevel) string {
	switch t {
	case LevelDebug:
		return cyan
	case LevelInfo:
		return green
	case LevelWarning:
		return yellow
	default:
		return red
	}
}

func level(t LogLevel) string {
	switch t {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarning:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "FATAL"
	}
}

//LogWriter 日志顺序输出
func LogWriter(lg *Logger) {
	for {
		p := <-lg.queuen
		if p.Level >= lg.Level {
			c := color(p.Level)
			l := level(p.Level)
			if lg.ConsoleOut != nil {
				fmt.Fprintf(lg.ConsoleOut, "[%s] %s: %s [%s] %s %s\n", lg.Prefix, p.timestamp.Format(lg.TimeFormat), c, l, reset, p.e.Error())
			}
			if lg.FileOut != nil {
				fmt.Fprintf(lg.FileOut, "[%s] %s: [%s] %s\n", lg.Prefix, p.timestamp.Format(lg.TimeFormat), l, p.e.Error())
			}
		}
	}
}

//NewLogger 建立新的日志记录器
func NewLogger(bufferSize uint) *Logger {
	lg := &Logger{
		Level:      defaultLevel,
		ConsoleOut: defaultConsoleOut,
		FileOut:    defaultFileOut,
		queuen: func() chan *LogMsg {
			if bufferSize == 0 {
				return make(chan *LogMsg)
			}
			return make(chan *LogMsg, bufferSize)
		}(),
		Prefix:     defaultPrefix,
		TimeFormat: defaultTimeFormat,
	}
	go LogWriter(lg)
	return lg
}

//Log 日志记录入口
func (lg *Logger) Log(l LogLevel, e ...string) {
	if len(e) > 0 {
		lg.queuen <- &LogMsg{
			timestamp: time.Now(),
			Level:     l,
			e: func() error {
				s := ""
				for _, t := range e {
					s = fmt.Sprintf("%s %s", s, t)
				}
				s = strings.TrimSpace(s)
				return errors.New(s)
			}(),
		}
	}
}

//Debug 记录调试级日志
func (lg *Logger) Debug(e ...string) {
	lg.Log(LevelDebug, e...)
}

//Info 记录信息级日志
func (lg *Logger) Info(e ...string) {
	lg.Log(LevelInfo, e...)
}

//Warn 记录警告级日志
func (lg *Logger) Warn(e ...string) {
	lg.Log(LevelWarning, e...)
}

//Error 记录错误级日志
func (lg *Logger) Error(e ...string) {
	lg.Log(LevelError, e...)
}

//Fatal 记录崩溃错误级日志
func (lg *Logger) Fatal(e ...string) {
	lg.Log(LevelFatal, e...)
}
