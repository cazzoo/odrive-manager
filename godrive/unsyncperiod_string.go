// Code generated by "stringer -type=unsyncPeriod -trimprefix=unsyncPeriod"; DO NOT EDIT.

package godrive

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[unsyncPeriodNever-0]
	_ = x[unsyncPeriodDay-1]
	_ = x[unsyncPeriodWeek-2]
	_ = x[unsyncPeriodMonth-3]
}

const _unsyncPeriod_name = "NeverDayWeekMonth"

var _unsyncPeriod_index = [...]uint8{0, 5, 8, 12, 17}

func (i unsyncPeriod) String() string {
	if i < 0 || i >= unsyncPeriod(len(_unsyncPeriod_index)-1) {
		return "unsyncPeriod(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _unsyncPeriod_name[_unsyncPeriod_index[i]:_unsyncPeriod_index[i+1]]
}