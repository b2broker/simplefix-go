package generator

import (
	"bytes"
	"io"
	"text/template"
)

type Formatter interface {
	Format() string
}

func Execute(writer io.Writer, formatter Formatter) error {
	t := template.Must(template.New("").Parse(formatter.Format()))
	return t.Execute(writer, formatter)
}

func ExecuteString(formatter Formatter) (string, error) {
	buffer := &bytes.Buffer{}
	t := template.Must(template.New("").Parse(formatter.Format()))
	err := t.Execute(buffer, formatter)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
