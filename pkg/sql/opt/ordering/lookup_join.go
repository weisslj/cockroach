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

package ordering

import (
	"github.com/weisslj/cockroach/pkg/sql/opt"
	"github.com/weisslj/cockroach/pkg/sql/opt/memo"
	"github.com/weisslj/cockroach/pkg/sql/opt/props"
	"github.com/weisslj/cockroach/pkg/sql/opt/props/physical"
)

func lookupOrIndexJoinCanProvideOrdering(
	expr memo.RelExpr, required *physical.OrderingChoice,
) bool {
	// LookupJoin and IndexJoin can pass through their ordering if the ordering
	// depends only on columns present in the input.
	return isOrderingBoundBy(expr.Child(0).(memo.RelExpr), required)
}

func lookupOrIndexJoinBuildChildReqOrdering(
	parent memo.RelExpr, required *physical.OrderingChoice, childIdx int,
) physical.OrderingChoice {
	if childIdx != 0 {
		return physical.OrderingChoice{}
	}

	// We may need to remove ordering columns that are not output by the input
	// expression.
	child := parent.Child(0).(memo.RelExpr)
	res := projectOrderingToInput(child, required)
	// It is in principle possible that the lookup join has an ON condition that
	// forces an equality on two columns in the input. In this case we need to
	// trim the column groups to keep the ordering valid w.r.t the child FDs
	// (similar to Select).
	//
	// This case indicates that we didn't do a good job pushing down equalities
	// (see #36219), but it should be handled correctly here nevertheless.
	return trimColumnGroups(&res, &child.Relational().FuncDeps)
}

func indexJoinBuildProvided(expr memo.RelExpr, required *physical.OrderingChoice) opt.Ordering {
	// If an index join has a requirement on some input columns, those columns
	// must be output columns (or equivalent to them). We may still need to remap
	// using column equivalencies.
	indexJoin := expr.(*memo.IndexJoinExpr)
	rel := indexJoin.Relational()
	return remapProvided(indexJoin.Input.ProvidedPhysical().Ordering, &rel.FuncDeps, rel.OutputCols)
}

func lookupJoinBuildProvided(expr memo.RelExpr, required *physical.OrderingChoice) opt.Ordering {
	lookupJoin := expr.(*memo.LookupJoinExpr)
	childProvided := lookupJoin.Input.ProvidedPhysical().Ordering

	// The lookup join includes an implicit projection (lookupJoin.Cols); some of
	// the input columns might not be output columns so we may need to remap them.
	// First check if we need to.
	needsRemap := false
	for i := range childProvided {
		if !lookupJoin.Cols.Contains(int(childProvided[i].ID())) {
			needsRemap = true
			break
		}
	}
	if !needsRemap {
		// Fast path: we don't need to remap any columns.
		return childProvided
	}

	// Because of the implicit projection, the FDs of the LookupJoin don't include
	// all the columns we care about; we have to recreate the FDs of the join
	// before the projection. These are the FDs of the input plus the equality
	// constraints implied by the lookup join.
	var fds props.FuncDepSet
	fds.CopyFrom(&lookupJoin.Input.Relational().FuncDeps)

	md := lookupJoin.Memo().Metadata()
	index := md.Table(lookupJoin.Table).Index(lookupJoin.Index)
	for i, colID := range lookupJoin.KeyCols {
		indexColID := lookupJoin.Table.ColumnID(index.Column(i).Ordinal)
		fds.AddEquivalency(colID, indexColID)
	}

	return remapProvided(childProvided, &fds, lookupJoin.Cols)
}
