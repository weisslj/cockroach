// Copyright 2017 The Cockroach Authors.
//
// Licensed as a CockroachDB Enterprise file under the Cockroach Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//     https://github.com/weisslj/cockroach/blob/master/licenses/CCL.txt

package cliccl_test

import (
	"os"
	"testing"

	"github.com/weisslj/cockroach/pkg/build"
	"github.com/weisslj/cockroach/pkg/ccl/utilccl"
	"github.com/weisslj/cockroach/pkg/server"
	"github.com/weisslj/cockroach/pkg/testutils/serverutils"
)

func TestMain(m *testing.M) {
	// CLI tests are sensitive to the server version, but test binaries don't have
	// a version injected. Pretend to be a very up-to-date version.
	defer build.TestingOverrideTag("v999.0.0")()

	defer utilccl.TestingEnableEnterprise()()
	serverutils.InitTestServerFactory(server.TestServerFactory)
	os.Exit(m.Run())
}

//go:generate ../../util/leaktest/add-leaktest.sh *_test.go
