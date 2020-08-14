// Code generated by "stringer -type=DataType"; DO NOT EDIT.

package tds

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TDS_VOID-31]
	_ = x[TDS_IMAGE-34]
	_ = x[TDS_TEXT-35]
	_ = x[TDS_BLOB-36]
	_ = x[TDS_VARBINARY-37]
	_ = x[TDS_INTN-38]
	_ = x[TDS_VARCHAR-39]
	_ = x[TDS_BINARY-45]
	_ = x[TDS_CHAR-47]
	_ = x[TDS_INT1-48]
	_ = x[TDS_DATE-49]
	_ = x[TDS_BIT-50]
	_ = x[TDS_TIME-51]
	_ = x[TDS_INT2-52]
	_ = x[TDS_INT4-56]
	_ = x[TDS_SHORTDATE-58]
	_ = x[TDS_FLT4-59]
	_ = x[TDS_MONEY-60]
	_ = x[TDS_DATETIME-61]
	_ = x[TDS_FLT8-62]
	_ = x[TDS_UINT2-65]
	_ = x[TDS_UINT4-66]
	_ = x[TDS_UINT8-67]
	_ = x[TDS_UINTN-68]
	_ = x[TDS_SENSITIVITY-103]
	_ = x[TDS_BOUNDARY-104]
	_ = x[TDS_DECN-106]
	_ = x[TDS_NUMN-108]
	_ = x[TDS_FLTN-109]
	_ = x[TDS_MONEYN-110]
	_ = x[TDS_DATETIMN-111]
	_ = x[TDS_SHORTMONEY-122]
	_ = x[TDS_DATEN-123]
	_ = x[TDS_TIMEN-147]
	_ = x[TDS_XML-163]
	_ = x[TDS_UNITEXT-174]
	_ = x[TDS_LONGCHAR-175]
	_ = x[TDS_BIGDATETIMEN-187]
	_ = x[TDS_BIGTIMEN-188]
	_ = x[TDS_INT8-191]
	_ = x[TDS_LONGBINARY-225]
	_ = x[TDS_INTERVAL-46]
	_ = x[TDS_SINT1-176]
	_ = x[TDS_DATETIMEN-111]
	_ = x[TDS_USER_TEXT-25]
	_ = x[TDS_USER_IMAGE-32]
	_ = x[TDS_USER_UNITEXT-54]
}

const _DataType_name = "TDS_USER_TEXTTDS_VOIDTDS_USER_IMAGETDS_IMAGETDS_TEXTTDS_BLOBTDS_VARBINARYTDS_INTNTDS_VARCHARTDS_BINARYTDS_INTERVALTDS_CHARTDS_INT1TDS_DATETDS_BITTDS_TIMETDS_INT2TDS_USER_UNITEXTTDS_INT4TDS_SHORTDATETDS_FLT4TDS_MONEYTDS_DATETIMETDS_FLT8TDS_UINT2TDS_UINT4TDS_UINT8TDS_UINTNTDS_SENSITIVITYTDS_BOUNDARYTDS_DECNTDS_NUMNTDS_FLTNTDS_MONEYNTDS_DATETIMNTDS_SHORTMONEYTDS_DATENTDS_TIMENTDS_XMLTDS_UNITEXTTDS_LONGCHARTDS_SINT1TDS_BIGDATETIMENTDS_BIGTIMENTDS_INT8TDS_LONGBINARY"

var _DataType_map = map[DataType]string{
	25:  _DataType_name[0:13],
	31:  _DataType_name[13:21],
	32:  _DataType_name[21:35],
	34:  _DataType_name[35:44],
	35:  _DataType_name[44:52],
	36:  _DataType_name[52:60],
	37:  _DataType_name[60:73],
	38:  _DataType_name[73:81],
	39:  _DataType_name[81:92],
	45:  _DataType_name[92:102],
	46:  _DataType_name[102:114],
	47:  _DataType_name[114:122],
	48:  _DataType_name[122:130],
	49:  _DataType_name[130:138],
	50:  _DataType_name[138:145],
	51:  _DataType_name[145:153],
	52:  _DataType_name[153:161],
	54:  _DataType_name[161:177],
	56:  _DataType_name[177:185],
	58:  _DataType_name[185:198],
	59:  _DataType_name[198:206],
	60:  _DataType_name[206:215],
	61:  _DataType_name[215:227],
	62:  _DataType_name[227:235],
	65:  _DataType_name[235:244],
	66:  _DataType_name[244:253],
	67:  _DataType_name[253:262],
	68:  _DataType_name[262:271],
	103: _DataType_name[271:286],
	104: _DataType_name[286:298],
	106: _DataType_name[298:306],
	108: _DataType_name[306:314],
	109: _DataType_name[314:322],
	110: _DataType_name[322:332],
	111: _DataType_name[332:344],
	122: _DataType_name[344:358],
	123: _DataType_name[358:367],
	147: _DataType_name[367:376],
	163: _DataType_name[376:383],
	174: _DataType_name[383:394],
	175: _DataType_name[394:406],
	176: _DataType_name[406:415],
	187: _DataType_name[415:431],
	188: _DataType_name[431:443],
	191: _DataType_name[443:451],
	225: _DataType_name[451:465],
}

func (i DataType) String() string {
	if str, ok := _DataType_map[i]; ok {
		return str
	}
	return "DataType(" + strconv.FormatInt(int64(i), 10) + ")"
}