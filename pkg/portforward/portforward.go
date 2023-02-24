package portforward

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// New creates a new port forwarder.
func New(restConfig *rest.Config, pod types.NamespacedName, localPort int, clusterPort int) *kubernetesPortForwarder {
	return &kubernetesPortForwarder{
		restConfig:  restConfig,
		pod:         pod,
		localPort:   localPort,
		clusterPort: clusterPort,
	}
}

// kubernetesPortForwarder establishes a kubernetes port-forward for the specified resource.
type kubernetesPortForwarder struct {
	// restConfig is the kubernetes config
	restConfig *rest.Config
	// pod references the selected pod for this port forwarding
	pod types.NamespacedName
	// localPort is the local port that will be selected to expose the clusterPort
	localPort int
	// clusterPort is the target port for the pod
	clusterPort int
}

// ExecuteWithPortForward establishes a port-forward until the given function returns.
func (kpf *kubernetesPortForwarder) ExecuteWithPortForward(fn func() error) error {
	apiAddress := fmt.Sprintf("%s/api/v1/namespaces/%s/pods/%s/portforward",
		kpf.restConfig.Host, kpf.pod.Namespace, kpf.pod.Name)
	apiUrl, err := url.Parse(apiAddress)
	if err != nil {
		return err
	}

	transport, upgrader, err := spdy.RoundTripperFor(kpf.restConfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, apiUrl)

	stopCh := make(chan struct{})
	defer close(stopCh)
	readyCh := make(chan struct{})

	stdOut := new(bytes.Buffer)
	errOut := new(bytes.Buffer)

	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", kpf.localPort, kpf.clusterPort)}, stopCh, readyCh, stdOut, errOut)
	if err != nil {
		return err
	}

	fmt.Printf("Starting port-forward %d:%d\n", kpf.localPort, kpf.clusterPort)

	go fw.ForwardPorts()
	// if err != nil {
	// 	return fmt.Errorf("could not forward port; stdOut: %s; errOut: %s: %w", stdOut, errOut, err)
	// }

	// wait for the port forward to be established
	<-readyCh

	err = fn()
	if err != nil {
		return fmt.Errorf("encoutered error during port-forward; stdOut: %s; errOut: %s: %w", stdOut, errOut, err)
	}

	return nil
}
