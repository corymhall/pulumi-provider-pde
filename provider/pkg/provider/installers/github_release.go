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
	AssetName   *string `pulumi:"assetName,optional"`
	Executable  *string `pulumi:"executable,optional"`
	Version     *string `pulumi:"version,optional"`
	BinLocation *string `pulumi:"binLocation,optional"`
	BinFolder   *string `pulumi:"binFolder,optional"`
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

	if news.InstallCommands != olds.InstallCommands {
		diff["installCommands"] = p.PropertyDiff{Kind: p.UpdateReplace}
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

	if (news.Version == nil && news.Version != olds.Version) ||
		(news.Version != nil && *news.Version != *olds.Version) {
		diff["version"] = p.PropertyDiff{Kind: p.UpdateReplace}
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
		downloadUrl, err := getReleaseDownloadURL(ctx, github.NewClient(nil), *input.Org, *input.Repo, *input.Version, *input.AssetName)
		if err != nil {
			return "", *state, err
		}
		state.DownloadURL = &downloadUrl
	} else {
		return "", *state, errors.New("assetName not defined, something went wrong!")
	}

	if preview {
		return name, *state, nil
	}

	commands := []string{}
	if input.InstallCommands != nil {
		commands = append(commands, *input.InstallCommands...)
	}
	if err := state.createOrUpdate(ctx, commands, &input); err != nil {
		return "", *state, err
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
	if input.Executable != nil {
		parts := strings.Split(*input.Executable, "/")
		exName = &parts[len(parts)-1]
		o.Executable = exName
	}
	ex := false
	if input.Executable != nil {
		ex = true
	}
	shellInputs := &ShellInputs{
		BaseInputs:      input.BaseInputs,
		BinLocation:     input.BinLocation,
		InstallCommands: input.InstallCommands,
		ProgramName:     exName,
		DownloadUrl:     o.DownloadURL,
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

func getReleaseAssetName(ctx p.Context, client *github.Client, org, repo, tag string) (string, error) {
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
	for _, ra := range release.Assets {
		assetName := strings.ToLower(*ra.Name)
		regx := fmt.Sprintf(".*%s.*%s", oss, arch)
		if ok, err := regexp.MatchString(regx, assetName); ok {
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

	if inputs.Version != nil {
		return id, inputs, state, nil
	}

	// TODO: implement import

	return "", inputs, state, nil
}

func (l *GitHubRelease) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (GitHubReleaseInputs, []p.CheckFailure, error) {
	client := github.NewClient(nil)
	if val, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		client.WithAuthToken(val)
	}
	org := newInputs["org"].StringValue()
	repo := newInputs["repo"].StringValue()
	failures := []p.CheckFailure{}
	if _, ok := newInputs["version"]; !ok {
		release, _, err := client.Repositories.GetLatestRelease(ctx, org, repo)
		if err != nil {
			failures = append(failures, p.CheckFailure{Property: "version", Reason: err.Error()})
		} else {
			newInputs["version"] = resource.NewStringProperty(*release.TagName)
			assetName, err := getReleaseAssetName(ctx, client, org, repo, *release.TagName)
			if err != nil {
				failures = append(failures, p.CheckFailure{Property: "version", Reason: err.Error()})
			} else {
				if _, ok := newInputs["assetName"]; !ok {
					newInputs["assetName"] = resource.NewStringProperty(assetName)
				} else if assetName != newInputs["assetName"].StringValue() {
					failures = append(failures, p.CheckFailure{
						Property: "assetName",
						Reason:   fmt.Sprintf("provided asset name %s not available for version %s", newInputs["assetName"].StringValue(), *release.TagName),
					})

				}
			}
		}
	} else {
		version := newInputs["version"].StringValue()
		assetName, err := getReleaseAssetName(ctx, client, org, repo, version)
		if err != nil {
			failures = append(failures, p.CheckFailure{
				Property: "version", Reason: err.Error(),
			})
		}
		if _, ok := newInputs["assetName"]; !ok {
			newInputs["assetName"] = resource.NewStringProperty(assetName)
		} else if assetName != newInputs["assetName"].StringValue() {
			failures = append(failures, p.CheckFailure{
				Property: "assetName",
				Reason:   fmt.Sprintf("provided asset name %s not available for version %s", newInputs["assetName"].StringValue(), version),
			})

		}
	}

	if _, ok := newInputs["binLocation"]; !ok {
		home, err := os.UserHomeDir()
		if err != nil {
			failures = append(failures, p.CheckFailure{Property: "binLocation", Reason: err.Error()})
		} else {
			binLocation := path.Join(home, ".local", "bin")
			newInputs["binLocation"] = resource.NewStringProperty(binLocation)
		}
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
		downloadUrl, err := getReleaseDownloadURL(ctx, github.NewClient(nil), *news.Org, *news.Repo, *news.Version, *news.AssetName)
		if err != nil {
			return *state, err
		}
		state.DownloadURL = &downloadUrl
	} else {
		return *state, errors.New("assetName not defined, something went wrong!")
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
		return *state, err
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
