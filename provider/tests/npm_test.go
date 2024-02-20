package tests

import (
	"os"
	"path"
	"testing"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNpmCommand(t *testing.T) {
	t.Parallel()
	cmd := provider()
	urn := urn("installers", "Npm")

	loc := path.Join(os.TempDir(), "npm-packages")
	os.MkdirAll(loc, 0777)
	t.Cleanup(func() {
		os.RemoveAll(loc)
	})
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
						"location": resource.PropertyValue{V: loc},
						"packages": resource.PropertyValue{V: []resource.PropertyValue{
							{V: "@cdk-cloudformation/alexa-ask-skill"},
						}},
					},
					Preview: preview,
				})
				require.NoError(t, err)
				return resp.Properties
			},
			preview: true,
			expected: resource.PropertyMap{
				"location": resource.PropertyValue{V: loc},
				"packages": resource.PropertyValue{V: []resource.PropertyValue{
					{V: "@cdk-cloudformation/alexa-ask-skill"},
				}},
				"deps": resource.MakeComputed(resource.PropertyValue{V: resource.PropertyMap{}}),
			},
		},
		{
			name: "create-actual",
			execute: func(t *testing.T, preview bool, props resource.PropertyValue) resource.PropertyMap {
				t.Helper()
				resp, err := cmd.Create(p.CreateRequest{
					Urn: urn,
					Properties: resource.PropertyMap{
						"location": resource.PropertyValue{V: loc},
						"packages": resource.PropertyValue{V: []resource.PropertyValue{
							{V: "@cdk-cloudformation/alexa-ask-skill"},
						}},
					},
					Preview: preview,
				})
				require.NoError(t, err)
				return resp.Properties
			},
			preview: false,
			expected: resource.PropertyMap{
				"location": resource.PropertyValue{V: loc},
				"packages": resource.PropertyValue{V: []resource.PropertyValue{
					{V: "@cdk-cloudformation/alexa-ask-skill"},
				}},
				"deps": resource.PropertyValue{V: resource.PropertyMap{
					"@cdk-cloudformation/alexa-ask-skill": resource.PropertyValue{V: "0.0.0-alpha.7"},
				}},
			},
		},
		{
			name: "update-preview-remove-packages",
			execute: func(t *testing.T, preview bool, props resource.PropertyValue) resource.PropertyMap {
				t.Helper()
				olds := resource.PropertyMap{
					"location": resource.PropertyValue{V: loc},
					"packages": resource.PropertyValue{V: []resource.PropertyValue{
						{V: "@cdk-cloudformation/alexa-ask-skill"},
					}},
					"deps": resource.PropertyValue{V: resource.PropertyMap{
						"@cdk-cloudformation/alexa-ask-skill": resource.PropertyValue{V: "0.0.0-alpha.7"},
					}},
				}
				news := resource.PropertyMap{
					"location": resource.PropertyValue{V: loc},
					"packages": resource.PropertyValue{V: []resource.PropertyValue{}},
				}
				dResp, err := cmd.Diff(p.DiffRequest{
					Urn:  urn,
					Olds: olds,
					News: news,
				})
				require.NoError(t, err)
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
				"location": resource.PropertyValue{V: loc},
				"packages": resource.PropertyValue{V: []resource.PropertyValue{}},
				"deps": resource.MakeComputed(resource.PropertyValue{V: resource.PropertyMap{
					"@cdk-cloudformation/alexa-ask-skill": resource.PropertyValue{V: "0.0.0-alpha.7"},
				}}),
			},
		},
		{
			name: "update-preview-add-packages",
			execute: func(t *testing.T, preview bool, props resource.PropertyValue) resource.PropertyMap {
				t.Helper()
				olds := resource.PropertyMap{
					"location": resource.PropertyValue{V: loc},
					"packages": resource.PropertyValue{V: []resource.PropertyValue{
						{V: "@cdk-cloudformation/alexa-ask-skill"},
					}},
					"deps": resource.PropertyValue{V: resource.PropertyMap{
						"@cdk-cloudformation/alexa-ask-skill": resource.PropertyValue{V: "0.0.0-alpha.7"},
					}},
				}
				news := resource.PropertyMap{
					"location": resource.PropertyValue{V: loc},
					"packages": resource.PropertyValue{V: []resource.PropertyValue{
						{V: "@cdk-cloudformation/alexa-ask-skill"},
						{V: "@cdk-cloudformation/registry-test-resource1-module"},
					}},
				}
				dResp, err := cmd.Diff(p.DiffRequest{
					Urn:  urn,
					Olds: olds,
					News: news,
				})
				require.NoError(t, err)
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
				"location": resource.PropertyValue{V: loc},
				"packages": resource.PropertyValue{V: []resource.PropertyValue{
					{V: "@cdk-cloudformation/alexa-ask-skill"},
					{V: "@cdk-cloudformation/registry-test-resource1-module"},
				}},
				"deps": resource.MakeComputed(resource.PropertyValue{V: resource.PropertyMap{
					"@cdk-cloudformation/alexa-ask-skill": resource.PropertyValue{V: "0.0.0-alpha.7"},
				}}),
			},
		},
		{
			name: "update-actual-remove-packages",
			execute: func(t *testing.T, preview bool, props resource.PropertyValue) resource.PropertyMap {
				t.Helper()
				olds := resource.PropertyMap{
					"location": resource.PropertyValue{V: loc},
					"packages": resource.PropertyValue{V: []resource.PropertyValue{
						{V: "@cdk-cloudformation/alexa-ask-skill"},
					}},
					"deps": resource.PropertyValue{V: resource.PropertyMap{
						"@cdk-cloudformation/alexa-ask-skill": resource.PropertyValue{V: "0.0.0-alpha.7"},
					}},
				}
				news := resource.PropertyMap{
					"location": resource.PropertyValue{V: loc},
					"packages": resource.PropertyValue{V: []resource.PropertyValue{}},
				}
				dResp, err := cmd.Diff(p.DiffRequest{
					Urn:  urn,
					Olds: olds,
					News: news,
				})
				require.NoError(t, err)
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
				"location": resource.PropertyValue{V: loc},
				"packages": resource.PropertyValue{V: []resource.PropertyValue{}},
				"deps": resource.PropertyValue{V: resource.PropertyMap{
					"@cdk-cloudformation/alexa-ask-skill": resource.PropertyValue{V: "latest"},
				}},
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
