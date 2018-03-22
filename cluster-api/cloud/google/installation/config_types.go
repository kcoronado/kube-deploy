package installation

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io"
	"io/ioutil"
	clustercommon "k8s.io/kube-deploy/cluster-api/pkg/apis/cluster/common"
	clusterv1 "k8s.io/kube-deploy/cluster-api/pkg/apis/cluster/v1alpha1"
	"k8s.io/kube-deploy/cluster-api/util"
	"os"
)

type ConfigWatch struct {
	Path string
}

type Config struct {
	infoList *infoList
}

type infoList struct {
	Items []info `json:"items"`
}

type info struct {
	OS       string                         `json:"os"`
	Roles    []clustercommon.MachineRole    `json:"roles"`
	Versions []clusterv1.MachineVersionInfo `json:"versions"`

	// This can either be a full projects path to an image/family,
	// or just the image name which is in the project.
	// If it's an image in the project, this field may be
	// identical to the OS field.
	Image    string   `json:"image"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	StartupScript string `json:"startupScript"`
}

type ConfigParams struct {
	OS       string
	Roles    []clustercommon.MachineRole
	Versions clusterv1.MachineVersionInfo
}

func NewConfigWatch(path string) (*ConfigWatch, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	return &ConfigWatch{Path: path}, nil
}

func (cw *ConfigWatch) Config() (*Config, error) {
	file, err := os.Open(cw.Path)
	if err != nil {
		return nil, err
	}
	return parseInstallationYaml(file)
}

// TODO(kcoronado): return an array of pointers to config
func parseInstallationYaml(reader io.Reader) (*Config, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	infoList := &infoList{}
	err = yaml.Unmarshal(bytes, infoList)
	if err != nil {
		return nil, err
	}

	return &Config{infoList}, nil
}

func (c *Config) GetImage(params *ConfigParams) (string, error) {
	installationConfig, err := c.matchInstallationConfig(params)
	if err != nil {
		return "", err
	}
	return installationConfig.Image, nil
}

func (c *Config) GetMetadata(params *ConfigParams) (*Metadata, error) {
	installationConfig, err := c.matchInstallationConfig(params)
	if err != nil {
		return nil, err
	}
	return &installationConfig.Metadata, nil
}

func (c *Config) matchInstallationConfig(params *ConfigParams) (*info, error) {
	for _, info := range c.infoList.Items {
		if params.OS != info.OS {
			continue
		}
		foundRoles := true
		for _, role := range params.Roles {
			if !util.RoleContains(role, info.Roles) {
				foundRoles = false
				break
			}
		}
		if !foundRoles {
			continue
		}
		foundVersion := false
		for _, versionSet := range info.Versions {
			if params.Versions == versionSet {
				foundVersion = true
			}
		}
		if foundVersion {
			return &info, nil
		}
	}
	return nil, fmt.Errorf("could not find a matching installation config for params %+v", params)
}
