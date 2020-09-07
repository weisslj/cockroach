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

package sql

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/weisslj/cockroach/pkg/settings"
	"github.com/weisslj/cockroach/pkg/sql/privilege"
	"github.com/weisslj/cockroach/pkg/sql/sqlbase"
	"github.com/weisslj/cockroach/pkg/util/log"
	"github.com/weisslj/cockroach/pkg/util/timeutil"
)

// This file contains facilities to report SQL activities to separate
// log files.
//
// The log format is currently as follows:
//
// Example audit log line:
// I180211 07:30:48.832004 317 sql/exec_log.go:90  [client=127.0.0.1:62503,user=root,n1] 13 exec "cockroach" {"ab"[53]:READ} "SELECT * FROM ab" {} 123.45 12 OK
// I180211 07:30:48.832004 317 sql/exec_log.go:90  [client=127.0.0.1:62503,user=root,n1] 13 exec "cockroach" {"ab"[53]:READ} "SELECT nonexistent FROM ab" {} 0.123 12 ERROR
// Example execution log:
// I180211 07:30:48.832004 317 sql/exec_log.go:90  [client=127.0.0.1:62503,user=root,n1] 13 exec "cockroach" {} "SELECT * FROM ab" {} 123.45 12 OK
// I180211 07:30:48.832004 317 sql/exec_log.go:90  [client=127.0.0.1:62503,user=root,n1] 13 exec "cockroach" {} "SELECT nonexistent FROM ab" {} 0.123 0 "column \"nonexistent\" not found"
//
// Explanation of fields:
// I180211 07:30:48.832004 317 sql/exec_log.go:90  [client=127.0.0.1:62503,user=root,n1] 13 exec "cockroach" {"ab"[53]:READ} "SELECT nonexistent FROM ab" {} 0.123 12 ERROR
// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//  \                                                                                                                                   .../
//   '- prefix generated by CockroachDB's standard log package. Contains:
//
// - date and time
//   - incidentally, this data is needed for auditing.
// - goroutine ID
// - where in this file the message was generated
// - the logging tags [...] - this contains the client address,
//   username, and node ID
//   - tags were populated in the logging context when setting up
//     the pgwire session.
//   - generally useful and actually a requirement for auditing.
// - a counter for the logging entry, monotonically increasing
//   from the point the process started.
//   - this is a requiredement for auditing too.
//
//  .-----------------------------------------------------------------------------------------------------------------------------.../
//  |
//  '- log message generated in this file. Includes:
//
//  - a label indicating where the data was generated - useful for troubleshooting.
//    - distinguishes e.g. exec, prepare, internal-exec, etc.
//  - the current value of `application_name`
//    - required for auditing, also helps filter out messages from a specific app.
//  - the logging trigger.
//    - "{}" for execution logs: any activity is worth logging in the exec log
//    - the list of triggering tables and access modes for audit
//      events. This is needed for auditing.
//  - the full text of the query.
//    - audit logs really ought to only "identify the data that's accessed but
//      without revealing PII". We don't know how to separate those two things
//      yet so the audit log contains the full SQL of the query even though
//      this may yield PII. This may need to be addressed later.
//  - the placeholder values. Useful for queries using placehodlers.
//    - "{}" when there are no placeholders.
//  - the query execution time in milliseconds. For troubleshooting.
//  - the number of rows that were produced. For troubleshooting.
//  - the status of the query (OK for success, ERROR or full error
//    message upon error). Needed for auditing and troubleshooting.

// logStatementsExecuteEnabled causes the Executor to log executed
// statements and, if any, resulting errors.
var logStatementsExecuteEnabled = settings.RegisterBoolSetting(
	"sql.trace.log_statement_execute",
	"set to true to enable logging of executed statements",
	false,
)

// maybeLogStatement conditionally records the current statement
// (p.curPlan) to the exec / audit logs.
func (p *planner) maybeLogStatement(ctx context.Context, lbl string, rows int, err error) {
	p.maybeLogStatementInternal(
		ctx, lbl, rows, err, p.statsCollector.PhaseTimes()[sessionQueryReceived])
}

func (p *planner) maybeLogStatementInternal(
	ctx context.Context, lbl string, rows int, err error, startTime time.Time,
) {
	// Note: if you find the code below crashing because p.execCfg == nil,
	// do not add a test "if p.execCfg == nil { do nothing }" !
	// Instead, make the logger work. This is critical for auditing - we
	// can't miss any statement.

	logV := log.V(2)
	logExecuteEnabled := logStatementsExecuteEnabled.Get(&p.execCfg.Settings.SV)
	auditEventsDetected := len(p.curPlan.auditEvents) != 0

	if !logV && !logExecuteEnabled && !auditEventsDetected {
		return
	}

	// Logged data, in order:

	// label passed as argument.

	appName := p.EvalContext().SessionData.ApplicationName

	logTrigger := "{}"
	if auditEventsDetected {
		var buf bytes.Buffer
		buf.WriteByte('{')
		sep := ""
		for _, ev := range p.curPlan.auditEvents {
			mode := "READ"
			if ev.writing {
				mode = "READWRITE"
			}
			fmt.Fprintf(&buf, "%s%q[%d]:%s", sep, ev.desc.GetName(), ev.desc.GetID(), mode)
			sep = ", "
		}
		buf.WriteByte('}')
		logTrigger = buf.String()
	}

	stmtStr := p.curPlan.AST.String()

	plStr := p.extendedEvalCtx.Placeholders.Values.String()

	age := float64(timeutil.Now().Sub(startTime).Nanoseconds()) / 1e6

	// rows passed as argument.

	execErrStr := ""
	auditErrStr := "OK"
	if err != nil {
		execErrStr = err.Error()
		auditErrStr = "ERROR"
	}

	// Now log!
	if auditEventsDetected {
		logger := p.execCfg.AuditLogger
		logger.Logf(ctx, "%s %q %s %q %s %.3f %d %s",
			lbl, appName, logTrigger, stmtStr, plStr, age, rows, auditErrStr)
	}
	if logExecuteEnabled {
		logger := p.execCfg.ExecLogger
		logger.Logf(ctx, "%s %q %s %q %s %.3f %d %q",
			lbl, appName, logTrigger, stmtStr, plStr, age, rows, execErrStr)
	}
	if logV {
		// Copy to the main log.
		log.VEventf(ctx, 2, "%s %q %s %q %s %.3f %d %q",
			lbl, appName, logTrigger, stmtStr, plStr, age, rows, execErrStr)
	}
}

// maybeAudit marks the current plan being constructed as flagged
// for auditing if the table being touched has an auditing mode set.
// This is later picked up by maybeLogStatement() above.
//
// It is crucial that this gets checked reliably -- we don't want to
// miss any statements! For now, we call this from CheckPrivilege(),
// as this is the function most likely to be called reliably from any
// caller that also uses a descriptor. Future changes that move the
// call to this method elsewhere must find a way to ensure that
// contributors who later add features do not have to remember to call
// this to get it right.
func (p *planner) maybeAudit(desc sqlbase.DescriptorProto, priv privilege.Kind) {
	wantedMode := desc.GetAuditMode()
	if wantedMode == sqlbase.TableDescriptor_DISABLED {
		return
	}

	switch priv {
	case privilege.INSERT, privilege.DELETE, privilege.UPDATE:
		p.curPlan.auditEvents = append(p.curPlan.auditEvents, auditEvent{desc: desc, writing: true})
	default:
		p.curPlan.auditEvents = append(p.curPlan.auditEvents, auditEvent{desc: desc, writing: false})
	}
}

// auditEvent represents an audit event for a single table.
type auditEvent struct {
	// The descriptor being audited.
	desc sqlbase.DescriptorProto
	// Whether the event was for INSERT/DELETE/UPDATE.
	writing bool
}
