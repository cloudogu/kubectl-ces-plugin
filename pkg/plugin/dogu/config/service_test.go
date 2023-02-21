package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func Test_doguConfigService_Edit(t *testing.T) {
	t.Run("should fail during delegate", func(t *testing.T) {
		// given
		delegator := newMockDelegator(t)
		delegator.EXPECT().Delegate(mock.Anything).Return(assert.AnError)
		sut := doguConfigService{delegator: delegator}

		// when
		err := sut.Edit("test-key", "test-value")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
