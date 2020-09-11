package hostwarning

import "fmt"

type HostWarning struct {
	host string
}

func NewHostWarning(host string) *HostWarning {
	return &HostWarning{
		host: host,
	}
}

func (hw HostWarning) Warn() string {
	if hw.host != "https://network.pivotal.io" && hw.host != ""  {
		return fmt.Sprintf( "\nWarning: You are currently targeting %s\n", hw.host)
	}
	 return ""
}
