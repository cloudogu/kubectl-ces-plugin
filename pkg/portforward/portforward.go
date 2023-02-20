package portforward

import (
	"bytes"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"net/http"
	"net/url"
	"strings"
)

type ResourceType string

const (
	PodType        ResourceType = "pods"
	DeploymentType ResourceType = "deployments"
	ReplicaSetType ResourceType = "replicasets"
	ServiceType    ResourceType = "services"
)

func NewPortForwarder(
	restConfig *rest.Config,
	resourceType ResourceType,
	name types.NamespacedName,
	localPort int,
	clusterPort int,
) *kubernetesPortForwarder {
	return &kubernetesPortForwarder{
		restConfig:   restConfig,
		resourceType: resourceType,
		name:         name,
		localPort:    localPort,
		clusterPort:  clusterPort,
	}
}

// kubernetesPortForwarder establishes a kubernetes port-forward for the specified resource.
type kubernetesPortForwarder struct {
	// restConfig is the kubernetes config
	restConfig *rest.Config
	// resourceType of the resource to port forward to
	resourceType ResourceType
	// name references the selected resource for this port forwarding
	name types.NamespacedName
	// localPort is the local port that will be selected to expose the clusterPort
	localPort int
	// clusterPort is the target port for the pod
	clusterPort int
}

// ExecuteWithPortForward establishes a port-forward until the given function returns.
func (spf kubernetesPortForwarder) ExecuteWithPortForward(fn func() error) error {
	path := fmt.Sprintf("/api/v1/namespaces/%s/%s/%s/portforward",
		spf.name.Namespace, spf.resourceType, spf.name.Name)
	hostIP := strings.TrimPrefix(spf.restConfig.Host, "https://")

	transport, upgrader, err := spdy.RoundTripperFor(spf.restConfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})

	stopCh := make(chan struct{})
	defer close(stopCh)
	readyCh := make(chan struct{})

	stdOut := new(bytes.Buffer)
	errOut := new(bytes.Buffer)

	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", spf.localPort, spf.clusterPort)}, stopCh, readyCh, stdOut, errOut)
	if err != nil {
		return err
	}

	err = fw.ForwardPorts()
	if err != nil {
		return fmt.Errorf("could not forward port; stdOut: %s; errOut: %s: %w", stdOut, errOut, err)
	}

	//wait for the port forward to be established
	<-readyCh

	err = fn()
	if err != nil {
		return fmt.Errorf("encoutered error during port-forward; stdOut: %s; errOut: %s: %w", stdOut, errOut, err)
	}

	return nil
}
