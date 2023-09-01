package provider

import (
	"github.com/corymhall/pulumi-provider-pde/provider/pkg/provider/installers"
	"github.com/corymhall/pulumi-provider-pde/provider/pkg/provider/local"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string = "0.0.1"

const (
	Name = "pde"
)

func NewProvider() p.Provider {

	// We tell the provider what resources it needs to support.
	// In this case, a single custom resource.
	return infer.Provider(infer.Options{
		Metadata: schema.Metadata{
			DisplayName: "pde",
			Description: "The pulumi pde provider...",
		},
		Resources: []infer.InferredResource{
			infer.Resource[*local.Link, local.LinkArgs, local.LinkState](),
			infer.Resource[*installers.GitHubRelease, installers.GitHubReleaseInputs, installers.GitHubReleaseOutputs](),
			infer.Resource[*installers.GitHubRepo, installers.GitHubRepoInputs, installers.GitHubRepoOutputs](),
		},
	})

}
