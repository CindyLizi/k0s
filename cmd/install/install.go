/*
Copyright 2021 k0s authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package install

import (
	"github.com/k0sproject/k0s/cmd/internal"
	"github.com/k0sproject/k0s/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type installFlags struct {
	force   bool
	envVars []string
}

func NewInstallCmd() *cobra.Command {
	var (
		debugFlags   internal.DebugFlags
		installFlags installFlags
	)

	cmd := &cobra.Command{
		Use:              "install",
		Short:            "Install k0s on a brand-new system. Must be run as root (or with sudo)",
		Args:             cobra.NoArgs,
		PersistentPreRun: debugFlags.Run,
		RunE:             func(*cobra.Command, []string) error { return pflag.ErrHelp }, // Enforce arg validation
	}

	pflags := cmd.PersistentFlags()
	debugFlags.AddToFlagSet(pflags)
	config.GetPersistentFlagSet().VisitAll(func(f *pflag.Flag) {
		f.Hidden = true
		f.Deprecated = "it has no effect and will be removed in a future release"
		pflags.AddFlag(f)
	})
	pflags.BoolVar(&installFlags.force, "force", false, "force init script creation")
	pflags.StringArrayVarP(&installFlags.envVars, "env", "e", nil, "set environment variable")

	cmd.AddCommand(installWorkerCmd(&installFlags))
	addPlatformSpecificCommands(cmd, &installFlags)

	return cmd
}
