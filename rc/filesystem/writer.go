package filesystem

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

const (
	fileModeUserReadWrite = 0600
)

type PivnetRCWriter struct {
}

func NewPivnetRCWriter() *PivnetRCWriter {
	return &PivnetRCWriter{}
}

func (h *PivnetRCWriter) WriteToFile(configFileLocation string, contents interface{}) error {
	yamlBytes, err := yaml.Marshal(contents)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFileLocation, yamlBytes, fileModeUserReadWrite)
	if err != nil {
		return err
	}

	err = os.Chmod(configFileLocation, fileModeUserReadWrite)
	if err != nil {
		return err
	}

	return nil
}
