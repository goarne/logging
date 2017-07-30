package logging

import (
	"io"
	"log"
	"os"
)

//Available loggers
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

//Datastructure for storing configurationvalues used to handle the file rotation.
type LogConfig struct {
	Filename         string
	Size             int64
	MaxNumberOfFiles int
}

func SetupFileAndConsoleLogging(lc LogConfig) {
	rotatingWriter := CreateRotatingWriter(lc)
	logger := CreateLogWriter(rotatingWriter)
	logger.Append(os.Stdout)

	InitLoggers(logger, logger, logger, logger)
}

func InitLoggers(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

}

func CreateLogWriter(w io.Writer) *LogWriter {
	return &LogWriter{w}
}

func CreateRotatingWriter(l LogConfig) io.Writer {
	rw := &RotatingFileWriter{FileName: l.Filename, Size: l.Size, MaxNumberOfFiles: l.MaxNumberOfFiles, File: &os.File{}}
	rw.OpenFile()
	return rw
}

//The structure support multiple writers.
type LogWriter struct {
	io.Writer
}

func (l *LogWriter) Append(w io.Writer) *LogWriter {
	if l.Writer != nil {
		l.Writer = io.MultiWriter(l.Writer, w)
	} else {
		l.Writer = io.MultiWriter(w)
	}

	return l
}
