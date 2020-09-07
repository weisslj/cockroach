// Copyright 2018 The Cockroach Authors.
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

package testcat

import "github.com/weisslj/cockroach/pkg/sql/sem/tree"

// DropTable is a partial implementation of the DROP TABLE statement.
func (tc *Catalog) DropTable(stmt *tree.DropTable) {
	for i := range stmt.Names {
		tn := &stmt.Names[i]

		// Update the table name to include catalog and schema if not provided.
		tc.qualifyTableName(tn)

		// Ensure that table with that name exists.
		tc.Table(tn)

		// Remove the table from the catalog.
		delete(tc.dataSources, tn.FQString())
	}
}
