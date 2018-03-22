package google

import (
	client "k8s.io/kube-deploy/ext-apiserver/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	"reflect"
	"testing"
	"k8s.io/kube-deploy/ext-apiserver/cloud/google/installation"
	"fmt"
)

func TestNewMachineActuator(t *testing.T) {
	testTables := []struct {
		token               string
		machineClient       client.MachineInterface
		path                string
		expectedConfigWatch *installation.ConfigWatch
		expectedErr         bool
	}{
		{
			token:               "token",
			machineClient:       nil,
			path:                "",
			expectedConfigWatch: &installation.ConfigWatch{},
			expectedErr:         false,
		},
	}

	for _, table := range testTables {
		gceClient, err := NewMachineActuator(table.token, table.machineClient, table.path)
		//config, err := gceClient.configWatch.Config()
		//fmt.Printf("Config: %v\n", config)
		fmt.Printf("Config == nil? %v\n", gceClient.configWatch == nil)
		if gceClient.kubeadmToken != table.token {
			t.Errorf("Kubeadm token was incorrect, got: %s, want %s.", gceClient.kubeadmToken, table.token)
		}
		if reflect.DeepEqual(gceClient.machineClient, table.machineClient) {
			t.Errorf("Machine client was incorrect, got: %v\nwant: %v.", gceClient.machineClient, table.machineClient)
		}
		if reflect.DeepEqual(gceClient.configWatch, table.expectedConfigWatch) {
			t.Errorf("Config watch was incorrect, got: %v\nwant: %v.", gceClient.configWatch, table.expectedConfigWatch)
		}
		if table.expectedErr && err == nil {
			t.Errorf("An error was not received as expected.")
		}
		if !table.expectedErr && err != nil {
			t.Errorf("Got unexpected error: %s", err)
		}
	}
}

//func TestGetImage(t *testing.T) {
//	defaultImg := "projects/ubuntu-os-cloud/global/images/family/ubuntu-1710"
//
//	testTables := []struct {
//		img string
//		project string
//		expectedPath string
//		expectedPreloaded bool
//	}{
//		{
//			img: "ubuntu-1710",
//			project: "my-project",
//			expectedPath: "projects/my-project/global/images/ubuntu-1710",
//			expectedPreloaded: false,
//		},
//	}
//
//	for _, table := range testTables {
//		path, isPreloaded :=
//	}
//}
