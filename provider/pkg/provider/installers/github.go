package installers

import (
	"github.com/pulumi/pulumi-go-provider/infer"
)

type GitHubBaseInputs struct {
	BaseInputs
	InstallCommands *[]string `pulumi:"installCommands,optional"`
	Org             *string   `pulumi:"org"`
	Repo            *string   `pulumi:"repo"`
}

func (g *GitHubBaseInputs) Annotate(a infer.Annotator) {
	a.Describe(&g.InstallCommands, "The commands to run to install the program")
	a.Describe(&g.Org, "The GitHub organization the repo belongs to")
	a.Describe(&g.Repo, "The GitHub repository name")
}
