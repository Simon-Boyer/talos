// Code generated by "enumer -type OutFormat -linecomment -text"; DO NOT EDIT.

package profile

import (
	"fmt"
	"strings"
)

const _OutFormatName = "unknownraw.tar.gz.xz.gz.zst"

var _OutFormatIndex = [...]uint8{0, 7, 10, 17, 20, 23, 27}

const _OutFormatLowerName = "unknownraw.tar.gz.xz.gz.zst"

func (i OutFormat) String() string {
	if i < 0 || i >= OutFormat(len(_OutFormatIndex)-1) {
		return fmt.Sprintf("OutFormat(%d)", i)
	}
	return _OutFormatName[_OutFormatIndex[i]:_OutFormatIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _OutFormatNoOp() {
	var x [1]struct{}
	_ = x[OutFormatUnknown-(0)]
	_ = x[OutFormatRaw-(1)]
	_ = x[OutFormatTar-(2)]
	_ = x[OutFormatXZ-(3)]
	_ = x[OutFormatGZ-(4)]
	_ = x[OutFormatZSTD-(5)]
}

var _OutFormatValues = []OutFormat{OutFormatUnknown, OutFormatRaw, OutFormatTar, OutFormatXZ, OutFormatGZ, OutFormatZSTD}

var _OutFormatNameToValueMap = map[string]OutFormat{
	_OutFormatName[0:7]:        OutFormatUnknown,
	_OutFormatLowerName[0:7]:   OutFormatUnknown,
	_OutFormatName[7:10]:       OutFormatRaw,
	_OutFormatLowerName[7:10]:  OutFormatRaw,
	_OutFormatName[10:17]:      OutFormatTar,
	_OutFormatLowerName[10:17]: OutFormatTar,
	_OutFormatName[17:20]:      OutFormatXZ,
	_OutFormatLowerName[17:20]: OutFormatXZ,
	_OutFormatName[20:23]:      OutFormatGZ,
	_OutFormatLowerName[20:23]: OutFormatGZ,
	_OutFormatName[23:27]:      OutFormatZSTD,
	_OutFormatLowerName[23:27]: OutFormatZSTD,
}

var _OutFormatNames = []string{
	_OutFormatName[0:7],
	_OutFormatName[7:10],
	_OutFormatName[10:17],
	_OutFormatName[17:20],
	_OutFormatName[20:23],
	_OutFormatName[23:27],
}

// OutFormatString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func OutFormatString(s string) (OutFormat, error) {
	if val, ok := _OutFormatNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _OutFormatNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to OutFormat values", s)
}

// OutFormatValues returns all values of the enum
func OutFormatValues() []OutFormat {
	return _OutFormatValues
}

// OutFormatStrings returns a slice of all String values of the enum
func OutFormatStrings() []string {
	strs := make([]string, len(_OutFormatNames))
	copy(strs, _OutFormatNames)
	return strs
}

// IsAOutFormat returns "true" if the value is listed in the enum definition. "false" otherwise
func (i OutFormat) IsAOutFormat() bool {
	for _, v := range _OutFormatValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalText implements the encoding.TextMarshaler interface for OutFormat
func (i OutFormat) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for OutFormat
func (i *OutFormat) UnmarshalText(text []byte) error {
	var err error
	*i, err = OutFormatString(string(text))
	return err
}
