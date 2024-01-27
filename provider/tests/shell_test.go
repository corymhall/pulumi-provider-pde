package tests

import (
	"testing"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellCommand(t *testing.T) {
	t.Parallel()
	cmd := provider()
	urn := urn("installers", "Shell")

	cases := []struct {
		name       string
		execute    func(t *testing.T, preview bool, props resource.PropertyValue) resource.PropertyMap
		expected   resource.PropertyMap
		inputProps resource.PropertyValue
		preview    bool
	}{
		{
			name: "create-preview",
			execute: func(t *testing.T, preview bool, props resource.PropertyValue) resource.PropertyMap {
				t.Helper()
				resp, err := cmd.Create(p.CreateRequest{
					Urn: urn,
					Properties: resource.PropertyMap{
						"installCommands": resource.PropertyValue{V: []resource.PropertyValue{}},
						"programName":     resource.PropertyValue{V: "cht.sh"},
						"downloadURL":     resource.PropertyValue{V: "https://cht.sh/:cht.sh"},
					},
					Preview: preview,
				})
				require.NoError(t, err)
				return resp.Properties
			},
			preview: true,
			expected: resource.PropertyMap{
				"installCommands": resource.PropertyValue{V: []resource.PropertyValue{}},
				"programName":     resource.PropertyValue{V: "cht.sh"},
				"downloadURL":     resource.PropertyValue{V: "https://cht.sh/:cht.sh"},
			},
		},
		{
			name: "create-actual",
			execute: func(t *testing.T, preview bool, props resource.PropertyValue) resource.PropertyMap {
				t.Helper()
				resp, err := cmd.Create(p.CreateRequest{
					Urn: urn,
					Properties: resource.PropertyMap{
						"installCommands": resource.PropertyValue{V: []resource.PropertyValue{}},
						"programName":     resource.PropertyValue{V: "cht.sh"},
						"downloadURL":     resource.PropertyValue{V: "https://cht.sh/:cht.sh"},
					},
					Preview: preview,
				})
				require.NoError(t, err)
				return resp.Properties
			},
			preview: false,
			expected: resource.PropertyMap{
				"installCommands": resource.PropertyValue{V: []resource.PropertyValue{}},
				"programName":     resource.PropertyValue{V: "cht.sh"},
				"downloadURL":     resource.PropertyValue{V: "https://cht.sh/:cht.sh"},
				"version":         resource.PropertyValue{V: "0.0.0"},
			},
		},
		{
			name: "update-preview",
			execute: func(t *testing.T, preview bool, props resource.PropertyValue) resource.PropertyMap {
				t.Helper()
				olds := resource.PropertyMap{
					"binLocation":     resource.PropertyValue{V: "/usr/local/bin"},
					"installCommands": resource.PropertyValue{V: []resource.PropertyValue{}},
					"programName":     resource.PropertyValue{V: "cht.sh"},
					"downloadURL":     resource.PropertyValue{V: "https://cht.sh/:cht.sh"},
					"version":         resource.PropertyValue{V: "0.0.0"},
				}
				news := resource.PropertyMap{
					"binLocation": resource.PropertyValue{V: "/usr/local/bin"},
					"installCommands": resource.PropertyValue{V: []resource.PropertyValue{
						{V: "echo 'hello world'"},
					}},
					"programName": resource.PropertyValue{V: "cht.sh"},
					"downloadURL": resource.PropertyValue{V: "https://cht.sh/:cht.sh"},
				}
				dResp, err := cmd.Diff(p.DiffRequest{
					Urn:  urn,
					Olds: olds,
					News: news,
				})
				assert.Equal(t, dResp.HasChanges, true)
				resp, err := cmd.Update(p.UpdateRequest{
					Urn:     urn,
					Olds:    olds,
					News:    news,
					Preview: preview,
				})
				require.NoError(t, err)
				return resp.Properties
			},
			preview: true,
			expected: resource.PropertyMap{
				"binLocation": resource.PropertyValue{V: "/usr/local/bin"},
				"installCommands": resource.PropertyValue{V: []resource.PropertyValue{
					{V: "echo 'hello world'"},
				}},
				"programName": resource.PropertyValue{V: "cht.sh"},
				"downloadURL": resource.PropertyValue{V: "https://cht.sh/:cht.sh"},
				"version":     resource.MakeComputed(resource.PropertyValue{V: "0.0.0"}),
			},
		},
		{
			name: "update-actual",
			execute: func(t *testing.T, preview bool, props resource.PropertyValue) resource.PropertyMap {
				t.Helper()
				olds := resource.PropertyMap{
					"binLocation":     resource.PropertyValue{V: "/usr/local/bin"},
					"installCommands": resource.PropertyValue{V: []resource.PropertyValue{}},
					"programName":     resource.PropertyValue{V: "cht.sh"},
					"downloadURL":     resource.PropertyValue{V: "https://cht.sh/:cht.sh"},
					"version":         resource.PropertyValue{V: "0.0.0"},
				}
				news := resource.PropertyMap{
					"binLocation": resource.PropertyValue{V: "/usr/local/bin"},
					"installCommands": resource.PropertyValue{V: []resource.PropertyValue{
						{V: "echo 'hello world'"},
					}},
					"programName": resource.PropertyValue{V: "cht.sh"},
					"downloadURL": resource.PropertyValue{V: "https://cht.sh/:cht.sh"},
				}
				dResp, err := cmd.Diff(p.DiffRequest{
					Urn:  urn,
					Olds: olds,
					News: news,
				})
				assert.Equal(t, dResp.HasChanges, true)
				resp, err := cmd.Update(p.UpdateRequest{
					Urn:     urn,
					Olds:    olds,
					News:    news,
					Preview: preview,
				})
				require.NoError(t, err)
				return resp.Properties
			},
			preview: false,
			expected: resource.PropertyMap{
				"binLocation": resource.PropertyValue{V: "/usr/local/bin"},
				"installCommands": resource.PropertyValue{V: []resource.PropertyValue{
					{V: "echo 'hello world'"},
				}},
				"programName": resource.PropertyValue{V: "cht.sh"},
				"downloadURL": resource.PropertyValue{V: "https://cht.sh/:cht.sh"},
				"version":     resource.PropertyValue{V: "0.0.0"},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t,
				tc.expected,
				tc.execute(t, tc.preview, tc.inputProps),
			)
		})
	}
}
