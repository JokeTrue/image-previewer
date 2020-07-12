package logging

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

var (
	DefaultLogger = newServiceLogger()
	NilLogger     = newServiceLogger()
)

func init() {
	DefaultLogger.SetFormatter(FormatterText)
	NilLogger.SetOutput(ioutil.Discard)
}

func SetLevel(level Level) {
	DefaultLogger.SetLevel(level)
}

func GetLevel() Level {
	return DefaultLogger.GetLevel()
}

func SetOutput(writer io.Writer) {
	DefaultLogger.SetOutput(writer)
}

func SetDefaultFields(fields Fields) {
	DefaultLogger.entry.Data = logrus.Fields(fields)
}

func SetFormatter(ftype Formatter) {
	DefaultLogger.SetFormatter(ftype)
}

func ParseLevel(lvl string) (Level, error) {
	level, err := logrus.ParseLevel(lvl)
	return Level(level), err
}

func WithField(key string, value interface{}) *ServiceLogger {
	return DefaultLogger.WithField(key, value)
}

func WithFields(fields Fields) *ServiceLogger {
	return DefaultLogger.WithFields(fields)
}

func WithError(err error) *ServiceLogger {
	return DefaultLogger.WithError(err)
}

func WithContext(ctx context.Context) *ServiceLogger {
	return DefaultLogger.WithContext(ctx)
}

func Debug(args ...interface{}) { DefaultLogger.Debug(args...) }
func Info(args ...interface{})  { DefaultLogger.Info(args...) }
func Print(args ...interface{}) { DefaultLogger.Print(args...) }
func Warn(args ...interface{})  { DefaultLogger.Warn(args...) }
func Error(args ...interface{}) { DefaultLogger.Error(args...) }
func Fatal(args ...interface{}) { DefaultLogger.Fatal(args...) }
func Panic(args ...interface{}) { DefaultLogger.Panic(args...) }

func Debugf(format string, args ...interface{}) { DefaultLogger.Debugf(format, args...) }
func Infof(format string, args ...interface{})  { DefaultLogger.Infof(format, args...) }
func Printf(format string, args ...interface{}) { DefaultLogger.Printf(format, args...) }
func Warnf(format string, args ...interface{})  { DefaultLogger.Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { DefaultLogger.Errorf(format, args...) }
func Fatalf(format string, args ...interface{}) { DefaultLogger.Fatalf(format, args...) }
func Panicf(format string, args ...interface{}) { DefaultLogger.Panicf(format, args...) }

func Debugln(args ...interface{}) { DefaultLogger.Debugln(args...) }
func Infoln(args ...interface{})  { DefaultLogger.Infoln(args...) }
func Println(args ...interface{}) { DefaultLogger.Println(args...) }
func Warnln(args ...interface{})  { DefaultLogger.Warnln(args...) }
func Errorln(args ...interface{}) { DefaultLogger.Errorln(args...) }
func Fatalln(args ...interface{}) { DefaultLogger.Fatalln(args...) }
func Panicln(args ...interface{}) { DefaultLogger.Panicln(args...) }
