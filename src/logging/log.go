package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

var Out *os.File = os.Stdout

type Logger struct {
	file *os.File
}

func NewLogger(f string) *Logger {
	file, err := os.OpenFile(f,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666)
	if err != nil {
		log.Fatalf("Can't attach log file %s", err)
	}

	l := &Logger{
		file: file,
	}

	return l
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	mw := io.MultiWriter(os.Stdout, l.file)
	log.SetOutput(mw)

	_, err := l.file.Write(l.formatArgs(format, args...))
	if err != nil {
		return
	}

	_, err = Out.Write(l.formatArgs(format, args...))
	if err != nil {
		return
	}

	os.Exit(1)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	mw := io.MultiWriter(os.Stdout, l.file)
	log.SetOutput(mw)

	_, err := l.file.Write(l.formatArgs(format, args...))
	if err != nil {
		return
	}

	_, err = Out.Write(l.formatArgs(format, args...))
	if err != nil {
		return
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	mw := io.MultiWriter(os.Stdout, l.file)
	log.SetOutput(mw)

	_, err := l.file.Write(l.formatArgs(format, args...))
	if err != nil {
		return
	}

	_, err = Out.Write(l.formatArgs(format, args...))
	if err != nil {
		return
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	mw := io.MultiWriter(os.Stdout, l.file)
	log.SetOutput(mw)

	_, err := l.file.Write(l.formatArgs(format, args...))
	if err != nil {
		return
	}

	_, err = Out.Write(l.formatArgs(format, args...))
	if err != nil {
		return
	}
}

func (l *Logger) formatArgs(format string, args ...interface{}) []byte {
	str := "\n" + fmt.Sprintf(format, args...)
	return []byte(str)
}
