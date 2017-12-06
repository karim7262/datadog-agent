// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

package app

import (
	"fmt"
	"io/ioutil"

	"github.com/DataDog/datadog-agent/cmd/agent/common"
	"github.com/DataDog/datadog-agent/cmd/agent/gui"
	"github.com/DataDog/datadog-agent/pkg/config"
	log "github.com/cihub/seelog"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var (
	launchCmd = &cobra.Command{
		Use:          "launch-gui",
		Short:        "starts the Datadog Agent GUI",
		Long:         ``,
		RunE:         launchGui,
		SilenceUsage: true,
	}
)

func init() {
	// attach the command to the root
	AgentCmd.AddCommand(launchCmd)

}

func launchGui(cmd *cobra.Command, args []string) error {
	err := common.SetupConfig(confFilePath)
	if err != nil {
		return fmt.Errorf("unable to set up global agent configuration: %v", err)
	}

	guiPort := config.Datadog.GetString("GUI_port")
	if guiPort == "-1" {
		log.Warnf("GUI not enabled: to enable, please set an appropriate port in your datadog.yaml file")
		return fmt.Errorf("GUI not enabled: to enable, please set an appropriate port in your datadog.yaml file")
	}

	// Read the authentication token DIRECTLY from the configuration file
	// If the user doesn't have authorization to read it, they can't launch the GUI
	data, err := ioutil.ReadFile(config.Datadog.ConfigFileUsed())
	if err != nil {
		return fmt.Errorf("unable to access GUI authentication token: " + err.Error())
	}
	cfg := make(map[string]interface{})
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}
	authTokenRaw := cfg["GUI_auth_token"]
	authToken, ok := authTokenRaw.(string)
	if !ok {
		return fmt.Errorf("error: 'GUI_auth_token' in %s is not a string", config.Datadog.ConfigFileUsed())
	}

	// TODO fix bug: gui.CsrfToken is not available from this process...
	// probably need to write this to the config file too, or some other location that can be accessed
	// doesn't need to persist between restarts, just needs to be available between the processes

	// Open the GUI in a browser, passing the authorization tokens as parameters
	err = open("http://127.0.0.1:" + guiPort + "/authenticate?authToken=" + authToken + ";csrf=" + gui.CsrfToken)
	if err != nil {
		log.Warnf("error opening GUI: " + err.Error())
		return fmt.Errorf("error opening GUI: " + err.Error())
	}

	log.Infof("GUI opened at 127.0.0.1:" + guiPort)
	return nil
}
