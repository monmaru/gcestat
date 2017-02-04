package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

type config struct {
	GCP gcpConfig
}

type gcpConfig struct {
	ProjectID string
	Zone      string
}

func getInstanceList(client *http.Client, config gcpConfig) (*compute.InstanceList, error) {
	service, err := compute.New(client)
	if err != nil {
		return nil, err
	}

	s := compute.NewInstancesService(service)
	list, err := s.List(config.ProjectID, config.Zone).Do()
	if err != nil {
		return nil, err
	}
	return list, err
}

func showInstanceInfo(client *http.Client, config gcpConfig) {
	list, err := getInstanceList(client, config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, item := range list.Items {
		fmt.Println("Name: ", item.Name)
		fmt.Println("Id: ", item.Id)
		fmt.Println("Zone: ", item.Zone)
		fmt.Println("Tags: ", item.Tags)
		fmt.Println("Kind: ", item.Kind)
		fmt.Println("MachineType: ", item.MachineType)
		fmt.Println("StatusMessage: ", item.StatusMessage)
		fmt.Println("Description: ", item.Description)
		fmt.Println("-------------------------------")
	}
}

func main() {
	var config config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	ctx := context.Background()
	client, err := google.DefaultClient(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	showInstanceInfo(client, config.GCP)
}
