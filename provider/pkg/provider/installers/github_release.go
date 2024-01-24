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

type GitHubReleaseInputs struct {
	GitHubBaseInputs
	AssetName      *string `pulumi:"assetName,optional"`
	Executable     *string `pulumi:"executable,optional"`
	ReleaseVersion *string `pulumi:"releaseVersion,optional"`
	BinLocation    *string `pulumi:"binLocation,optional"`
	BinFolder      *string `pulumi:"binFolder,optional"`
}

type GitHubReleaseOutputs struct {
	GitHubReleaseInputs
	CommandOutputs
	DownloadURL *string   `pulumi:"downloadURL"`
	Locations   *[]string `pulumi:"locations,optional"`
}

var _ = (infer.CustomRead[GitHubReleaseInputs, GitHubReleaseOutputs])((*GitHubRelease)(nil))
var _ = (infer.CustomUpdate[GitHubReleaseInputs, GitHubReleaseOutputs])((*GitHubRelease)(nil))
var _ = (infer.CustomDiff[GitHubReleaseInputs, GitHubReleaseOutputs])((*GitHubRelease)(nil))
var _ = (infer.CustomDelete[GitHubReleaseOutputs])((*GitHubRelease)(nil))
var _ = (infer.CustomCheck[GitHubReleaseInputs])((*GitHubRelease)(nil))

func (l *GitHubRelease) Diff(ctx p.Context, id string, olds GitHubReleaseOutputs, news GitHubReleaseInputs) (p.DiffResponse, error) {
	diff := map[string]p.PropertyDiff{}

	if *news.AssetName != *news.AssetName {
		diff["assetName"] = p.PropertyDiff{Kind: p.UpdateReplace}
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
	var newUninstall string
	var oldUninstall string
	if news.UninstallCommands != nil {
		newUninstall = strings.Join(*news.UninstallCommands, " && ")
	}
	if olds.UninstallCommands != nil {
		oldUninstall = strings.Join(*olds.UninstallCommands, " && ")
	}
	if newUninstall != oldUninstall {
		diff["uninstallCommands"] = p.PropertyDiff{Kind: p.UpdateReplace}
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

	if *news.Org != *olds.Org {
		diff["org"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}

	if *news.Repo != *olds.Repo {
		diff["repo"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}

	if (news.ReleaseVersion == nil && news.ReleaseVersion != olds.ReleaseVersion) ||
		(news.ReleaseVersion != nil && *news.ReleaseVersion != *olds.ReleaseVersion) {
		diff["releaseVersion"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}

	return p.DiffResponse{
		DeleteBeforeReplace: true,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// All resources must implement Create at a minumum.
func (l *GitHubRelease) Create(ctx p.Context, name string, input GitHubReleaseInputs, preview bool) (string, GitHubReleaseOutputs, error) {
	state := &GitHubReleaseOutputs{
		GitHubReleaseInputs: input,
	}

	if input.AssetName != nil {
		client := github.NewClient(nil)
		if val, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
			client.WithAuthToken(val)
		}
		downloadUrl, err := getReleaseDownloadURL(ctx, client, *input.Org, *input.Repo, *input.ReleaseVersion, *input.AssetName)
		if err != nil {
			return "", GitHubReleaseOutputs{}, err
		}
		state.DownloadURL = &downloadUrl
	} else {
		return "", GitHubReleaseOutputs{}, errors.New("assetName not defined, something went wrong!")
	}

	if preview {
		return name, *state, nil
	}

	commands := []string{}
	if input.InstallCommands != nil {
		commands = append(commands, *input.InstallCommands...)
	}
	if err := state.createOrUpdate(ctx, commands, &input); err != nil {
		return "", GitHubReleaseOutputs{}, err
	}

	return name, *state, nil
}

func (o *GitHubReleaseOutputs) createOrUpdate(ctx p.Context, commands []string, input *GitHubReleaseInputs) error {
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
		exName = &parts[len(parts)-1]
		ex = true
	}
	shellInputs := &ShellInputs{
		BaseInputs:      input.BaseInputs,
		BinLocation:     input.BinLocation,
		InstallCommands: input.InstallCommands,
		ProgramName:     exName,
		DownloadURL:     o.DownloadURL,
		Executable:      &ex,
	}

	shellOutputs := &ShellOutputs{
		ShellInputs: *shellInputs,
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
		if ok, err := regexp.MatchString(assetName, name); ok && err == nil {
			return *ra.Name, nil
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

func (l *GitHubRelease) Read(ctx p.Context, id string, inputs GitHubReleaseInputs, state GitHubReleaseOutputs) (
	canonicalID string, normalizedInputs GitHubReleaseInputs, normalizedState GitHubReleaseOutputs, err error) {

	if inputs.ReleaseVersion != nil {
		return id, inputs, state, nil
	}
	client := github.NewClient(nil)
	if val, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		client.WithAuthToken(val)
	}

	release, _, err := client.Repositories.GetLatestRelease(ctx, *inputs.Org, *inputs.Repo)
	if err != nil {
		return "", GitHubReleaseInputs{}, GitHubReleaseOutputs{}, err
	}
	state.ReleaseVersion = release.TagName

	return id, inputs, state, nil
}

func (l *GitHubRelease) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (GitHubReleaseInputs, []p.CheckFailure, error) {
	client := github.NewClient(nil)
	if val, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		client.WithAuthToken(val)
	}
	org := newInputs["org"].StringValue()
	repo := newInputs["repo"].StringValue()
	failures := []p.CheckFailure{}
	var version string
	if _, ok := newInputs["releaseVersion"]; !ok {
		release, _, err := client.Repositories.GetLatestRelease(ctx, org, repo)
		if err != nil {
			failures = append(failures, p.CheckFailure{Property: "releaseVersion", Reason: err.Error()})
			return GitHubReleaseInputs{}, failures, nil
		}
		newInputs["releaseVersion"] = resource.NewStringProperty(*release.TagName)
		version = *release.TagName
	} else {
		version = newInputs["releaseVersion"].StringValue()
	}
	var inputAsssetName string
	if val, ok := newInputs["assetName"]; ok {
		inputAsssetName = val.StringValue()
	}
	assetName, err := getReleaseAssetName(ctx, client, org, repo, version, inputAsssetName)
	if err != nil {
		failures = append(failures, p.CheckFailure{Property: "releaseVersion", Reason: err.Error()})
		return GitHubReleaseInputs{}, failures, nil
	}
	newInputs["assetName"] = resource.NewStringProperty(assetName)

	if _, ok := newInputs["binLocation"]; !ok {
		home, err := os.UserHomeDir()
		if err != nil {
			failures = append(failures, p.CheckFailure{Property: "binLocation", Reason: err.Error()})
			return GitHubReleaseInputs{}, failures, nil
		}
		binLocation := path.Join(home, ".local", "bin")
		newInputs["binLocation"] = resource.NewStringProperty(binLocation)
	}

	inputs, fails, err := infer.DefaultCheck[GitHubReleaseInputs](newInputs)
	return inputs, append(failures, fails...), err
}

func (l *GitHubRelease) Update(ctx p.Context, name string, olds GitHubReleaseOutputs, news GitHubReleaseInputs, preview bool) (GitHubReleaseOutputs, error) {
	state := &GitHubReleaseOutputs{
		GitHubReleaseInputs: news,
		DownloadURL:         olds.DownloadURL,
		Locations:           olds.Locations,
	}
	if news.AssetName != nil {
		client := github.NewClient(nil)
		if val, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
			client.WithAuthToken(val)
		}
		downloadUrl, err := getReleaseDownloadURL(ctx, client, *news.Org, *news.Repo, *news.ReleaseVersion, *news.AssetName)
		if err != nil {
			return GitHubReleaseOutputs{}, err
		}
		state.DownloadURL = &downloadUrl
	} else {
		return GitHubReleaseOutputs{}, errors.New("assetName not defined, something went wrong!")
	}

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
		return GitHubReleaseOutputs{}, err
	}

	return *state, nil
}

func (l *GitHubRelease) Delete(ctx p.Context, id string, props GitHubReleaseOutputs) error {
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
