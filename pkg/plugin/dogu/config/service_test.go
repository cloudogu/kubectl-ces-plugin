package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

func TestNew(t *testing.T) {
	t.Run("should at least throw an error", func(t *testing.T) {
		// given

		// when
		_, err := New(testDoguName, testNameSpace, &rest.Config{})

		// then
		require.Error(t, err)
	})
}
