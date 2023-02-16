package cli

import (
	"bytes"
	"github.com/cloudogu/kubectl-ces-plugin/cmd/plugin/cli/dogu/config"
	config2 "github.com/cloudogu/kubectl-ces-plugin/pkg/plugin/dogu/config"
	"github.com/spf13/viper"
	"testing"
)

func Test_getAllForDoguCmd(t *testing.T) {

	actual := new(bytes.Buffer)
	RootCmd := RootCmd()
	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	RootCmd.SetArgs([]string{"dogu", "config", "list", "redmine"})
	//TODO: mock config service, so that only the CLI is tested, inject mock by overriding the factory
	config.DoguConfigServiceFactory = func(viper *viper.Viper) (*config2.DoguConfigService, error) {
		return nil, nil //inject mock here
	}
	//err := RootCmd.Execute()

	//assert.NoError(t, err, "command should be successful")
	//TODO: test specific output
	//assert.Equal(t, "config for redmine\n", actual.String())
}
