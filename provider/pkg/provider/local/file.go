package local

import (
	"fmt"
	"os"
	"path"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type File struct{}

var _ = (infer.CustomDelete[FileState])((*File)(nil))
var _ = (infer.CustomCheck[FileArgs])((*File)(nil))
var _ = (infer.CustomUpdate[FileArgs, FileState])((*File)(nil))
var _ = (infer.CustomDiff[FileArgs, FileState])((*File)(nil))
var _ = (infer.ExplicitDependencies[FileArgs, FileState])((*File)(nil))
var _ = (infer.Annotated)((*File)(nil))
var _ = (infer.Annotated)((*FileArgs)(nil))
var _ = (infer.Annotated)((*FileState)(nil))

func (f *File) Annotate(a infer.Annotator) {
	a.Describe(&f, "A file projected into a pulumi resource")
}

type FileArgs struct {
	Path    string   `pulumi:"path"`
	Content []string `pulumi:"content"`
	Force   bool     `pulumi:"force,optional"`
}

func (f *FileArgs) Annotate(a infer.Annotator) {
	a.Describe(&f.Content, "The content of the file.")
	a.Describe(&f.Force, "If an already existing file should be deleted if it exists.")
	a.Describe(&f.Path, "The path of the file. This defaults to the name of the pulumi resource.")
}

type FileState struct {
	Path    string   `pulumi:"path"`
	Force   bool     `pulumi:"force"`
	Content []string `pulumi:"content"`
}

func (f *FileState) Annotate(a infer.Annotator) {
	a.Describe(&f.Content, "The content of the file.")
	a.Describe(&f.Force, "If an already existing file should be deleted if it exists.")
	a.Describe(&f.Path, "The path of the file.")
}

func (*File) Create(ctx p.Context, name string, input FileArgs, preview bool) (id string, output FileState, err error) {
	if !input.Force {
		_, err := os.Stat(input.Path)
		if !os.IsNotExist(err) {
			return "", FileState{}, fmt.Errorf("file already exists; pass force=true to override")
		}
	}
	contentString := strings.Join(input.Content, "\n")
	state := &FileState{
		Path:    input.Path,
		Force:   input.Force,
		Content: input.Content,
	}

	if preview { // Don't do the actual creating if in preview
		return name, *state, nil
	}

	_, err = os.Stat(path.Dir(input.Path))
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(input.Path), 0755); err != nil {
			return "", FileState{}, err
		}
	}

	f, err := os.Create(input.Path)
	if err != nil {
		return "", FileState{}, err
	}
	defer f.Close()
	n, err := f.WriteString(contentString)
	if err != nil {
		return "", FileState{}, err
	}
	if n != len([]byte(contentString)) {
		return "", FileState{}, fmt.Errorf("only wrote %d/%d bytes", n, len(contentString))
	}
	return name, *state, nil
}

func (*File) Delete(ctx p.Context, id string, props FileState) error {
	err := os.Remove(props.Path)
	if os.IsNotExist(err) {
		ctx.Logf(diag.Warning, "file %q already deleted", props.Path)
		err = nil
	}
	return err
}

func (*File) Check(ctx p.Context, name string, oldInputs, newInputs resource.PropertyMap) (FileArgs, []p.CheckFailure, error) {
	if _, ok := newInputs["path"]; !ok {
		newInputs["path"] = resource.NewStringProperty(name)
	}
	return infer.DefaultCheck[FileArgs](newInputs)
}

func (*File) Update(ctx p.Context, id string, olds FileState, news FileArgs, preview bool) (FileState, error) {
	state := &FileState{
		Path:    news.Path,
		Force:   news.Force,
		Content: news.Content,
	}
	if preview {
		return *state, nil
	}
	newContentString := strings.Join(news.Content, "\n")
	f, err := os.Create(news.Path)
	if err != nil {
		return FileState{}, err
	}
	defer f.Close()
	n, err := f.WriteString(newContentString)
	if err != nil {
		return FileState{}, err
	}
	if n != len([]byte(newContentString)) {
		return FileState{}, fmt.Errorf("only wrote %d/%d bytes", n, len(news.Content))
	}

	return *state, nil
}

func (*File) Diff(ctx p.Context, id string, olds FileState, news FileArgs) (p.DiffResponse, error) {
	diff := map[string]p.PropertyDiff{}
	byteContent, err := os.ReadFile(olds.Path)
	if err != nil {
		return p.DiffResponse{}, err
	}
	content := string(byteContent)
	newContentString := strings.Join(news.Content, "\n")
	oldContentString := strings.Join(olds.Content, "\n")
	if newContentString != oldContentString || newContentString != content {
		diff["content"] = p.PropertyDiff{Kind: p.Update}
	}
	if news.Force != olds.Force {
		diff["force"] = p.PropertyDiff{Kind: p.Update}
	}
	if news.Path != olds.Path {
		diff["path"] = p.PropertyDiff{Kind: p.UpdateReplace}
	}
	return p.DiffResponse{
		DeleteBeforeReplace: true,
		HasChanges:          len(diff) > 0,
		DetailedDiff:        diff,
	}, nil
}

func (*File) WireDependencies(f infer.FieldSelector, args *FileArgs, state *FileState) {
	f.OutputField(&state.Content).DependsOn(f.InputField(&args.Content))
	f.OutputField(&state.Force).DependsOn(f.InputField(&args.Force))
	f.OutputField(&state.Path).DependsOn(f.InputField(&args.Path))
}
