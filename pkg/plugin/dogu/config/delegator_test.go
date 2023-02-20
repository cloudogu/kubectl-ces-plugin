package config

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry/mocks"
)

const testNameSpace = "test-namespace"
const testDoguName = "official/ldap"

//go:embed testdata/ldap-dogu.json
var ldapBytes []byte

func Test_doguConfigurationDelegator_Delegate(t *testing.T) {
	t.Run("should return any error during port forward", func(t *testing.T) {
		// given
		portForwarderMock := newMockPortForwarder(t)
		portForwarderMock.EXPECT().ExecuteWithPortForward(mocks.Anything).RunAndReturn(func(func() error) error {
			return assert.AnError
		})

		sut := &doguConfigurationDelegator{
			doguName:  testDoguName,
			forwarder: portForwarderMock,
		}

		// when
		err := sut.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
			return nil
		})

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should return error during getting dogu json", func(t *testing.T) {
		// given
		portForwarderMock := newMockPortForwarder(t)
		portForwarderMock.EXPECT().ExecuteWithPortForward(mocks.Anything).RunAndReturn(func(payload func() error) error {
			return payload()
		})

		doguRegMock := newMockDoguRegistry(t)
		doguRegMock.EXPECT().Get(testDoguName).Return(nil, assert.AnError)
		sut := &doguConfigurationDelegator{
			doguName:  testDoguName,
			forwarder: portForwarderMock,
			doguReg:   doguRegMock,
		}

		// when
		err := sut.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
			return nil
		})

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not get dogu")
		assert.ErrorIs(t, err, assert.AnError)
	})
	// t.Run("should return error during getting dogu json", func(t *testing.T) {
	// 	// given
	// 	portForwarderMock := newMockPortForwarder(t)
	// 	portForwarderMock.EXPECT().ExecuteWithPortForward(mocks.Anything).RunAndReturn(func(payload func() error) error {
	// 		return payload()
	// 	})
	// 	dogu := readDoguResource(t, ldapBytes)
	// 	doguRegMock := newMockDoguRegistry(t)
	// 	doguRegMock.EXPECT().Get(testDoguName).Return(dogu, nil)
	// 	sut := &doguConfigurationDelegator{
	// 		doguName:  testDoguName,
	// 		forwarder: portForwarderMock,
	// 		doguReg:   doguRegMock,
	// 	}
	//
	// 	// when
	// 	err := sut.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
	// 		return nil
	// 	})
	//
	// 	// then
	// 	require.NoError(t, err)
	// })
}

func readDoguResource(t *testing.T, doguResourceBytes []byte) *core.Dogu {
	t.Helper()

	data := &core.Dogu{}
	err := json.Unmarshal(doguResourceBytes, data)
	if err != nil {
		t.Fatal(err.Error())
	}

	return data
}
