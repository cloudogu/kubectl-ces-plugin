package portforward

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	"github.com/cloudogu/kubectl-ces-plugin/pkg/logger"
)

// New creates a new executive port forwarder.
func New(restConfig *rest.Config, pod types.NamespacedName, localPort, clusterPort int) *executivePortForwarder {
	return &executivePortForwarder{
		dialerFactory: &portForwardDialerFactory{
			restConfig: restConfig,
			pod:        pod,
		},
		forwarderFactory: &kubernetesPortForwarderFactory{
			localPort:   localPort,
			clusterPort: clusterPort,
		},
	}
}

// executivePortForwarder establishes a kubernetes port-forward to the specified resource for the duration of a function call.
type executivePortForwarder struct {
	// dialerFactory creates
	dialerFactory dialerFactory
	// forwarderFactory creates the native port forward in use by this
	forwarderFactory portForwarderFactory
}

// ExecuteWithPortForward establishes a port-forward until the given function returns.
func (epf *executivePortForwarder) ExecuteWithPortForward(fn func() error) error {
	dialer, err := epf.dialerFactory.create()
	if err != nil {
		return fmt.Errorf("failed to create dialer for port-forward: %w", err)
	}

	stopCh := make(chan struct{})
	defer func() {
		close(stopCh)
		logger.GetInstance().Info("Closing port-forward")
	}()
	readyCh := make(chan struct{})

	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)

	fw, err := epf.forwarderFactory.create(dialer, stopCh, readyCh, out, errOut)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes port forwarder: %w", err)
	}

	errCh := make(chan error)
	go func() {
		err2 := fw.ForwardPorts()
		errCh <- err2
	}()

	// wait for the port forward to be established
	select {
	case err := <-errCh:
		return fmt.Errorf("could not create port-forward: %w", err)
	case <-readyCh:
		logger.GetInstance().Info("Port forward is ready")
	}

	err = fn()
	if err != nil {
		logger.GetInstance().Debugf("encountered error during port-forward; out: %s; errOut: %s", out, errOut)
		return fmt.Errorf("encountered error during port-forward: %w", err)
	}

	return nil
}

type portForwardDialerFactory struct {
	// restConfig is the kubernetes config
	restConfig *rest.Config
	// pod references the selected pod for this port forwarding
	pod types.NamespacedName
}

func (pfdf *portForwardDialerFactory) create() (dialer, error) {
	apiUrl, err := createApiUrl(pfdf.restConfig.Host, pfdf.pod)
	if err != nil {
		return nil, err
	}

	transport, upgrader, err := spdy.RoundTripperFor(pfdf.restConfig)
	if err != nil {
		return nil, err
	}

	return spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, apiUrl), nil
}

func createApiUrl(host string, pod types.NamespacedName) (*url.URL, error) {
	apiAddress := fmt.Sprintf("%s/api/v1/namespaces/%s/pods/%s/portforward",
		host, pod.Namespace, pod.Name)
	apiUrl, err := url.Parse(apiAddress)
	if err != nil {
		return nil, err
	}

	return apiUrl, nil
}

type kubernetesPortForwarderFactory struct {
	// localPort is the local port that will be selected to expose the clusterPort
	localPort int
	// clusterPort is the target port for the pod
	clusterPort int
}

func (pff *kubernetesPortForwarderFactory) create(dialer dialer, stopCh, readyCh chan struct{}, out, errOut io.Writer) (portForwarder, error) {
	pf, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", pff.localPort, pff.clusterPort)}, stopCh, readyCh, out, errOut)
	if err != nil {
		return nil, err
	}

	logger.GetInstance().Infof("Starting port-forward %d:%d\n", pff.localPort, pff.clusterPort)
	return pf, nil
}
