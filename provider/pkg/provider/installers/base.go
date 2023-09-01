package installers

type BaseInputs struct {
	UpdateCommands    *[]string `pulumi:"updateCommands,optional"`
	UninstallCommands *[]string `pulumi:"uninstallCommands,optional"`
}

type CommandInputs struct {
	Interpreter *[]string          `pulumi:"interpreter"`
	Environment *map[string]string `pulumi:"environment,optional"`
}

type CommandOutputs struct {
	CommandInputs
}

type BaseOutputs struct {
	Version *string `pulumi:"version"`
}
