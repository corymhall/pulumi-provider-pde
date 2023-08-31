package installers

import (
	// "errors"
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
var _ = (infer.CustomDelete[ShellOutputs])((*Shell)(nil))
var _ = (infer.CustomCheck[ShellInputs])((*Shell)(nil))

type ShellInputs struct {
	BaseInputs
	CommandInputs
	VersionCommand  *string   `pulumi:"version_command"`
	InstallCommands *[]string `pulumi:"install_commands"`
	ProgramName     *string   `pulumi:"program_name"`
	DownloadUrl     *string   `pulumi:"download_url"`
	BinLocation     *string   `pulumi:"bin_location,optional"`
	Executable      *bool     `pulumi:"executable,optional"`
}

type ShellOutputs struct {
	ShellInputs
	CommandOutputs
	BaseOutputs
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

func (s *ShellOutputs) createOrUpdate(ctx p.Context, input ShellInputs, commands []string) error {
	file, err := downloadTmpFile(*input.DownloadUrl)
	defer os.Remove(file)
	if err != nil {
		return err
	}
	dir := os.TempDir()
	_, err = s.run(ctx, strings.Join(commands, " && "), dir)
	if err != nil {
		return err
	}
	var binLocation string
	if input.BinLocation != nil {
		binLocation = *input.BinLocation
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		binLocation = path.Join(home, ".local", "bin")
	}

	if input.Executable != nil && *input.Executable {
		target := path.Join(binLocation, *input.ProgramName)
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
	out, err := os.Create(name)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return name, err
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
	return infer.DefaultCheck[ShellInputs](newInputs)
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
	file := path.Join(*props.BinLocation, *props.ProgramName)
	if props.UninstallCommands != nil {
		_, err := props.run(ctx, strings.Join(*props.UninstallCommands, " && "), "")
		if err != nil {
			return err
		}

	}
	if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
