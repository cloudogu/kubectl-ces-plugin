package config

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"
)

type servicePortForward struct {
	// RestConfig is the kubernetes config
	RestConfig *rest.Config
	// Service is the selected service for this port forwarding
	Service v1.Service
	// LocalPort is the local port that will be selected to expose the PodPort
	LocalPort int
	// PodPort is the target port for the pod
	PodPort int
	// Streams configures where to write or read input from
	Streams genericclioptions.IOStreams
}

func (spf servicePortForward) ExecuteWithPortForward(fn func() error) error {
	path := fmt.Sprintf("/api/v1/namespaces/%s/services/%s/portforward",
		spf.Service.Namespace, spf.Service.Name)
	hostIP := strings.TrimPrefix(spf.RestConfig.Host, "https://")

	transport, upgrader, err := spdy.RoundTripperFor(spf.RestConfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	stopCh := make(chan struct{})
	defer close(stopCh)
	readyCh := make(chan struct{})
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", spf.LocalPort, spf.PodPort)}, stopCh, readyCh, spf.Streams.Out, spf.Streams.ErrOut)
	if err != nil {
		return err
	}

	err = fw.ForwardPorts()
	if err != nil {
		return err
	}

	<-readyCh

	return fn()
}

func NewDoguConfigService(namespace string) (*DoguConfigService, error) {
	endpoint := fmt.Sprintf("http://etcd.%s.svc.cluster.local:4001", namespace)
	reg, err := registry.New(core.Registry{
		Type:      "etcd",
		Endpoints: []string{endpoint},
	})
	if err != nil {
		return nil, err
	}

	return &DoguConfigService{
		registry:      reg,
		portForwarder: servicePortForward{}, //TODO
	}, nil
}

type DoguConfigService struct {
	registry      registry.Registry
	portForwarder servicePortForward
}

func (s DoguConfigService) Edit(doguName string, registryKey string, registryValue string) error {
	return s.portForwarder.ExecuteWithPortForward(func() error {
		err := s.checkInstallStatus(doguName)
		if err != nil {
			return err
		}

		err = s.registry.DoguConfig(doguName).Set(registryKey, registryValue)
		if err != nil {
			return fmt.Errorf("error while editing key '%s' for dogu %s: %w", registryKey, doguName, err)
		}

		return nil
	})
}

func (s DoguConfigService) Delete(doguName string, registryKey string) error {
	return s.portForwarder.ExecuteWithPortForward(func() error {
		err := s.checkInstallStatus(doguName)
		if err != nil {
			return err
		}

		err = s.registry.DoguConfig(doguName).Delete(registryKey)
		if err != nil {
			return fmt.Errorf("error while deleting key '%s' for dogu %s: %w", registryKey, doguName, err)
		}

		return nil
	})
}

func (s DoguConfigService) getAllForDogu(doguName string) (map[string]string, error) {
	var configEntries map[string]string
	err := s.portForwarder.ExecuteWithPortForward(func() error {
		err := s.checkInstallStatus(doguName)
		if err != nil {
			return err
		}

		configEntries, err = s.registry.DoguConfig(doguName).GetAll()
		if err != nil {
			return fmt.Errorf("error while reading all keys for dogu %s: %w", doguName, err)
		}

		return nil
	})
	return configEntries, err
}

func (s DoguConfigService) GetValue(doguName string, registryKey string) (string, error) {
	var configValue string
	err := s.portForwarder.ExecuteWithPortForward(func() error {
		err := s.checkInstallStatus(doguName)
		if err != nil {
			return err
		}

		configValue, err = s.registry.DoguConfig(doguName).Get(registryKey)
		if err != nil {
			return fmt.Errorf("error while reading key '%s' for dogu %s", registryKey, doguName)
		}
		return nil
	})
	return configValue, err
}

func (s DoguConfigService) checkInstallStatus(wantedDogu string) error {
	enabled, err := s.registry.DoguRegistry().IsEnabled(wantedDogu)
	if err != nil {
		return fmt.Errorf("cannot check if dogu '%s' is installed: %w", wantedDogu, err)
	}

	if !enabled {
		return fmt.Errorf("dogu %s is not installed", wantedDogu)
	}
	return nil
}
