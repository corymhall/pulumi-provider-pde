package provider

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/corymhall/pulumi-provider-pde/components/pkg/local"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
)

type module struct {
	version semver.Version
}

func (m *module) Version() semver.Version {
	return m.version
}

func (m *module) Construct(ctx *pulumi.Context, name, typ, urn string) (r pulumi.Resource, err error) {
	switch typ {
	case "pdec:local:Profile":
		r = &local.Profile{}
		err = ctx.RegisterResource(typ, name, nil, r, pulumi.URN_(urn))
		return
	default:
		return nil, fmt.Errorf("unknown resource type: %s", typ)
	}
}

func Serve(providerName, version string, schema []byte) {
	// Start the gRPC service
	pulumi.RegisterResourceModule("pdec", "local", &module{semver.MustParse(version)})
	if err := provider.MainWithOptions(provider.Options{
		Name:    providerName,
		Version: version,
		Schema:  schema,
		Construct: func(ctx *pulumi.Context, typ, name string, inputs pulumiprovider.ConstructInputs,
			options pulumi.ResourceOption) (*pulumiprovider.ConstructResult, error) {
			switch typ {
			case "pdec:local:Profile":
				args := &local.ProfileArgs{}
				if err := inputs.CopyTo(args); err != nil {
					return nil, fmt.Errorf("setting args: %w", err)
				}
				profile, err := local.NewProfile(ctx, name, args, options)
				if err != nil {
					return nil, fmt.Errorf("creating component: %w", err)
				}
				return pulumiprovider.NewConstructResult(profile)

			default:
				return nil, fmt.Errorf("unknown resource type: %s", typ)
			}

		},
		Call: func(ctx *pulumi.Context, tok string, args pulumiprovider.CallArgs) (*pulumiprovider.CallResult, error) {
			switch tok {
			case "pdec:local:Profile/addLines":
				methodArgs := &local.AddLinesArgs{}
				res, err := args.CopyTo(methodArgs)
				if err != nil {
					return nil, fmt.Errorf("settings args: %w", err)
				}
				profile, ok := res.(*local.Profile)
				if !ok {
					fmt.Println("res not of type profile")
				}
				profile.AddLines(methodArgs)
				return pulumiprovider.NewCallResult(&local.Empty{})
			case "pdec:local:Profile/getFileName":
				methodArgs := &local.GetFileNameArgs{}
				res, err := args.CopyTo(methodArgs)
				if err != nil {
					return nil, fmt.Errorf("settings args: %w", err)
				}
				profile, ok := res.(*local.Profile)
				if !ok {
					fmt.Println("res not of type profile")
				}
				result := profile.GetFileName(methodArgs)
				return pulumiprovider.NewCallResult(&local.GetFileNameResult{Result: result})
			default:
				return nil, fmt.Errorf("unknown method %s", tok)
			}

		},
	}); err != nil {
		cmdutil.ExitError(err.Error())
	}
}
