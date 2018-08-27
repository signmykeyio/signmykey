package logging

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// TextFormatter represents a formatter logging type
type TextFormatter struct{}

// Format log entry and return it as a slice of bytes
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

	_, err := fmt.Fprintf(buf, "%s %s %s %s",
		entry.Time.Format("2006-01-02 15:04:05"),
		color.New(color.Bold).Sprint(app),
		color.New(colorLevel[entry.Level]).Sprint(strings.ToUpper(entry.Level.String()[0:4])),
		entry.Message,
	)
	if err != nil {
		return []byte(""), errors.Wrap(err, "error formating log entry")
	}

	for field, value := range entry.Data {
		if field != "app" {
			fmt.Fprintf(buf, " %v=%v", color.CyanString(field), value)
		}
	}

	_, err = fmt.Fprintln(buf)
	if err != nil {
		return []byte(""), err
	}

	return buf.Bytes(), nil
}
