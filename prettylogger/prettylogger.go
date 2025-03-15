package plog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

//// ========================================================
//// ========================================================
//! Output requirements constants

// Requirements are bitmasks that showed are needed to be outputed in logs
const (
	RequireTimestamp = 1 << iota
	RequireLevel
	RequireCaller
	RequireAll = RequireCaller | RequireLevel | RequireTimestamp
)

//// ========================================================
//// ========================================================
//! Information levels

const (
	// Output all logs
	LevelDebug = iota

	// Output all logs except Debug logs
	LevelInfo

	// Only bad logs: Warnings, Errors and Fatal errors
	LevelWarn

	// Only critical messages: Errors and Fatal errors
	LevelError

	// Only fatal errors are outputed
	LevelFatal
)

// Return if given level is in range of possible values
func levelRange(level int) bool {
	return level >= 0 && level <= LevelFatal
}

//// ========================================================
//// ========================================================
//! Colors constants

// ANSI sybmols for a colored output
const (
	colorReset   = "\033[0m"
	colorBlue    = "\033[34m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorRed     = "\033[31m"
	colorBoldRed = "\033[1;31m"
)

//// ========================================================
//// ========================================================
//! Logger basic implementation

// Customly written logger struct for a pretty output
type Logger struct {
	// Used to block any modifications with same logger
	mu sync.Mutex

	// Depending on the level, certain log functions will do nothing
	level int

	// Logger will use fmt.Fprint() using an this output attribute
	output io.Writer

	// What logger will output additionally to message
	require int

	// If set to false, logger won't work
	enabled bool

	// If set to true, output will contain ansi symbols
	colored bool
}

/*
Return new logger.

Require a valid output level and output writer.

Requirements - additional info that will be outputed with each log (see Require constants).

ContainAnsi - if true log will use ansi symbols for coloring the output (might be incompatable with windows cmd or others consoles/files).
*/
func NewLogger(level int, output io.Writer, outputRequirements int, containAnsi bool) (*Logger, error) {
	if levelRange(level) {
		return &Logger{mu: sync.Mutex{}, level: level, output: output, require: outputRequirements, enabled: true, colored: containAnsi}, nil
	}
	return nil, fmt.Errorf("level value is out of range")
}

/*
Change logger's level to a given value.

In case of level being out of range, return error and do not change logger in any way.
*/
func (l *Logger) SetLevel(level int) error {
	if levelRange(level) {
		return fmt.Errorf("level value is out of range")
	}
	l.mu.Lock()
	l.level = level
	l.mu.Unlock()
	return nil
}

// Change logger's writer interface to a given one
func (l *Logger) SetWriter(output io.Writer) {
	l.mu.Lock()
	l.output = output
	l.mu.Unlock()
}

/*
Change logger's output requirements to a given one.

Expects usage of Require bitmask constants.
*/
func (l *Logger) SetRequirements(outputRequirements int) {
	l.mu.Lock()
	l.require = outputRequirements
	l.mu.Unlock()
}

// Change logger's enabled field to a given value
func (l *Logger) SetActivity(enabled bool) {
	l.mu.Lock()
	l.enabled = enabled
	l.mu.Unlock()
}

/*
Change logger's coloring ability.

If true, will use ANSI characters.
*/
func (l *Logger) SetColoring(containAnsi bool) {
	l.mu.Lock()
	l.colored = containAnsi
	l.mu.Unlock()
}

//// ========================================================
//// ========================================================
//! Information string getters

// Return string name of a current level
func getColoredLevel(level int) string {
	switch level {
	case LevelDebug:
		return colorBlue + "DEBUG" + colorReset
	case LevelInfo:
		return colorGreen + "INFO" + colorReset
	case LevelWarn:
		return colorYellow + "WARN" + colorReset
	case LevelError:
		return colorRed + "ERROR" + colorReset
	case LevelFatal:
		return colorBoldRed + "FATAL" + colorReset
	default:
		return "INVALID"
	}
}

// Return string name of a current level
func getRegularLevel(level int) string {
	switch level {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "INVALID"
	}
}

// Return current time as a string in date-time format
func getTimestamp() string {
	return time.Now().Format(time.DateTime)
}

/*
Return string or current caller function.

Skip - how many functions must be skipped before accessing the caller function information.
*/
func getCall(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s:%d (%s)", file, line, funcName)
}

//// ========================================================
//// ========================================================
//! Logger output implementation

// Using given values print them in the logger's output writer attribute
func (l *Logger) print(level int, format string, attrs ...any) error {
	if !l.enabled {
		return fmt.Errorf("logger is disabled")
	}
	if l.level > level {
		return fmt.Errorf("logger ingores given level logs")
	}
	var stats string
	if l.require&RequireAll != 0 {
		stats = "-"
		if l.require&RequireCaller != 0 {
			stats = fmt.Sprintf("%s %s", getCall(3), stats)
		}
		if l.require&RequireTimestamp != 0 {
			stats = fmt.Sprintf("[%s] %s", getTimestamp(), stats)
		}
		if l.require&RequireLevel != 0 {
			if l.colored {
				stats = fmt.Sprintf("[%s] %s", getColoredLevel(level), stats)
			} else {
				stats = fmt.Sprintf("[%s] %s", getRegularLevel(level), stats)
			}
		}
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(
		l.output,
		"%s %s\n",
		stats,
		fmt.Sprintf(format, attrs...),
	)
	return nil
}

/*
Print log into logger's output writer with any selected level.

Takes formated string and output it with additional info.

Additional info depended on require value of a logger: current time, level, caller function, caller file and given message.
*/
func (l *Logger) Log(level int, format string, attrs ...any) error {
	return l.print(level, format, attrs...)
}

/*
Print log into logger's output writer with Debug level

Takes formated string and output it with additional info.

Additional info depended on require value of a logger: current time, level, caller function, caller file and given message
*/
func (l *Logger) Debug(format string, attrs ...any) error {
	return l.print(LevelDebug, format, attrs...)
}

/*
Print log into logger's output writer with Info level.

Takes formated string and output it with additional info.

Additional info depended on require value of a logger: current time, level, caller function, caller file and given message.
*/
func (l *Logger) Info(format string, attrs ...any) error {
	return l.print(LevelInfo, format, attrs...)
}

/*
Print log into logger's output writer with Warn level.

Takes formated string and output it with additional info.

Additional info depended on require value of a logger: current time, level, caller function, caller file and given message.
*/
func (l *Logger) Warn(format string, attrs ...any) error {
	return l.print(LevelWarn, format, attrs...)
}

/*
Print log into logger's output writer with Error level.

Takes formated string and output it with additional info.

Additional info depended on require value of a logger: current time, level, caller function, caller file and given message.
*/
func (l *Logger) Error(format string, attrs ...any) error {
	return l.print(LevelError, format, attrs...)
}

/*
Stops the program and before that print log into logger's output writer with Fatal level.

Takes formated string and output it with additional info.

Additional info depended on require value of a logger: current time, level, caller function, caller file and given message.
*/
func (l *Logger) Fatal(format string, attrs ...any) {
	l.print(LevelFatal, format, attrs...)
	os.Exit(0)
}

// Print empty line in the output stream
func (l *Logger) Line() {
	l.mu.Lock()
	fmt.Fprintln(l.output)
	l.mu.Unlock()
}

//// ========================================================
