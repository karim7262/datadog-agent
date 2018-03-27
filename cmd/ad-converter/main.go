// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

package main

import (
	"errors"
	"os"

	log "github.com/cihub/seelog"
	"github.com/spf13/cobra"
)

var (
	convertCmd = &cobra.Command{
		Use:               "ad-convert [command]",
		Short:             "Converts yaml check configurations to Docker / Kubernetes labels",
		PersistentPreRunE: parseToJson,
		SilenceUsage:      true,
	}
	toDockerCmd = &cobra.Command{
		Use:     "dockerfile",
		Aliases: []string{"df"},
		Short:   "Convert to Dockerfile LABEL commands",
		Args:    cobra.MinimumNArgs(1),
		RunE:    toDocker,
	}
	toComposeCmd = &cobra.Command{
		Use:     "compose",
		Aliases: []string{"dc"},
		Short:   "Convert to docker compose labels",
		Args:    cobra.MinimumNArgs(1),
		RunE:    toCompose,
	}
	toKubeCmd = &cobra.Command{
		Use:     "kubernetes",
		Aliases: []string{"kube", "k"},
		Short:   "Convert to Kubernetes pod annotations",
		Args:    cobra.MinimumNArgs(1),
		RunE:    toKube,
	}
	kubeContainerName string
	outputJson        jsonConfigs
)

func init() {
	convertCmd.AddCommand(toDockerCmd)
	convertCmd.AddCommand(toComposeCmd)
	convertCmd.AddCommand(toKubeCmd)
	toKubeCmd.Flags().StringVarP(&kubeContainerName, "name", "n", "", "Container name to target")
}

func parseToJson(cmd *cobra.Command, args []string) error {
	configs, err := parseFiles(args)
	if err != nil {
		return err
	}
	outputJson, err = jsonizeConfig(configs)
	if err != nil {
		return err
	}
	return nil
}

func toDocker(cmd *cobra.Command, args []string) error {
	printDockerLabels(outputJson)
	return nil
}

func toCompose(cmd *cobra.Command, args []string) error {
	printComposeLabels(outputJson)
	return nil
}
func toKube(cmd *cobra.Command, args []string) error {
	if kubeContainerName == "" {
		return errors.New("Please provide a target container name with -n")
	}
	printKubeAnnotations(outputJson, kubeContainerName)
	return nil
}

func main() {
	if err := convertCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
