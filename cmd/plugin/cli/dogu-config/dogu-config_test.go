package dogu_config

import "bytes"

func (s *DoguConfigCLITestSuite) TestCmd() {
	s.Run("should print usage", func() {
		// given
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		sut := Cmd()
		sut.SetOut(outBuf)
		sut.SetErr(errBuf)

		// when
		sut.SetArgs([]string{})
		err := sut.Execute()

		// then
		s.NoError(err)
		s.Contains(outBuf.String(), "Usage:", "should have usage output")
	})
}
