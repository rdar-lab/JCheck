package commands

import (
	"github.com/rdar-lab/JCheck/checks"
	"github.com/rdar-lab/JCheck/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type CheckCommandSuite struct {
	suite.Suite
}

func (suite *CheckCommandSuite) TestSanityCheckPass() {
	err := os.Setenv("PanicTest", "0")
	if err != nil {
		assert.Fail(suite.T(), "was unable to set env variable")
	}

	err = os.Setenv("FailureTest", "0")
	if err != nil {
		assert.Fail(suite.T(), "was unable to set env variable")
	}

	conf := &checkConfiguration{
		what: "SelfCheck",
		loop: 1,
	}

	err = doCheck(conf)
	assert.Nil(suite.T(), err)
}

func (suite *CheckCommandSuite) TestSanityCheckPanic() {
	err := os.Setenv("PanicTest", "1")
	if err != nil {
		assert.Fail(suite.T(), "was unable to set env variable")
	}

	err = os.Setenv("FailureTest", "0")
	if err != nil {
		assert.Fail(suite.T(), "was unable to set env variable")
	}

	conf := &checkConfiguration{
		what: "SelfCheck",
		loop: 1,
	}

	err = doCheck(conf)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), err.Error(), "failures detected")
}

func (suite *CheckCommandSuite) TestSanityCheckFailure() {
	err := os.Setenv("PanicTest", "0")
	if err != nil {
		assert.Fail(suite.T(), "was unable to set env variable")
	}

	err = os.Setenv("FailureTest", "1")
	if err != nil {
		assert.Fail(suite.T(), "was unable to set env variable")
	}

	conf := &checkConfiguration{
		what: "SelfCheck",
		loop: 1,
	}

	err = doCheck(conf)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), err.Error(), "failures detected")
}

func (suite *CheckCommandSuite) TestWhenNoChecksFound() {
	err := os.Setenv("PanicTest", "0")
	if err != nil {
		assert.Fail(suite.T(), "was unable to set env variable")
	}

	err = os.Setenv("FailureTest", "1")
	if err != nil {
		assert.Fail(suite.T(), "was unable to set env variable")
	}

	conf := &checkConfiguration{
		what: "NoChecksFound",
		loop: 1,
	}

	err = doCheck(conf)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), err.Error(), "no checks performed. Check 'what' argument")
}

func TestCheckCommandSuite(t *testing.T) {
	common.GetRegistry().Register(checks.GetSelfCheck())
	suite.Run(t, new(CheckCommandSuite))
}
