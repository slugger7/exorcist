package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/slugger7/exorcist/internal/environment"
)

type ILogger interface {
	Debug(message string)
	Debugf(format string, args ...any)
	Info(message string)
	Infof(format string, args ...any)
	Warning(message string)
	Warningf(format string, args ...any)
	Error(message string)
	Errorf(format string, args ...any)
}

type logger struct {
	env           *environment.EnvironmentVariables
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

var loggerInstance *logger

func New(env *environment.EnvironmentVariables) ILogger {
	if loggerInstance == nil {
		loggerInstance = &logger{
			env:           env,
			debugLogger:   log.New(os.Stdout, "[DEBUG]", log.Default().Flags()),
			infoLogger:    log.New(os.Stdout, "[INFO] ", log.Default().Flags()),
			warningLogger: log.New(os.Stdout, "[WARN]", log.Default().Flags()),
			errorLogger:   log.New(os.Stdout, "[ERROR]", log.Default().Flags()),
		}
	}
	return loggerInstance
}

type callerInfo struct {
	file     string
	funcName string
	lineNo   int
}

func getCallerInformation(skip int) callerInfo {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		log.Println("runtime.Caller() failed")
	}
	funcName := runtime.FuncForPC(pc).Name()

	return callerInfo{file: file, funcName: funcName, lineNo: lineNo}
}

func (l *logger) logDebug(message string) {
	ci := getCallerInformation(3)
	l.debugLogger.Printf("%v@%v(%v): %v", ci.file, ci.lineNo, ci.funcName, message)
}

func (l *logger) Debug(message string) {
	l.logDebug(message)
}
func (l *logger) Debugf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.logDebug(message)
}

func (l *logger) Info(message string) {
	l.infoLogger.Println(message)
}

func (l *logger) Infof(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.Info(message)
}

func (l *logger) logWarning(message string) {
	ci := getCallerInformation(3)
	l.warningLogger.Printf("%v: %v", ci.funcName, message)
}

func (l *logger) Warning(message string) {
	l.logWarning(message)
}

func (l *logger) Warningf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.logWarning(message)
}

func (l *logger) logErorr(message string) {
	ci := getCallerInformation(3)
	l.errorLogger.Printf("%v@%v(%v): %v", ci.file, ci.lineNo, ci.funcName, message)
}

func (l *logger) Error(message string) {
	l.logErorr(message)
}

func (l *logger) Errorf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.logErorr(message)
}
