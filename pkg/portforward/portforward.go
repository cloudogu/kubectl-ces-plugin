package portforward

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// ResourceType contains the type of the Kubernetes resource kind to which the port forward shall be created.
type ResourceType string

const (
	// PodType determines the port forward type that goes towards pods.
	PodType ResourceType = "pods"
	// DeploymentType determines the port forward type that goes towards deployments.
	DeploymentType ResourceType = "deployments"
	// ReplicaSetType determines the port forward type that goes towards replicasets.
	ReplicaSetType ResourceType = "replicasets"
	// ServiceType determines the port forward type that goes towards services.
	ServiceType ResourceType = "services"
)

// New creates a new port forwarder.
func New(
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
func (kpf *kubernetesPortForwarder) ExecuteWithPortForward(fn func() error) error {
	path := fmt.Sprintf("/api/v1/namespaces/%s/%s/%s/portforward",
		kpf.name.Namespace, kpf.resourceType, kpf.name.Name)
	hostIP := strings.TrimPrefix(kpf.restConfig.Host, "https://")

	transport, upgrader, err := spdy.RoundTripperFor(kpf.restConfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})

	stopCh := make(chan struct{})
	defer close(stopCh)
	readyCh := make(chan struct{})

	stdOut := new(bytes.Buffer)
	errOut := new(bytes.Buffer)

	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", kpf.localPort, kpf.clusterPort)}, stopCh, readyCh, stdOut, errOut)
	if err != nil {
		return err
	}

	err = fw.ForwardPorts()
	if err != nil {
		return fmt.Errorf("could not forward port; stdOut: %s; errOut: %s: %w", stdOut, errOut, err)
	}

	// wait for the port forward to be established
	<-readyCh

	err = fn()
	if err != nil {
		return fmt.Errorf("encoutered error during port-forward; stdOut: %s; errOut: %s: %w", stdOut, errOut, err)
	}

	return nil
}
