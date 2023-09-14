package local

import (
	"github.com/corymhall/pulumi-provider-pde/sdk/go/pde/local"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Profile struct {
	fileName pulumi.StringOutput `pulumi:"fileName"`
	content  pulumi.StringArray  `pulumi:"content"`
}

type profile struct{}

type ProfileArgs struct {
	FileName pulumi.StringInput `pulumi:"fileName"`
}

type ProfileState struct {
	pulumi.ResourceState
	FileName pulumi.StringOutput `pulumi:"fileName"`
}

type GetFileNameArgs struct{}

func (p *Profile) Construct(ctx *pulumi.Context, name, typ string, args ProfileArgs, opts pulumi.ResourceOption) (*ProfileState, error) {
	profile := &ProfileState{}
	if err := ctx.RegisterComponentResource(typ, name, profile, opts); err != nil {
		return nil, err
	}
	profile.FileName = args.FileName.ToStringOutput()
	p.content = pulumi.StringArray{pulumi.String("hello")}
	local.NewFile(ctx, name, &local.FileArgs{
		Path:    args.FileName,
		Force:   pulumi.BoolPtr(false),
		Content: p.content,
	})

	if err := ctx.RegisterResourceOutputs(profile, pulumi.Map{
		"fileName": args.FileName,
	}); err != nil {
		return nil, err
	}
	return profile, nil
}

type GetFileNameResult struct {
	Result pulumi.StringOutput `pulumi:"result"`
}

func (p *Profile) Call(ctx p.Context, input GetFileNameArgs) (GetFileNameResult, error) {
	return GetFileNameResult{
		Result: p.GetFileName(&input),
	}, nil

}

func (p *Profile) GetFileName(args *GetFileNameArgs) pulumi.StringOutput {
	return p.fileName
}

type Empty struct{}

func (p *Profile) AddLines(lines ...pulumi.StringInput) Empty {
	p.content = append(p.content, lines...)
	return Empty{}
}
