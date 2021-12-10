package analyze

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type textFileWriter struct {
	path         string
	file         *os.File
	stringBuffer *bufio.Writer
}

func (r *textFileWriter) openFileForWrite(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Errorln(err)
		return err
	}

	r.file = file
	r.stringBuffer = bufio.NewWriter(file)

	return nil
}

func (r *textFileWriter) close() error {
	errWriter := r.stringBuffer.Flush()
	if errWriter != nil {
		return errWriter
	}

	errFile := r.file.Close()
	if errFile != nil {
		return errFile
	}

	return nil
}

func (r *textFileWriter) writeText(text string) {
	_, _ = r.stringBuffer.WriteString(text + "\n")
}

func (r *textFileWriter) writeTextWithBottomPadding(text string, padding int) {
	r.writeText(text + strings.Repeat("\n", padding))
}
