package tool

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Oc struct {
	ProcessExecutorer
}

func NewOc(exec ProcessExecutorer) Oc {
	return Oc{
		ProcessExecutorer: exec,
	}
}

const osVersionKey = "serverVersion"

// Based on https://access.redhat.com/solutions/4870701
var kubeOpenShiftVersionMap map[string]string = map[string]string{
	"1.22": "4.9",
	"1.21": "4.8",
	"1.20": "4.7",
	"1.19": "4.6",
	"1.18": "4.5",
	"1.17": "4.4",
	"1.16": "4.3",
	"1.14": "4.2",
	"1.13": "4.1",
}

func (o Oc) GetVersion() (string, error) {
	rawOutput, err := o.RunProcessAndCaptureOutput("oc", "version", "-o", "yaml")
	if err != nil {
		return "", err
	}
	out := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(rawOutput), &out)
	if err != nil {
		return "", err
	}
	// Relying on Kubernetes version can be replaced after fixing this issue:
	// https://bugzilla.redhat.com/show_bug.cgi?id=1850656
	kubeServerVersion := out[osVersionKey].(map[string]interface{})
	kubeVersion := fmt.Sprintf("%s.%s", kubeServerVersion["major"], kubeServerVersion["minor"])
	osVersion, ok := kubeOpenShiftVersionMap[kubeVersion]
	if !ok {
		return "", fmt.Errorf("Internal error: %q not found in Kubernetes-OpenShift version map", kubeVersion)
	}

	return osVersion, nil
}

func GetKubeOpenShiftVersionMap() map[string]string {
	return kubeOpenShiftVersionMap
}
