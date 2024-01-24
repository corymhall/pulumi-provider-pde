package installers

import (
	// "errors"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"

	"github.com/pulumi/pulumi-command/provider/pkg/provider/util"
)

func (c *CommandOutputs) run(ctx p.Context, command, dir string) (string, error) {
	var args []string
	if c.Interpreter != nil && len(*c.Interpreter) > 0 {
		args = append(args, *c.Interpreter...)
	} else {
		if runtime.GOOS == "windows" {
			args = []string{"cmd", "/C"}
		} else {
			args = []string{"/bin/sh", "-c"}
		}
	}
	args = append(args, command)

	var err error
	var stdoutbuf, stderrbuf, stdouterrbuf bytes.Buffer
	stdouterrwriter := util.ConcurrentWriter{Writer: &stdouterrbuf}
	r, w := io.Pipe()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = io.MultiWriter(&stdoutbuf, &stdouterrwriter, w)
	cmd.Stderr = io.MultiWriter(&stderrbuf, &stdouterrwriter, w)
	cmd.Env = os.Environ()
	// if c.Environment != nil {
	// 	for k, v := range *c.Environment {
	// 		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	// 	}
	// }

	stdouterrch := make(chan struct{})
	go util.CopyOutput(ctx, r, stdouterrch, diag.Info)

	err = cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	w.Close()
	<-stdouterrch

	if err != nil {
		return "", fmt.Errorf("%w: running %q:\n%s", err, command, stdouterrbuf.String())
	}

	return strings.TrimSuffix(stdoutbuf.String(), "\n"), nil
	// return nil
}
