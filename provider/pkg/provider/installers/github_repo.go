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
	fetchCmd       = "git fetch --all"
)

type GitHubRepo struct{}

var _ = (infer.CustomRead[GitHubRepoInputs, GitHubRepoOutputs])((*GitHubRepo)(nil))
var _ = (infer.CustomUpdate[GitHubRepoInputs, GitHubRepoOutputs])((*GitHubRepo)(nil))
var _ = (infer.CustomDelete[GitHubRepoOutputs])((*GitHubRepo)(nil))
var _ = (infer.CustomCheck[GitHubRepoInputs])((*GitHubRepo)(nil))

type GitHubRepoInputs struct {
	GitHubBaseInputs
	FolderName *string `pulumi:"folderName,optional"`
	Branch     *string `pulumi:"branch,optional"`
}

type GitHubRepoOutputs struct {
	CommandOutputs
	GitHubRepoInputs
	BaseOutputs
	OutputBranch  *string `pulumi:"outputBranch"`
	AbsFolderName *string `pulumi:"absFolderName"`
}

// All resources must implement Create at a minumum.
func (l *GitHubRepo) Create(ctx p.Context, name string, input GitHubRepoInputs, preview bool) (string, GitHubRepoOutputs, error) {
	state := &GitHubRepoOutputs{
		GitHubRepoInputs: input,
	}

	if input.Branch == nil {
		branch := "main"
		state.OutputBranch = &branch
	}

	absPath, err := state.getLocation(ctx, &input)
	if err != nil {
		return "", *state, err
	}
	if preview {
		return name, *state, nil
	}

	if err = state.clone(ctx, absPath, input); err != nil {
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
	state.AbsFolderName = &absPath

	return name, *state, nil
}

func (o *GitHubRepoOutputs) getLocation(ctx p.Context, inputs *GitHubRepoInputs) (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if inputs.FolderName == nil {
		inputs.FolderName = inputs.Repo
	}
	absPath := path.Join(dir, *inputs.FolderName)
	_, err = os.Lstat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(absPath, 0777)
		} else {
			return "", err
		}
	}
	return absPath, nil
}

func (o *GitHubRepoOutputs) clone(ctx p.Context, absPath string, inputs GitHubRepoInputs) error {

	command := fmt.Sprintf("git clone -b %s https://github.com/%s/%s %s", *o.Branch, *inputs.Org, *inputs.Repo, *inputs.FolderName)

	// clone the repo
	_, err := o.run(ctx, command, path.Dir(absPath))
	if err != nil {
		return err
	}
	return nil
}

func (l *GitHubRepo) Read(ctx p.Context, id string, inputs GitHubRepoInputs, state GitHubRepoOutputs) (
	canonicalID string, normalizedInputs GitHubRepoInputs, normalizedState GitHubRepoOutputs, err error) {

	versionCmd := fmt.Sprintf("git rev-parse origin/%s", *state.Branch)
	_, err = state.run(ctx, fetchCmd, *state.AbsFolderName)
	if err != nil {
		return "", inputs, state, nil
	}
	if inputs.Branch == nil {
		branch := "main"
		state.OutputBranch = &branch
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
		BaseOutputs:      olds.BaseOutputs,
	}
	if news.Branch == nil {
		branch := "main"
		state.OutputBranch = &branch
	} else {
		state.OutputBranch = news.Branch
	}
	if preview {
		return *state, nil
	}

	absPath, err := state.getLocation(ctx, &news)
	if err != nil {
		return *state, err
	}

	// clone location has changed so we need to remove the old location
	if absPath != *olds.AbsFolderName {
		if err := os.RemoveAll(*olds.AbsFolderName); err != nil && !os.IsNotExist(err) {
			return *state, err
		}
		if err = state.clone(ctx, absPath, news); err != nil {
			return *state, err
		}
		if news.InstallCommands != nil {
			_, err = state.run(ctx, strings.Join(*news.InstallCommands, " && "), absPath)
			if err != nil {
				return *state, err
			}
		}
	} else {
		_, err = state.run(ctx, fetchCmd, absPath)
		if err != nil {
			return *state, err
		}
		if *news.Branch != *olds.Branch {
			// switch branch
			_, err = state.run(ctx, fmt.Sprintf("git checkout %s", *news.Branch), absPath)
			if err != nil {
				return *state, err
			}
		} else {
			// pull
			_, err = state.run(ctx, "git pull", absPath)
			if err != nil {
				return *state, err
			}
		}
	}

	version, err := state.run(ctx, versionCommand, absPath)
	if err != nil {
		return *state, err
	}
	state.Version = &version
	state.AbsFolderName = &absPath

	return *state, nil
}

func (l *GitHubRepo) Delete(ctx p.Context, id string, props GitHubRepoOutputs) error {
	return nil
}
