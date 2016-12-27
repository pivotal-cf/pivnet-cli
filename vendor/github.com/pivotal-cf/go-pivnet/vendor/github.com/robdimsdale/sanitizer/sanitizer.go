package sanitizer

import (
	"io"
	"strings"
)

type Sanitizer interface {
	io.Writer
}

type sanitizer struct {
	sanitized map[string]string
	sink      io.Writer
}

func NewSanitizer(sanitized map[string]string, sink io.Writer) Sanitizer {
	if _, ok := sanitized[""]; ok {
		delete(sanitized, "")
	}
	return &sanitizer{
		sanitized: sanitized,
		sink:      sink,
	}
}

func (s sanitizer) Write(p []byte) (n int, err error) {
	input := string(p)

	for k, v := range s.sanitized {
		input = strings.Replace(input, k, v, -1)
	}

	scrubbed := []byte(input)

	return s.sink.Write(scrubbed)
}
