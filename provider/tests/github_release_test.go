package tests

import (
	"fmt"
	"os"
	"path"
	"testing"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitHubReleaseCommand(t *testing.T) {
	t.Parallel()
	cmd := provider()
	urn := urn("installers", "GitHubRelease", "repo")

	// The state that we expect a non-preview create to return
	//
	// We use this as the final expect for create and the old state during update
	bin := path.Join(os.TempDir(), ".bin")
	os.MkdirAll(bin, 0777)
	defer os.Remove(bin)
	t.Cleanup(func() {
		os.Remove(bin)
	})

	locations := resource.PropertyValue{
		V: []resource.PropertyValue{
			{V: path.Join(bin, "pulumi")},
			{V: path.Join(bin, "pulumi-analyzer-policy")},
			{V: path.Join(bin, "pulumi-analyzer-policy-python")},
			{V: path.Join(bin, "pulumi-language-dotnet")},
			{V: path.Join(bin, "pulumi-language-go")},
			{V: path.Join(bin, "pulumi-language-java")},
			{V: path.Join(bin, "pulumi-language-nodejs")},
			{V: path.Join(bin, "pulumi-language-python")},
			{V: path.Join(bin, "pulumi-language-python-exec")},
			{V: path.Join(bin, "pulumi-language-yaml")},
			{V: path.Join(bin, "pulumi-resource-pulumi-nodejs")},
			{V: path.Join(bin, "pulumi-resource-pulumi-python")},
			{V: path.Join(bin, "pulumi-watch")},
		},
	}
	base := resource.PropertyMap{
		"org":         resource.PropertyValue{V: "pulumi"},
		"repo":        resource.PropertyValue{V: "pulumi"},
		"binLocation": resource.PropertyValue{V: bin},
		"binFolder":   resource.PropertyValue{V: "pulumi"},
		"version":     resource.PropertyValue{V: "v3.81.0"},
		"downloadURL": resource.PropertyValue{V: "https://github.com/pulumi/pulumi/releases/download/v3.81.0/pulumi-v3.81.0-darwin-arm64.tar.gz"},
	}
	// Run a create against an in-memory provider, assert it succeeded, and return the
	// created property map
	create := func(preview bool, props resource.PropertyValue) resource.PropertyMap {
		cResp, err := cmd.Check(p.CheckRequest{
			Urn: urn,
			News: resource.PropertyMap{
				"org":         resource.PropertyValue{V: "pulumi"},
				"repo":        resource.PropertyValue{V: "pulumi"},
				"binLocation": resource.PropertyValue{V: bin},
				"binFolder":   resource.PropertyValue{V: "pulumi"},
				"version":     props,
			},
		})
		require.NoError(t, err)
		resp, err := cmd.Create(p.CreateRequest{
			Urn:        urn,
			Properties: cResp.Inputs.Copy(),
			Preview:    preview,
		})
		require.NoError(t, err)
		return resp.Properties
	}

	del := func(location resource.PropertyValue) {
		err := cmd.Delete(p.DeleteRequest{
			Urn: urn,
			Properties: resource.PropertyMap{
				"downloadURL": resource.PropertyValue{V: "https://github.com/pulumi/pulumi/releases/download/v3.81.0/pulumi-v3.81.0-darwin-arm64.tar.gz"},
				"org":         resource.PropertyValue{V: "pulumi"},
				"repo":        resource.PropertyValue{V: "pulumi"},
				"binLocation": resource.PropertyValue{V: bin},
				"version":     resource.PropertyValue{V: "abc"},
				"locations":   location,
			},
		})
		require.NoError(t, err)
	}

	update := func(preview bool, version resource.PropertyValue) resource.PropertyMap {
		cResp, err := cmd.Check(p.CheckRequest{
			Urn: urn,
			News: resource.PropertyMap{
				"org":         resource.PropertyValue{V: "pulumi"},
				"repo":        resource.PropertyValue{V: "pulumi"},
				"binLocation": resource.PropertyValue{V: bin},
				"binFolder":   resource.PropertyValue{V: "pulumi"},
				"version":     version,
			},
			Olds: base,
		})
		require.NoError(t, err)
		resp, err := cmd.Update(p.UpdateRequest{
			ID:      "echo1234",
			Urn:     urn,
			Olds:    base,
			News:    cResp.Inputs.Copy(),
			Preview: preview,
		})
		require.NoError(t, err)
		return resp.Properties
	}

	t.Run("create-preview", func(t *testing.T) {
		assert.Equal(t, resource.PropertyMap{
			"org":         resource.PropertyValue{V: "pulumi"},
			"repo":        resource.PropertyValue{V: "pulumi"},
			"version":     resource.PropertyValue{V: "v3.81.0"},
			"assetName":   resource.PropertyValue{V: "pulumi-v3.81.0-darwin-arm64.tar.gz"},
			"binLocation": resource.PropertyValue{V: bin},
			"binFolder":   resource.PropertyValue{V: "pulumi"},
			"downloadURL": resource.MakeComputed(resource.PropertyValue{V: "https://github.com/pulumi/pulumi/releases/download/v3.81.0/pulumi-v3.81.0-darwin-arm64.tar.gz"}),
		}, create(true /* preview */, resource.PropertyValue{V: "v3.81.0"}))
	})

	t.Run("create-actual", func(t *testing.T) {
		assert.Equal(t, expectedProps(bin, "v3.81.0", resource.PropertyMap{
			"downloadURL": resource.PropertyValue{V: "https://github.com/pulumi/pulumi/releases/download/v3.81.0/pulumi-v3.81.0-darwin-arm64.tar.gz"},
			"locations":   locations,
		}),
			create(false /*preview*/, resource.PropertyValue{V: "v3.81.0"}),
		)
	})
	//
	t.Run("update-preview", func(t *testing.T) {
		assert.Equal(t, expectedProps(bin, "v3.80.0", resource.PropertyMap{
			"version":     resource.PropertyValue{V: "v3.80.0"},
			"downloadURL": resource.MakeComputed(resource.PropertyValue{V: "https://github.com/pulumi/pulumi/releases/download/v3.80.0/pulumi-v3.80.0-darwin-arm64.tar.gz"}),
		}), update(true /*preview*/, resource.PropertyValue{V: "v3.80.0"}))
	})

	t.Run("update-replace-actual", func(t *testing.T) {
		assert.Equal(t, expectedProps(bin, "v3.80.0", resource.PropertyMap{
			"version":     resource.PropertyValue{V: "v3.80.0"},
			"downloadURL": resource.PropertyValue{V: "https://github.com/pulumi/pulumi/releases/download/v3.80.0/pulumi-v3.80.0-darwin-arm64.tar.gz"},
			"locations":   locations,
		}), update(false /*preview*/, resource.PropertyValue{V: "v3.80.0"}))
	})

	t.Run("delete-actual", func(t *testing.T) {
		locations := resource.PropertyValue{
			V: []resource.PropertyValue{
				{V: path.Join(bin, "pulumi")},
				{V: path.Join(bin, "pulumi-analyzer-policy")},
				{V: path.Join(bin, "pulumi-analyzer-policy-python")},
				{V: path.Join(bin, "pulumi-language-dotnet")},
				{V: path.Join(bin, "pulumi-language-go")},
				{V: path.Join(bin, "pulumi-language-java")},
				{V: path.Join(bin, "pulumi-language-nodejs")},
				{V: path.Join(bin, "pulumi-language-python")},
				{V: path.Join(bin, "pulumi-language-python-exec")},
				{V: path.Join(bin, "pulumi-language-yaml")},
				{V: path.Join(bin, "pulumi-resource-pulumi-nodejs")},
				{V: path.Join(bin, "pulumi-resource-pulumi-python")},
				{V: path.Join(bin, "pulumi-watch")},
			},
		}
		del(locations)
	})
}

func expectedProps(
	bin, version string,
	overrides resource.PropertyMap,
) resource.PropertyMap {
	base := resource.PropertyMap{
		"org":         resource.PropertyValue{V: "pulumi"},
		"repo":        resource.PropertyValue{V: "pulumi"},
		"version":     resource.PropertyValue{V: version},
		"assetName":   resource.PropertyValue{V: fmt.Sprintf("pulumi-%s-darwin-arm64.tar.gz", version)},
		"binLocation": resource.PropertyValue{V: bin},
		"binFolder":   resource.PropertyValue{V: "pulumi"},
	}
	for pk, v := range overrides {
		base[pk] = v
	}
	return base
}
