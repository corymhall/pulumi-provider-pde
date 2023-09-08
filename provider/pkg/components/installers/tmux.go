package installers

import (
	"fmt"

	"github.com/corymhall/pulumi-provider-pde/sdk/go/pde/installers"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Tmux struct {
	pulumi.ResourceState
}

func NewTmux(ctx *pulumi.Context, name string, opts ...pulumi.ResourceOption) (*Tmux, error) {
	tmux := &Tmux{}
	if err := ctx.RegisterComponentResource("pde:installers:Tmux", name, tmux); err != nil {
		return nil, err
	}
	_, err := installers.NewGitHubRelease(ctx, fmt.Sprintf("%s-tmux", name), &installers.GitHubReleaseArgs{
		Version: pulumi.String("3.3a"),
		Org:     pulumi.String("tmux"),
		Repo:    pulumi.String("tmux"),
		InstallCommands: pulumi.StringArray{
			pulumi.String("tar -zxf tmux-3.3a.tar.gz"),
			pulumi.String(`cd tmux-3.3a &&
sudo ./configure --prefix $HOME/.local &&
sudo make && sudo make install &&
sudo chmod +x $HOME/.local/bin/tmux &&`),
		},
	}, pulumi.Parent(tmux))
	if err != nil {
		return nil, err
	}

	return tmux, nil
}
