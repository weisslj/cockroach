// Copyright 2016 The Cockroach Authors.
//
// Licensed as a CockroachDB Enterprise file under the Cockroach Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//     https://github.com/weisslj/cockroach/blob/master/licenses/CCL.txt

package buildccl

import "github.com/weisslj/cockroach/pkg/build"

func init() {
	build.Distribution = "CCL"
}
