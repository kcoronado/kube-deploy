package installation

import (
	"io"
	"k8s.io/kube-deploy/ext-apiserver/pkg/apis/cluster/common"
	"k8s.io/kube-deploy/ext-apiserver/pkg/apis/cluster/v1alpha1"
	"reflect"
	"strings"
	"testing"
)

func TestParseInstallationYaml(t *testing.T) {
	testTables := []struct {
		reader      io.Reader
		expectedErr bool
	}{
		{
			reader: strings.NewReader(`items:
- os: ubuntu-1710
  roles:
  - Master
  versions:
  - kubelet: 1.9.3
    controlPlane: 1.9.3
    containerRuntime:
      name: docker
      version: 1.12.0
  - kubelet: 1.9.4
    controlPlane: 1.9.4
    containerRuntime:
      name: docker
      version: 1.12.0
  image: projects/ubuntu-os-cloud/global/images/family/ubuntu-1710
  metadata:
    startupScript: |
      #!/bin/bash
- os: ubuntu-1710
  roles:
  - Node
  versions:
  - kubelet: 1.9.3
    containerRuntime:
      name: docker
      version: 1.12.0
  - kubelet: 1.9.4
    containerRuntime:
      name: docker
      version: 1.12.0
  image: projects/ubuntu-os-cloud/global/images/family/ubuntu-1710
  metadata:
    startupScript: |
      #!/bin/bash
      echo this is the node config.`),
			expectedErr: false,
		},
		{
			reader:      strings.NewReader("Not valid yaml"),
			expectedErr: true,
		},
	}

	for _, table := range testTables {
		config, err := parseInstallationYaml(table.reader)
		if table.expectedErr {
			if err == nil {
				t.Errorf("An error was not received as expected.")
			}
			if config != nil {
				t.Errorf("Config should be nil, got %v", config)
			}
		}
		if !table.expectedErr {
			if err != nil {
				t.Errorf("Got unexpected error: %s", err)
			}
			if config == nil {
				t.Errorf("Config should have been parsed, but was nil")
			}
		}
	}
}

func TestMatchInstallationConfig(t *testing.T) {
	masterInstallationInfo := info{
		OS:    "ubuntu-1710",
		Roles: []common.MachineRole{common.MasterRole},
		Versions: []v1alpha1.MachineVersionInfo{
			{
				Kubelet:      "1.9.3",
				ControlPlane: "1.9.3",
				ContainerRuntime: v1alpha1.ContainerRuntimeInfo{
					Name:    "docker",
					Version: "1.12.0",
				},
			}, {
				Kubelet:      "1.9.4",
				ControlPlane: "1.9.4",
				ContainerRuntime: v1alpha1.ContainerRuntimeInfo{
					Name:    "docker",
					Version: "1.12.0",
				},
			},
		},
		Image: "projects/ubuntu-os-cloud/global/images/family/ubuntu-1710",
		Metadata: Metadata{
			StartupScript: "Master startup script",
		},
	}
	nodeInstallationInfo := info{
		OS:    "ubuntu-1710",
		Roles: []common.MachineRole{common.NodeRole},
		Versions: []v1alpha1.MachineVersionInfo{
			{
				Kubelet: "1.9.3",
				ContainerRuntime: v1alpha1.ContainerRuntimeInfo{
					Name:    "docker",
					Version: "1.12.0",
				},
			}, {
				Kubelet: "1.9.4",
				ContainerRuntime: v1alpha1.ContainerRuntimeInfo{
					Name:    "docker",
					Version: "1.12.0",
				},
			},
		},
		Image: "projects/ubuntu-os-cloud/global/images/family/ubuntu-1710",
		Metadata: Metadata{
			StartupScript: "Node startup script",
		},
	}

	config := Config{
		infoList: &infoList{
			Items: []info{masterInstallationInfo, nodeInstallationInfo},
		},
	}

	testTables := []struct {
		params        ConfigParams
		expectedMatch *info
		expectedErr   bool
	}{
		{
			params: ConfigParams{
				OS:    "ubuntu-1710",
				Roles: []common.MachineRole{common.MasterRole},
				Versions: v1alpha1.MachineVersionInfo{
					Kubelet:      "1.9.4",
					ControlPlane: "1.9.4",
					ContainerRuntime: v1alpha1.ContainerRuntimeInfo{
						Name:    "docker",
						Version: "1.12.0",
					},
				},
			},
			expectedMatch: &masterInstallationInfo,
			expectedErr:   false,
		},
		{
			params: ConfigParams{
				OS:    "ubuntu-1710",
				Roles: []common.MachineRole{common.NodeRole},
				Versions: v1alpha1.MachineVersionInfo{
					Kubelet: "1.9.4",
					ContainerRuntime: v1alpha1.ContainerRuntimeInfo{
						Name:    "docker",
						Version: "1.12.0",
					},
				},
			},
			expectedMatch: &nodeInstallationInfo,
			expectedErr:   false,
		},
		{
			params: ConfigParams{
				OS:    "ubuntu-1710",
				Roles: []common.MachineRole{common.NodeRole},
				Versions: v1alpha1.MachineVersionInfo{
					Kubelet:      "1.9.4",
					ControlPlane: "1.9.4",
					ContainerRuntime: v1alpha1.ContainerRuntimeInfo{
						Name:    "docker",
						Version: "1.13.0",
					},
				},
			},
			expectedMatch: nil,
			expectedErr:   true,
		},
	}

	for _, table := range testTables {
		matched, err := config.matchInstallationConfig(&table.params)
		if !reflect.DeepEqual(matched, table.expectedMatch) {
			t.Errorf("Matched installation info was incorrect, got: %+v,\n want %+v.", matched, table.expectedMatch)
		}
		if table.expectedErr && err == nil {
			t.Errorf("An error was not received as expected.")
		}
		if !table.expectedErr && err != nil {
			t.Errorf("Got unexpected error: %s", err)
		}
	}
}
