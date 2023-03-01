package go_abi_param

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

// 老版本的数组解析只支持一维，且格式为 aaa,bbb,ccc
// 当前版本引入多维数组，格式变更为：
// [1,2,3]        1,2,3    !strings.contains("[") strings.split(",")
// [[1,2],[3,4]]   [1,2],[3,4]   strings.contains("[") 以 ],[为分割符 index%2 == 1 加 ] , index%2==0 加 [
// [[[[1,2],[11,22]],[3,4]]]
func fmtToBlock(originVal string) string {
	if !hasBlock(originVal) {
		//return originVal
		return fmt.Sprintf("%s%s%s", "[", originVal, "]")
	}
	return originVal
}

func hasBlock(originVal string) bool {
	return originVal[:1] == "[" && originVal[len(originVal)-1:] == "]"
}

func parseUnpackString(value string) ([]interface{}, error) {
	value = fmtToBlock(value)
	if strings.Count(value, "[") != strings.Count(value, "]") {
		return nil, fmt.Errorf("unpaired block")
	}

	// 可能传入可能值，不一定是非标准的字符串，所以需要将其转换为字符串，再使用 json 包解析
	var _value []byte
	for i, _v := range []byte(value) {
		if _v == '[' {
			if []byte(value)[i+1] == '[' {
				_value = append(_value, _v)
				continue
			}
			_value = append(_value, _v, '"')
			continue
		} else if _v == '"' {
			_value = append(_value, _v)
			continue
		} else if _v == ']' {
			if len([]byte(value)) > i+1 && []byte(value)[i+1] == ']' {
				_value = append(_value, _v)
				continue
			}
			_value = append(_value, '"', _v)
			continue
		} else if _v == ',' {
			_value = append(_value, '"', _v, '"')
			continue
		} else {
			if len([]byte(value)) > i+1 && []byte(value)[i+1] == ']' {
				_value = append(_value, _v, '"')
				continue
			}
			_value = append(_value, _v)
		}
	}
	value = strings.ReplaceAll(string(_value), `""`, `"`)
	value = strings.ReplaceAll(value, `]",`, `],`)
	value = strings.ReplaceAll(value, `,"[`, `,[`)
	value = strings.ReplaceAll(value, `]"]`, `]]`)

	var parsedData interface{}
	err := json.Unmarshal([]byte(value), &parsedData)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(parsedData).String() != "[]interface {}" {
		return nil, fmt.Errorf("parsed data type is interface should be []interface {]")
	}
	return parsedData.([]interface{}), err
}

func (ap *AbiParam) forEachUnpackForString(t abi.Type, originVal string) (interface{}, error) {
	output, err := parseUnpackString(originVal)
	if err != nil {
		return nil, err
	}

	if t.Size < 0 {
		return nil, fmt.Errorf("cannot parse input array, size is negative (%d)", t.Size)
	}

	if t.Size > len(output) {
		return nil, fmt.Errorf("abi: cannot marshal in to go array: offset %d would go over slice boundary (len=%d)", len(output), t.Size)
	}

	// this value will become our slice or our array, depending on the type
	var refSlice reflect.Value

	if t.T == abi.SliceTy {
		// reset slice size
		t.Size = len(output)
		// declare our slice
		refSlice = reflect.MakeSlice(t.GetType(), t.Size, t.Size)
	} else if t.T == abi.ArrayTy {
		// declare our array
		refSlice = reflect.New(t.GetType()).Elem()
	} else {
		return nil, fmt.Errorf("abi: invalid type in array/slice unpacking stage")
	}

	for i := 0; i < t.Size; i++ {
		ap.logger.Debugf("nest type: %s", getType(*t.Elem))

		//opVal := unpackDynamicData(t.Size, output)
		opVal := ""
		if reflect.TypeOf(output).Kind() == reflect.String {
			opVal = output[i].(string)
		} else {
			opVal = unpackDynamicData(output[i])
		}
		inter, err := ap.parseParam(getType(*t.Elem), opVal)
		if err != nil {
			return nil, err
		}
		refSlice.Index(i).Set(reflect.ValueOf(inter))
	}
	return refSlice.Interface(), nil
}

func unpackDynamicData(ov interface{}) string {
	var _val string

	iValType := reflect.TypeOf(ov).String()
	switch iValType {
	case "string":
		return ov.(string)
	case "[]interface {}":
		var _ovArr []string
		for i, _ov := range ov.([]interface{}) {
			_dData := unpackDynamicData(_ov)
			// replace to the original value
			if i == 0 {
				_dData = "[" + _dData
			}

			// replace to the original value
			if i == len(ov.([]interface{}))-1 {
				_dData = _dData + "]"
			}
			_ovArr = append(_ovArr, _dData)
		}
		_val = strings.Join(_ovArr, ",")
	}
	return _val
}

// readFixedBytes uses reflection to create a fixed array to be read from.
func readFixedBytes(t abi.Type, word []byte) (interface{}, error) {
	if t.T != abi.FixedBytesTy {
		return nil, fmt.Errorf("abi: invalid type in call to make fixed byte array")
	}
	// convert
	array := reflect.New(t.GetType()).Elem()

	reflect.Copy(array, reflect.ValueOf(word[0:t.Size]))
	return array.Interface(), nil
}

// readFunctionType enforces that standard by always presenting it as a 24-array (address + sig = 24 bytes)
func readFunctionType(t abi.Type, word []byte) (funcTy [24]byte, err error) {
	if t.T != abi.FunctionTy {
		return [24]byte{}, fmt.Errorf("abi: invalid type in call to make function type byte array")
	}
	if garbage := binary.BigEndian.Uint64(word[24:32]); garbage != 0 {
		err = fmt.Errorf("abi: got improperly encoded function type, got %v", word)
	} else {
		copy(funcTy[:], word[0:24])
	}
	return
}

func readBigInt(value string) (*big.Int, error) {
	v, ok := big.NewInt(0).SetString(value, 10)
	if !ok {
		return nil, fmt.Errorf("param %s can not convent to int", value)
	}
	return v, nil
}

func readAddress(value string) (common.Address, error) {
	if value == "" {
		return common.Address{}, fmt.Errorf("can't convent param %s to address", value)
	}
	return common.HexToAddress(value), nil
}

// readBool reads a bool.
func readBool(word string) (bool, error) {
	switch word {
	case "0", "false":
		return false, nil
	case "1", "true":
		return true, nil
	default:
		return false, errBadBool
	}
}

func readInt8(value string) (int8, error) {
	bv, err := readBigInt(value)
	if err != nil {
		return 0, err
	}
	return int8(bv.Int64()), nil
}

func readUint8(value string) (uint8, error) {
	bv, err := readBigInt(value)
	if err != nil {
		return 0, err
	}
	return uint8(bv.Int64()), nil
}

func readInt16(value string) (int16, error) {
	v, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

func readUint16(value string) (uint16, error) {
	v, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(v), nil
}

func readInt32(value string) (int32, error) {
	v, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

func readUint32(value string) (uint32, error) {
	v, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

func readInt64(value string) (int64, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func readUint64(value string) (uint64, error) {
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func readBytes(value string) ([]byte, error) {
	return hexutil.Decode(value)
}

func readString(value string) (string, error) {
	return value, nil
}
