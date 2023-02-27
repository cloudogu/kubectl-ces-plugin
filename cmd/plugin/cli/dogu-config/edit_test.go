package dogu_config

import (
	"bytes"

	"github.com/stretchr/testify/assert"
)

func (s *DoguConfigCLITestSuite) Test_editCmd() {
	s.Run("should set config value", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		sut := editCmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)
		doguName := "redmine"

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		const configKey = "redmineKey"
		const deleteOnEdit = false
		doguConfigServiceFactoryMock.EXPECT().Edit(configKey, deleteOnEdit).Return(nil).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		sut.SetArgs([]string{doguName, configKey})
		err := sut.Execute()

		// then
		s.NoError(err, "command should be successful")
		s.Empty(outBuf.String())
		s.Empty(errBuf.String())
	})

	s.Run("should return error from configService", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		sut := editCmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)
		doguName := "redmine"

		doguConfigServiceFactoryMock := newMockDoguConfigService(s.T())
		const configKey = "redmineKey"
		const deleteOnEdit = false
		doguConfigServiceFactoryMock.EXPECT().Edit(configKey, deleteOnEdit).Return(assert.AnError).Once()
		doguConfigServiceFactory = noopDoguConfigServiceFactory(doguConfigServiceFactoryMock)

		// when
		sut.SetArgs([]string{doguName, configKey})
		err := sut.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.ErrorContains(err, "cannot edit config keys")
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should return error that the config service cannot be created", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		sut := editCmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)
		doguName := "redmine"

		configKey := "redmineKey"
		doguConfigServiceFactory = errorDoguConfigServiceFactory(assert.AnError)

		// when
		sut.SetArgs([]string{doguName, configKey})
		err := sut.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.ErrorContains(err, "cannot create config service")
		s.ErrorIs(err, assert.AnError)
	})

	s.Run("should fail with too few Arguments", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		sut := editCmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		err := sut.Execute()

		// then
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
		s.Contains(errBuf.String(), err.Error(), "should contain error output")
		s.EqualError(err, "accepts between 1 and 2 arg(s), received 0")
	})
}
