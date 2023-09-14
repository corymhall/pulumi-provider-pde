package main

import (
	"os"
	"path"

	"github.com/corymhall/pulumi-provider-pde/sdk/go/pdec/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		prof, err := local.NewProfile(ctx, "zsh", &local.ProfileArgs{
			FileName: pulumi.String(path.Join(wd, "zshbkp")),
		})
		prof.AddLines(ctx, &local.ProfileAddLinesArgs{
			Lines: pulumi.StringArray{
				pulumi.String("Hello"),
				pulumi.String("World"),
			},
		})
		name, err := prof.GetFileName(ctx)
		if err != nil {
			return err
		}
		ctx.Export("name", name)
		return nil
	})
}
