// Code generated by "stringer -type=trashCleanFrequency -trimprefix=trashCleanFrequency"; DO NOT EDIT.

package godrive

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[trashCleanFrequencyNever-0]
	_ = x[trashCleanFrequencyImmediately-1]
	_ = x[trashCleanFrequencyFifteen-2]
	_ = x[trashCleanFrequencyHour-3]
	_ = x[trashCleanFrequencyDay-4]
}

const _trashCleanFrequency_name = "NeverImmediatelyFifteenHourDay"

var _trashCleanFrequency_index = [...]uint8{0, 5, 16, 23, 27, 30}

func (i trashCleanFrequency) String() string {
	if i < 0 || i >= trashCleanFrequency(len(_trashCleanFrequency_index)-1) {
		return "trashCleanFrequency(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _trashCleanFrequency_name[_trashCleanFrequency_index[i]:_trashCleanFrequency_index[i+1]]
}