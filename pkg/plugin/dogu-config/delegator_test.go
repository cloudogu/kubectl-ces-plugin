package dogu_config

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/keys"
	"github.com/cloudogu/cesapp-lib/registry/mocks"
)

const (
	testNameSpace = "test-namespace"
	testDoguName  = "ldap"
)

//go:embed testdata/ldap-dogu.json
var ldapBytes []byte

func Test_newDelegator(t *testing.T) {
	t.Run("should create delegator correctly", func(t *testing.T) {
		// given
		dogu := testDoguName
		mockForwarder := newMockPortForwarder(t)
		mockDoguReg := newMockDoguRegistry(t)
		mockReg := newMockCesRegistry(t)
		mockReg.EXPECT().DoguRegistry().Return(mockDoguReg).Once()

		// when
		actual := newDelegator(dogu, mockForwarder, mockReg)

		// then
		require.IsType(t, &doguConfigurationDelegator{}, actual)
		assert.NotNil(t, actual)
	})
}

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
	t.Run("should return without error because there are no config keys available for the dogu", func(t *testing.T) {
		// given
		portForwarderMock := newMockPortForwarder(t)
		portForwarderMock.EXPECT().ExecuteWithPortForward(mocks.Anything).RunAndReturn(func(payload func() error) error {
			return payload()
		})
		dogu := readDoguResource(t, ldapBytes)
		dogu.Configuration = []core.ConfigurationField{} // mock zero config keys
		doguRegMock := newMockDoguRegistry(t)
		doguRegMock.EXPECT().Get(testDoguName).Return(dogu, nil)
		sut := &doguConfigurationDelegator{
			doguName:  testDoguName,
			forwarder: portForwarderMock,
			doguReg:   doguRegMock,
		}

		realStdout := os.Stdout
		defer restoreOriginalStdout(realStdout)
		fakeReaderPipe, fakeWriterPipe := routeStdoutToReplacement()

		// when
		err := sut.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
			return nil
		})

		// then
		actual := captureOutput(fakeReaderPipe, fakeWriterPipe, realStdout)
		require.NoError(t, err)
		assert.Equal(t, "dogu 'ldap' has no configuration fields\n", actual)

	})
	t.Run("should return error when creating config editor", func(t *testing.T) {
		// given
		portForwarderMock := newMockPortForwarder(t)
		portForwarderMock.EXPECT().ExecuteWithPortForward(mocks.Anything).RunAndReturn(func(payload func() error) error {
			return payload()
		})
		dogu := readDoguResource(t, ldapBytes)
		doguRegMock := newMockDoguRegistry(t)
		doguRegMock.EXPECT().Get(testDoguName).Return(dogu, nil).Once()
		edFactoryMock := newMockEditorFactory(t)
		edFactoryMock.EXPECT().create().Return(nil, assert.AnError).Once()
		sut := &doguConfigurationDelegator{
			doguName:    testDoguName,
			forwarder:   portForwarderMock,
			doguReg:     doguRegMock,
			editFactory: edFactoryMock,
		}

		// when
		err := sut.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
			return nil
		})

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "could not create configuration editor for dogu 'ldap'")
	})
	t.Run("should return error from payload function", func(t *testing.T) {
		// given
		portForwarderMock := newMockPortForwarder(t)
		portForwarderMock.EXPECT().ExecuteWithPortForward(mocks.Anything).RunAndReturn(func(payload func() error) error {
			return payload()
		})
		dogu := readDoguResource(t, ldapBytes)
		doguRegMock := newMockDoguRegistry(t)
		doguRegMock.EXPECT().Get(testDoguName).Return(dogu, nil)
		editorMock := newMockDoguConfigurationEditor(t)
		edFactoryMock := newMockEditorFactory(t)
		edFactoryMock.EXPECT().create().Return(editorMock, nil).Once()
		sut := &doguConfigurationDelegator{
			doguName:    testDoguName,
			forwarder:   portForwarderMock,
			doguReg:     doguRegMock,
			editFactory: edFactoryMock,
		}

		// when
		err := sut.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
			return assert.AnError
		})

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error during registry interaction: ")
	})
	t.Run("should succeed for a reasonable dogu", func(t *testing.T) {
		// given
		portForwarderMock := newMockPortForwarder(t)
		portForwarderMock.EXPECT().ExecuteWithPortForward(mocks.Anything).RunAndReturn(func(payload func() error) error {
			return payload()
		})
		dogu := readDoguResource(t, ldapBytes)
		doguRegMock := newMockDoguRegistry(t)
		doguRegMock.EXPECT().Get(testDoguName).Return(dogu, nil)
		editorMock := newMockDoguConfigurationEditor(t)
		edFactoryMock := newMockEditorFactory(t)
		edFactoryMock.EXPECT().create().Return(editorMock, nil).Once()
		sut := &doguConfigurationDelegator{
			doguName:    testDoguName,
			forwarder:   portForwarderMock,
			doguReg:     doguRegMock,
			editFactory: edFactoryMock,
		}

		// when
		err := sut.Delegate(func(dogu *core.Dogu, editor doguConfigurationEditor) error {
			return nil
		})

		// then
		require.NoError(t, err)
	})
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

func routeStdoutToReplacement() (readerPipe, writerPipe *os.File) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	return r, w
}

func captureOutput(fakeReaderPipe, fakeWriterPipe, originalStdout *os.File) string {
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, fakeReaderPipe)
		outC <- buf.String()
	}()

	// back to normal state
	_ = fakeWriterPipe.Close()
	restoreOriginalStdout(originalStdout)

	actualOutput := <-outC

	return actualOutput
}

func restoreOriginalStdout(stdout *os.File) {
	os.Stdout = stdout
}

func Test_defaultEditorFactory_create(t *testing.T) {
	t.Run("should return error from key manager factory creation", func(t *testing.T) {
		// given
		keyManagerFactoryMock := newMockKeyManagerFactory(t)
		keyManagerFactoryMock.EXPECT().create(mocks.Anything, testDoguName).Return(nil, assert.AnError)

		sut := &defaultEditorFactory{
			doguName:          testDoguName,
			keyManagerFactory: keyManagerFactoryMock,
		}

		// when
		_, err := sut.create()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "could not create key manager for dogu '"+testDoguName+"'")
	})
	t.Run("should return error from key manager", func(t *testing.T) {
		// given
		keyManagerMock := newMockKeyManager(t)
		keyManagerMock.EXPECT().GetPublicKey().Return(nil, assert.AnError)
		keyManagerFactoryMock := newMockKeyManagerFactory(t)
		keyManagerFactoryMock.EXPECT().create(mocks.Anything, testDoguName).Return(keyManagerMock, nil)

		sut := &defaultEditorFactory{
			doguName:          testDoguName,
			keyManagerFactory: keyManagerFactoryMock,
		}

		// when
		_, err := sut.create()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "could not get public key for dogu '"+testDoguName+"'")
	})
	t.Run("should return error from config editor", func(t *testing.T) {
		// given
		mockReg := newMockCesRegistry(t)
		mockReg.EXPECT().DoguConfig(testDoguName).Return(&mockNoopDoguConfig{}).Once()

		keyManagerMock := newMockKeyManager(t)
		keyManagerMock.EXPECT().GetPublicKey().Return(&keys.PublicKey{}, nil)
		keyManagerFactoryMock := newMockKeyManagerFactory(t)
		keyManagerFactoryMock.EXPECT().create(mocks.Anything, testDoguName).Return(keyManagerMock, nil)

		sut := &defaultEditorFactory{
			doguName:          testDoguName,
			keyManagerFactory: keyManagerFactoryMock,
			registry:          mockReg,
		}

		// when
		actual, err := sut.create()

		// then
		require.NoError(t, err)
		assert.NotNil(t, actual)
	})
	t.Run("should return config editor", func(t *testing.T) {
		// given
		mockReg := newMockCesRegistry(t)
		mockReg.EXPECT().DoguConfig(testDoguName).Return(&mockNoopDoguConfig{}).Once()

		keyManagerMock := newMockKeyManager(t)
		keyManagerMock.EXPECT().GetPublicKey().Return(&keys.PublicKey{}, nil)
		keyManagerFactoryMock := newMockKeyManagerFactory(t)
		keyManagerFactoryMock.EXPECT().create(mocks.Anything, testDoguName).Return(keyManagerMock, nil)

		sut := &defaultEditorFactory{
			doguName:          testDoguName,
			keyManagerFactory: keyManagerFactoryMock,
			registry:          mockReg,
		}

		// when
		actual, err := sut.create()

		// then
		require.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

type mockNoopDoguConfig struct{}

func (m *mockNoopDoguConfig) Set(_, _ string) error {
	return nil
}

func (m *mockNoopDoguConfig) SetWithLifetime(_, _ string, _ int) error {
	return nil
}

func (m *mockNoopDoguConfig) Refresh(_ string, _ int) error {
	return nil
}

func (m *mockNoopDoguConfig) Get(_ string) (string, error) {
	return "", nil
}

func (m *mockNoopDoguConfig) GetAll() (map[string]string, error) {
	return nil, nil
}

func (m *mockNoopDoguConfig) Delete(_ string) error {
	return nil
}

func (m *mockNoopDoguConfig) DeleteRecursive(_ string) error {
	return nil
}

func (m *mockNoopDoguConfig) Exists(_ string) (bool, error) {
	return false, nil
}

func (m *mockNoopDoguConfig) RemoveAll() error {
	return nil
}

func (m *mockNoopDoguConfig) GetOrFalse(_ string) (bool, string, error) {
	return false, "", nil
}
