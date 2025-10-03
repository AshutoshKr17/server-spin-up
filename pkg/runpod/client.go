package runpod

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gpu-cloud-manager/pkg/types"
)

const (
	BaseURL = "https://api.runpod.io/graphql"
)

// Client represents a RunPod API client
type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewClient creates a new RunPod API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: BaseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RunPodTemplate represents a GPU template from RunPod
type RunPodTemplate struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	DockerImage   string  `json:"dockerImage"`
	ContainerDisk int     `json:"containerDiskInGb"`
	VolumeInGb    int     `json:"volumeInGb"`
	VolumeMountPath string `json:"volumeMountPath"`
	Ports         []string `json:"ports"`
	Env           []map[string]string `json:"env"`
}

// RunPodPod represents a GPU instance from RunPod
type RunPodPod struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Runtime         *PodRuntime `json:"runtime"`
	Machine         *Machine    `json:"machine"`
	ImageName       string      `json:"imageName"`
	ContainerDisk   int         `json:"containerDiskInGb"`
	VolumeInGb      int         `json:"volumeInGb"`
	CostPerHr       float64     `json:"costPerHr"`
	DesiredStatus   string      `json:"desiredStatus"`
	LastStatusChange string     `json:"lastStatusChange"`
}

// PodRuntime represents runtime information
type PodRuntime struct {
	UptimeInSeconds int `json:"uptimeInSeconds"`
	Ports           []PortInfo `json:"ports"`
	GPUs            []GPUInfo  `json:"gpus"`
}

// PortInfo represents port mapping information
type PortInfo struct {
	IP                 string `json:"ip"`
	IsIpPublic         bool   `json:"isIpPublic"`
	PrivatePort        int    `json:"privatePort"`
	PublicPort         int    `json:"publicPort"`
	Type              string `json:"type"`
}

// GPUInfo represents GPU information
type GPUInfo struct {
	ID               string `json:"id"`
	GPUUtilPercent   int    `json:"gpuUtilPercent"`
	MemoryUtilPercent int   `json:"memoryUtilPercent"`
}

// Machine represents machine specifications
type Machine struct {
	PodHostID       string `json:"podHostId"`
	GPUCount        int    `json:"gpuCount"`
	CPUCount        int    `json:"cpuCount"`
	MemoryInGb      int    `json:"memoryInGb"`
	GPUDisplayName  string `json:"gpuDisplayName"`
	SecureCloud     bool   `json:"secureCloud"`
	Location        string `json:"location"`
}

// SearchPods searches for available GPU pods
func (c *Client) SearchPods() ([]RunPodPod, error) {
	query := `
	query {
		myself {
			pods {
				id
				name
				runtime {
					uptimeInSeconds
					ports {
						ip
						isIpPublic
						privatePort
						publicPort
						type
					}
					gpus {
						id
						gpuUtilPercent
						memoryUtilPercent
					}
				}
				machine {
					podHostId
					gpuCount
					cpuCount
					memoryInGb
					gpuDisplayName
					secureCloud
					location
				}
				imageName
				containerDiskInGb
				volumeInGb
				costPerHr
				desiredStatus
				lastStatusChange
			}
		}
	}`

	var response struct {
		Data struct {
			Myself struct {
				Pods []RunPodPod `json:"pods"`
			} `json:"myself"`
		} `json:"data"`
	}

	err := c.makeGraphQLRequest(query, nil, &response)
	if err != nil {
		return nil, err
	}

	return response.Data.Myself.Pods, nil
}

// GetGPUTypes retrieves available GPU types
func (c *Client) GetGPUTypes() ([]GPUType, error) {
	query := `
	query {
		gpuTypes {
			id
			displayName
			memoryInGb
			secureCloud
			communityCloud
			lowestPrice {
				minimumBidPrice
				uninterruptablePrice
			}
		}
	}`

	var response struct {
		Data struct {
			GPUTypes []GPUType `json:"gpuTypes"`
		} `json:"data"`
	}

	err := c.makeGraphQLRequest(query, nil, &response)
	if err != nil {
		return nil, err
	}

	return response.Data.GPUTypes, nil
}

// GPUType represents a GPU type available on RunPod
type GPUType struct {
	ID           string  `json:"id"`
	DisplayName  string  `json:"displayName"`
	MemoryInGb   int     `json:"memoryInGb"`
	SecureCloud  bool    `json:"secureCloud"`
	CommunityCloud bool  `json:"communityCloud"`
	LowestPrice  *Price  `json:"lowestPrice"`
}

// Price represents pricing information
type Price struct {
	MinimumBidPrice     float64 `json:"minimumBidPrice"`
	UninterruptablePrice float64 `json:"uninterruptablePrice"`
}

// CreatePodRequest represents parameters for creating a pod
type CreatePodRequest struct {
	Name            string            `json:"name"`
	ImageName       string            `json:"imageName"`
	GPUTypeID       string            `json:"gpuTypeId"`
	CloudType       string            `json:"cloudType"` // SECURE, COMMUNITY
	SupportPublicIp bool              `json:"supportPublicIp"`
	StartJupyter    bool              `json:"startJupyter"`
	StartSsh        bool              `json:"startSsh"`
	ContainerDisk   int               `json:"containerDiskInGb"`
	VolumeInGb      int               `json:"volumeInGb"`
	VolumeMountPath string            `json:"volumeMountPath"`
	Ports           string            `json:"ports"`
	Env             []EnvVar          `json:"env"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreatePod creates a new GPU pod
func (c *Client) CreatePod(req *CreatePodRequest) (*RunPodPod, error) {
	mutation := `
	mutation createPod($input: PodCreateInput!) {
		podCreate(input: $input) {
			id
			name
			imageName
			containerDiskInGb
			volumeInGb
			costPerHr
			desiredStatus
			machine {
				podHostId
				gpuCount
				cpuCount
				memoryInGb
				gpuDisplayName
				secureCloud
				location
			}
		}
	}`

	variables := map[string]interface{}{
		"input": req,
	}

	var response struct {
		Data struct {
			PodCreate RunPodPod `json:"podCreate"`
		} `json:"data"`
	}

	err := c.makeGraphQLRequest(mutation, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.PodCreate, nil
}

// StopPod stops a running pod
func (c *Client) StopPod(podID string) error {
	mutation := `
	mutation stopPod($input: PodStopInput!) {
		podStop(input: $input) {
			id
			desiredStatus
		}
	}`

	variables := map[string]interface{}{
		"input": map[string]string{
			"podId": podID,
		},
	}

	var response struct {
		Data struct {
			PodStop struct {
				ID            string `json:"id"`
				DesiredStatus string `json:"desiredStatus"`
			} `json:"podStop"`
		} `json:"data"`
	}

	return c.makeGraphQLRequest(mutation, variables, &response)
}

// ResumePod resumes a stopped pod
func (c *Client) ResumePod(podID string) error {
	mutation := `
	mutation resumePod($input: PodResumeInput!) {
		podResume(input: $input) {
			id
			desiredStatus
		}
	}`

	variables := map[string]interface{}{
		"input": map[string]string{
			"podId": podID,
		},
	}

	var response struct {
		Data struct {
			PodResume struct {
				ID            string `json:"id"`
				DesiredStatus string `json:"desiredStatus"`
			} `json:"podResume"`
		} `json:"data"`
	}

	return c.makeGraphQLRequest(mutation, variables, &response)
}

// TerminatePod terminates a pod
func (c *Client) TerminatePod(podID string) error {
	mutation := `
	mutation terminatePod($input: PodTerminateInput!) {
		podTerminate(input: $input) {
			id
		}
	}`

	variables := map[string]interface{}{
		"input": map[string]string{
			"podId": podID,
		},
	}

	var response struct {
		Data struct {
			PodTerminate struct {
				ID string `json:"id"`
			} `json:"podTerminate"`
		} `json:"data"`
	}

	return c.makeGraphQLRequest(mutation, variables, &response)
}

// makeGraphQLRequest performs GraphQL requests to RunPod API
func (c *Client) makeGraphQLRequest(query string, variables interface{}, result interface{}) error {
	requestBody := map[string]interface{}{
		"query": query,
	}
	
	if variables != nil {
		requestBody["variables"] = variables
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Make the request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("error parsing response: %v", err)
		}
	}

	return nil
}

// Helper functions to convert between RunPod types and our internal types

// ConvertPodToGPUInstance converts a RunPod pod to our internal GPU instance type
func ConvertPodToGPUInstance(pod RunPodPod) types.GPUInstance {
	status := types.StatusOffline
	switch pod.DesiredStatus {
	case "RUNNING":
		status = types.StatusRunning
	case "STOPPED":
		status = types.StatusOffline
	case "EXITED":
		status = types.StatusOffline
	}

	instance := types.GPUInstance{
		ID:           fmt.Sprintf("runpod_%s", pod.ID),
		Provider:     types.RunPod,
		ProviderID:   pod.ID,
		Name:         pod.Name,
		Status:       status,
		PricePerHour: pod.CostPerHr,
		ProviderData: map[string]interface{}{
			"image_name":        pod.ImageName,
			"container_disk":    pod.ContainerDisk,
			"volume_gb":         pod.VolumeInGb,
			"desired_status":    pod.DesiredStatus,
			"last_status_change": pod.LastStatusChange,
		},
	}

	if pod.Machine != nil {
		instance.GPUModel = pod.Machine.GPUDisplayName
		instance.GPUCount = pod.Machine.GPUCount
		instance.CPUCount = pod.Machine.CPUCount
		instance.RAM = pod.Machine.MemoryInGb
		instance.Region = pod.Machine.Location
		instance.Storage = pod.ContainerDisk + pod.VolumeInGb

		// Add machine info to provider data
		machineData := map[string]interface{}{
			"pod_host_id":    pod.Machine.PodHostID,
			"secure_cloud":   pod.Machine.SecureCloud,
		}
		for k, v := range machineData {
			instance.ProviderData[k] = v
		}
	}

	if pod.Runtime != nil {
		instance.ProviderData["uptime_seconds"] = pod.Runtime.UptimeInSeconds
		
		if len(pod.Runtime.Ports) > 0 {
			instance.ProviderData["ports"] = pod.Runtime.Ports
		}
		
		if len(pod.Runtime.GPUs) > 0 {
			instance.ProviderData["gpu_utilization"] = pod.Runtime.GPUs
		}
	}

	return instance
}

// ConvertGPUTypeToGPUInstance converts a RunPod GPU type to a searchable instance
func ConvertGPUTypeToGPUInstance(gpuType GPUType) types.GPUInstance {
	instance := types.GPUInstance{
		ID:           fmt.Sprintf("runpod_type_%s", gpuType.ID),
		Provider:     types.RunPod,
		ProviderID:   gpuType.ID,
		Name:         fmt.Sprintf("RunPod %s", gpuType.DisplayName),
		Status:       types.StatusOffline,
		GPUModel:     gpuType.DisplayName,
		GPUCount:     1,
		RAM:          gpuType.MemoryInGb,
		Region:       "Global", // RunPod has multiple regions
		ProviderData: map[string]interface{}{
			"secure_cloud":    gpuType.SecureCloud,
			"community_cloud": gpuType.CommunityCloud,
		},
	}

	if gpuType.LowestPrice != nil {
		instance.PricePerHour = gpuType.LowestPrice.UninterruptablePrice
		instance.ProviderData["minimum_bid_price"] = gpuType.LowestPrice.MinimumBidPrice
		instance.ProviderData["uninterruptable_price"] = gpuType.LowestPrice.UninterruptablePrice
	}

	return instance
}
