// Code generated by "stringer -type=Method"; DO NOT EDIT.

package roachpb

import "strconv"

const _Method_name = "GetPutConditionalPutIncrementDeleteDeleteRangeClearRangeScanReverseScanBeginTransactionEndTransactionAdminSplitAdminMergeAdminTransferLeaseAdminChangeReplicasAdminRelocateRangeHeartbeatTxnGCPushTxnQueryTxnQueryIntentResolveIntentResolveIntentRangeMergeTruncateLogRequestLeaseTransferLeaseLeaseInfoComputeChecksumCheckConsistencyInitPutWriteBatchExportImportAdminScatterAddSSTableRecomputeStatsRefreshRefreshRangeSubsumeRangeStats"

var _Method_index = [...]uint16{0, 3, 6, 20, 29, 35, 46, 56, 60, 71, 87, 101, 111, 121, 139, 158, 176, 188, 190, 197, 205, 216, 229, 247, 252, 263, 275, 288, 297, 312, 328, 335, 345, 351, 357, 369, 379, 393, 400, 412, 419, 429}

func (i Method) String() string {
	if i < 0 || i >= Method(len(_Method_index)-1) {
		return "Method(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Method_name[_Method_index[i]:_Method_index[i+1]]
}