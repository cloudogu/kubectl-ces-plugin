package util

// These constants identify arguments passed from the viper CLI to any service in order to avoid parameter pollution.
const (
	// CliTransportParamK8sArgs passes K8s kube config infos globally.
	CliTransportParamK8sArgs = "globalK8sArgs"
)
