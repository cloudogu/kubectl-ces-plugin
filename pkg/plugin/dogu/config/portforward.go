package config

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

type PortForwardType string

const (
	PodType        PortForwardType = "pods"
	DeploymentType PortForwardType = "deployments"
	ReplicaSetType PortForwardType = "replicasets"
	ServiceType    PortForwardType = "services"
)

// PortForwarder establishes a port-forward until the given function returns.
type PortForwarder interface {
	ExecuteWithPortForward(fn func() error) error
}

// KubernetesPortForwarder establishes a kubernetes port-forward for the specified resource.
type KubernetesPortForwarder struct {
	// RestConfig is the kubernetes config
	RestConfig *rest.Config
	// Type of the resource to port forward to
	Type PortForwardType
	// NamespacedName references the selected resource for this port forwarding
	NamespacedName types.NamespacedName
	// LocalPort is the local port that will be selected to expose the ClusterPort
	LocalPort int
	// ClusterPort is the target port for the pod
	ClusterPort int
}

// ExecuteWithPortForward establishes a port-forward until the given function returns.
func (spf KubernetesPortForwarder) ExecuteWithPortForward(fn func() error) error {
	path := fmt.Sprintf("/api/v1/namespaces/%s/%s/%s/portforward",
		spf.NamespacedName.Namespace, spf.Type, spf.NamespacedName.Name)
	hostIP := strings.TrimPrefix(spf.RestConfig.Host, "https://")

	transport, upgrader, err := spdy.RoundTripperFor(spf.RestConfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})

	stopCh := make(chan struct{})
	defer close(stopCh)
	readyCh := make(chan struct{})

	stdOut := new(bytes.Buffer)
	errOut := new(bytes.Buffer)

	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", spf.LocalPort, spf.ClusterPort)}, stopCh, readyCh, stdOut, errOut)
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
