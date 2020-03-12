// Copyright 2020 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/okteto/okteto/pkg/errors"
	"github.com/okteto/okteto/pkg/model"
	"go.undefinedlabs.com/scopeagent"
)

func Test_isWatchesConfigurationTooLow(t *testing.T) {
	var tests = []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "too-low",
			value:    "2",
			expected: true,
		},
		{
			name:     "too-low-trim",
			value:    "2\n",
			expected: true,
		},
		{
			name:     "ok",
			value:    "20000",
			expected: false,
		},
		{
			name:     "ok-trim",
			value:    "20000\n",
			expected: false,
		},
		{
			name:     "wrong",
			value:    "2a4d",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isWatchesConfigurationTooLow(tt.value)
			if tt.expected != result {
				t.Errorf("expected %t got %t in test %s", tt.expected, result, tt.name)
			}
		})
	}
}

func Test_loadDevOrDefault(t *testing.T) {
	var tests = []struct {
		name       string
		deployment string
		expectErr  bool
		dev        *model.Dev
	}{
		{
			name:       "default",
			deployment: "default-deployment",
			expectErr:  false,
		},
		{
			name:       "default-no-name",
			deployment: "",
			expectErr:  true,
		},
		{
			name:       "load-dev",
			deployment: "test-deployment",
			expectErr:  false,
			dev: &model.Dev{
				Name:  "loaded",
				Image: "okteto/test:1.0",
			},
		},
	}

	test := scopeagent.GetTest(t)
	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			def, err := loadDevOrDefault("/tmp/a-path", tt.deployment)
			if tt.expectErr {
				if err == nil {
					t.Fatal("expected error when loading")
				}

				if !errors.IsNotExist(err) {
					t.Fatalf("expected not found got: %s", err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if def.Name != tt.deployment {
				t.Errorf("expected default name, got %s", def.Name)
			}

			if tt.dev == nil {
				return
			}

			f, err := ioutil.TempFile("", "")
			if err != nil {
				t.Fatal(err)
			}
			f.Close()
			defer os.Remove(f.Name())

			if err := saveManifest(tt.dev, f.Name()); err != nil {
				t.Fatal(err)
			}

			loaded, err := loadDevOrDefault(f.Name(), "foo")
			if err != nil {
				t.Fatal("expected error when loading existing manifest")
			}

			if tt.dev.Image != loaded.Image {
				t.Fatalf("expected %s got %s", tt.dev.Image, loaded.Image)
			}

			if tt.dev.Name != loaded.Name {
				t.Fatalf("expected %s got %s", tt.dev.Name, loaded.Name)
			}

		})
	}
	name := "demo-deployment"
	def, err := loadDevOrDefault("/tmp/bad-path", name)
	if err != nil {
		t.Fatal("default dev was not returned")
	}

	if def.Name != name {
		t.Errorf("expected %s, got %s", name, def.Name)
	}

	_, err = loadDevOrDefault("/tmp/bad-path", "")
	if err == nil {
		t.Error("expected error with empty deployment name")
	}

}
