package local

import (
	// "errors"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// Each resource has a controlling struct.
// Resource behavior is determined by implementing methods on the controlling struct.
// The `Create` method is mandatory, but other methods are optional.
// - Check: Remap inputs before they are typed.
// - Diff: Change how instances of a resource are compared.
// - Update: Mutate a resource in place.
// - Read: Get the state of a resource from the backing provider.
// - Delete: Custom logic when the resource is deleted.
// - Annotate: Describe fields and set defaults for a resource.
// - WireDependencies: Control how outputs and secrets flows through values.
type Link struct{}

var _ = (infer.CustomRead[LinkArgs, LinkState])((*Link)(nil))
var _ = (infer.CustomUpdate[LinkArgs, LinkState])((*Link)(nil))
var _ = (infer.CustomDelete[LinkState])((*Link)(nil))
var _ = (infer.CustomCheck[LinkArgs])((*Link)(nil))

// Each resource has in input struct, defining what arguments it accepts.
type LinkArgs struct {
	Source    *string `pulumi:"source"`
	Target    *string `pulumi:"target"`
	Overwrite *bool   `pulumi:"overwrite,optional"`
	Retain    *bool   `pulumi:"retain,optional"`
	Recursive *bool   `pulumi:"recursive,optional"`
	// Fields projected into Pulumi must be public and have a `pulumi:"..."` tag.
	// The pulumi tag doesn't need to match the field name, but its generally a
	// good idea.
	// Length int `pulumi:"length"`
}

// Each resource has a state, describing the fields that exist on the created resource.
type LinkState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	LinkArgs
	Linked  *bool     `pulumi:"linked"`
	IsDir   *bool     `pulumi:"is_dir"`
	Targets *[]string `pulumi:"targets"`
}

// All resources must implement Create at a minumum.
func (l *Link) Create(ctx p.Context, name string, input LinkArgs, preview bool) (string, LinkState, error) {
	state := &LinkState{
		LinkArgs: input,
		Targets:  &[]string{},
	}
	if preview {
		return name, *state, nil
	}

	source, err := os.Lstat(*input.Source)
	if err != nil {
		return "", *state, err
	}
	if input.Recursive != nil && *input.Recursive == true && source.IsDir() {
		_, err := os.Lstat(*input.Target)
		if err != nil {
			return "", *state, err
		}
		if err := state.linkFile(source, *input.Target, ""); err != nil {
			return "", *state, err
		}
	} else {
		if err := os.Link(*input.Source, *input.Target); err != nil {
			return "", *state, err
		}
		*state.Targets = append(*state.Targets, *input.Target)
	}

	if err := state.stats(input); err != nil {
		return name, *state, nil
	}

	return name, *state, nil
}

func (l *Link) Read(ctx p.Context, id string, inputs LinkArgs, state LinkState) (
	canonicalID string, normalizedInputs LinkArgs, normalizedState LinkState, err error) {

	sExists, err := os.Lstat(*inputs.Source)
	if err != nil {
		return "", LinkArgs{}, LinkState{}, err
	}
	if sExists.IsDir() {
		b := true
		state.IsDir = &b
	} else {
		b := false
		state.IsDir = &b
	}
	exists, err := os.Lstat(*inputs.Target)
	if err != nil {
		state.Target = nil
	} else {
		if exists.Mode() != os.ModeSymlink {
			b := false
			state.Linked = &b
		} else {
			b := true
			state.Linked = &b
		}
	}
	return *inputs.Source, inputs, state, nil

}

func (l *Link) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (LinkArgs, []p.CheckFailure, error) {
	return infer.DefaultCheck[LinkArgs](newInputs)
}

func (l *Link) Update(ctx p.Context, name string, olds LinkState, news LinkArgs, preview bool) (LinkState, error) {
	state := &LinkState{
		LinkArgs: news,
		Linked:   olds.Linked,
		IsDir:    olds.IsDir,
		Targets:  olds.Targets,
	}
	recursive := news.Recursive != nil && *news.Recursive

	if preview {
		return *state, nil
	}

	oldTargets := *state.Targets

	// if we are recursively linking files in a directory
	// then we need to do the update
	if *olds.IsDir && recursive {
		source, err := os.Lstat(*news.Source)
		_, err = os.Lstat(*news.Target)
		if err != nil {
			return *state, err
		}
		if err := state.linkFile(source, *news.Target, ""); err != nil {
			return *state, err
		}
		state.stats(news)
		removeOld(oldTargets, *state.Targets)
		return *state, nil
	}

	// nothing has changed
	if *news.Source == *olds.Source &&
		*news.Target == *olds.Target {
		return *state, nil
	}

	// remove old target
	if *news.Target != *olds.Target {
		if err := os.Remove(*olds.Target); err != nil {
			return LinkState{}, errors.New(fmt.Sprintf("Error removing target %s: %s", *olds.Target, err))
		}
	}

	// make new link
	if *news.Source != *olds.Source ||
		*news.Target != *olds.Target {
		// first remove the old link
		if err := os.Remove(*news.Target); err != nil && !os.IsNotExist(err) {
			return LinkState{}, err
		}
		if err := os.Link(*news.Source, *news.Target); err != nil {
			return LinkState{}, err
		}
	}
	state.stats(news)
	return *state, nil
}

func (l *Link) Delete(ctx p.Context, id string, props LinkState) error {
	if props.Retain != nil && *props.Retain == true {
		return nil
	}

	if err := props.removeTargets(); err != nil {
		return err
	}

	return nil
}

// TODO: currently this doesn't handle errors very well. If we fail on one file
// we end up in a bad state where some things have been linked while others are not
// TODO: handle name collisions
func (l *LinkState) linkFile(source fs.FileInfo, target string, parent string) error {
	p := path.Join(parent, source.Name())
	if source.IsDir() {
		entries, err := os.ReadDir(p)
		if err != nil {
			return err
		}
		for _, de := range entries {
			info, err := de.Info()
			if err != nil {
				return err
			}
			return l.linkFile(info, target, p)
		}
	}
	if err := os.Link(source.Name(), p); err != nil {
		return err
	}
	*l.Targets = append(*l.Targets, p)
	return nil
}

func (l *LinkState) removeTargets() error {
	for _, v := range *l.Targets {
		if err := os.Remove(v); err != nil && !os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("Error unlinking target %s: %s", v, err))
		}
	}
	return nil
}

func removeOld(olds, news []string) error {
	for _, v := range olds {
		if !contains(v, news) {
			if err := os.Remove(v); err != nil && !os.IsNotExist(err) {
				return err
			}

		}
	}
	return nil
}

func contains(a string, b []string) bool {
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}

func (l *LinkState) stats(input LinkArgs) error {
	t := true
	f := false
	l.Linked = &t
	file, err := os.Lstat(*input.Source)
	if err != nil {
		return err
	}
	if file.IsDir() {
		l.IsDir = &t
	} else {
		l.IsDir = &f
	}
	return nil
}
