package installers

import (
	// "errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type Shell struct{}

var _ = (infer.CustomRead[ShellInputs, ShellOutputs])((*Shell)(nil))
var _ = (infer.CustomUpdate[ShellInputs, ShellOutputs])((*Shell)(nil))
var _ = (infer.CustomDiff[ShellInputs, ShellOutputs])((*Shell)(nil))
var _ = (infer.CustomDelete[ShellOutputs])((*Shell)(nil))
var _ = (infer.CustomCheck[ShellInputs])((*Shell)(nil))

type ShellInputs struct {
	BaseInputs
	CommandInputs
	InstallCommands *[]string `pulumi:"installCommands"`
	ProgramName     *string   `pulumi:"programName"`
	DownloadUrl     *string   `pulumi:"downloadUrl"`
	VersionCommand  *string   `pulumi:"versionCommand,optional"`
	BinLocation     *string   `pulumi:"binLocation,optional"`
	Executable      *bool     `pulumi:"executable,optional"`
}

type ShellOutputs struct {
	ShellInputs
	CommandOutputs
	BaseOutputs
	Location *string `pulumi:"location"`
}

func (l *Shell) Diff(ctx p.Context, id string, olds ShellOutputs, news ShellInputs) (p.DiffResponse, error) {
	diff := map[string]p.PropertyDiff{}

	if *news.BinLocation != *olds.BinLocation {
		diff["binLocation"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	if *news.DownloadUrl != *olds.DownloadUrl {
		diff["downloadUrl"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	if news.Executable != olds.Executable {
		diff["executable"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	if news.InstallCommands != olds.InstallCommands {
		diff["installCommands"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	if news.UpdateCommands != olds.UpdateCommands {
		diff["updateCommands"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	if *news.ProgramName != *olds.ProgramName {
		diff["programName"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}

	return p.DiffResponse{
		DeleteBeforeReplace: true,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

// All resources must implement Create at a minumum.
func (l *Shell) Create(ctx p.Context, name string, input ShellInputs, preview bool) (string, ShellOutputs, error) {
	state := &ShellOutputs{
		ShellInputs: input,
	}
	if preview {
		return name, *state, nil
	}

	if err := state.createOrUpdate(ctx, input, *input.InstallCommands); err != nil {
		return "", *state, err
	}
	return name, *state, nil
}

func (l *Shell) Read(ctx p.Context, id string, inputs ShellInputs, state ShellOutputs) (
	canonicalID string, normalizedInputs ShellInputs, normalizedState ShellOutputs, err error) {

	if inputs.VersionCommand != nil {
		output, err := state.run(ctx, *inputs.VersionCommand, os.TempDir())
		if err != nil {
			return "", inputs, state, err
		}
		state.Version = &output
	}

	return *inputs.ProgramName, inputs, state, nil
}

func (l *Shell) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (ShellInputs, []p.CheckFailure, error) {
	fails := []p.CheckFailure{}
	if _, ok := newInputs["binLocation"]; !ok {
		home, err := os.UserHomeDir()
		if err != nil {
			fails = append(fails, p.CheckFailure{Property: "binLocation", Reason: err.Error()})
		} else {
			binLocation := path.Join(home, ".local", "bin")
			newInputs["binLocation"] = resource.NewStringProperty(binLocation)
		}
	}

	inputs, failures, err := infer.DefaultCheck[ShellInputs](newInputs)
	return inputs, append(failures, fails...), err
}

func (l *Shell) Update(ctx p.Context, name string, olds ShellOutputs, news ShellInputs, preview bool) (ShellOutputs, error) {
	state := &ShellOutputs{
		ShellInputs: news,
		BaseOutputs: BaseOutputs{
			Version: olds.Version,
		},
	}
	if preview {
		return *state, nil
	}
	var commands []string
	if news.UpdateCommands != nil {
		commands = *news.UpdateCommands
	} else {
		commands = *news.InstallCommands
	}
	if err := state.createOrUpdate(ctx, news, commands); err != nil {
		return *state, err
	}
	return *state, nil
}

func (l *Shell) Delete(ctx p.Context, id string, props ShellOutputs) error {
	if props.UninstallCommands != nil {
		_, err := props.run(ctx, strings.Join(*props.UninstallCommands, " && "), "")
		if err != nil {
			return err
		}

	}
	if err := os.Remove(*props.Location); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *ShellOutputs) createOrUpdate(ctx p.Context, input ShellInputs, commands []string) error {
	dir := os.TempDir()
	_, err := s.run(ctx, fmt.Sprintf("curl -OL %s", *input.DownloadUrl), dir)
	if err != nil {
		return err
	}
	_, err = s.run(ctx, strings.Join(commands, " && "), dir)
	if err != nil {
		return err
	}

	if input.Executable != nil && *input.Executable {
		target := path.Join(*input.BinLocation, *input.ProgramName)
		s.Location = &target
		// move it
		if err = os.Rename(path.Join(dir, *input.ProgramName), target); err != nil {
			return err
		}
		// make executable
		if err = os.Chmod(target, 0777); err != nil {
			return err
		}
	}

	if input.VersionCommand != nil {
		output, err := s.run(ctx, *input.VersionCommand, dir)
		if err != nil {
			return err
		}
		s.Version = &output
	} else {
		dv := "0.0.0"
		s.Version = &dv
	}
	return nil
}

func downloadTmpFile(downloadUrl string) (string, error) {
	fileURL, err := url.Parse(downloadUrl)
	p := fileURL.Path
	segments := strings.Split(p, "/")
	fileName := segments[len(segments)-1]
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	name := path.Join(os.TempDir(), fileName)
	fmt.Println("tmp file", name)
	out, err := os.Create(name)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println()
	return name, err
}
