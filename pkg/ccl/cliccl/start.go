// Copyright 2017 The Cockroach Authors.
//
// Licensed as a CockroachDB Enterprise file under the Cockroach Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//     https://github.com/weisslj/cockroach/blob/master/licenses/CCL.txt

package cliccl

import (
	"github.com/weisslj/cockroach/pkg/ccl/baseccl"
	"github.com/weisslj/cockroach/pkg/ccl/cliccl/cliflagsccl"
	"github.com/weisslj/cockroach/pkg/cli"
	"github.com/spf13/cobra"
)

// This does not define a `start` command, only modifications to the existing command
// in `pkg/cli/start.go`.

var storeEncryptionSpecs baseccl.StoreEncryptionSpecList

func init() {
	cli.VarFlag(cli.StartCmd.Flags(), &storeEncryptionSpecs, cliflagsccl.EnterpriseEncryption)

	// Add a new pre-run command to match encryption specs to store specs.
	cli.AddPersistentPreRunE(cli.StartCmd, func(cmd *cobra.Command, _ []string) error {
		return populateStoreSpecsEncryption()
	})
}

// populateStoreSpecsEncryption is a PreRun hook that matches store encryption specs with the
// parsed stores and populates some fields in the StoreSpec.
func populateStoreSpecsEncryption() error {
	return baseccl.PopulateStoreSpecWithEncryption(cli.GetServerCfgStores(), storeEncryptionSpecs)
}
