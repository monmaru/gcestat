package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	crm "google.golang.org/api/cloudresourcemanager/v1"
	compute "google.golang.org/api/compute/v1"
)

func main() {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	exitIfError(err)

	resourceManager, err := crm.New(client)
	exitIfError(err)

	projectsResponse, err := resourceManager.Projects.List().Do()
	exitIfError(err)

	compute, err := compute.New(client)
	exitIfError(err)

	for _, project := range projectsResponse.Projects {
		fmt.Printf("[%s]\n", project.ProjectId)
		zonesResponse, err := compute.Zones.List(project.ProjectId).Do()
		if err != nil {
			fmt.Println("Could not access this project.")
			continue
		}

		for _, zone := range zonesResponse.Items {
			instancesResponse, err := compute.Instances.List(project.ProjectId, zone.Name).Do()
			if err != nil {
				fmt.Println("Could not get instance information.")
				continue
			}

			if len(instancesResponse.Items) == 0 {
				continue
			}

			showInstances(instancesResponse.Items)
		}
	}
}

func showInstances(instances []*compute.Instance) {
	for _, instance := range instances {
		fmt.Println("-------------------------------")
		fmt.Println("Name        :", instance.Name)
		fmt.Println("Id          :", instance.Id)
		fmt.Println("Zone        :", instance.Zone)
		fmt.Println("Tags        :", instance.Tags)
		fmt.Println("MachineType :", instance.MachineType)
		fmt.Println("Status      :", instance.Status)
		fmt.Println("Description :", instance.Description)
	}
}

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
