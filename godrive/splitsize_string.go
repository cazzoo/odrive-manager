// Code generated by "stringer -type=splitSize -trimprefix=splitSize"; DO NOT EDIT.

package godrive

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[splitSizeSmall-0]
	_ = x[splitSizeMedium-1]
	_ = x[splitSizeLarge-2]
	_ = x[splitSizeXlarge-3]
}

const _splitSize_name = "SmallMediumLargeXlarge"

var _splitSize_index = [...]uint8{0, 5, 11, 16, 22}

func (i splitSize) String() string {
	if i < 0 || i >= splitSize(len(_splitSize_index)-1) {
		return "splitSize(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _splitSize_name[_splitSize_index[i]:_splitSize_index[i+1]]
}