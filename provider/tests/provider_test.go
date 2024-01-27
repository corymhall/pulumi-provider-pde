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
func urn(typ, res string) resource.URN {
	return resource.NewURN("stack", "proj", "",
		tokens.Type(fmt.Sprintf("test:%s:%s", typ, res)), "name")
}
