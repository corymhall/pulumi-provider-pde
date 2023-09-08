package tests

import (
	"fmt"

	"github.com/blang/semver"
	pde "github.com/corymhall/pulumi-provider-pde/provider/pkg/provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

func provider() integration.Server {
	v := semver.MustParse(pde.Version)
	return integration.NewServer(pde.Name, v, pde.NewProvider())

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
			tokens.NewModuleToken(pde.Name, m),
			r),
		n)
}
