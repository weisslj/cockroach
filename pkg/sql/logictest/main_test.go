// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package logictest

import (
	"os"
	"testing"

	"github.com/weisslj/cockroach/pkg/security"
	"github.com/weisslj/cockroach/pkg/security/securitytest"
	"github.com/weisslj/cockroach/pkg/server"
	_ "github.com/weisslj/cockroach/pkg/sql/tests"
	"github.com/weisslj/cockroach/pkg/testutils/serverutils"
	"github.com/weisslj/cockroach/pkg/testutils/testcluster"
	"github.com/weisslj/cockroach/pkg/util/randutil"
)

//go:generate ../../util/leaktest/add-leaktest.sh *_test.go

func TestMain(m *testing.M) {
	security.SetAssetLoader(securitytest.EmbeddedAssets)
	randutil.SeedForTests()
	serverutils.InitTestServerFactory(server.TestServerFactory)
	serverutils.InitTestClusterFactory(testcluster.TestClusterFactory)
	os.Exit(m.Run())
}
