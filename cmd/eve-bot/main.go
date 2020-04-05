package main

import (
	"time"

	"gitlab.unanet.io/devops/eve-bot/internal/api"
	"gitlab.unanet.io/devops/eve-bot/internal/servicefactory"
	"gitlab.unanet.io/devops/eve-bot/internal/version"
	"go.uber.org/zap"
)

// Public/Global Variables Passed in dynamically during Build time
// used to add build metadata into the binary
var (
	// GitCommit is the Full Git Commit SHA
	GitCommit string
	// GitCommitAuthor is the author of the Git Commit
	GitCommitAuthor string
	// GitBranch is the Full Git Branch Name
	GitBranch string
	// BuildDate is the DateTimeStamp during build
	BuildDate string
	// GitDescribe is a way to intentionally describe the version
	GitDescribe string
	// Version is the Full Semantic Version
	Version string
	// VersionPrerelease is the pre-release name (dev,rc-1,alpha,beta,nightly,etc.)
	VersionPrerelease string
	// VersionMetaData is the optional metadata to attach to a version
	VersionMetaData string
	// Builder is the name of the user that builds the artifact (i.e whoami)
	Builder string
	// BuildHost is the name of the host that builds the artifact
	BuildHost string
)

func main() {
	svcFactory := servicefactory.Initialize(
		version.New(
			Version,
			GitCommit,
			GitBranch,
			GitCommitAuthor,
			GitDescribe,
			BuildDate,
			VersionPrerelease,
			VersionMetaData,
			Builder,
			BuildHost,
			time.Now(),
		),
	)

	// Serve up the main server API
	if err := api.New(svcFactory).Serve(); err != nil {
		svcFactory.Logger.Bg().Fatal("api serve failed", zap.Error(err))
	}
}
