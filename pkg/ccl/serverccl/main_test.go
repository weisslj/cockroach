// Copyright 2018 The Cockroach Authors.
//
// Licensed as a CockroachDB Enterprise file under the Cockroach Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//     https://github.com/weisslj/cockroach/blob/master/licenses/CCL.txt

package serverccl

import (
	"os"
	"testing"

	_ "github.com/weisslj/cockroach/pkg/ccl"
	"github.com/weisslj/cockroach/pkg/ccl/utilccl"
	"github.com/weisslj/cockroach/pkg/security"
	"github.com/weisslj/cockroach/pkg/security/securitytest"
	"github.com/weisslj/cockroach/pkg/server"
	"github.com/weisslj/cockroach/pkg/testutils/serverutils"
	"github.com/weisslj/cockroach/pkg/testutils/testcluster"
	"github.com/weisslj/cockroach/pkg/util/randutil"
)

func TestMain(m *testing.M) {
	defer utilccl.TestingEnableEnterprise()()
	security.SetAssetLoader(securitytest.EmbeddedAssets)
	randutil.SeedForTests()
	serverutils.InitTestServerFactory(server.TestServerFactory)
	serverutils.InitTestClusterFactory(testcluster.TestClusterFactory)
	os.Exit(m.Run())
}

//go:generate ../../util/leaktest/add-leaktest.sh *_test.go
