package main

import (
	"fmt"
	"github.com/callebjorkell/huectl/cmd"
	log "github.com/sirupsen/logrus"
)

type colorFormatter struct {
	log.TextFormatter
}

func (f *colorFormatter) Format(entry *log.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case log.DebugLevel, log.TraceLevel:
		levelColor = 90 // dark grey
	case log.WarnLevel:
		levelColor = 33 // yellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		levelColor = 91 // bright red
	default:
		levelColor = 39 // default
	}
	return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m\n", levelColor, entry.Message)), nil
}

func main() {
	log.SetFormatter(&colorFormatter{})
	log.SetLevel(log.DebugLevel)

	if err := cmd.Huectl().Execute(); err != nil {
		log.Fatal(err)
	}
}
