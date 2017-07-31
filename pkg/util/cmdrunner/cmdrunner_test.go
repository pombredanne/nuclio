package cmdrunner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nuclio/nuclio/pkg/zap"

	"github.com/nuclio/nuclio-sdk"
	"github.com/stretchr/testify/suite"
)

type CmdRunnerTestSuite struct {
	suite.Suite
	logger        nuclio.Logger
	commandRunner *CmdRunner
}

func (suite *CmdRunnerTestSuite) SetupSuite() {
	var err error

	suite.logger, _ = nucliozap.NewNuclioZap("test", nucliozap.ErrorLevel)
	suite.commandRunner, err = NewCmdRunner(suite.logger)
	if err != nil {
		panic("Failed to create command runner")
	}
}

func (suite *CmdRunnerTestSuite) TestWorkingDir() {
	currentDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		suite.Fail("Failed to get current directory")
	}

	options := RunOptions{
		WorkingDir: &currentDirectory,
	}

	output, err := suite.commandRunner.Run(&options, "pwd")
	suite.NoError(err)

	// remove "private" on OSX
	privatePrefix := "/private"
	if strings.HasPrefix(output, privatePrefix) {
		output = output[len(privatePrefix):]
	}

	suite.True(strings.HasPrefix(output, currentDirectory))
}

func (suite *CmdRunnerTestSuite) TestFormattedCommand() {
	output, err := suite.commandRunner.Run(nil, `echo "%s %d"`, "hello", 1)
	suite.NoError(err)

	// ignore newlines, if any
	suite.True(strings.HasPrefix(output, "hello 1"))
}

func (suite *CmdRunnerTestSuite) TestEnv() {
	options := RunOptions{
		Env: map[string]string{
			"ENV1": "env1",
			"ENV2": "env2",
		},
	}

	output, err := suite.commandRunner.Run(&options, `echo $ENV1 && echo $ENV2`)
	suite.NoError(err)

	// ignore newlines, if any
	suite.True(strings.HasPrefix(output, "env1\nenv2"))
}

func (suite *CmdRunnerTestSuite) TestStdin() {
	stdinValue := "from stdin"

	options := RunOptions{
		Stdin: &stdinValue,
	}

	output, err := suite.commandRunner.Run(&options, "more")
	suite.NoError(err)

	// ignore newlines, if any
	suite.True(strings.HasPrefix(output, stdinValue))
}

func (suite *CmdRunnerTestSuite) TestBadShell() {
	commandRunner, err := NewCmdRunner(suite.logger)
	if err != nil {
		panic("Failed to create command runner")
	}

	commandRunner.SetShell("/bin/definitelynotashell")

	_, err = commandRunner.Run(nil, `pwd`)
	suite.Error(err)
}

func TestCmdRunnerTestSuite(t *testing.T) {
	suite.Run(t, new(CmdRunnerTestSuite))
}