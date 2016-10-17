package printer

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

const (
	PrintAsTable = "table"
	PrintAsJSON  = "json"
	PrintAsYAML  = "yaml"
)

//go:generate counterfeiter . Printer

type Printer interface {
	PrintYAML(interface{}) error
	PrintJSON(interface{}) error
	Println(message string) error
}

type printer struct {
	outputWriter io.Writer
}

func NewPrinter(outputWriter io.Writer) Printer {
	return &printer{
		outputWriter: outputWriter,
	}
}

func (p printer) PrintYAML(object interface{}) (err error) {
	// We have to do the recovery ourselves here because go-yaml panics
	// when it fails to marshal, unlike JSON which returns an error.
	// This logic is heavily inspired by the equivalent in the
	// json package in the standard library.
	defer func() {
		if r := recover(); r != nil {
			if rString, ok := r.(string); ok {
				err = fmt.Errorf(rString)
				return
			}
			panic(r)
		}
	}()

	b, err := yaml.Marshal(object)
	if err != nil {
		return err
	}

	output := fmt.Sprintf("---\n%s\n", string(b))
	_, err = p.outputWriter.Write([]byte(output))
	return err
}

func (p printer) PrintJSON(object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}

	_, err = p.outputWriter.Write(b)
	return err
}

func (p printer) Println(message string) error {
	_, err := p.outputWriter.Write([]byte(fmt.Sprintln(message)))
	return err
}
