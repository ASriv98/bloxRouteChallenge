// Package config handles ingestion of run time configuration needed on a tester by tester basis
package config

import "github.com/pkg/errors"

type Data struct {
	file []byte

	ServerLogPath      string `yaml:"serverLogPath"`
	DefaultLogPath     string `yaml:"defaultLogPath"`
	SQSUrl             string `yaml:"sqsURL"`
	AccessKeyID        string `yaml:"awsAccessKeyID"`
	SecretKeyID        string `yaml:"awsSecretKeyID"`
	Region             string `yaml:"awsRegion"`
	ServerIdleInterval string `yaml:"serverIdleInterval"`
	NumServerWorkers   string `yaml:"numServerWorkers"`
}

// IsConfigFileValid checks to see whether the config file has all fields logically usable in the test
func (d *Data) IsConfigFileValid() error {

	// TODO make all these functions return an error
	if d.AccessKeyID == "" || d.SecretKeyID == "" || d.Region == "" {
		return errors.New("missing aws config parameters")
	}

	return nil
}
