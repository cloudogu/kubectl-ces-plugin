package dogu_config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
	"testing"
)

const testNamespace = "ecosystem"

func Test_defaultServiceFactory_create(t *testing.T) {
	t.Run("should fail to create rest config", func(t *testing.T) {
		// given
		cfgMock := newMockRestClientGetter(t)
		cfgMock.EXPECT().ToRESTConfig().Return(nil, assert.AnError).Once()
		sut := defaultServiceFactory{
			namespace:   testNamespace,
			configFlags: cfgMock,
		}

		// when
		actual, err := sut.create(testDoguName)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "could not create rest config")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		cfgMock := newMockRestClientGetter(t)
		cfgMock.EXPECT().ToRESTConfig().Return(&rest.Config{}, nil).Once()
		sut := defaultServiceFactory{
			namespace:   testNamespace,
			configFlags: cfgMock,
		}

		// when
		actual, err := sut.create(testDoguName)

		// then
		require.NoError(t, err)
		assert.NotNil(t, actual)
	})
}
