package manager

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

func (m *EipMetadata) getMDFilepath() string {
	return fmt.Sprintf("%s/eip-%s.metadata", MDDir, m.ExternalIP)
}

func (m *EipMetadata) dumpMD() error {
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

func (m *EipMetadata) deleteMD() error {
	err := os.Remove(m.getMDFilepath())
	if err != nil && !errhandle.IsNoSuchFileError(err) {
		return err
	}

	return nil
}

func parseMD(eip string) (*EipMetadata, error) {
	md := new(EipMetadata)
	filepath := fmt.Sprintf("%s/eip-%s.metadata", MDDir, eip)

	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
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
