// Package config handles ingestion of run time configuration needed on a tester by tester basis
package config

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
