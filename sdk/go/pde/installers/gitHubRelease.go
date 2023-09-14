// Code generated by pulumi-gen-pde DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package installers

import (
	"context"
	"reflect"

	"errors"
	"github.com/corymhall/pulumi-provider-pde/sdk/go/pde/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type GitHubRelease struct {
	pulumi.CustomResourceState

	AssetName         pulumi.StringPtrOutput   `pulumi:"assetName"`
	BinFolder         pulumi.StringPtrOutput   `pulumi:"binFolder"`
	DownloadURL       pulumi.StringOutput      `pulumi:"downloadURL"`
	Environment       pulumi.StringMapOutput   `pulumi:"environment"`
	Executable        pulumi.StringPtrOutput   `pulumi:"executable"`
	InstallCommands   pulumi.StringArrayOutput `pulumi:"installCommands"`
	Interpreter       pulumi.StringArrayOutput `pulumi:"interpreter"`
	Locations         pulumi.StringArrayOutput `pulumi:"locations"`
	Org               pulumi.StringOutput      `pulumi:"org"`
	Repo              pulumi.StringOutput      `pulumi:"repo"`
	UninstallCommands pulumi.StringArrayOutput `pulumi:"uninstallCommands"`
	UpdateCommands    pulumi.StringArrayOutput `pulumi:"updateCommands"`
	Version           pulumi.StringPtrOutput   `pulumi:"version"`
}

// NewGitHubRelease registers a new resource with the given unique name, arguments, and options.
func NewGitHubRelease(ctx *pulumi.Context,
	name string, args *GitHubReleaseArgs, opts ...pulumi.ResourceOption) (*GitHubRelease, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Org == nil {
		return nil, errors.New("invalid value for required argument 'Org'")
	}
	if args.Repo == nil {
		return nil, errors.New("invalid value for required argument 'Repo'")
	}
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource GitHubRelease
	err := ctx.RegisterResource("pde:installers:GitHubRelease", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetGitHubRelease gets an existing GitHubRelease resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetGitHubRelease(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *GitHubReleaseState, opts ...pulumi.ResourceOption) (*GitHubRelease, error) {
	var resource GitHubRelease
	err := ctx.ReadResource("pde:installers:GitHubRelease", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering GitHubRelease resources.
type gitHubReleaseState struct {
}

type GitHubReleaseState struct {
}

func (GitHubReleaseState) ElementType() reflect.Type {
	return reflect.TypeOf((*gitHubReleaseState)(nil)).Elem()
}

type gitHubReleaseArgs struct {
	AssetName         *string  `pulumi:"assetName"`
	BinFolder         *string  `pulumi:"binFolder"`
	Executable        *string  `pulumi:"executable"`
	InstallCommands   []string `pulumi:"installCommands"`
	Org               string   `pulumi:"org"`
	Repo              string   `pulumi:"repo"`
	UninstallCommands []string `pulumi:"uninstallCommands"`
	UpdateCommands    []string `pulumi:"updateCommands"`
}

// The set of arguments for constructing a GitHubRelease resource.
type GitHubReleaseArgs struct {
	AssetName         pulumi.StringPtrInput
	BinFolder         pulumi.StringPtrInput
	Executable        pulumi.StringPtrInput
	InstallCommands   pulumi.StringArrayInput
	Org               pulumi.StringInput
	Repo              pulumi.StringInput
	UninstallCommands pulumi.StringArrayInput
	UpdateCommands    pulumi.StringArrayInput
}

func (GitHubReleaseArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*gitHubReleaseArgs)(nil)).Elem()
}

type GitHubReleaseInput interface {
	pulumi.Input

	ToGitHubReleaseOutput() GitHubReleaseOutput
	ToGitHubReleaseOutputWithContext(ctx context.Context) GitHubReleaseOutput
}

func (*GitHubRelease) ElementType() reflect.Type {
	return reflect.TypeOf((**GitHubRelease)(nil)).Elem()
}

func (i *GitHubRelease) ToGitHubReleaseOutput() GitHubReleaseOutput {
	return i.ToGitHubReleaseOutputWithContext(context.Background())
}

func (i *GitHubRelease) ToGitHubReleaseOutputWithContext(ctx context.Context) GitHubReleaseOutput {
	return pulumi.ToOutputWithContext(ctx, i).(GitHubReleaseOutput)
}

// GitHubReleaseArrayInput is an input type that accepts GitHubReleaseArray and GitHubReleaseArrayOutput values.
// You can construct a concrete instance of `GitHubReleaseArrayInput` via:
//
//	GitHubReleaseArray{ GitHubReleaseArgs{...} }
type GitHubReleaseArrayInput interface {
	pulumi.Input

	ToGitHubReleaseArrayOutput() GitHubReleaseArrayOutput
	ToGitHubReleaseArrayOutputWithContext(context.Context) GitHubReleaseArrayOutput
}

type GitHubReleaseArray []GitHubReleaseInput

func (GitHubReleaseArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*GitHubRelease)(nil)).Elem()
}

func (i GitHubReleaseArray) ToGitHubReleaseArrayOutput() GitHubReleaseArrayOutput {
	return i.ToGitHubReleaseArrayOutputWithContext(context.Background())
}

func (i GitHubReleaseArray) ToGitHubReleaseArrayOutputWithContext(ctx context.Context) GitHubReleaseArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(GitHubReleaseArrayOutput)
}

// GitHubReleaseMapInput is an input type that accepts GitHubReleaseMap and GitHubReleaseMapOutput values.
// You can construct a concrete instance of `GitHubReleaseMapInput` via:
//
//	GitHubReleaseMap{ "key": GitHubReleaseArgs{...} }
type GitHubReleaseMapInput interface {
	pulumi.Input

	ToGitHubReleaseMapOutput() GitHubReleaseMapOutput
	ToGitHubReleaseMapOutputWithContext(context.Context) GitHubReleaseMapOutput
}

type GitHubReleaseMap map[string]GitHubReleaseInput

func (GitHubReleaseMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*GitHubRelease)(nil)).Elem()
}

func (i GitHubReleaseMap) ToGitHubReleaseMapOutput() GitHubReleaseMapOutput {
	return i.ToGitHubReleaseMapOutputWithContext(context.Background())
}

func (i GitHubReleaseMap) ToGitHubReleaseMapOutputWithContext(ctx context.Context) GitHubReleaseMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(GitHubReleaseMapOutput)
}

type GitHubReleaseOutput struct{ *pulumi.OutputState }

func (GitHubReleaseOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**GitHubRelease)(nil)).Elem()
}

func (o GitHubReleaseOutput) ToGitHubReleaseOutput() GitHubReleaseOutput {
	return o
}

func (o GitHubReleaseOutput) ToGitHubReleaseOutputWithContext(ctx context.Context) GitHubReleaseOutput {
	return o
}

func (o GitHubReleaseOutput) AssetName() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringPtrOutput { return v.AssetName }).(pulumi.StringPtrOutput)
}

func (o GitHubReleaseOutput) BinFolder() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringPtrOutput { return v.BinFolder }).(pulumi.StringPtrOutput)
}

func (o GitHubReleaseOutput) DownloadURL() pulumi.StringOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringOutput { return v.DownloadURL }).(pulumi.StringOutput)
}

func (o GitHubReleaseOutput) Environment() pulumi.StringMapOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringMapOutput { return v.Environment }).(pulumi.StringMapOutput)
}

func (o GitHubReleaseOutput) Executable() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringPtrOutput { return v.Executable }).(pulumi.StringPtrOutput)
}

func (o GitHubReleaseOutput) InstallCommands() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringArrayOutput { return v.InstallCommands }).(pulumi.StringArrayOutput)
}

func (o GitHubReleaseOutput) Interpreter() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringArrayOutput { return v.Interpreter }).(pulumi.StringArrayOutput)
}

func (o GitHubReleaseOutput) Locations() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringArrayOutput { return v.Locations }).(pulumi.StringArrayOutput)
}

func (o GitHubReleaseOutput) Org() pulumi.StringOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringOutput { return v.Org }).(pulumi.StringOutput)
}

func (o GitHubReleaseOutput) Repo() pulumi.StringOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringOutput { return v.Repo }).(pulumi.StringOutput)
}

func (o GitHubReleaseOutput) UninstallCommands() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringArrayOutput { return v.UninstallCommands }).(pulumi.StringArrayOutput)
}

func (o GitHubReleaseOutput) UpdateCommands() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringArrayOutput { return v.UpdateCommands }).(pulumi.StringArrayOutput)
}

func (o GitHubReleaseOutput) Version() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *GitHubRelease) pulumi.StringPtrOutput { return v.Version }).(pulumi.StringPtrOutput)
}

type GitHubReleaseArrayOutput struct{ *pulumi.OutputState }

func (GitHubReleaseArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*GitHubRelease)(nil)).Elem()
}

func (o GitHubReleaseArrayOutput) ToGitHubReleaseArrayOutput() GitHubReleaseArrayOutput {
	return o
}

func (o GitHubReleaseArrayOutput) ToGitHubReleaseArrayOutputWithContext(ctx context.Context) GitHubReleaseArrayOutput {
	return o
}

func (o GitHubReleaseArrayOutput) Index(i pulumi.IntInput) GitHubReleaseOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *GitHubRelease {
		return vs[0].([]*GitHubRelease)[vs[1].(int)]
	}).(GitHubReleaseOutput)
}

type GitHubReleaseMapOutput struct{ *pulumi.OutputState }

func (GitHubReleaseMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*GitHubRelease)(nil)).Elem()
}

func (o GitHubReleaseMapOutput) ToGitHubReleaseMapOutput() GitHubReleaseMapOutput {
	return o
}

func (o GitHubReleaseMapOutput) ToGitHubReleaseMapOutputWithContext(ctx context.Context) GitHubReleaseMapOutput {
	return o
}

func (o GitHubReleaseMapOutput) MapIndex(k pulumi.StringInput) GitHubReleaseOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *GitHubRelease {
		return vs[0].(map[string]*GitHubRelease)[vs[1].(string)]
	}).(GitHubReleaseOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*GitHubReleaseInput)(nil)).Elem(), &GitHubRelease{})
	pulumi.RegisterInputType(reflect.TypeOf((*GitHubReleaseArrayInput)(nil)).Elem(), GitHubReleaseArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*GitHubReleaseMapInput)(nil)).Elem(), GitHubReleaseMap{})
	pulumi.RegisterOutputType(GitHubReleaseOutput{})
	pulumi.RegisterOutputType(GitHubReleaseArrayOutput{})
	pulumi.RegisterOutputType(GitHubReleaseMapOutput{})
}
