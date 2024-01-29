package local

import (
	// "errors"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
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
var _ = (infer.CustomDiff[LinkArgs, LinkState])((*Link)(nil))
var _ = (infer.CustomDelete[LinkState])((*Link)(nil))
var _ = (infer.CustomCheck[LinkArgs])((*Link)(nil))

// Each resource has in input struct, defining what arguments it accepts.
type LinkArgs struct {
	Source    *string `pulumi:"source"`
	Target    *string `pulumi:"target"`
	Overwrite *bool   `pulumi:"overwrite,optional"`
	Retain    *bool   `pulumi:"retain,optional"`
	Recursive *bool   `pulumi:"recursive,optional"`
}

// Each resource has a state, describing the fields that exist on the created resource.
type LinkState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	LinkArgs
	Linked  *bool     `pulumi:"linked"`
	IsDir   *bool     `pulumi:"isDir"`
	Targets *[]string `pulumi:"targets"`
}

func (l *Link) Annotate(a infer.Annotator) {
	a.Describe(&l, "Create a symlink for a file or directory")
}

func (l *LinkState) Annotate(a infer.Annotator) {
	a.Describe(&l.Linked, "Whether the symlink has been created")
	a.Describe(&l.IsDir, "Whether the source is a directory")
	a.Describe(&l.Targets, "The targets locations of the symlink")
}

func (l *LinkArgs) Annotate(a infer.Annotator) {
	a.Describe(&l.Source, "The source file or directory to create a link to")
	a.Describe(&l.Target, "The target file or directory to create a link at")
	a.Describe(&l.Overwrite, "Whether to overwrite the target if it exists")
	a.Describe(&l.Retain, "Whether to retain the link if the resource is deleted")
	a.Describe(&l.Recursive, "Whether to recursively create links for directories")
}

func (l *Link) Diff(ctx p.Context, id string, olds LinkState, news LinkArgs) (p.DiffResponse, error) {
	diff := map[string]p.PropertyDiff{}
	if (news.Recursive == nil && olds.Recursive != nil) || (news.Recursive != nil && olds.Recursive == nil) {
		diff["recursive"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	if (news.Recursive != nil && olds.Recursive != nil) && *news.Recursive != *olds.Recursive {
		diff["recursive"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	if *news.Source != *olds.Source {
		diff["source"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	if *news.Target != *olds.Target {
		diff["target"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}

	if olds.Linked == nil || *olds.Linked == false {
		diff["linked"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}

	return p.DiffResponse{
		DeleteBeforeReplace: true,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// All resources must implement Create at a minumum.
func (l *Link) Create(ctx p.Context, name string, input LinkArgs, preview bool) (string, LinkState, error) {
	f := false
	state := &LinkState{
		LinkArgs: input,
		Targets:  &[]string{},
		Linked:   &f,
		IsDir:    &f,
	}
	if preview {
		return name, *state, nil
	}

	if err := state.link(ctx); err != nil {
		return name, LinkState{}, err
	}

	return name, *state, nil
}

func (l *Link) Read(ctx p.Context, id string, inputs LinkArgs, state LinkState) (
	canonicalID string, normalizedInputs LinkArgs, normalizedState LinkState, err error) {

	sExists, err := os.Lstat(*inputs.Source)
	if err != nil {
		return "", LinkArgs{}, LinkState{}, err
	}
	f := false
	state.Linked = &f
	isDir := sExists.IsDir()
	state.IsDir = &isDir
	exists, err := os.Lstat(*inputs.Target)
	ctx.Logf(diag.Warning, "target file", exists.Mode())
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
	return id, inputs, state, nil

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

	if preview {
		return *state, nil
	}

	state.Targets = &[]string{}

	if err := state.link(ctx); err != nil {
		return LinkState{}, err
	}

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

func (l *LinkState) link(ctx p.Context) error {
	source, err := os.Lstat(*l.Source)
	if err != nil {
		return err
	}
	if l.Recursive != nil && *l.Recursive && source.IsDir() {
		_, err := os.Lstat(*l.Target)
		if err != nil {
			return err
		}
		if err := l.linkFile(ctx, *l.Source, *l.Target); err != nil {
			return err
		}
	} else {
		if l.Overwrite != nil && *l.Overwrite {
			os.Remove(*l.Target)
		}
		if err := os.Symlink(*l.Source, *l.Target); err != nil {
			return err
		}
		tfTargets := append(*l.Targets, *l.Target)
		l.Targets = &tfTargets
	}

	if err := l.stats(l.LinkArgs); err != nil {
		return err
	}
	return nil
}

// TODO: currently this doesn't handle errors very well. If we fail on one file
// we end up in a bad state where some things have been linked while others are not
// TODO: handle name collisions
func (l *LinkState) linkFile(ctx p.Context, source string, target string) error {
	s, err := os.Lstat(source)
	if err != nil {
		return err
	}
	var tfTargets []string
	if s.IsDir() {
		entries, err := os.ReadDir(source)
		if err != nil {
			return err
		}
		for _, de := range entries {
			p := filepath.Join(source, de.Name())
			t := filepath.Join(target, de.Name())
			if err := os.Symlink(p, t); err != nil {
				return err
			}
			tfTargets = append(tfTargets, *l.Targets...)
		}
	} else {
		if err := os.Symlink(source, target); err != nil {
			return err
		}

		tfTargets = append(tfTargets, *l.Targets...)
	}
	l.Targets = &tfTargets
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
	file, err = os.Lstat(*input.Target)
	if err != nil {
		return err
	}
	if file.Mode() == os.ModeSymlink {
		l.Linked = &t
	}
	return nil
}
