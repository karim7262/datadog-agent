package main

import (
	"fmt"

	"github.com/DataDog/datadog-agent/pkg/util/docker"
)

func main() {
	du, err := docker.GetDockerUtil()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	listConfig := &docker.ContainerListConfig{
		IncludeExited: false,
		FlagExcluded:  false,
	}

	du.ListContainers(listConfig)
}
