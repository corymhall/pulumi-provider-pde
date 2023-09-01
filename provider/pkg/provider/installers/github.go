package installers

type GitHubBaseInputs struct {
	BaseInputs
	InstallCommands *[]string `pulumi:"install_commands,optional"`
	Org             *string   `pulumi:"org"`
	Repo            *string   `pulumi:"repo"`
}
