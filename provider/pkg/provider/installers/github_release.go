package installers

import (
	// "errors"

	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/google/go-github/v54/github"
	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type GitHubRelease struct{}

type GitHubReleaseInputs struct {
	AssetName      *string `pulumi:"assetName,optional"`
	Executable     *string `pulumi:"executable,optional"`
	ReleaseVersion *string `pulumi:"releaseVersion,optional"`
	GitHubBaseInputs
}

type GitHubReleaseOutputs struct {
	GitHubReleaseInputs
	CommandOutputs
	BaseOutputs
	DownloadURL *string `pulumi:"download_url"`
}

var _ = (infer.CustomRead[GitHubReleaseInputs, GitHubReleaseOutputs])((*GitHubRelease)(nil))
var _ = (infer.CustomUpdate[GitHubReleaseInputs, GitHubReleaseOutputs])((*GitHubRelease)(nil))
var _ = (infer.CustomDelete[GitHubReleaseOutputs])((*GitHubRelease)(nil))
var _ = (infer.CustomCheck[GitHubReleaseInputs])((*GitHubRelease)(nil))

// All resources must implement Create at a minumum.
func (l *GitHubRelease) Create(ctx p.Context, name string, input GitHubReleaseInputs, preview bool) (string, GitHubReleaseOutputs, error) {
	state := &GitHubReleaseOutputs{
		GitHubReleaseInputs: input,
	}

	if preview {
		return name, *state, nil
	}

	client := github.NewClient(nil)
	if input.ReleaseVersion != nil {
		state.Version = input.ReleaseVersion
	} else {
		if err := state.getLatestRelease(ctx, client, *input.Org, *input.Repo); err != nil {
			return "", *state, err
		}
	}
	if input.AssetName != nil {
		downloadUrl := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s", *input.Org, *input.Repo, *input.AssetName)
		state.DownloadURL = &downloadUrl
	} else {
		if err := state.getReleaseAsset(ctx, client, *input.Org, *input.Repo, *state.Version); err != nil {
			return "", *state, err
		}
	}

	if err := state.createOrUpdate(ctx, *input.InstallCommands, &input); err != nil {
		return "", *state, err
	}

	return name, *state, nil

}

func (o *GitHubReleaseOutputs) createOrUpdate(ctx p.Context, commands []string, input *GitHubReleaseInputs) error {
	if o.DownloadURL == nil {
		return errors.New("Couldn't find a GitHub release to use")
	}

	loc := path.Join(os.TempDir(), *o.AssetName)
	ext := path.Ext(*o.AssetName)
	switch ext {
	case ".gz":
		commands = append(commands, fmt.Sprintf("tar -xzvf %s", loc))
	case ".zip":
		commands = append(commands, fmt.Sprintf("unzip -o %s", loc))
	}

	exName := input.GitHubBaseInputs.Repo
	if input.Executable != nil {
		parts := strings.Split(*input.Executable, "/")
		exName = &parts[len(parts)-1]
	}
	ex := false
	if input.Executable != nil {
		ex = true
	}
	shellInputs := &ShellInputs{
		BaseInputs:      input.BaseInputs,
		InstallCommands: input.InstallCommands,
		ProgramName:     exName,
		DownloadUrl:     o.DownloadURL,
		Executable:      &ex,
	}

	o.Executable = exName
	shellOutputs := &ShellOutputs{
		ShellInputs: *shellInputs,
	}
	if err := shellOutputs.createOrUpdate(ctx, *shellInputs, *input.InstallCommands); err != nil {
		return err
	}

	return nil
}

func (o *GitHubReleaseOutputs) getReleaseAsset(ctx p.Context, client *github.Client, org, repo, tag string) error {
	release, _, err := client.Repositories.GetReleaseByTag(ctx, org, repo, tag)
	if err != nil {
		return err
	}
	o.Version = release.TagName
	// darwin/amd64
	// darwin/arm64
	// linux/amd64
	// linux/arm64
	oss := runtime.GOOS
	arch := runtime.GOARCH
	if o.AssetName == nil {
		for _, ra := range release.Assets {
			assetName := strings.ToLower(*ra.Name)
			regx := fmt.Sprintf(".*%s.*%s", oss, arch)
			if ok, err := regexp.MatchString(regx, assetName); ok {
				if err != nil {
					return err
				}
				o.DownloadURL = ra.BrowserDownloadURL
				o.AssetName = ra.Name
			}

		}
	}
	return nil
}

func (o *GitHubReleaseOutputs) getLatestRelease(ctx p.Context, client *github.Client, org, repo string) error {
	release, _, err := client.Repositories.GetLatestRelease(ctx, org, repo)
	if err != nil {
		return err
	}
	o.Version = release.TagName
	return nil
}

func (l *GitHubRelease) Read(ctx p.Context, id string, inputs GitHubReleaseInputs, state GitHubReleaseOutputs) (
	canonicalID string, normalizedInputs GitHubReleaseInputs, normalizedState GitHubReleaseOutputs, err error) {

	if inputs.ReleaseVersion != nil {
		return id, inputs, state, nil
	}

	client := github.NewClient(nil)
	if err := state.getLatestRelease(ctx, client, *inputs.Org, *inputs.Repo); err != nil {
		return "", inputs, state, nil
	}

	if err := state.getReleaseAsset(ctx, client, *inputs.Org, *inputs.Repo, *state.Version); err != nil {
		return "", inputs, state, nil
	}

	return "", inputs, state, nil
}

func (l *GitHubRelease) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (GitHubReleaseInputs, []p.CheckFailure, error) {
	return infer.DefaultCheck[GitHubReleaseInputs](newInputs)
}

func (l *GitHubRelease) Update(ctx p.Context, name string, olds GitHubReleaseOutputs, news GitHubReleaseInputs, preview bool) (GitHubReleaseOutputs, error) {
	state := &GitHubReleaseOutputs{
		GitHubReleaseInputs: news,
		BaseOutputs:         olds.BaseOutputs,
		DownloadURL:         olds.DownloadURL,
	}
	if preview {
		return *state, nil
	}

	var commands []string
	if news.UpdateCommands != nil {
		commands = *news.UninstallCommands
	} else {
		commands = *news.InstallCommands
	}
	if err := state.createOrUpdate(ctx, commands, &news); err != nil {
		return *state, err
	}

	return *state, nil
}

// TODO: implement this
func (l *GitHubRelease) Delete(ctx p.Context, id string, props GitHubReleaseOutputs) error {
	return nil
}
