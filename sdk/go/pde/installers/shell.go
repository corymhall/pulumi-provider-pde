// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package installers

import (
	"context"
	"reflect"

	"errors"
	"github.com/corymhall/pulumi-provider-pde/sdk/go/pde/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Install something from a URL using shell commands
type Shell struct {
	pulumi.CustomResourceState

	// The location to put the program. Defaults to $HOME/.local/bin
	BinLocation pulumi.StringPtrOutput `pulumi:"binLocation"`
	// The URL to download the program from
	DownloadURL pulumi.StringOutput `pulumi:"downloadURL"`
	// The environment variables to set when running the commands
	Environment pulumi.StringMapOutput `pulumi:"environment"`
	// Whether the program that is download is an executable
	Executable pulumi.BoolPtrOutput `pulumi:"executable"`
	// The commands to run to install the program
	InstallCommands pulumi.StringArrayOutput `pulumi:"installCommands"`
	// The interpreter to use to run the commands. Defaults to ['/bin/sh', '-c']
	Interpreter pulumi.StringArrayOutput `pulumi:"interpreter"`
	// The location the program was installed to
	Location pulumi.StringPtrOutput `pulumi:"location"`
	// The name of the program. This is the name you would use to execute the program
	ProgramName pulumi.StringOutput `pulumi:"programName"`
	// Optional Commands to run to uninstall the program
	UninstallCommands pulumi.StringArrayOutput `pulumi:"uninstallCommands"`
	// Optional Commands to run to update the program
	UpdateCommands pulumi.StringArrayOutput `pulumi:"updateCommands"`
	// The version of the program
	Version pulumi.StringOutput `pulumi:"version"`
	// The command to run to get the version of the program. This is needed if you want to keep track of the version in state
	VersionCommand pulumi.StringPtrOutput `pulumi:"versionCommand"`
}

// NewShell registers a new resource with the given unique name, arguments, and options.
func NewShell(ctx *pulumi.Context,
	name string, args *ShellArgs, opts ...pulumi.ResourceOption) (*Shell, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.DownloadURL == nil {
		return nil, errors.New("invalid value for required argument 'DownloadURL'")
	}
	if args.InstallCommands == nil {
		return nil, errors.New("invalid value for required argument 'InstallCommands'")
	}
	if args.ProgramName == nil {
		return nil, errors.New("invalid value for required argument 'ProgramName'")
	}
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource Shell
	err := ctx.RegisterResource("pde:installers:Shell", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetShell gets an existing Shell resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetShell(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *ShellState, opts ...pulumi.ResourceOption) (*Shell, error) {
	var resource Shell
	err := ctx.ReadResource("pde:installers:Shell", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering Shell resources.
type shellState struct {
}

type ShellState struct {
}

func (ShellState) ElementType() reflect.Type {
	return reflect.TypeOf((*shellState)(nil)).Elem()
}

type shellArgs struct {
	// The location to put the program. Defaults to $HOME/.local/bin
	BinLocation *string `pulumi:"binLocation"`
	// The URL to download the program from
	DownloadURL string `pulumi:"downloadURL"`
	// The environment variables to set when running the commands
	Environment map[string]string `pulumi:"environment"`
	// Whether the program that is download is an executable
	Executable *bool `pulumi:"executable"`
	// The commands to run to install the program
	InstallCommands []string `pulumi:"installCommands"`
	// The interpreter to use to run the install commands. Defaults to ['/bin/sh', '-c']
	Interpreter []string `pulumi:"interpreter"`
	// The name of the program. This is the name you would use to execute the program
	ProgramName string `pulumi:"programName"`
	// Optional Commands to run to uninstall the program
	UninstallCommands []string `pulumi:"uninstallCommands"`
	// Optional Commands to run to update the program
	UpdateCommands []string `pulumi:"updateCommands"`
	// The command to run to get the version of the program. This is needed if you want to keep track of the version in state
	VersionCommand *string `pulumi:"versionCommand"`
}

// The set of arguments for constructing a Shell resource.
type ShellArgs struct {
	// The location to put the program. Defaults to $HOME/.local/bin
	BinLocation pulumi.StringPtrInput
	// The URL to download the program from
	DownloadURL pulumi.StringInput
	// The environment variables to set when running the commands
	Environment pulumi.StringMapInput
	// Whether the program that is download is an executable
	Executable pulumi.BoolPtrInput
	// The commands to run to install the program
	InstallCommands pulumi.StringArrayInput
	// The interpreter to use to run the install commands. Defaults to ['/bin/sh', '-c']
	Interpreter pulumi.StringArrayInput
	// The name of the program. This is the name you would use to execute the program
	ProgramName pulumi.StringInput
	// Optional Commands to run to uninstall the program
	UninstallCommands pulumi.StringArrayInput
	// Optional Commands to run to update the program
	UpdateCommands pulumi.StringArrayInput
	// The command to run to get the version of the program. This is needed if you want to keep track of the version in state
	VersionCommand pulumi.StringPtrInput
}

func (ShellArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*shellArgs)(nil)).Elem()
}

type ShellInput interface {
	pulumi.Input

	ToShellOutput() ShellOutput
	ToShellOutputWithContext(ctx context.Context) ShellOutput
}

func (*Shell) ElementType() reflect.Type {
	return reflect.TypeOf((**Shell)(nil)).Elem()
}

func (i *Shell) ToShellOutput() ShellOutput {
	return i.ToShellOutputWithContext(context.Background())
}

func (i *Shell) ToShellOutputWithContext(ctx context.Context) ShellOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ShellOutput)
}

// ShellArrayInput is an input type that accepts ShellArray and ShellArrayOutput values.
// You can construct a concrete instance of `ShellArrayInput` via:
//
//	ShellArray{ ShellArgs{...} }
type ShellArrayInput interface {
	pulumi.Input

	ToShellArrayOutput() ShellArrayOutput
	ToShellArrayOutputWithContext(context.Context) ShellArrayOutput
}

type ShellArray []ShellInput

func (ShellArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*Shell)(nil)).Elem()
}

func (i ShellArray) ToShellArrayOutput() ShellArrayOutput {
	return i.ToShellArrayOutputWithContext(context.Background())
}

func (i ShellArray) ToShellArrayOutputWithContext(ctx context.Context) ShellArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ShellArrayOutput)
}

// ShellMapInput is an input type that accepts ShellMap and ShellMapOutput values.
// You can construct a concrete instance of `ShellMapInput` via:
//
//	ShellMap{ "key": ShellArgs{...} }
type ShellMapInput interface {
	pulumi.Input

	ToShellMapOutput() ShellMapOutput
	ToShellMapOutputWithContext(context.Context) ShellMapOutput
}

type ShellMap map[string]ShellInput

func (ShellMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*Shell)(nil)).Elem()
}

func (i ShellMap) ToShellMapOutput() ShellMapOutput {
	return i.ToShellMapOutputWithContext(context.Background())
}

func (i ShellMap) ToShellMapOutputWithContext(ctx context.Context) ShellMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ShellMapOutput)
}

type ShellOutput struct{ *pulumi.OutputState }

func (ShellOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**Shell)(nil)).Elem()
}

func (o ShellOutput) ToShellOutput() ShellOutput {
	return o
}

func (o ShellOutput) ToShellOutputWithContext(ctx context.Context) ShellOutput {
	return o
}

// The location to put the program. Defaults to $HOME/.local/bin
func (o ShellOutput) BinLocation() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Shell) pulumi.StringPtrOutput { return v.BinLocation }).(pulumi.StringPtrOutput)
}

// The URL to download the program from
func (o ShellOutput) DownloadURL() pulumi.StringOutput {
	return o.ApplyT(func(v *Shell) pulumi.StringOutput { return v.DownloadURL }).(pulumi.StringOutput)
}

// The environment variables to set when running the commands
func (o ShellOutput) Environment() pulumi.StringMapOutput {
	return o.ApplyT(func(v *Shell) pulumi.StringMapOutput { return v.Environment }).(pulumi.StringMapOutput)
}

// Whether the program that is download is an executable
func (o ShellOutput) Executable() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *Shell) pulumi.BoolPtrOutput { return v.Executable }).(pulumi.BoolPtrOutput)
}

// The commands to run to install the program
func (o ShellOutput) InstallCommands() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *Shell) pulumi.StringArrayOutput { return v.InstallCommands }).(pulumi.StringArrayOutput)
}

// The interpreter to use to run the commands. Defaults to ['/bin/sh', '-c']
func (o ShellOutput) Interpreter() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *Shell) pulumi.StringArrayOutput { return v.Interpreter }).(pulumi.StringArrayOutput)
}

// The location the program was installed to
func (o ShellOutput) Location() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Shell) pulumi.StringPtrOutput { return v.Location }).(pulumi.StringPtrOutput)
}

// The name of the program. This is the name you would use to execute the program
func (o ShellOutput) ProgramName() pulumi.StringOutput {
	return o.ApplyT(func(v *Shell) pulumi.StringOutput { return v.ProgramName }).(pulumi.StringOutput)
}

// Optional Commands to run to uninstall the program
func (o ShellOutput) UninstallCommands() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *Shell) pulumi.StringArrayOutput { return v.UninstallCommands }).(pulumi.StringArrayOutput)
}

// Optional Commands to run to update the program
func (o ShellOutput) UpdateCommands() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *Shell) pulumi.StringArrayOutput { return v.UpdateCommands }).(pulumi.StringArrayOutput)
}

// The version of the program
func (o ShellOutput) Version() pulumi.StringOutput {
	return o.ApplyT(func(v *Shell) pulumi.StringOutput { return v.Version }).(pulumi.StringOutput)
}

// The command to run to get the version of the program. This is needed if you want to keep track of the version in state
func (o ShellOutput) VersionCommand() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Shell) pulumi.StringPtrOutput { return v.VersionCommand }).(pulumi.StringPtrOutput)
}

type ShellArrayOutput struct{ *pulumi.OutputState }

func (ShellArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*Shell)(nil)).Elem()
}

func (o ShellArrayOutput) ToShellArrayOutput() ShellArrayOutput {
	return o
}

func (o ShellArrayOutput) ToShellArrayOutputWithContext(ctx context.Context) ShellArrayOutput {
	return o
}

func (o ShellArrayOutput) Index(i pulumi.IntInput) ShellOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *Shell {
		return vs[0].([]*Shell)[vs[1].(int)]
	}).(ShellOutput)
}

type ShellMapOutput struct{ *pulumi.OutputState }

func (ShellMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*Shell)(nil)).Elem()
}

func (o ShellMapOutput) ToShellMapOutput() ShellMapOutput {
	return o
}

func (o ShellMapOutput) ToShellMapOutputWithContext(ctx context.Context) ShellMapOutput {
	return o
}

func (o ShellMapOutput) MapIndex(k pulumi.StringInput) ShellOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *Shell {
		return vs[0].(map[string]*Shell)[vs[1].(string)]
	}).(ShellOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*ShellInput)(nil)).Elem(), &Shell{})
	pulumi.RegisterInputType(reflect.TypeOf((*ShellArrayInput)(nil)).Elem(), ShellArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*ShellMapInput)(nil)).Elem(), ShellMap{})
	pulumi.RegisterOutputType(ShellOutput{})
	pulumi.RegisterOutputType(ShellArrayOutput{})
	pulumi.RegisterOutputType(ShellMapOutput{})
}
