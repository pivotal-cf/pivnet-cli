package commands_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/printer"
)

const (
	apiPrefix = "/api/v2"
	apiToken  = "some-api-token"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var (
	fakeAuthenticator *commandsfakes.FakeAuthenticator
	authErr           error

	initInvocationArg bool
	initErr           error

	origInitFunc func(bool) error
)

var _ = BeforeSuite(func() {
	origInitFunc = commands.Init

	commands.Pivnet = commands.PivnetCommand{
		Format: printer.PrintAsJSON,
	}
})

var _ = BeforeEach(func() {
	fakeAuthenticator = &commandsfakes.FakeAuthenticator{}
	commands.Auth = fakeAuthenticator
	authErr = nil

	initErr = nil
})

var _ = JustBeforeEach(func() {
	commands.Init = func(arg bool) error {
		initInvocationArg = arg
		return initErr
	}

	fakeAuthenticator.AuthenticateClientReturns(authErr)
})

func fieldFor(command interface{}, name string) reflect.StructField {
	field, success := reflect.TypeOf(command).FieldByName(name)
	Expect(success).To(BeTrue(), "Expected %s field to exist on command", name)
	return field
}

func longTag(f reflect.StructField) string {
	return f.Tag.Get("long")
}

func shortTag(f reflect.StructField) string {
	return f.Tag.Get("short")
}

var alias = func(f reflect.StructField) string {
	return f.Tag.Get("alias")
}

var command = func(f reflect.StructField) string {
	return f.Tag.Get("command")
}

var isRequired = func(f reflect.StructField) bool {
	return f.Tag.Get("required") == "true"
}

var defaultVal = func(f reflect.StructField) string {
	return f.Tag.Get("default")
}
