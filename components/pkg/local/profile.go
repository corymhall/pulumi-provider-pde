package local

import (
	"github.com/corymhall/pulumi-provider-pde/sdk/go/pde/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Profile struct {
	pulumi.ResourceState
	FileName pulumi.StringOutput `pulumi:"fileName"`
	content  pulumi.StringArray
}

type ProfileArgs struct {
	FileName pulumi.StringInput `pulumi:"fileName"`
}

type GetFileNameArgs struct{}

func NewProfile(ctx *pulumi.Context, name string, args *ProfileArgs, opts ...pulumi.ResourceOption) (*Profile, error) {
	profile := &Profile{
		content: pulumi.StringArray{},
	}
	if err := ctx.RegisterComponentResource("pdec:local:Profile", name, profile, opts...); err != nil {
		return nil, err
	}
	profile.FileName = args.FileName.ToStringOutput()
	content := profile.content.ToStringArrayOutput().ApplyT(func(c []string) []string {
		return c
	}).(pulumi.StringArrayOutput)

	if err := ctx.RegisterResourceOutputs(profile, pulumi.Map{
		"fileName": profile.FileName,
		"content":  content,
	}); err != nil {
		return nil, err
	}
	local.NewFile(ctx, name, &local.FileArgs{
		Path:    profile.FileName,
		Force:   pulumi.BoolPtr(false),
		Content: content,
	})
	return profile, nil
}

type GetFileNameResult struct {
	Result pulumi.StringOutput `pulumi:"result"`
}

func (p *Profile) GetFileName(args *GetFileNameArgs) pulumi.StringOutput {
	p.FileName = pulumi.String("something_else").ToStringOutput()
	return p.FileName
}

type Empty struct{}

type AddLinesArgs struct {
	Lines []pulumi.StringInput `pulumi:"lines"`
}

func (p *Profile) AddLines(args *AddLinesArgs) Empty {
	p.content = append(p.content, args.Lines...)
	return Empty{}
}
