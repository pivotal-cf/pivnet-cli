package filesystem

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type PivnetRCReadWriter struct {
	configFilepath string
}

func NewPivnetRCReadWriter(configFilepath string) *PivnetRCReadWriter {
	return &PivnetRCReadWriter{
		configFilepath: configFilepath,
	}
}

// loadPivnetRC does not return an error if the file does not exist
// but will return an error for other reasons e.g. the file cannot be read.
func (h *PivnetRCReadWriter) ReadFromFile() ([]byte, error) {
	_, err := os.Stat(h.configFilepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	pivnetRCBytes, err := ioutil.ReadFile(h.configFilepath)
	if err != nil {
		return nil, err
	}

	return pivnetRCBytes, nil
}

const (
	fileModeUserReadWrite = 0600
)

type PivnetRCWriter struct {
	configFilepath string
}

func (h *PivnetRCReadWriter) WriteToFile(contents interface{}) error {
	yamlBytes, err := yaml.Marshal(contents)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(h.configFilepath, yamlBytes, fileModeUserReadWrite)
	if err != nil {
		return err
	}

	err = os.Chmod(h.configFilepath, fileModeUserReadWrite)
	if err != nil {
		return err
	}

	return nil
}
