package installers

type GitHubBaseInputs struct {
	BaseInputs
	InstallCommands *[]string `pulumi:"installCommands,optional"`
	Org             *string   `pulumi:"org"`
	Repo            *string   `pulumi:"repo"`
}
