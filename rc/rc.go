package rc

import yaml "gopkg.in/yaml.v2"

//go:generate counterfeiter . PivnetRCReadWriter
type PivnetRCReadWriter interface {
	WriteToFile(contents []byte) error
	ReadFromFile() ([]byte, error)
}

type PivnetRC struct {
	Profiles []PivnetProfile `yaml:"profiles"`
}

type RCHandler struct {
	rcReadWriter PivnetRCReadWriter
}

func NewRCHandler(rcReadWriter PivnetRCReadWriter) *RCHandler {
	return &RCHandler{
		rcReadWriter: rcReadWriter,
	}
}

func (h *RCHandler) SaveProfile(
	profileName string,
	apiToken string,
	host string,
) error {
	pivnetRC, err := h.loadPivnetRC()
	if err != nil {
		return err
	}

	if pivnetRC == nil {
		pivnetRC = &PivnetRC{}
	}

	var newInfo *PivnetProfile
	var index int
	for i, p := range pivnetRC.Profiles {
		if p.Name == profileName {
			newInfo = &p
			index = i
			break
		}
	}

	if newInfo == nil {
		newInfo = &PivnetProfile{
			Name:     profileName,
			APIToken: apiToken,
			Host:     host,
		}
		index = len(pivnetRC.Profiles)
		pivnetRC.Profiles = append(pivnetRC.Profiles, PivnetProfile{})
	}

	newInfo.APIToken = apiToken
	newInfo.Host = host

	pivnetRC.Profiles[index] = *newInfo

	yamlBytes, err := yaml.Marshal(pivnetRC)
	if err != nil {
		// untested as we cannot force yaml unmarshal to return an error
		return err
	}

	return h.rcReadWriter.WriteToFile(yamlBytes)
}

// ProfileForName will return (nil,nil) if the file does not exist,
// or if the profile does not exist,
// but will return (nil,err) for other reasons e.g. the file cannot be read.
func (h *RCHandler) ProfileForName(profileName string) (*PivnetProfile, error) {
	pivnetRC, err := h.loadPivnetRC()
	if err != nil {
		return nil, err
	}

	if pivnetRC == nil {
		return nil, nil
	}

	for _, p := range pivnetRC.Profiles {
		if p.Name == profileName {
			return &p, nil
		}
	}

	return nil, nil
}

// RemoveProfileWithName will return error for all errors except if file does not exist
func (h *RCHandler) RemoveProfileWithName(profileName string) error {
	pivnetRC, err := h.loadPivnetRC()
	if err != nil {
		return err
	}

	// Nothing to logout from if we didn't find an RC file
	if pivnetRC == nil {
		return nil
	}

	foundIndex := -1
	for i, p := range pivnetRC.Profiles {
		if p.Name == profileName {
			foundIndex = i
		}
	}

	if foundIndex >= 0 {
		pivnetRC.Profiles = append(pivnetRC.Profiles[:foundIndex], pivnetRC.Profiles[foundIndex+1:]...)
	}

	yamlBytes, err := yaml.Marshal(pivnetRC)
	if err != nil {
		// untested as we cannot force yaml unmarshal to return an error
		return err
	}

	return h.rcReadWriter.WriteToFile(yamlBytes)
}

// loadPivnetRC does not return an error if the file does not exist
// but will return an error for other reasons e.g. the file cannot be read.
func (h *RCHandler) loadPivnetRC() (*PivnetRC, error) {
	b, err := h.rcReadWriter.ReadFromFile()
	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, nil
	}

	var pivnetRC PivnetRC
	err = yaml.Unmarshal(b, &pivnetRC)
	if err != nil {
		return nil, err
	}

	return &pivnetRC, nil
}
