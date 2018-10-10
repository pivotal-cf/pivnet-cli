package commands_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pivnet-cli/commands"
	"github.com/pivotal-cf/pivnet-cli/commands/commandsfakes"
	"github.com/pivotal-cf/pivnet-cli/rc"
)

var _ = Describe("who ami command", func() {
	var (
		cmd           commands.WhoAmICommand
		profileName   string
		host          string
		fakeRCHandler *commandsfakes.FakeRCHandler
		err           error
		outBuffer     bytes.Buffer
	)

	BeforeEach(func() {
		profileName = "default"
		host = "test.host.com"

		commands.Pivnet.ProfileName = profileName

		fakeRCHandler = &commandsfakes.FakeRCHandler{}
		commands.RC = fakeRCHandler

		outBuffer = bytes.Buffer{}
		commands.OutputWriter = &outBuffer
	})

	Context("with good params", func() {
		BeforeEach(func() {
			profile := &rc.PivnetProfile{
				Name:     "some-profile",
				APIToken: "token",
				Host:     host,
			}
			fakeRCHandler.ProfileForNameReturns(profile, nil)
		})

		It("invokes the Init function with 'false'", func() {
			err = cmd.Execute(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(initInvocationArg).To(BeFalse())
		})

		It("prints profileName", func() {
			err = cmd.Execute(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(outBuffer.String()).Should(ContainSubstring(profileName))
		})

		It("prints host", func() {
			err = cmd.Execute(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(outBuffer.String()).Should(ContainSubstring(host))
		})
	})

	Context("when retrieving the profile returns an error", func() {
		BeforeEach(func() {
			fakeRCHandler.ProfileForNameReturns(nil, fmt.Errorf(""))
		})

		It("returns an error", func() {
			err = cmd.Execute(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(outBuffer.String()).Should(ContainSubstring("No user is logged in"))
		})
	})

	Context("when retrieving the profile doesn't exist", func() {
		BeforeEach(func() {
			fakeRCHandler.ProfileForNameReturns(nil, nil)
		})

		It("returns an error", func() {
			err = cmd.Execute(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(outBuffer.String()).Should(ContainSubstring("No user is logged in"))
		})
	})
})
