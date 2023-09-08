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

func TestGitHubCommand(t *testing.T) {
	t.Parallel()
	cmd := provider()
	urn := urn("installers", "GitHubRepo", "repo")

	// The state that we expect a non-preview create to return
	//
	// We use this as the final expect for create and the old state during update
	absPath := absFolder("pulumi-provider-pde")
	defer os.Remove(absPath)
	t.Cleanup(func() {
		os.Remove(absPath)
	})
	// Run a create against an in-memory provider, assert it succeeded, and return the
	// created property map
	create := func(preview bool, props resource.PropertyValue) resource.PropertyMap {
		cResp, err := cmd.Check(p.CheckRequest{
			Urn: urn,
			News: resource.PropertyMap{
				"org":    resource.PropertyValue{V: "corymhall"},
				"repo":   resource.PropertyValue{V: "pulumi-provider-pde"},
				"branch": props,
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

	del := func(name resource.PropertyValue) {
		err := cmd.Delete(p.DeleteRequest{
			Urn: urn,
			Properties: resource.PropertyMap{
				"absFolderName": name,
				"version":       resource.PropertyValue{V: "abc"},
				"repo":          resource.PropertyValue{V: "pulumi-provider-pde"},
				"org":           resource.PropertyValue{V: "corymhall"},
			},
		})
		require.NoError(t, err)
	}

	update := func(preview bool, props resource.PropertyValue) resource.PropertyMap {
		olds := resource.PropertyMap{
			"org":           resource.PropertyValue{V: "corymhall"},
			"repo":          resource.PropertyValue{V: "pulumi-provider-pde"},
			"branch":        resource.PropertyValue{V: "testing"},
			"absFolderName": resource.PropertyValue{V: absPath},
			"folderName":    resource.PropertyValue{V: "pulumi-provider-pde"},
			"version":       resource.PropertyValue{V: "f9a0bfe30df2f36d677240f811468ec27ac78446"},
		}
		cResp, err := cmd.Check(p.CheckRequest{
			Urn: urn,
			News: resource.PropertyMap{
				"org":        resource.PropertyValue{V: "corymhall"},
				"repo":       resource.PropertyValue{V: "pulumi-provider-pde"},
				"branch":     resource.PropertyValue{V: "testing"},
				"folderName": props,
			},
			Olds: olds,
		})
		fmt.Println("foldername: ", props)
		require.NoError(t, err)
		resp, err := cmd.Update(p.UpdateRequest{
			ID:      "echo1234",
			Urn:     urn,
			Olds:    olds,
			News:    cResp.Inputs.Copy(),
			Preview: preview,
		})
		require.NoError(t, err)
		return resp.Properties
	}

	t.Run("create-preview", func(t *testing.T) {
		assert.Equal(t, resource.PropertyMap{
			"org":           resource.PropertyValue{V: "corymhall"},
			"repo":          resource.PropertyValue{V: "pulumi-provider-pde"},
			"absFolderName": resource.MakeComputed(resource.PropertyValue{V: absPath}),
			"folderName":    resource.PropertyValue{V: "pulumi-provider-pde"},
		}, create(true /* preview */, resource.PropertyValue{}))
	})

	t.Run("create-actual", func(t *testing.T) {
		assert.Equal(t, resource.PropertyMap{
			"org":           resource.PropertyValue{V: "corymhall"},
			"repo":          resource.PropertyValue{V: "pulumi-provider-pde"},
			"branch":        resource.PropertyValue{V: "testing"},
			"absFolderName": resource.PropertyValue{V: absPath},
			"folderName":    resource.PropertyValue{V: "pulumi-provider-pde"},
			"version":       resource.PropertyValue{V: "f9a0bfe30df2f36d677240f811468ec27ac78446"},
		},
			create(false /*preview*/, resource.PropertyValue{V: "testing"}),
		)
	})

	t.Run("update-preview", func(t *testing.T) {
		assert.Equal(t, resource.PropertyMap{
			"org":           resource.PropertyValue{V: "corymhall"},
			"repo":          resource.PropertyValue{V: "pulumi-provider-pde"},
			"absFolderName": resource.MakeComputed(resource.PropertyValue{V: absPath}),
			"version":       resource.MakeComputed(resource.PropertyValue{V: "f9a0bfe30df2f36d677240f811468ec27ac78446"}),
			"branch":        resource.PropertyValue{V: "testing"},
		}, update(true /*preview*/, resource.NewNullProperty()))
	})

	t.Run("update-replace-actual", func(t *testing.T) {
		assert.Equal(t, resource.PropertyMap{
			"org":           resource.PropertyValue{V: "corymhall"},
			"repo":          resource.PropertyValue{V: "pulumi-provider-pde"},
			"folderName":    resource.PropertyValue{V: "pulumi-provider-tmp"},
			"absFolderName": resource.PropertyValue{V: absFolder("pulumi-provider-tmp")},
			"version":       resource.PropertyValue{V: "f9a0bfe30df2f36d677240f811468ec27ac78446"},
			"branch":        resource.PropertyValue{V: "testing"},
		}, update(false /*preview*/, resource.PropertyValue{V: "pulumi-provider-tmp"}))
	})

	t.Run("delete-actual", func(t *testing.T) {
		folder := absFolder("pulumi-provider-tmp")
		del(resource.PropertyValue{V: folder})
		_, err := os.Lstat(folder)
		if err != nil && !os.IsNotExist(err) {
			t.Fatalf("file not cleaned up!")
		}
	})
}

func absFolder(name string) string {
	homeDir, _ := os.UserHomeDir()
	return path.Join(homeDir, name)
}
