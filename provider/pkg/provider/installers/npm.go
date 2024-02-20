package installers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"
	"golang.org/x/exp/slices"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type Npm struct{}

var _ = (infer.CustomUpdate[NpmArgs, NpmState])((*Npm)(nil))
var _ = (infer.CustomDiff[NpmArgs, NpmState])((*Npm)(nil))
var _ = (infer.CustomDelete[NpmState])((*Npm)(nil))
var _ = (infer.CustomCheck[NpmArgs])((*Npm)(nil))

type NpmArgs struct {
	Location *string   `pulumi:"location"`
	Packages *[]string `pulumi:"packages"`
}

type NpmState struct {
	NpmArgs
	Deps *map[string]string `pulumi:"deps"`
}

func (s *Npm) Annotate(a infer.Annotator) {
	a.Describe(&s, `
Install global npm packages.

This resource will create a local node project at the location you specify
and will then symlink the node_modules/.bin directory so that all the executables
are available globally.`)
}

func (s *NpmArgs) Annotate(a infer.Annotator) {
	a.Describe(&s.Location, "The location to create the local node project")
	a.Describe(&s.Packages, "The npm packages to install")
}

func (s *NpmState) Annotate(a infer.Annotator) {
	a.Describe(&s.Deps, "The npm packages that have been installed and their versions")
}

func (s *Npm) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (NpmArgs, []p.CheckFailure, error) {
	return infer.DefaultCheck[NpmArgs](newInputs)
}

func (s *Npm) Diff(ctx p.Context, id string, olds NpmState, news NpmArgs) (p.DiffResponse, error) {
	diff := map[string]p.PropertyDiff{}
	if *news.Location != *olds.Location {
		diff["location"] = p.PropertyDiff{Kind: p.UpdateReplace, InputDiff: true}
	}

	for _, o := range *olds.Packages {
		if !slices.Contains(*news.Packages, o) {
			diff["packages"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
		}
	}
	for _, n := range *news.Packages {
		if !slices.Contains(*olds.Packages, n) {
			diff["packages"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
		}
	}

	c := &CommandOutputs{}
	deps := *olds.Deps
	for k := range deps {
		cmd := fmt.Sprintf("npm view %s version", k)
		v, err := c.run(ctx, cmd, *olds.Location)
		if err != nil {
			return p.DiffResponse{}, err
		}
		if strings.TrimSpace(v) != deps[k] {
			ctx.Logf(diag.Debug, "%s: old version: %s, new version: %s", k, deps[k], v)
			diff["deps"] = p.PropertyDiff{Kind: p.Update, InputDiff: false}
		}
	}

	return p.DiffResponse{
		DeleteBeforeReplace: true,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

func (s *Npm) Create(ctx p.Context, name string, input NpmArgs, preview bool) (string, NpmState, error) {
	state := &NpmState{
		NpmArgs: input,
		Deps:    &map[string]string{},
	}

	if preview {
		return name, *state, nil
	}

	_, err := os.Lstat(*input.Location)
	if err != nil {
		return "", NpmState{}, err
	}

	content := map[string]interface{}{
		"name":         "npm",
		"scripts":      map[string]string{},
		"dependencies": map[string]string{},
		"main":         "lib/index.js",
		"license":      "Apache-2.0",
		"version":      "0.0.0",
	}

	c, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		return "", NpmState{}, err
	}
	if err := os.WriteFile(filepath.Join(*input.Location, "package.json"), c, 0755); err != nil {
		return "", NpmState{}, err
	}
	if err := state.install(ctx); err != nil {
		return "", NpmState{}, err
	}
	return name, *state, nil
}

func (s *Npm) Update(ctx p.Context, name string, olds NpmState, news NpmArgs, preview bool) (NpmState, error) {
	state := &NpmState{
		NpmArgs: news,
		Deps:    olds.Deps,
	}

	if preview {
		return *state, nil
	}
	deps := *olds.Deps
	for k := range *olds.Deps {
		deps[k] = "latest"
	}

	if err := state.install(ctx); err != nil {
		return NpmState{}, err
	}
	return *state, nil
}

func (s *Npm) Delete(ctx p.Context, id string, props NpmState) error {
	if err := os.Remove(filepath.Join(*props.Location, "package.json")); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.RemoveAll(filepath.Join(*props.Location, "node_modules")); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Install a npm package to a local directory, get the version of that package that
// was just installed and save it to the state.
// We need to get the version by executing the package's binary because if this is the first time
// we are running this, the location will not be available globally yet
func (n *NpmState) install(ctx p.Context) error {
	c := &CommandOutputs{}

	deps := *n.Deps
	for _, p := range *n.Packages {
		pkgVersion := "latest"
		if v, ok := deps[p]; ok {
			pkgVersion = v
		}
		if _, err := c.run(ctx, fmt.Sprintf("npm install %s@%s", p, pkgVersion), *n.Location); err != nil {
			return err
		}
		versionCmd := fmt.Sprintf(`npm ls --depth 0 --json -l | jq '.dependencies."%s".version' -r\`, p)
		version, err := c.run(ctx, versionCmd, *n.Location)
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}
		deps[p] = version
	}
	n.Deps = &deps

	return nil
}
