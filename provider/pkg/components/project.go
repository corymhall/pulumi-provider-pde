package components

import (
	"os"
	"path"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Project struct {
	pulumi.ResourceState
	Dir *string
}

func NewProject(ctx *pulumi.Context, name string, opts ...pulumi.ResourceOption) (*Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	projectDir := path.Join(cwd, name)
	project := &Project{
		Dir: &projectDir,
	}
	if err := ctx.RegisterComponentResource("pde:components:Project", name, project); err != nil {
		return nil, err
	}
	return project, nil
}
