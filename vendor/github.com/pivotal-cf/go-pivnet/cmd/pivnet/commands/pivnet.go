package commands

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/errorhandler"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/gp"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/logshim"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/printer"
	"github.com/pivotal-cf/go-pivnet/cmd/pivnet/version"
	"github.com/pivotal-cf/go-pivnet/logger"
	"github.com/robdimsdale/sanitizer"
)

const (
	DefaultHost = "https://network.pivotal.io"
)

var (
	OutputWriter io.Writer
	LogWriter    io.Writer

	ErrorHandler errorhandler.ErrorHandler
	Printer      printer.Printer
)

type PivnetCommand struct {
	VersionFunc func() `short:"v" long:"version" description:"Print the version of this CLI and exit"`

	Format  string `long:"format" description:"Format to print as" default:"table" choice:"table" choice:"json" choice:"yaml"`
	Verbose bool   `long:"verbose" description:"Display verbose output"`

	APIToken string `long:"api-token" description:"Pivnet API token"`
	Host     string `long:"host" description:"Pivnet API Host"`

	Help    HelpCommand    `command:"help" alias:"h" description:"Print this help message"`
	Version VersionCommand `command:"version" alias:"v" description:"Print the version of this CLI and exit"`

	Curl CurlCommand `command:"curl" alias:"c" description:"Curl an endpoint"`

	ReleaseTypes ReleaseTypesCommand `command:"release-types" alias:"rts" description:"List release types"`

	EULAs      EULAsCommand      `command:"eulas" alias:"es" description:"List EULAs"`
	EULA       EULACommand       `command:"eula" alias:"e" description:"Show EULA"`
	AcceptEULA AcceptEULACommand `command:"accept-eula" alias:"ae" description:"Accept EULA"`

	Products ProductsCommand `command:"products" alias:"ps" description:"List products"`
	Product  ProductCommand  `command:"product" alias:"p" description:"Show product"`

	ProductFiles      ProductFilesCommand      `command:"product-files" alias:"pfs" description:"List product files"`
	ProductFile       ProductFileCommand       `command:"product-file" alias:"pf" description:"Show product file"`
	CreateProductFile CreateProductFileCommand `command:"create-product-file" alias:"cpf" description:"Create product file"`
	UpdateProductFile UpdateProductFileCommand `command:"update-product-file" alias:"upf" description:"Update product file"`
	AddProductFile    AddProductFileCommand    `command:"add-product-file" alias:"apf" description:"Add product file to release"`
	RemoveProductFile RemoveProductFileCommand `command:"remove-product-file" alias:"rpf" description:"Remove product file from release"`
	DeleteProductFile DeleteProductFileCommand `command:"delete-product-file" alias:"dpf" description:"Delete product file"`

	DownloadProductFile DownloadProductFileCommand `command:"download-product-file" alias:"dlpf" description:"Download product file"`

	FileGroups                 FileGroupsCommand                 `command:"file-groups" alias:"fgs" description:"List file groups"`
	FileGroup                  FileGroupCommand                  `command:"file-group" alias:"fg" description:"Show file group"`
	CreateFileGroup            CreateFileGroupCommand            `command:"create-file-group" alias:"cfg" description:"Create file group"`
	UpdateFileGroup            UpdateFileGroupCommand            `command:"update-file-group" alias:"ufg" description:"Update file group"`
	DeleteFileGroup            DeleteFileGroupCommand            `command:"delete-file-group" alias:"dfg" description:"Delete file group"`
	AddFileGroupToRelease      AddFileGroupToReleaseCommand      `command:"add-file-group" alias:"afg" description:"Add file group to release"`
	RemoveFileGroupFromRelease RemoveFileGroupFromReleaseCommand `command:"remove-file-group" alias:"rfg" description:"Remove file group from release"`

	Releases      ReleasesCommand      `command:"releases" alias:"rs" description:"List releases"`
	Release       ReleaseCommand       `command:"release" alias:"r" description:"Show release"`
	CreateRelease CreateReleaseCommand `command:"create-release" alias:"cr" description:"Create release"`
	DeleteRelease DeleteReleaseCommand `command:"delete-release" alias:"dr" description:"Delete release"`

	UserGroups      UserGroupsCommand      `command:"user-groups" alias:"ugs" description:"List user groups"`
	UserGroup       UserGroupCommand       `command:"user-group" alias:"ug" description:"Show user group"`
	AddUserGroup    AddUserGroupCommand    `command:"add-user-group" alias:"aug" description:"Add user group to release"`
	RemoveUserGroup RemoveUserGroupCommand `command:"remove-user-group" alias:"rug" description:"Remove user group from release"`
	CreateUserGroup CreateUserGroupCommand `command:"create-user-group" alias:"cug" description:"Create user group"`
	UpdateUserGroup UpdateUserGroupCommand `command:"update-user-group" alias:"uug" description:"Update user group"`
	DeleteUserGroup DeleteUserGroupCommand `command:"delete-user-group" alias:"dug" description:"Delete user group"`

	AddUserGroupMember    AddUserGroupMemberCommand    `command:"add-user-group-member" alias:"augm" description:"Add user group member to group"`
	RemoveUserGroupMember RemoveUserGroupMemberCommand `command:"remove-user-group-member" alias:"rugm" description:"Remove user group member from group"`

	ReleaseDependencies     ReleaseDependenciesCommand     `command:"release-dependencies" alias:"rds" description:"List release dependencies"`
	AddReleaseDependency    AddReleaseDependencyCommand    `command:"add-release-dependency" alias:"ard" description:"Add release dependency"`
	RemoveReleaseDependency RemoveReleaseDependencyCommand `command:"remove-release-dependency" alias:"rrd" description:"Remove release dependency"`

	ReleaseUpgradePaths      ReleaseUpgradePathsCommand      `command:"release-upgrade-paths" alias:"rups" description:"List release upgrade paths"`
	AddReleaseUpgradePath    AddReleaseUpgradePathCommand    `command:"add-release-upgrade-path" alias:"arup" description:"Add release upgrade path"`
	RemoveReleaseUpgradePath RemoveReleaseUpgradePathCommand `command:"remove-release-upgrade-path" alias:"rrup" description:"Remove release upgrade path"`

	Logger    logger.Logger
	userAgent string
}

var Pivnet PivnetCommand

func init() {
	Pivnet.VersionFunc = func() {
		fmt.Println(version.Version)
		os.Exit(0)
	}

	if Pivnet.Host == "" {
		Pivnet.Host = DefaultHost
	}
}

func NewPivnetClient() *gp.CompositeClient {
	return gp.NewCompositeClient(
		pivnet.ClientConfig{
			Token:     Pivnet.APIToken,
			Host:      Pivnet.Host,
			UserAgent: Pivnet.userAgent,
		},
		Pivnet.Logger,
	)
}

func Init() {
	if OutputWriter == nil {
		OutputWriter = os.Stdout
	}

	if LogWriter == nil {
		switch Pivnet.Format {
		case printer.PrintAsJSON, printer.PrintAsYAML:
			LogWriter = os.Stderr
			break
		default:
			LogWriter = os.Stdout
		}
	}

	if ErrorHandler == nil {
		ErrorHandler = errorhandler.NewErrorHandler(Pivnet.Format, OutputWriter, LogWriter)
	}

	if Printer == nil {
		Printer = printer.NewPrinter(OutputWriter)
	}

	sanitized := map[string]string{
		Pivnet.APIToken: "*** redacted api token ***",
	}

	OutputWriter = sanitizer.NewSanitizer(sanitized, OutputWriter)
	LogWriter = sanitizer.NewSanitizer(sanitized, LogWriter)

	infoLogger := log.New(LogWriter, "", log.LstdFlags)
	debugLogger := log.New(LogWriter, "", log.LstdFlags)

	Pivnet.userAgent = fmt.Sprintf(
		"go-pivnet/%s",
		version.Version,
	)

	Pivnet.Logger = logshim.NewLogShim(infoLogger, debugLogger, Pivnet.Verbose)
}
