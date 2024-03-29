package installers

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/google/go-github/v55/github"
	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type GitHubRelease struct{}

type GitHubReleaseArgs struct {
	GitHubBaseInputs
	AssetName      *string `pulumi:"assetName,optional"`
	Executable     *string `pulumi:"executable,optional"`
	ReleaseVersion *string `pulumi:"releaseVersion,optional"`
	BinLocation    *string `pulumi:"binLocation,optional"`
	BinFolder      *string `pulumi:"binFolder,optional"`
}

type GitHubReleaseState struct {
	GitHubReleaseArgs
	CommandOutputs
	DownloadURL *string   `pulumi:"downloadURL"`
	Locations   *[]string `pulumi:"locations,optional"`
}

func (l *GitHubRelease) Annotate(a infer.Annotator) {
	a.Describe(&l, "Install a program from a GitHub release")
}

func (l *GitHubReleaseArgs) Annotate(a infer.Annotator) {
	a.Describe(&l.AssetName, `The name of the release asset to install. If this is not provided then
				the resource will try and find the correct asset name to install. Supports regex`)
	a.Describe(&l.Executable, "The name of the executable to create a symlink for. If not provided then the executable name will be the same as the repo name")
	a.Describe(&l.ReleaseVersion, `The release version to install. If this is not provided then
				the resource will try and find the latest release version to install.`)
	a.Describe(&l.BinLocation, "The location to put the program. Defaults to $HOME/.local/bin")
	a.Describe(&l.BinFolder, `Sometimes release assets contain a folder containing
				program binaries which can just be copied. If that is the case, then provide the
				location here. This will copy all files in the directory to the bin_location`)
}

func (l *GitHubReleaseState) Annotate(a infer.Annotator) {
	a.Describe(&l.DownloadURL, "The URL of the GitHub release asset")
	a.Describe(&l.Locations, "The locations the program was installed to")
}

var _ = (infer.CustomUpdate[GitHubReleaseArgs, GitHubReleaseState])((*GitHubRelease)(nil))
var _ = (infer.CustomDiff[GitHubReleaseArgs, GitHubReleaseState])((*GitHubRelease)(nil))
var _ = (infer.CustomDelete[GitHubReleaseState])((*GitHubRelease)(nil))
var _ = (infer.CustomCheck[GitHubReleaseArgs])((*GitHubRelease)(nil))

func (l *GitHubRelease) Diff(ctx p.Context, id string, olds GitHubReleaseState, news GitHubReleaseArgs) (p.DiffResponse, error) {
	diff := map[string]p.PropertyDiff{}

	var newInstall string
	var oldInstall string
	if news.InstallCommands != nil {
		newInstall = strings.Join(*news.InstallCommands, " && ")
	}
	if olds.InstallCommands != nil {
		oldInstall = strings.Join(*olds.InstallCommands, " && ")
	}

	if newInstall != oldInstall {
		diff["installCommands"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
	}
	var newUninstall string
	var oldUninstall string
	if news.UninstallCommands != nil {
		newUninstall = strings.Join(*news.UninstallCommands, " && ")
	}
	if olds.UninstallCommands != nil {
		oldUninstall = strings.Join(*olds.UninstallCommands, " && ")
	}
	if newUninstall != oldUninstall {
		diff["uninstallCommands"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
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
		diff["updateCommands"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
	}

	pdiff := p.PropertyDiff{Kind: p.UpdateReplace, InputDiff: true}
	if newUpdate != "" {
		pdiff = p.PropertyDiff{Kind: p.Update, InputDiff: true}
	}

	if *news.AssetName != *news.AssetName {
		diff["assetName"] = pdiff
	}

	if news.Org != olds.Org {
		diff["org"] = p.PropertyDiff{Kind: p.UpdateReplace, InputDiff: true}
	}

	if news.Repo != olds.Repo {
		diff["repo"] = p.PropertyDiff{Kind: p.UpdateReplace, InputDiff: true}
	}

	if (news.ReleaseVersion == nil && news.ReleaseVersion != olds.ReleaseVersion) ||
		(news.ReleaseVersion != nil && *news.ReleaseVersion != *olds.ReleaseVersion) {
		diff["releaseVersion"] = pdiff
	}

	return p.DiffResponse{
		DeleteBeforeReplace: true,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// All resources must implement Create at a minumum.
func (l *GitHubRelease) Create(ctx p.Context, name string, input GitHubReleaseArgs, preview bool) (string, GitHubReleaseState, error) {
	state := &GitHubReleaseState{
		GitHubReleaseArgs: input,
	}

	if input.AssetName != nil {
		client := github.NewClient(nil)
		if val, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
			client.WithAuthToken(val)
		}
		downloadUrl, err := getReleaseDownloadURL(ctx, client, input.Org, input.Repo, *input.ReleaseVersion, *input.AssetName)
		if err != nil {
			return "", GitHubReleaseState{}, err
		}
		state.DownloadURL = &downloadUrl
	} else {
		return "", GitHubReleaseState{}, errors.New("assetName not defined, something went wrong!")
	}

	if preview {
		return name, *state, nil
	}

	commands := []string{}
	if input.InstallCommands != nil {
		commands = append(commands, *input.InstallCommands...)
	}
	if err := state.createOrUpdate(ctx, commands, &input); err != nil {
		return "", GitHubReleaseState{}, err
	}

	return name, *state, nil
}

func (o *GitHubReleaseState) createOrUpdate(ctx p.Context, commands []string, input *GitHubReleaseArgs) error {
	if o.DownloadURL == nil {
		return errors.New("Couldn't find a GitHub release to use")
	}

	ext := path.Ext(*o.AssetName)
	switch ext {
	case ".gz":
		commands = append(commands, fmt.Sprintf("tar -xzvf %s", *o.AssetName))
	case ".zip":
		commands = append(commands, fmt.Sprintf("unzip -o %s", *o.AssetName))
	}

	exName := input.Repo
	ex := false
	// if the user provided a path to the executable, then use the last part of that path as the executable name
	// e.g. if they provided /usr/local/bin/terraform, then use terraform as the executable name
	if input.Executable != nil {
		parts := strings.Split(*input.Executable, "/")
		exName = parts[len(parts)-1]
		ex = true
	}
	shellInputs := &ShellArgs{
		BaseInputs:      input.BaseInputs,
		BinLocation:     input.BinLocation,
		InstallCommands: *input.InstallCommands,
		ProgramName:     exName,
		DownloadURL:     *o.DownloadURL,
		Executable:      &ex,
	}

	shellOutputs := &ShellState{
		ShellArgs: *shellInputs,
	}
	locations := []string{}
	if input.BinFolder != nil {
		commands = append(commands, fmt.Sprintf("cp -r %s/* %s", *input.BinFolder, *input.BinLocation))
	}
	if err := shellOutputs.createOrUpdate(ctx, *shellInputs, commands); err != nil {
		return err
	}

	if input.BinFolder != nil {
		ls, err := shellOutputs.run(ctx, fmt.Sprintf("ls -l -1 %s", *input.BinFolder), os.TempDir())
		if err != nil {
			return err
		}
		for _, l := range strings.Split(ls, "\n") {
			locations = append(locations, path.Join(*input.BinLocation, l))
		}

	}
	if shellOutputs.Location != nil {
		locations = append(locations, *shellOutputs.Location)
	}
	o.Locations = &locations

	return nil
}

func getReleaseDownloadURL(ctx p.Context, client *github.Client, org, repo, tag, assetName string) (string, error) {
	release, _, err := client.Repositories.GetReleaseByTag(ctx, org, repo, tag)
	if err != nil {
		return "", err
	}
	for _, ra := range release.Assets {
		if *ra.Name == assetName {
			return *ra.BrowserDownloadURL, nil
		}

	}
	return "", nil
}

func getReleaseAssetName(
	ctx p.Context,
	client *github.Client,
	org,
	repo,
	tag,
	assetName string,
) (string, error) {
	release, _, err := client.Repositories.GetReleaseByTag(ctx, org, repo, tag)
	if err != nil {
		return "", err
	}
	// darwin/amd64
	// darwin/arm64
	// linux/amd64
	// linux/arm64
	oss := runtime.GOOS
	arch := runtime.GOARCH
	// loop over the release assets and try to find the correct one based on
	// the runtime and arch
	for _, ra := range release.Assets {
		name := strings.ToLower(*ra.Name)
		// if the user has provided their own regex
		if assetName != "" {
			if ok, err := regexp.MatchString(assetName, name); ok && err == nil {
				return *ra.Name, nil
			}
		}

		// TODO: make a better regex
		regx := fmt.Sprintf(".*%s.*(%s)?", oss, arch)
		if ok, err := regexp.MatchString(regx, name); ok {
			if err != nil {
				return "", err
			}
			return *ra.Name, nil
		}
	}
	return "", nil
}

func (l *GitHubRelease) Read(ctx p.Context, id string, inputs GitHubReleaseArgs, state GitHubReleaseState) (
	canonicalID string, normalizedInputs GitHubReleaseArgs, normalizedState GitHubReleaseState, err error) {

	if inputs.ReleaseVersion != nil {
		return id, inputs, state, nil
	}
	client := github.NewClient(nil)
	if val, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		client.WithAuthToken(val)
	}

	release, _, err := client.Repositories.GetLatestRelease(ctx, inputs.Org, inputs.Repo)
	if err != nil {
		return "", GitHubReleaseArgs{}, GitHubReleaseState{}, err
	}
	inputs.ReleaseVersion = release.TagName

	assetName, err := getReleaseAssetName(ctx, client, inputs.Org, inputs.Repo, *state.ReleaseVersion, *inputs.AssetName)
	inputs.AssetName = &assetName

	return id, inputs, state, nil
}

func (l *GitHubRelease) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (GitHubReleaseArgs, []p.CheckFailure, error) {
	client := github.NewClient(nil)
	if val, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		client.WithAuthToken(val)
	}
	failures := []p.CheckFailure{}

	// then this is a create operation
	assetName := newInputs["assetName"].StringValue()
	releaseVersion := newInputs["releaseVersion"].StringValue()
	if _, ok := oldInputs["org"]; !ok {
		_, inputs, _, err := l.Read(ctx, name, GitHubReleaseArgs{
			GitHubBaseInputs: GitHubBaseInputs{
				Org:  newInputs["org"].StringValue(),
				Repo: newInputs["repo"].StringValue(),
			},
			AssetName:      &assetName,
			ReleaseVersion: &releaseVersion,
		}, GitHubReleaseState{})
		if err != nil {
			return GitHubReleaseArgs{}, nil, err
		}
		newInputs["assetName"] = resource.NewStringProperty(*inputs.AssetName)
		newInputs["releaseVersion"] = resource.NewStringProperty(*inputs.ReleaseVersion)
	} else {
		// this is an update operation
		if _, ok := newInputs["releaseVersion"]; !ok {
			newInputs["releaseVersion"] = oldInputs["releaseVersion"]
		}
		if _, ok := newInputs["assetName"]; !ok {
			newInputs["assetName"] = oldInputs["assetName"]
		}
		if _, ok := newInputs["binLocation"]; !ok {
			newInputs["binLocation"] = oldInputs["binLocation"]
		}
	}
	inputs, fails, err := infer.DefaultCheck[GitHubReleaseArgs](newInputs)
	return inputs, append(failures, fails...), err
}

func (l *GitHubRelease) Update(ctx p.Context, name string, olds GitHubReleaseState, news GitHubReleaseArgs, preview bool) (GitHubReleaseState, error) {
	state := &GitHubReleaseState{
		GitHubReleaseArgs: news,
		DownloadURL:       olds.DownloadURL,
		Locations:         olds.Locations,
	}

	if news.AssetName == nil {
		return GitHubReleaseState{}, errors.New("assetName not defined, something went wrong! Try running a refresh")
	}
	client := github.NewClient(nil)
	if val, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		client.WithAuthToken(val)
	}
	downloadUrl, err := getReleaseDownloadURL(ctx, client, news.Org, news.Repo, *news.ReleaseVersion, *news.AssetName)
	if err != nil {
		return GitHubReleaseState{}, err
	}
	state.DownloadURL = &downloadUrl

	if preview {
		return *state, nil
	}

	var commands []string
	if news.UpdateCommands != nil {
		commands = *news.UninstallCommands
	} else if news.InstallCommands != nil {
		commands = *news.InstallCommands
	}
	if err := state.createOrUpdate(ctx, commands, &news); err != nil {
		return GitHubReleaseState{}, err
	}

	return *state, nil
}

func (l *GitHubRelease) Delete(ctx p.Context, id string, props GitHubReleaseState) error {
	if props.UninstallCommands != nil {
		_, err := props.run(ctx, strings.Join(*props.UninstallCommands, " && "), "")
		if err != nil {
			return err
		}

	}
	for _, l := range *props.Locations {
		if err := os.Remove(l); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
