package profile

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SystemProfiles struct {
	pulumi.ResourceState
	Lines       []string
	SystemPaths map[string]bool
	SystemEnv   map[string]string
}

func NewSystemProfiles(ctx *pulumi.Context, name string, opts ...pulumi.ResourceOption) (*SystemProfiles, error) {
	profile := &SystemProfiles{
		Lines:       []string{},
		SystemPaths: map[string]bool{},
		SystemEnv: map[string]string{
			"GOPROXY":         "direct",
			"XDG_CONFIG_HOME": "",
		},
	}
	if err := ctx.RegisterComponentResource("pde:profile:SystemProfiles", name, profile); err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *SystemProfiles) AddToEnv(key, value string) {
	s.SystemEnv[key] = value
}

func (s *SystemProfiles) AddToSystemPath(path string) {
	s.SystemPaths[path] = true
}

func (s *SystemProfiles) AddLines(lines []string) {
	s.Lines = append(s.Lines, lines...)
}

func (s *SystemProfiles) AddAlias(name, command string) {
	s.Lines = append(s.Lines, fmt.Sprintf("alias %s='%s'", name, command))
}

func (s *SystemProfiles) Register(name, command string) {

}

type Profile struct {
	pulumi.ResourceState
}

func NewProfile(ctx *pulumi.Context, name string, opts ...pulumi.ResourceOption) (*Profile, error) {
	profile := &Profile{}
	if err := ctx.RegisterComponentResource("pde:profile:Profile", name, profile); err != nil {
		return nil, err
	}
	return profile, nil
}
