package dslog

import (
    "errors"
    "fmt"
    "io"
    "os"
    "time"

    "github.com/sirupsen/logrus"
)

type logFormatter struct {
}

func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    var err error
    timeStamp := time.Now().Format(time.RFC3339)
    levelString := entry.Level.String()
    message := fmt.Sprintf("%s %s %s\n", timeStamp, levelString, entry.Message)
    return []byte(message), err
}

const DebugLevel int = 1
const ErrorLevel int = 2

func init() {
    logrus.SetOutput(os.Stdout)
    logrus.SetLevel(logrus.DebugLevel)
    logrus.SetFormatter(new(logFormatter))
}

func SetDebugMode(debugMode bool) {
    if debugMode {
        logrus.SetLevel(logrus.DebugLevel)
        return
    }
    logrus.SetLevel(logrus.ErrorLevel)
}

func SetLevel(level int) error {
    var err error
    switch level {
        case DebugLevel:
            logrus.SetLevel(logrus.DebugLevel)
        case ErrorLevel:
            logrus.SetLevel(logrus.ErrorLevel)
        default:
            return errors.New("wrong log level")
    }
    return err
}

func SetOutput(writer io.Writer) error {
    var err error
    logrus.SetOutput(writer)
    return err
}

func LogDebug(message ...interface{}) {
    logrus.Debug(message)
}

func LogError(message ...interface{}) {
    logrus.Error(message)
}

func LogWarning(message ...interface{}) {
    logrus.Warning(message)
}

func LogInfo(message ...interface{}) {
    logrus.Info(message)
}


func LogDebugf(format string, args ...interface{}) {
    format = "[" + format + "]"
    message := fmt.Sprintf(format, args...)
    logrus.Debug(message)
}

func LogErrorf(format string, args ...interface{}) {
    format = "[" + format + "]"
    message := fmt.Sprintf(format, args...)
    logrus.Error(message)
}

func LogWarningf(format string, args ...interface{}) {
    format = "[" + format + "]"
    message := fmt.Sprintf(format, args...)
    logrus.Warning(message)
}

func LogInfof(format string, args ...interface{}) {
    format = "[" + format + "]"
    message := fmt.Sprintf(format, args...)
    logrus.Info(message)
}
