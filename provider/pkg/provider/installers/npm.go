package installers

import (
	"fmt"
	"os"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type Npm struct{}

var _ = (infer.CustomUpdate[NpmArgs, NpmState])((*Npm)(nil))
var _ = (infer.CustomDiff[NpmArgs, NpmState])((*Npm)(nil))
var _ = (infer.CustomDelete[NpmState])((*Npm)(nil))
var _ = (infer.CustomCheck[NpmArgs])((*Npm)(nil))

type NpmArgs struct {
	Location string  `pulumi:"location"`
	Package  string  `pulumi:"package"`
	Version  *string `pulumi:"version,optional"`
}

type NpmState struct {
	NpmArgs
}

func (s *Npm) Annotate(a infer.Annotator) {
	a.Describe(&s, `
Install global npm packages.

This resource will create a local node project at the location you specify
and will then symlink the node_modules/.bin directory so that all the executables
are available globally.`)
}

func (s *NpmArgs) Annotate(a infer.Annotator) {
	a.Describe(&s.Location, "The location of the node project")
	a.Describe(&s.Package, "The npm package to install")
	a.Describe(&s.Version, "The version of the package to install")
}

func (s *NpmState) Annotate(a infer.Annotator) {}

func (s *Npm) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (NpmArgs, []p.CheckFailure, error) {
	if _, ok := newInputs["version"]; !ok {
		// if package is not in oldInputs, then this is a create operation and the read method is not
		// called
		if _, ok := oldInputs["package"]; !ok {
			_, inputs, _, err := s.Read(ctx, name, NpmArgs{
				Location: newInputs["location"].StringValue(),
				Package:  newInputs["package"].StringValue(),
			}, NpmState{})
			if err != nil {
				return NpmArgs{}, nil, err
			}
			newInputs["version"] = resource.NewStringProperty(*inputs.Version)
		} else {
			newInputs["version"] = oldInputs["version"]
		}
	}
	return infer.DefaultCheck[NpmArgs](newInputs)
}

// Read is only called during --refresh operations
func (s *Npm) Read(ctx p.Context, id string, inputs NpmArgs, state NpmState) (
	canonicalID string, normalizedInputs NpmArgs, normalizedState NpmState, err error,
) {
	c := &CommandOutputs{}
	if inputs.Version == nil {
		cmd := fmt.Sprintf("npm view %s version", inputs.Package)
		v, err := c.run(ctx, cmd, state.Location)
		if err != nil {
			return "", NpmArgs{}, NpmState{}, err
		}
		inputs.Version = &v
	}
	return id, inputs, state, nil
}

func (s *Npm) Diff(ctx p.Context, id string, olds NpmState, news NpmArgs) (p.DiffResponse, error) {
	diff := map[string]p.PropertyDiff{}
	if news.Location != olds.Location {
		diff["location"] = p.PropertyDiff{Kind: p.UpdateReplace, InputDiff: true}
	}

	if news.Package != olds.Package {
		diff["packages"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
	}

	if *news.Version != *olds.Version {
		diff["version"] = p.PropertyDiff{Kind: p.Update, InputDiff: true}
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
	}

	if preview {
		return name, *state, nil
	}

	_, err := os.Lstat(input.Location)
	if err != nil {
		return "", NpmState{}, fmt.Errorf("location %s does not exist", input.Location)
	}

	if err := state.install(ctx); err != nil {
		return "", NpmState{}, err
	}
	return name, *state, nil
}

func (s *Npm) Update(ctx p.Context, name string, olds NpmState, news NpmArgs, preview bool) (NpmState, error) {
	state := &NpmState{
		NpmArgs: news,
	}

	if preview {
		return *state, nil
	}

	if err := state.install(ctx); err != nil {
		return NpmState{}, err
	}
	return *state, nil
}

func (s *Npm) Delete(ctx p.Context, id string, props NpmState) error {
	c := &CommandOutputs{}

	if _, err := c.run(ctx, fmt.Sprintf("npm uninstall %s", props.Package), props.Location); err != nil {
		return err
	}
	return nil
}

// Install a npm package to a local directory
func (n *NpmState) install(ctx p.Context) error {
	c := &CommandOutputs{}

	if _, err := c.run(ctx, fmt.Sprintf("npm install %s@%s", n.Package, *n.Version), n.Location); err != nil {
		return err
	}

	return nil
}
