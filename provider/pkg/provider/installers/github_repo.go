package installers

import (
	// "errors"

	"fmt"
	"os"
	"path"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

const (
	versionCommand = "git log --pretty=format:\"%H\" -1"
)

type GitHubRepo struct{}

type GitHubRepoBaseInputs struct {
	BaseInputs
	InstallCommands *[]string `pulumi:"install_commands,optional"`
	Org             *string   `pulumi:"org"`
	Repo            *string   `pulumi:"repo"`
}

var _ = (infer.CustomRead[GitHubRepoInputs, GitHubRepoOutputs])((*GitHubRepo)(nil))
var _ = (infer.CustomUpdate[GitHubRepoInputs, GitHubRepoOutputs])((*GitHubRepo)(nil))
var _ = (infer.CustomDelete[GitHubRepoOutputs])((*GitHubRepo)(nil))
var _ = (infer.CustomCheck[GitHubRepoInputs])((*GitHubRepo)(nil))

type GitHubRepoInputs struct {
	GitHubRepoBaseInputs
	FolderName *string `pulumi:"folder_name,optional"`
	Branch     *string `pulumi:"branch,optional"`
}

type GitHubRepoOutputs struct {
	CommandOutputs
	GitHubRepoInputs
	BaseOutputs
	AbsFolderName *string `pulumi:"abs_folder_name"`
}

// All resources must implement Create at a minumum.
func (l *GitHubRepo) Create(ctx p.Context, name string, input GitHubRepoInputs, preview bool) (string, GitHubRepoOutputs, error) {
	state := &GitHubRepoOutputs{
		GitHubRepoInputs: input,
	}

	var baseName string
	if state.Branch == nil {
		branch := "main"
		state.Branch = &branch
	}
	if preview {
		return name, *state, nil
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		return "", *state, err
	}
	if input.FolderName != nil {
		baseName = *input.FolderName
	} else {
		baseName = *input.Repo
	}
	absPath := path.Join(dir, baseName)
	_, err = os.Lstat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(absPath, 0777)
		} else {
			return "", *state, err
		}
	}

	command := fmt.Sprintf("git clone -b %s https://github.com/%s/%s %s", *state.Branch, *input.Org, *input.Repo, baseName)

	// clone the repo
	_, err = state.run(ctx, command, path.Dir(absPath))
	if err != nil {
		return "", *state, err
	}

	if input.InstallCommands != nil {
		_, err = state.run(ctx, strings.Join(*input.InstallCommands, " && "), absPath)
		if err != nil {
			return "", *state, err
		}
	}

	version, err := state.run(ctx, versionCommand, absPath)
	if err != nil {
		return "", *state, err
	}
	state.Version = &version
	state.FolderName = &baseName
	state.AbsFolderName = &absPath

	return name, *state, nil
}

func (l *GitHubRepo) Read(ctx p.Context, id string, inputs GitHubRepoInputs, state GitHubRepoOutputs) (
	canonicalID string, normalizedInputs GitHubRepoInputs, normalizedState GitHubRepoOutputs, err error) {

	fetch := "git fetch --all"
	versionCmd := fmt.Sprintf("git rev-parse origin/%s", *state.Branch)
	_, err = state.run(ctx, fetch, *state.AbsFolderName)
	if err != nil {
		return "", inputs, state, nil
	}
	version, err := state.run(ctx, versionCmd, *state.AbsFolderName)
	state.Version = &version
	return "", inputs, state, nil

}

func (l *GitHubRepo) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (GitHubRepoInputs, []p.CheckFailure, error) {
	return infer.DefaultCheck[GitHubRepoInputs](newInputs)
}

func (l *GitHubRepo) Update(ctx p.Context, name string, olds GitHubRepoOutputs, news GitHubRepoInputs, preview bool) (GitHubRepoOutputs, error) {
	state := &GitHubRepoOutputs{
		GitHubRepoInputs: news,
		BaseOutputs: BaseOutputs{
			Version: olds.Version,
		},
	}
	if preview {
		return *state, nil
	}
	return *state, nil
}

func (l *GitHubRepo) Delete(ctx p.Context, id string, props GitHubRepoOutputs) error {
	return nil
}
