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

func versionCommand(branch string) string {
	return fmt.Sprintf("git log origin/%s --pretty=format:\"%%H\" -1", branch)
}

const (
	fetchCmd = "git fetch --all"
)

type GitHubRepo struct{}

var _ = (infer.CustomRead[GitHubRepoArgs, GitHubRepoState])((*GitHubRepo)(nil))
var _ = (infer.CustomUpdate[GitHubRepoArgs, GitHubRepoState])((*GitHubRepo)(nil))
var _ = (infer.CustomDelete[GitHubRepoState])((*GitHubRepo)(nil))
var _ = (infer.CustomDiff[GitHubRepoArgs, GitHubRepoState])((*GitHubRepo)(nil))
var _ = (infer.CustomCheck[GitHubRepoArgs])((*GitHubRepo)(nil))

type GitHubRepoArgs struct {
	GitHubBaseInputs
	FolderName *string `pulumi:"folderName,optional"`
	Branch     *string `pulumi:"branch,optional"`
}

type GitHubRepoState struct {
	CommandOutputs
	GitHubRepoArgs
	BaseOutputs
	AbsFolderName *string `pulumi:"absFolderName"`
}

func (l *GitHubRepo) Annotate(a infer.Annotator) {
	a.Describe(&l, "Install a program from a GitHub repository")
}

func (l *GitHubRepoArgs) Annotate(a infer.Annotator) {
	a.Describe(&l.Branch, "The branch to clone from. Default to main")
	a.Describe(&l.FolderName, "The folder to clone the repo to. By default this is will be $HOME/$REPO_NAME")
}

func (l *GitHubRepoState) Annotate(a infer.Annotator) {
	a.Describe(&l.AbsFolderName, "The absolute path to the folder the repo was cloned to")
}

func (l *GitHubRepo) Diff(ctx p.Context, id string, olds GitHubRepoState, news GitHubRepoArgs) (p.DiffResponse, error) {
	diff := map[string]p.PropertyDiff{}
	if news.Branch == nil || *news.Branch != *olds.Branch {
		diff["branch"] = p.PropertyDiff{Kind: p.Update}
	}
	if news.FolderName == nil || *news.FolderName != *olds.FolderName {
		diff["folderName"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	var newInstall string
	var oldInstall string
	if news.InstallCommands != nil {
		newInstall = strings.Join(*news.InstallCommands, " && ")
	}
	if olds.InstallCommands != nil {
		oldInstall = strings.Join(*olds.InstallCommands, " && ")
	}

	if newInstall != oldInstall {
		diff["installCommands"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	if *news.Org != *olds.Org {
		diff["org"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}

	if *news.Repo != *olds.Repo {
		diff["repo"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}

	var newUpdate string
	var oldUpdate string
	if news.UpdateCommands != nil {
		newUpdate = strings.Join(*news.UpdateCommands, " && ")
	}
	if olds.UpdateCommands != nil {
		oldUpdate = strings.Join(*olds.UpdateCommands, " && ")
	}
	if newUpdate != oldUpdate {
		diff["updateCommands"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}

	return p.DiffResponse{
		DeleteBeforeReplace: true,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// All resources must implement Create at a minumum.
func (l *GitHubRepo) Create(ctx p.Context, name string, input GitHubRepoArgs, preview bool) (string, GitHubRepoState, error) {
	state := &GitHubRepoState{
		GitHubRepoArgs: input,
	}

	if err := state.getLocation(ctx, &input); err != nil {
		return "", GitHubRepoState{}, err
	}
	if preview {
		return name, GitHubRepoState{}, nil
	}

	if err := state.clone(ctx, input); err != nil {
		return "", GitHubRepoState{}, err
	}

	if input.InstallCommands != nil {
		_, err := state.run(ctx, strings.Join(*input.InstallCommands, " && "), *state.AbsFolderName)
		if err != nil {
			return "", GitHubRepoState{}, err
		}
	}

	version, err := state.run(ctx, versionCommand(*input.Branch), *state.AbsFolderName)
	if err != nil {
		return "", GitHubRepoState{}, err
	}
	state.Version = &version

	return name, *state, nil
}

// used for import operation
func (l *GitHubRepo) Read(ctx p.Context, id string, inputs GitHubRepoArgs, state GitHubRepoState) (
	canonicalID string, normalizedInputs GitHubRepoArgs, normalizedState GitHubRepoState, err error) {

	_, err = state.run(ctx, fetchCmd, *state.AbsFolderName)
	if err != nil {
		return "", GitHubRepoArgs{}, GitHubRepoState{}, err
	}

	versionCmd := fmt.Sprintf("git rev-parse origin/%s", *state.Branch)
	version, err := state.run(ctx, versionCmd, *state.AbsFolderName)
	state.Version = &version
	return "", inputs, state, nil

}

func (l *GitHubRepo) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (GitHubRepoArgs, []p.CheckFailure, error) {
	if _, ok := newInputs["branch"]; !ok {
		newInputs["branch"] = resource.NewStringProperty("main")
	}
	repo := newInputs["repo"].StringValue()
	if _, ok := newInputs["folderName"]; !ok {
		newInputs["folderName"] = resource.NewStringProperty(repo)
	}
	return infer.DefaultCheck[GitHubRepoArgs](newInputs)
}

func (l *GitHubRepo) Update(ctx p.Context, name string, olds GitHubRepoState, news GitHubRepoArgs, preview bool) (GitHubRepoState, error) {
	state := &GitHubRepoState{
		GitHubRepoArgs: news,
		AbsFolderName:  olds.AbsFolderName,
		BaseOutputs:    olds.BaseOutputs,
		CommandOutputs: olds.CommandOutputs,
	}

	if preview {
		return *state, nil
	}

	if err := state.getLocation(ctx, &news); err != nil {
		return GitHubRepoState{}, err
	}

	_, err := state.run(ctx, fetchCmd, *state.AbsFolderName)
	if err != nil {
		return GitHubRepoState{}, err
	}
	// switch branch
	_, err = state.run(ctx, fmt.Sprintf("git checkout %s", *state.Branch), *state.AbsFolderName)
	if err != nil {
		return GitHubRepoState{}, err
	}
	// pull
	_, err = state.run(ctx, "git pull", *state.AbsFolderName)
	if err != nil {
		return GitHubRepoState{}, err
	}

	version, err := state.run(ctx, versionCommand(*state.Branch), *state.AbsFolderName)
	if err != nil {
		return GitHubRepoState{}, err
	}
	state.Version = &version

	if state.UpdateCommands != nil {
		_, err := state.run(ctx, strings.Join(*state.UpdateCommands, " && "), *state.AbsFolderName)
		if err != nil {
			return GitHubRepoState{}, err
		}
	}

	return *state, nil
}

func (l *GitHubRepo) Delete(ctx p.Context, id string, props GitHubRepoState) error {
	if err := os.RemoveAll(*props.AbsFolderName); err != nil {
		return err
	}
	if props.UninstallCommands != nil {
		_, err := props.run(ctx, strings.Join(*props.UninstallCommands, " && "), "")
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *GitHubRepoState) getLocation(ctx p.Context, inputs *GitHubRepoArgs) error {
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

func (o *GitHubRepoState) clone(ctx p.Context, inputs GitHubRepoArgs) error {

	command := fmt.Sprintf(
		"git clone -b %s https://github.com/%s/%s %s",
		*inputs.Branch,
		*inputs.Org,
		*inputs.Repo,
		*inputs.FolderName,
	)

	// clone the repo
	_, err := o.run(ctx, command, path.Dir(*o.AbsFolderName))
	if err != nil {
		return err
	}
	return nil
}
