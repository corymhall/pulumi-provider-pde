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
			LanguageMap: map[string]any{
				"go": map[string]any{
					"generateResourceContainerTypes": true,
					"importBasePath":                 "github.com/corymhall/pulumi-provider-pde/sdk/go/pde",
				},
			},
		},
		Resources: []infer.InferredResource{
			infer.Resource[*local.Link, local.LinkArgs, local.LinkState](),
			infer.Resource[*local.File, local.FileArgs, local.FileState](),
			infer.Resource[*installers.GitHubRelease, installers.GitHubReleaseArgs, installers.GitHubReleaseState](),
			infer.Resource[*installers.GitHubRepo, installers.GitHubRepoArgs, installers.GitHubRepoState](),
			infer.Resource[*installers.Shell, installers.ShellArgs, installers.ShellState](),
			infer.Resource[*installers.Npm, installers.NpmArgs, installers.NpmState](),
		},
	})

}
