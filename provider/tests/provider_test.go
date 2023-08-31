package tests

import (
	"fmt"
	"os"
	"testing"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/blang/semver"
	pde "github.com/corymhall/pulumi-provider-pde/provider/pkg/provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func provider() integration.Server {
	v := semver.MustParse(pde.Version)
	return integration.NewServer(pde.Name, v, pde.NewProvider())

}
func urn(mod, res, name string) resource.URN {
	m := tokens.ModuleName(mod)
	r := tokens.TypeName(res)
	if !tokens.IsQName(name) {
		panic(fmt.Sprintf("invalid resource name: %q", name))
	}
	n := tokens.QName(name)
	return resource.NewURN("test", "command", "",
		tokens.NewTypeToken(
			tokens.NewModuleToken(pde.Name, m),
			r),
		n)
}

func defaultState(fromFile, toFile string, makeComputed bool, overwrites resource.PropertyMap) resource.PropertyMap {
	d := resource.PropertyMap{
		"source": resource.PropertyValue{V: fromFile},
		"target": resource.PropertyValue{V: toFile},
		"targets": resource.PropertyValue{
			V: []resource.PropertyValue{{V: toFile}},
		},
		"linked":    resource.PropertyValue{V: true},
		"is_dir":    resource.PropertyValue{V: false},
		"overwrite": resource.PropertyValue{V: false},
	}
	if makeComputed {
		d["linked"] = resource.PropertyValue{V: false}
		d["exists"] = resource.PropertyValue{V: false}
		for pk, v := range d {
			d[pk] = resource.MakeComputed(v)
		}
	}
	for pk, v := range overwrites {
		d[pk] = v
	}
	return d
}

func TestLinkCommand(t *testing.T) {
	t.Parallel()
	cmd := provider()
	urn := urn("local", "Link", "file")

	// The state that we expect a non-preview create to return
	//
	// We use this as the final expect for create and the old state during update
	file, _ := os.CreateTemp(os.TempDir(), "abc")
	toFile := "/tmp/xyz"
	t.Cleanup(func() {
		os.Remove(file.Name())
		os.Remove(toFile)
	})
	// Run a create against an in-memory provider, assert it succeeded, and return the
	// created property map
	create := func(preview bool, props resource.PropertyValue) resource.PropertyMap {
		resp, err := cmd.Create(p.CreateRequest{
			Urn: urn,
			Properties: resource.PropertyMap{
				"source":    resource.PropertyValue{V: file.Name()},
				"target":    resource.PropertyValue{V: toFile},
				"overwrite": props,
			},
			Preview: preview,
		})
		require.NoError(t, err)
		return resp.Properties
	}

	del := func() {
		err := cmd.Delete(p.DeleteRequest{
			Urn:        urn,
			Properties: defaultState(file.Name(), toFile, false, nil),
		})
		require.NoError(t, err)
	}

	update := func(preview bool, props resource.PropertyValue) resource.PropertyMap {
		resp, err := cmd.Update(p.UpdateRequest{
			ID:   "echo1234",
			Urn:  urn,
			Olds: defaultState(file.Name(), toFile, false, resource.PropertyMap{}),
			News: resource.PropertyMap{
				"source":    resource.PropertyValue{V: file.Name()},
				"target":    resource.PropertyValue{V: toFile},
				"overwrite": props,
			},
		})
		require.NoError(t, err)
		return resp.Properties
	}

	t.Run("create-preview", func(t *testing.T) {
		assert.Equal(t, resource.PropertyMap{
			"source":    resource.PropertyValue{V: file.Name()},
			"target":    resource.PropertyValue{V: toFile},
			"overwrite": resource.MakeComputed(resource.PropertyValue{V: false}),
			"targets": resource.MakeComputed(resource.PropertyValue{
				V: []resource.PropertyValue{},
			}),
		}, create(true /* preview */, resource.MakeComputed(resource.PropertyValue{V: false})))
	})

	t.Run("create-actual", func(t *testing.T) {
		assert.Equal(t, defaultState(file.Name(), toFile, false, resource.PropertyMap{
			"source":    resource.PropertyValue{V: file.Name()},
			"target":    resource.PropertyValue{V: toFile},
			"overwrite": resource.PropertyValue{V: false},
		}),
			create(false /*preview*/, resource.PropertyValue{V: false}),
		)
	})

	t.Run("update-preview", func(t *testing.T) {
		assert.Equal(t, defaultState(file.Name(), toFile, false, resource.PropertyMap{
			"source":    resource.PropertyValue{V: file.Name()},
			"target":    resource.PropertyValue{V: toFile},
			"overwrite": resource.PropertyValue{V: true},
			"targets": resource.PropertyValue{
				V: []resource.PropertyValue{{V: toFile}},
			},
		}), update(true /*preview*/, resource.PropertyValue{V: true}))
	})

	t.Run("update-no-op", func(t *testing.T) {
		assert.Equal(t, defaultState(file.Name(), toFile, false, resource.PropertyMap{
			"source":    resource.PropertyValue{V: file.Name()},
			"target":    resource.PropertyValue{V: toFile},
			"overwrite": resource.PropertyValue{V: true},
		}), update(false /*preview*/, resource.PropertyValue{V: true}))
	})

	t.Run("delete-actual", func(t *testing.T) {
		del()
		_, err := os.Lstat("/tmp/xyz")
		if !os.IsNotExist(err) {
			t.Fatalf("file not cleaned up!")
		}
	})
}
