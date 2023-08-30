package tests

import (
	"fmt"
	"testing"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/blang/semver"
	pde "github.com/corymhall/pulumi-provider-pde/pkg/provider"
	"github.com/hashicorp/go-version"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

func provider() integration.Server {
	v := semver.MustParse(version.Version)
	return integration.NewServer(sde.Name, v, pde.NewProvider())

}
func urn(mod, res, name string) resource.URN {
	m := tokens.ModuleName(mod)
	r := tokens.TypeName(res)
	if !tokens.IsQName(name) {
		panic(fmt.Sprintf("invalid resource name: %q", name))
	}
	n := tokens.QName(name)
	return resource.NewURN("test", "command", "",
		tokens.NewTypeToken(
			tokens.NewModuleToken(command.Name, m),
			r),
		n)
}

func TestLinkCommand(t *testing.T) {
	t.Parallel()
	cmd := provider()
	urn := urn("local", "link", "file")

	// Run a create against an in-memory provider, assert it succeeded, and return the
	// created property map
	create := func(preview bool, env resource.PropertyValue) resource.PropertyMap {
		resp, err := cmd.Create(p.CreateRequest{})
	}

}
