package logging

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

type TextFormatter struct{}

func (f *TextFormatter) Format(entry *log.Entry) ([]byte, error) {

	var buf *bytes.Buffer
	if entry.Buffer == nil {
		buf = &bytes.Buffer{}
	} else {
		buf = entry.Buffer
	}

	colorLevel := map[log.Level]color.Attribute{
		log.DebugLevel: color.FgWhite,
		log.InfoLevel:  color.FgBlue,
		log.WarnLevel:  color.FgYellow,
		log.ErrorLevel: color.FgRed,
		log.FatalLevel: color.FgRed,
		log.PanicLevel: color.FgRed,
	}

	app := "MAIN"
	if val, ok := entry.Data["app"]; ok {
		app = strings.ToUpper(val.(string))
		if len(app) > 4 {
			app = app[0:4]
		}
	}

	fmt.Fprintf(buf, "%s %s %s %s",
		entry.Time.Format("2006-01-02 15:04:05"),
		color.New(color.Bold).Sprint(app),
		color.New(colorLevel[entry.Level]).Sprint(strings.ToUpper(entry.Level.String()[0:4])),
		entry.Message,
	)

	for field, value := range entry.Data {
		if field != "app" {
			fmt.Fprintf(buf, " %s=%s", color.CyanString(field), value)
		}
	}

	fmt.Fprintln(buf)
	return buf.Bytes(), nil
}
