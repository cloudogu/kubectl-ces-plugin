package portforward

import (
	"io"
	"k8s.io/apimachinery/pkg/util/httpstream"
)

type portForwarder interface {
	// ForwardPorts formats and executes a port forwarding request.
	ForwardPorts() error
}

type portForwarderFactory interface {
	create(portForwardDialer dialer, stopCh, readyCh chan struct{}, out, errOut io.Writer) (portForwarder, error)
}

type dialerFactory interface {
	create() (dialer, error)
}

type dialer interface {
	httpstream.Dialer
}
