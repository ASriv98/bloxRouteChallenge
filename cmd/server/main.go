package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"sqsclientserver/config"
	aws "sqsclientserver/src/aws"
	"sqsclientserver/src/logging"
	"sqsclientserver/src/queue"
)

func main() {
	var (
		configData config.Data
	)

	configData, err := readConfig()
	if err != nil {
		panic(fmt.Sprintf("read config %s", err))
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sess, err = aws.NewSession(aws.Config{ID: configData.AccessKeyID,
		Secret: configData.SecretKeyID,
		Region: configData.Region})
	if err != nil {
		log.Fatalf("cannot create aws session")
	}

	queue := queue.NewQueue(sess, configData.SQSUrl)
	stdOutLog := logging.NewLogger(configData.ServerLogPath)

	server := NewServer(queue, stdOutLog, configData)
	server.Start(context.Background())
}

func readConfig() (config.Data, error) {
	var (
		data []byte
		ret  config.Data
		err  error
	)

	if data, err = ioutil.ReadFile(config.DefaultPath); err != nil {
		return ret, errors.Wrap(err, "reading config file")
	}

	if err = yaml.Unmarshal(data, &ret); err != nil {
		return ret, errors.Wrap(err, "parsing config file")
	}

	if err = ret.IsConfigFileValid(); err != nil {
		return ret, errors.Wrap(err, "is config file valid")
	}

	return ret, nil
}
