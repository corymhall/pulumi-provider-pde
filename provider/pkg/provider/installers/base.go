package installers

type BaseInputs struct {
	InstallCommands   *[]string `pulumi:"install_commands"`
	UpdateCommands    *[]string `pulumi:"update_commands,optional"`
	UninstallCommands *[]string `pulumi:"uninstall_commands,optional"`
	VersionCommand    *string   `pulumi:"version_command"`
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
