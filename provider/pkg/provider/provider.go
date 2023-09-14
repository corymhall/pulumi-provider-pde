package provider

import (
	"strings"

	"github.com/blang/semver"
	"github.com/corymhall/pulumi-provider-pde/provider/pkg/provider/installers"
	"github.com/corymhall/pulumi-provider-pde/provider/pkg/provider/local"
	// "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
	// "github.com/pulumi/pulumi/pkg/v3/resource/provider"
	// pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string = "0.0.1"

const (
	Name = "pde"
)

func NewProvider() p.Provider {

	// provider.Main(Name, func(hc *provider.HostClient) (pulumirpc.ResourceProviderServer, error) {
	//
	// })

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
		// Components: []infer.InferredComponent{
		// 	infer.Component[*local.Profile, local.ProfileArgs, *local.ProfileState](),
		// },
		// Functions: []infer.InferredFunction{
		// 	infer.Function[*local.Profile, local.GetFileNameArgs, local.GetFileNameResult](),
		// 	infer.Function[*local.Profile, pulumi.StringInput, local.GetFileNameResult](),
		// },
		Resources: []infer.InferredResource{
			infer.Resource[*local.Link, local.LinkArgs, local.LinkState](),
			infer.Resource[*local.File, local.FileArgs, local.FileState](),
			infer.Resource[*installers.GitHubRelease, installers.GitHubReleaseInputs, installers.GitHubReleaseOutputs](),
			infer.Resource[*installers.GitHubRepo, installers.GitHubRepoInputs, installers.GitHubRepoOutputs](),
		},
	})

}
func Schema(version string) (string, error) {
	version = strings.TrimPrefix(version, "v")
	s, err := integration.NewServer(Name, semver.MustParse(version), NewProvider()).
		GetSchema(p.GetSchemaRequest{})
	return s.Schema, err
}
