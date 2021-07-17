// Copyright 2021 The Okteto Authors
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

package preview

import (
	"context"

	"github.com/spf13/cobra"
)

//Preview preview management commands
func Preview(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preview",
		Short: "Preview environment management commands",
	}
	cmd.AddCommand(Deploy(ctx))
	cmd.AddCommand(Destroy(ctx))
	cmd.AddCommand(List(ctx))
	return cmd
}
