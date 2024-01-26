package installers

import (
	"github.com/pulumi/pulumi-go-provider/infer"
)

type BaseInputs struct {
	UpdateCommands    *[]string `pulumi:"updateCommands,optional"`
	UninstallCommands *[]string `pulumi:"uninstallCommands,optional"`
}

func (b *BaseInputs) Annotate(a infer.Annotator) {
	a.Describe(&b.UpdateCommands, "Optional Commands to run to update the program")
	a.Describe(&b.UninstallCommands, "Optional Commands to run to uninstall the program")
}

type CommandInputs struct {
	Interpreter *[]string          `pulumi:"interpreter,optional"`
	Environment *map[string]string `pulumi:"environment,optional"`
}

func (c *CommandInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Interpreter, "The interpreter to use to run the commands. Defaults to ['/bin/sh', '-c']")
	a.Describe(&c.Environment, "The environment variables to set when running the commands")
}

type CommandOutputs struct {
	CommandInputs
}

type BaseOutputs struct {
	Version *string `pulumi:"version"`
}

func (b *BaseOutputs) Annotate(a infer.Annotator) {
	a.Describe(&b.Version, "The version of the program")
}
