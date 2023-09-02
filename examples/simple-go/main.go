package main

import (
	"github.com/corymhall/pulumi-provider-pde/sdk/go/pde/installers"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		res, err := installers.NewGitHubRepo(ctx, "pulumi-provider-local", &installers.GitHubRepoArgs{
			Org:  pulumi.String("corymhall"),
			Repo: pulumi.String("pulumi-provider-local"),
		})
		ctx.Export("abspath", res.AbsFolderName)
		ctx.Export("version", res.Version)
		return err
	})
}
