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
var _ = (infer.CustomDiff[GitHubRepoInputs, GitHubRepoOutputs])((*GitHubRepo)(nil))
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
	AbsFolderName *string `pulumi:"absFolderName"`
}

func (l *GitHubRepo) Diff(ctx p.Context, id string, olds GitHubRepoOutputs, news GitHubRepoInputs) (p.DiffResponse, error) {
	diff := map[string]p.PropertyDiff{}
	if news.Branch == nil || *news.Branch != *olds.Branch {
		diff["branch"] = p.PropertyDiff{Kind: p.Update}
	}
	if news.FolderName == nil || *news.FolderName != *olds.FolderName {
		diff["folderName"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	if news.InstallCommands != olds.InstallCommands {
		diff["installCommands"] = p.PropertyDiff{Kind: p.Update}
	}
	if *news.Org != *olds.Org {
		diff["org"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}

	if *news.Repo != *olds.Repo {
		diff["repo"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}

	if news.UpdateCommands != olds.UpdateCommands {
		diff["updateCommands"] = p.PropertyDiff{Kind: p.Update}
	}

	return p.DiffResponse{
		DeleteBeforeReplace: true,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// All resources must implement Create at a minumum.
func (l *GitHubRepo) Create(ctx p.Context, name string, input GitHubRepoInputs, preview bool) (string, GitHubRepoOutputs, error) {
	state := &GitHubRepoOutputs{
		GitHubRepoInputs: input,
	}

	if err := state.getLocation(ctx, &input); err != nil {
		return "", *state, err
	}
	if preview {
		return name, *state, nil
	}

	if err := state.clone(ctx, input); err != nil {
		return "", *state, err
	}

	if input.InstallCommands != nil {
		_, err := state.run(ctx, strings.Join(*input.InstallCommands, " && "), *state.AbsFolderName)
		if err != nil {
			return "", *state, err
		}
	}

	version, err := state.run(ctx, versionCommand, *state.AbsFolderName)
	if err != nil {
		return "", *state, err
	}
	state.Version = &version

	return name, *state, nil
}

// used for import operation
func (l *GitHubRepo) Read(ctx p.Context, id string, inputs GitHubRepoInputs, state GitHubRepoOutputs) (
	canonicalID string, normalizedInputs GitHubRepoInputs, normalizedState GitHubRepoOutputs, err error) {

	versionCmd := fmt.Sprintf("git rev-parse origin/%s", *state.Branch)
	_, err = state.run(ctx, fetchCmd, *state.AbsFolderName)
	if err != nil {
		return "", inputs, state, nil
	}

	version, err := state.run(ctx, versionCmd, *state.AbsFolderName)
	state.Version = &version
	return "", inputs, state, nil

}

func (l *GitHubRepo) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (GitHubRepoInputs, []p.CheckFailure, error) {
	if _, ok := newInputs["branch"]; !ok {
		newInputs["branch"] = resource.NewStringProperty("main")
	}
	repo := newInputs["repo"].StringValue()
	if _, ok := newInputs["folderName"]; !ok {
		newInputs["folderName"] = resource.NewStringProperty(repo)
	}
	fmt.Println(newInputs["folderName"])
	return infer.DefaultCheck[GitHubRepoInputs](newInputs)
}

func (l *GitHubRepo) Update(ctx p.Context, name string, olds GitHubRepoOutputs, news GitHubRepoInputs, preview bool) (GitHubRepoOutputs, error) {
	state := &GitHubRepoOutputs{
		GitHubRepoInputs: news,
		AbsFolderName:    olds.AbsFolderName,
		BaseOutputs:      olds.BaseOutputs,
		CommandOutputs:   olds.CommandOutputs,
	}

	if preview {
		return *state, nil
	}

	if err := state.getLocation(ctx, &news); err != nil {
		return *state, err
	}

	// clone location has changed so we need to remove the old location
	if *state.AbsFolderName != *olds.AbsFolderName {
		if err := os.RemoveAll(*olds.AbsFolderName); err != nil && !os.IsNotExist(err) {
			return *state, err
		}
		if err := state.clone(ctx, news); err != nil {
			return *state, err
		}
		if news.InstallCommands != nil {
			_, err := state.run(ctx, strings.Join(*news.InstallCommands, " && "), *state.AbsFolderName)
			if err != nil {
				return *state, err
			}
		}
	} else {
		_, err := state.run(ctx, fetchCmd, *state.AbsFolderName)
		if err != nil {
			return *state, err
		}
		if *news.Branch != *olds.Branch {
			// switch branch
			_, err = state.run(ctx, fmt.Sprintf("git checkout %s", *state.Branch), *state.AbsFolderName)
			if err != nil {
				return *state, err
			}
		} else {
			// pull
			_, err = state.run(ctx, "git pull", *state.AbsFolderName)
			if err != nil {
				return *state, err
			}
		}
	}

	version, err := state.run(ctx, versionCommand, *state.AbsFolderName)
	if err != nil {
		return *state, err
	}
	state.Version = &version

	return *state, nil
}

func (l *GitHubRepo) Delete(ctx p.Context, id string, props GitHubRepoOutputs) error {
	if err := os.RemoveAll(*props.AbsFolderName); err != nil {
		return err
	}
	return nil
}

func (o *GitHubRepoOutputs) getLocation(ctx p.Context, inputs *GitHubRepoInputs) error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	absPath := path.Join(dir, *inputs.FolderName)
	_, err = os.Lstat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(absPath, 0777)
		} else {
			return err
		}
	}
	o.AbsFolderName = &absPath
	return nil
}

func (o *GitHubRepoOutputs) clone(ctx p.Context, inputs GitHubRepoInputs) error {

	command := fmt.Sprintf("git clone -b %s https://github.com/%s/%s %s", *inputs.Branch, *inputs.Org, *inputs.Repo, *inputs.FolderName)

	// clone the repo
	_, err := o.run(ctx, command, path.Dir(*o.AbsFolderName))
	if err != nil {
		return err
	}
	return nil
}
