package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/lucheng0127/kube-eip/pkg/utils/errhandle"
)

const (
	MDDir              string = "/var/run/eip_agent"
	MD_STATUS_FINISHED string = "Finished"
	MD_STATUS_FAILED   string = "Failed"
)

type EipMetadata struct {
	Status     string `json:"status"`
	ExternalIP string `json:"exip"`
	InternalIP string `json:"inip"`
	Phase      int    `json:"phase"`
}

func formatMDPath(eip, iip string) string {
	return fmt.Sprintf("%s/eip-%s-%s.metadata", MDDir, eip, iip)
}

func (m *EipMetadata) getMDFilepath() string {
	return formatMDPath(m.ExternalIP, m.InternalIP)
}

func (m *EipMetadata) DumpMD() error {
	if _, err := os.Stat(MDDir); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(MDDir, os.ModePerm); err != nil {
			return err
		}
	}

	content, err := json.MarshalIndent(m, "", "")
	if err != nil {
		return err
	}

	err = os.WriteFile(m.getMDFilepath(), content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (m *EipMetadata) DeleteMD() error {
	err := os.Remove(m.getMDFilepath())
	if err != nil && !errhandle.IsNoSuchFileError(err) {
		return err
	}

	return nil
}

func ParseMD(eip, iip string) (*EipMetadata, error) {
	md := new(EipMetadata)
	filepath := formatMDPath(eip, iip)

	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		// If metadata file not exist, do not return error
		return nil, nil
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, md)
	if err != nil {
		return nil, err
	}

	return md, nil
}
