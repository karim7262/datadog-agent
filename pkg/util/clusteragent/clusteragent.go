// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

package clusteragent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"strings"
	"time"

	"github.com/DataDog/datadog-agent/pkg/util/log"

	"github.com/DataDog/datadog-agent/pkg/api/security"
	"github.com/DataDog/datadog-agent/pkg/api/util"
	apiv1 "github.com/DataDog/datadog-agent/pkg/clusteragent/api/v1"
	"github.com/DataDog/datadog-agent/pkg/clusteragent/clusterchecks/types"
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util/retry"
	"github.com/DataDog/datadog-agent/pkg/version"
)

/*
Client to query the Datadog Cluster Agent (DCA) API.
*/

const (
	authorizationHeaderKey = "Authorization"
	// RealIPHeader refers to the cluster level check runner ip passed in the request headers
	RealIPHeader = "X-Real-Ip"
)

var globalClusterAgentClient *DCAClient

type metadataNames []string

// DCAClientInterface  is required to query the API of Datadog cluster agent
type DCAClientInterface interface {
	Version() version.Version
	ClusterAgentAPIEndpoint() string

	GetVersion() (version.Version, error)
	GetNodeLabels(nodeName string) (map[string]string, error)
	GetPodsMetadataForNode(nodeName string) (apiv1.NamespacesPodsStringsSet, error)
	GetKubernetesMetadataNames(nodeName, ns, podName string) ([]string, error)

	PostClusterCheckStatus(nodeName string, status types.NodeStatus) (types.StatusResponse, error)
	GetClusterCheckConfigs(nodeName string) (types.ConfigResponse, error)
	GetEndpointsCheckConfigs(nodeName string) (types.ConfigResponse, error)
}

// DCAClient is required to query the API of Datadog cluster agent
type DCAClient struct {
	// used to setup the DCAClient
	initRetry retry.Retrier

	clusterAgentAPIEndpoints      []string        // ${SCHEME}://${clusterAgentHost}:${PORT}
	ClusterAgentVersion           version.Version // Version of the cluster-agent we're connected to
	clusterAgentAPIClient         *http.Client
	clusterAgentAPIRequestHeaders http.Header
	lastSuccessfulAPIEndpoint     int
	leaderClient                  *leaderClient
}

// resetGlobalClusterAgentClient is a helper to remove the current DCAClient global
// It is ONLY to be used for tests
func resetGlobalClusterAgentClient() {
	globalClusterAgentClient = nil
}

// GetClusterAgentClient returns or init the DCAClient
func GetClusterAgentClient() (DCAClientInterface, error) {
	if globalClusterAgentClient == nil {
		globalClusterAgentClient = &DCAClient{}
		globalClusterAgentClient.initRetry.SetupRetrier(&retry.Config{
			Name:              "clusterAgentClient",
			AttemptMethod:     globalClusterAgentClient.init,
			Strategy:          retry.Backoff,
			InitialRetryDelay: 1 * time.Second,
			MaxRetryDelay:     5 * time.Minute,
		})
	}
	if err := globalClusterAgentClient.initRetry.TriggerRetry(); err != nil {
		log.Debugf("Cluster Agent init error: %v", err)
		return nil, err
	}
	return globalClusterAgentClient, nil
}

func (c *DCAClient) init() error {
	var err error

	c.clusterAgentAPIEndpoints, err = getClusterAgentEndpoints()
	if err != nil {
		return err
	}
	c.lastSuccessfulAPIEndpoint = 0

	authToken, err := security.GetClusterAgentAuthToken()
	if err != nil {
		return err
	}

	c.clusterAgentAPIRequestHeaders = http.Header{}
	c.clusterAgentAPIRequestHeaders.Set(authorizationHeaderKey, fmt.Sprintf("Bearer %s", authToken))
	podIP := config.Datadog.GetString("clc_runner_host")
	c.clusterAgentAPIRequestHeaders.Set(RealIPHeader, podIP)

	// TODO remove insecure
	c.clusterAgentAPIClient = util.GetClient(false)
	c.clusterAgentAPIClient.Timeout = 2 * time.Second

	// Validate the cluster-agent client by checking the version
	c.ClusterAgentVersion, err = c.GetVersion()
	if err != nil {
		return err
	}
	log.Infof("Successfully connected to the Datadog Cluster Agent %s", c.ClusterAgentVersion.String())

	// Clone the http client in a new client with built-in redirect handler
	// TODO: is this useful in any way?
	c.leaderClient = newLeaderClient(c.clusterAgentAPIClient, c.clusterAgentAPIEndpoints[0])

	return nil
}

// Version returns ClusterAgentVersion already stored in the DCAClient
func (c *DCAClient) Version() version.Version {
	return c.ClusterAgentVersion
}

// ClusterAgentAPIEndpoint returns the last successfully contacted Agent API Endpoint URL as a string
// if no endpoint has been contacted successfully, it will return the first configured
func (c *DCAClient) ClusterAgentAPIEndpoint() string {
	return c.clusterAgentAPIEndpoints[c.lastSuccessfulAPIEndpoint]
}

// getClusterAgentEndpoints provides list of validated https endpoints from configuration keys in datadog.yaml:
// 1st. configuration key "cluster_agent.url" (or the DD_CLUSTER_AGENT_URL environment variable),
//      add the https prefix if the scheme isn't specified, multiple URLs can be separated by semicolon
// 2nd. environment variables associated with "cluster_agent.kubernetes_service_name"
//      ${dcaServiceName}_SERVICE_HOST and ${dcaServiceName}_SERVICE_PORT (not possible to specify multiple URLs)
func getClusterAgentEndpoints() ([]string, error) {
	const configDcaURL = "cluster_agent.url"
	const configDcaSvcName = "cluster_agent.kubernetes_service_name"

	dcaURL := config.Datadog.GetString(configDcaURL)
	// The following part is primarily (but not only) used on Cloud Foundry where we might have multiple DCA endpoints
	if dcaURL != "" {
		urls := strings.Split(dcaURL, ";")
		results := make([]string, 0, len(urls))
		for i, oneURL := range urls {
			if strings.HasPrefix(oneURL, "http://") {
				return nil, fmt.Errorf("cannot get cluster agent endpoint on position %d, not a https scheme: %s", i, oneURL)
			}
			if strings.Contains(oneURL, "://") == false {
				log.Tracef("Adding https scheme to URL %s on position %d: https://%s", oneURL, i, oneURL)
				oneURL = fmt.Sprintf("https://%s", oneURL)
			}
			u, err := url.Parse(oneURL)
			if err != nil {
				return nil, fmt.Errorf("failed parsing URL %s on position %d: %s", oneURL, i, err.Error())
			}
			if u.Scheme != "https" {
				return nil, fmt.Errorf("cannot get cluster agent endpoint on position %d, not a https scheme: %s", i, u.Scheme)
			}
			log.Debugf("Will attempt connecting to the configured URL for the Datadog Cluster Agent: %s", oneURL)
			results = append(results, u.String())
		}
		return results, nil
	}

	// The following part is Kubernetes specific and will always return exactly one endpoint (or none if there's error)
	// Construct the URL with the Kubernetes service environment variables
	// *_SERVICE_HOST and *_SERVICE_PORT
	dcaSvc := config.Datadog.GetString(configDcaSvcName)
	log.Debugf("Identified service for the Datadog Cluster Agent: %s", dcaSvc)
	if dcaSvc == "" {
		return nil, fmt.Errorf("cannot get a cluster agent endpoint, both %s and %s are empty", configDcaURL, configDcaSvcName)
	}

	dcaSvc = strings.ToUpper(dcaSvc)
	dcaSvc = strings.Replace(dcaSvc, "-", "_", -1) // Kubernetes replaces "-" with "_" in the service names injected in the env var.

	// host
	dcaSvcHostEnv := fmt.Sprintf("%s_SERVICE_HOST", dcaSvc)
	dcaSvcHost := os.Getenv(dcaSvcHostEnv)
	if dcaSvcHost == "" {
		return nil, fmt.Errorf("cannot get a cluster agent endpoint for kubernetes service %s, env %s is empty", dcaSvc, dcaSvcHostEnv)
	}

	// port
	dcaSvcPort := os.Getenv(fmt.Sprintf("%s_SERVICE_PORT", dcaSvc))
	if dcaSvcPort == "" {
		return nil, fmt.Errorf("cannot get a cluster agent endpoint for kubernetes service %s, env %s is empty", dcaSvc, dcaSvcPort)
	}

	// validate the URL
	dcaURL = fmt.Sprintf("https://%s:%s", dcaSvcHost, dcaSvcPort)
	u, err := url.Parse(dcaURL)
	if err != nil {
		return nil, err
	}

	return []string{u.String()}, nil
}

func (c *DCAClient) httpDo(method, urlPath string) (*http.Response, error) {
	var response *http.Response
	var err error
	lastSuccess := c.clusterAgentAPIEndpoints[c.lastSuccessfulAPIEndpoint]
	response, err = c.httpDoOne(method, lastSuccess, urlPath)
	if err == nil && response.StatusCode != http.StatusMisdirectedRequest {
		return response, err
	}

	for i, url := range c.clusterAgentAPIEndpoints {
		if url == lastSuccess {
			continue
		} else if response != nil {
			// close the response body from last loop iteration, as we'll get a new response
			// (and either return that or close it here in the next iteration)
			response.Body.Close()
		}
		response, err = c.httpDoOne(method, url, urlPath)
		// only break if there's no error *and* we didn't get StatusMisdirectedRequest, IOW the request went ok and
		// was aimed at the current leader; otherwise just continue the loop
		if err == nil && response.StatusCode != http.StatusMisdirectedRequest {
			log.Infof("Successfully requested data from DCA at %s, will use it next time", url)
			c.lastSuccessfulAPIEndpoint = i
			break
		} else {
			errStr := "not a leader (code 421)"
			if err != nil {
				errStr = err.Error()
			}
			log.Debugf("Failed requesting data from DCA at %s: %s", url, errStr)
		}
	}
	return response, err
}

func (c *DCAClient) httpDoOne(method, url, urlPath string) (*http.Response, error) {
	log.Debugf("Trying DCA at %s", url)
	request, err := http.NewRequest(method, fmt.Sprintf("%s/%s", url, urlPath), nil)
	if err != nil {
		return nil, err
	}
	request.Header = c.clusterAgentAPIRequestHeaders
	return c.clusterAgentAPIClient.Do(request)
}

// GetVersion fetches the version of the Cluster Agent. Used in the agent status command.
func (c *DCAClient) GetVersion() (version.Version, error) {
	const dcaVersionPath = "version"
	var version version.Version
	var err error

	// https://host:port/version
	resp, err := c.httpDo("GET", dcaVersionPath)
	if err != nil {
		return version, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return version, fmt.Errorf("unexpected status code from cluster agent: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return version, err
	}

	err = json.Unmarshal(body, &version)

	return version, err
}

// GetNodeLabels returns the node labels from the Cluster Agent.
func (c *DCAClient) GetNodeLabels(nodeName string) (map[string]string, error) {
	const dcaNodeMeta = "api/v1/tags/node"
	var err error
	var labels map[string]string

	// https://host:port/api/v1/tags/node/{nodeName}
	rawURL := fmt.Sprintf("%s/%s", dcaNodeMeta, nodeName)
	resp, err := c.httpDo("GET", rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from cluster agent: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &labels)
	return labels, err
}

// GetPodsMetadataForNode queries the datadog cluster agent to get nodeName registered
// Kubernetes pods metadata.
func (c *DCAClient) GetPodsMetadataForNode(nodeName string) (apiv1.NamespacesPodsStringsSet, error) {
	const dcaMetadataPath = "api/v1/tags/pod"
	var err error

	if c == nil {
		return nil, fmt.Errorf("cluster agent's client is not properly initialized")
	}
	/* https://host:port/api/v1/tags/pod/{nodeName}
	response example:
	{
		"Nodes": {
			"node1": {
				"services": {
					"default": {
						"datadog-monitoring-cluster-agent-58f45b9b44-pkxrv": {
							"datadog-monitoring-cluster-agent": {},
							"datadog-monitoring-cluster-agent-metrics-api": {}
						}
					},
					"kube-system": {
						"kube-dns-6b98c9c9bf-ts7gc": {
							"kube-dns": {}
						}
					}
				}
			}
		}
	}
	*/
	rawURL := fmt.Sprintf("%s/%s", dcaMetadataPath, nodeName)
	resp, err := c.httpDo("GET", rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from cluster agent: %d", resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	metadataPodPayload := apiv1.NewMetadataResponse()
	if err = json.Unmarshal(b, metadataPodPayload); err != nil {
		return nil, err
	}

	if _, ok := metadataPodPayload.Nodes[nodeName]; !ok {
		return nil, fmt.Errorf("cluster agent didn't return pods metadata for node: %s", nodeName)
	}
	return metadataPodPayload.Nodes[nodeName].Services, nil
}

// GetKubernetesMetadataNames queries the datadog cluster agent to get nodeName/podName registered
// Kubernetes metadata.
func (c *DCAClient) GetKubernetesMetadataNames(nodeName, ns, podName string) ([]string, error) {
	const dcaMetadataPath = "api/v1/tags/pod"
	var metadataNames metadataNames
	var err error

	if c == nil {
		return nil, fmt.Errorf("cluster agent's client is not properly initialized")
	}
	if ns == "" {
		return nil, fmt.Errorf("namespace is empty")
	}

	// https://host:port/api/v1/metadata/{nodeName}/{ns}/{pod-[0-9a-z]+}
	rawURL := fmt.Sprintf("%s/%s/%s/%s", dcaMetadataPath, nodeName, ns, podName)
	resp, err := c.httpDo("GET", rawURL)
	if err != nil {
		return metadataNames, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return metadataNames, fmt.Errorf("unexpected status code from cluster agent: %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return metadataNames, err
	}
	err = json.Unmarshal(b, &metadataNames)
	if err != nil {
		return metadataNames, err
	}

	return metadataNames, nil
}
