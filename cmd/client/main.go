package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"sqsclientserver/config"
	aws "sqsclientserver/src/aws"
	"sqsclientserver/src/data"
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

	ctx := context.Background()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sess, err = aws.NewSession(aws.Config{ID: configData.AccessKeyID,
		Secret: configData.SecretKeyID,
		Region: configData.Region})
	if err != nil {
		log.Fatalf("cant create AWS session")
	}

	fmt.Println("sqs url: ", configData.SQSUrl)
	queue := queue.NewQueue(sess, configData.SQSUrl)
	client := NewClient(queue)

	s := promptui.Select{
		Label: "Choose client option",
		Items: []string{_addItem, _deleteItem, _getItem, _getAllItems},
	}

	var choice string
	_, choice, err = s.Run()
	if err != nil {
		fmt.Println("failed to run promptui %s", err.Error())
		return
	}

	for {
		switch choice {
		case _addItem:
			fmt.Println("adding a new item on server")

			prompt := promptui.Prompt{
				Label:    "Crate Key (must be string)",
				Validate: validate,
			}

			key, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			prompt = promptui.Prompt{
				Label: "Value",
			}

			value, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			d := data.NewData(key, value)
			err = client.AddItem(ctx, d)
			if err != nil {
				fmt.Printf("client add item %v\n", err)
			}

			continue

		case _deleteItem:
			fmt.Println("deleting item from server")
			prompt := promptui.Prompt{
				Label:    "Key (must be string)",
				Validate: validate,
			}

			key, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			err = client.DeleteItem(ctx, key)
			if err != nil {
				fmt.Printf("client delete item %v\n", err)
			}

			continue

		case _getItem:
			fmt.Println("getting single item from server")
			prompt := promptui.Prompt{
				Label:    "Key (must be string)",
				Validate: validate,
			}

			key, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			err = client.GetItem(ctx, key)
			if err != nil {
				fmt.Printf("client get item %v\n", err)
			}

			continue

		case _getAllItems:
			fmt.Println("get all items")
			err := client.GetAllItems(ctx)
			if err != nil {
				fmt.Printf("client get all items %v\n", err)
			}
		}

		continue
	}
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

	return ret, nil
}

func validate(input string) error {
	if reflect.ValueOf(input).Kind() != reflect.String {
		return errors.New("entered value is not a string")
	}

	return nil
}
