// Copyright 2019 The Cockroach Authors.
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

package distsqlrun

import (
	"context"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/weisslj/cockroach/pkg/base"
	"github.com/weisslj/cockroach/pkg/settings/cluster"
	"github.com/weisslj/cockroach/pkg/sql/distsqlpb"
	"github.com/weisslj/cockroach/pkg/sql/sem/tree"
	"github.com/weisslj/cockroach/pkg/sql/sqlbase"
	"github.com/weisslj/cockroach/pkg/storage/engine"
	"github.com/weisslj/cockroach/pkg/util/encoding"
	"github.com/weisslj/cockroach/pkg/util/leaktest"
	"github.com/weisslj/cockroach/pkg/util/mon"
)

const noFilterIdx = -1

func TestWindowerAccountingForResults(t *testing.T) {
	defer leaktest.AfterTest(t)()
	ctx := context.Background()
	st := cluster.MakeTestingClusterSettings()
	monitor := mon.MakeMonitorWithLimit(
		"test-monitor",
		mon.MemoryResource,
		100000,        /* limit */
		nil,           /* curCount */
		nil,           /* maxHist */
		5000,          /* increment */
		math.MaxInt64, /* noteworthy */
		st,
	)
	evalCtx := tree.MakeTestingEvalContextWithMon(st, &monitor)
	defer evalCtx.Stop(ctx)
	diskMonitor := makeTestDiskMonitor(ctx, st)
	defer diskMonitor.Stop(ctx)
	tempEngine, err := engine.NewTempEngine(base.DefaultTestTempStorageConfig(st), base.DefaultTestStoreSpec)
	if err != nil {
		t.Fatal(err)
	}
	defer tempEngine.Close()

	flowCtx := &FlowCtx{
		Settings:    st,
		EvalCtx:     &evalCtx,
		diskMonitor: diskMonitor,
		TempStorage: tempEngine,
	}

	post := &distsqlpb.PostProcessSpec{}
	input := NewRepeatableRowSource(sqlbase.OneIntCol, sqlbase.MakeIntRows(1000, 1))
	aggSpec := distsqlpb.AggregatorSpec_ARRAY_AGG
	spec := distsqlpb.WindowerSpec{
		PartitionBy: []uint32{},
		WindowFns: []distsqlpb.WindowerSpec_WindowFn{{
			Func:         distsqlpb.WindowerSpec_Func{AggregateFunc: &aggSpec},
			ArgIdxStart:  0,
			ArgCount:     1,
			Ordering:     distsqlpb.Ordering{Columns: []distsqlpb.Ordering_Column{{ColIdx: 0}}},
			FilterColIdx: noFilterIdx,
			Frame: &distsqlpb.WindowerSpec_Frame{
				Mode: distsqlpb.WindowerSpec_Frame_ROWS,
				Bounds: distsqlpb.WindowerSpec_Frame_Bounds{
					Start: distsqlpb.WindowerSpec_Frame_Bound{
						BoundType: distsqlpb.WindowerSpec_Frame_OFFSET_PRECEDING,
						IntOffset: 100,
					},
				},
			},
		}},
	}
	output := NewRowBuffer(
		sqlbase.OneIntCol, nil, RowBufferArgs{},
	)

	d, err := newWindower(flowCtx, 0 /* processorID */, &spec, input, post, output)
	if err != nil {
		t.Fatal(err)
	}
	d.Run(ctx)
	for {
		row, meta := output.Next()
		if row != nil {
			t.Fatalf("unexpectedly received row %+v", row)
		}
		if meta == nil {
			t.Fatalf("unexpectedly didn't receive an OOM error")
		}
		if meta.Err != nil {
			if !strings.Contains(meta.Err.Error(), "memory budget exceeded") {
				t.Fatalf("unexpectedly received an error other than OOM")
			}
			break
		}
	}
}

type windowTestSpec struct {
	// The column indices of PARTITION BY clause.
	partitionBy []uint32
	// Window function to be computed.
	windowFn windowFnTestSpec
}

type windowFnTestSpec struct {
	funcName       string
	argIdxStart    uint32
	argCount       uint32
	columnOrdering sqlbase.ColumnOrdering
}

func windows(windowTestSpecs []windowTestSpec) ([]distsqlpb.WindowerSpec, error) {
	windows := make([]distsqlpb.WindowerSpec, len(windowTestSpecs))
	for i, spec := range windowTestSpecs {
		windows[i].PartitionBy = spec.partitionBy
		windows[i].WindowFns = make([]distsqlpb.WindowerSpec_WindowFn, 1)
		windowFnSpec := distsqlpb.WindowerSpec_WindowFn{}
		fnSpec, err := CreateWindowerSpecFunc(spec.windowFn.funcName)
		if err != nil {
			return nil, err
		}
		windowFnSpec.Func = fnSpec
		windowFnSpec.ArgIdxStart = spec.windowFn.argIdxStart
		windowFnSpec.ArgCount = spec.windowFn.argCount
		if spec.windowFn.columnOrdering != nil {
			ordCols := make([]distsqlpb.Ordering_Column, 0, len(spec.windowFn.columnOrdering))
			for _, column := range spec.windowFn.columnOrdering {
				ordCols = append(ordCols, distsqlpb.Ordering_Column{
					ColIdx: uint32(column.ColIdx),
					// We need this -1 because encoding.Direction has extra value "_"
					// as zeroth "entry" which its proto equivalent doesn't have.
					Direction: distsqlpb.Ordering_Column_Direction(column.Direction - 1),
				})
			}
			windowFnSpec.Ordering = distsqlpb.Ordering{Columns: ordCols}
		}
		windowFnSpec.FilterColIdx = noFilterIdx
		windows[i].WindowFns[0] = windowFnSpec
	}
	return windows, nil
}

func BenchmarkWindower(b *testing.B) {
	const numCols = 3
	const numRows = 100000

	specs, err := windows([]windowTestSpec{
		{ // sum(@1) OVER ()
			partitionBy: []uint32{},
			windowFn: windowFnTestSpec{
				funcName:    "SUM",
				argIdxStart: uint32(0),
				argCount:    uint32(1),
			},
		},
		{ // sum(@1) OVER (ORDER BY @3)
			windowFn: windowFnTestSpec{
				funcName:       "SUM",
				argIdxStart:    uint32(0),
				argCount:       uint32(1),
				columnOrdering: sqlbase.ColumnOrdering{{ColIdx: 2, Direction: encoding.Ascending}},
			},
		},
		{ // sum(@1) OVER (PARTITION BY @2)
			partitionBy: []uint32{1},
			windowFn: windowFnTestSpec{
				funcName:    "SUM",
				argIdxStart: uint32(0),
				argCount:    uint32(1),
			},
		},
		{ // sum(@1) OVER (PARTITION BY @2 ORDER BY @3)
			partitionBy: []uint32{1},
			windowFn: windowFnTestSpec{
				funcName:       "SUM",
				argIdxStart:    uint32(0),
				argCount:       uint32(1),
				columnOrdering: sqlbase.ColumnOrdering{{ColIdx: 2, Direction: encoding.Ascending}},
			},
		},
	})
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	st := cluster.MakeTestingClusterSettings()
	evalCtx := tree.MakeTestingEvalContext(st)
	defer evalCtx.Stop(ctx)

	flowCtx := &FlowCtx{
		Settings: st,
		EvalCtx:  &evalCtx,
	}

	rowsGenerators := []func(int, int) sqlbase.EncDatumRows{
		sqlbase.MakeIntRows,
		func(numRows, numCols int) sqlbase.EncDatumRows {
			return sqlbase.MakeRepeatedIntRows(numRows/100, numRows, numCols)
		},
	}
	skipRepeatedSpecs := map[int]bool{0: true, 1: true}

	for j, rowsGenerator := range rowsGenerators {
		for i, spec := range specs {
			if skipRepeatedSpecs[i] && j == 1 {
				continue
			}
			runName := fmt.Sprintf("%s() OVER (", spec.WindowFns[0].Func.AggregateFunc.String())
			if len(spec.PartitionBy) > 0 {
				runName = runName + "PARTITION BY"
				if j == 0 {
					runName = runName + "/* SINGLE ROW PARTITIONS */"
				} else {
					runName = runName + "/* MULTIPLE ROWS PARTITIONS */"
				}
			}
			if len(spec.PartitionBy) > 0 && spec.WindowFns[0].Ordering.Columns != nil {
				runName = runName + " "
			}
			if spec.WindowFns[0].Ordering.Columns != nil {
				runName = runName + "ORDER BY"
			}
			runName = runName + ")"

			b.Run(runName, func(b *testing.B) {
				post := &distsqlpb.PostProcessSpec{}
				disposer := &RowDisposer{}
				input := NewRepeatableRowSource(sqlbase.ThreeIntCols, rowsGenerator(numRows, numCols))

				b.SetBytes(int64(8 * numRows * numCols))
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					d, err := newWindower(flowCtx, 0 /* processorID */, &spec, input, post, disposer)
					if err != nil {
						b.Fatal(err)
					}
					d.Run(context.TODO())
					input.Reset()
				}
				b.StopTimer()
			})
		}
	}
}
