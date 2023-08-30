package local

import (
	"errors"
	"os"

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
var _ = (infer.CustomCheck[LinkArgs])((*Link)(nil))

// Each resource has in input struct, defining what arguments it accepts.
type LinkArgs struct {
	Source    string `pulumi:"source"`
	Target    string `pulumi:"target"`
	Linked    bool   `pulumi:"linked"`
	IsDir     bool   `pulumi:"is_dir"`
	Overwrite bool   `pulumi:"overwrite"`
	Exists    bool   `pulumi:"exists"`
	// Fields projected into Pulumi must be public and hava a `pulumi:"..."` tag.
	// The pulumi tag doesn't need to match the field name, but its generally a
	// good idea.
	// Length int `pulumi:"length"`
}

// Each resource has a state, describing the fields that exist on the created resource.
type LinkState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	LinkArgs
	// Here we define a required output called result.
	Result string `pulumi:"result"`
}

// All resources must implement Create at a minumum.
func (Link) Create(ctx p.Context, name string, input LinkArgs, preview bool) (string, LinkState, error) {
	if input.Exists && !input.Overwrite {
		return "", LinkState{}, errors.New("")
	}

	if err := os.Link(input.Source, input.Target); err != nil {
		return "", LinkState{}, err
	}
	file, err := os.Lstat(input.Source)
	if err != nil {
		if file.IsDir() {
			input.IsDir = true
		} else {
			input.IsDir = false
		}
	}

	input.Linked = true
	input.Exists = true

	return name, LinkState{LinkArgs: input}, nil
}

func (Link) Read(ctx p.Context, id string, inputs LinkArgs, state LinkState) (
	canonicalID string, normalizedInputs LinkArgs, normalizedState LinkState, err error) {

	sExists, err := os.Lstat(inputs.Source)
	if err != nil {
		return "", LinkArgs{}, LinkState{}, err
	}
	if sExists.IsDir() {
		inputs.IsDir = true
	} else {
		inputs.IsDir = false
	}
	exists, err := os.Lstat(inputs.Target)
	if err != nil {
		inputs.Exists = false
	} else {
		if exists.Mode() != os.ModeSymlink {
			inputs.Linked = false
		} else {
			inputs.Linked = true
		}
	}
	return inputs.Source, inputs, LinkState{LinkArgs: inputs}, nil

}

func (Link) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (LinkArgs, []p.CheckFailure, error) {
	return infer.DefaultCheck[LinkArgs](newInputs)
}
