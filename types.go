package go_abi_param

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"reflect"
	"strings"
)

// 老版本的数组解析只支持一维，且格式为 aaa,bbb,ccc
// 当前版本引入多维数组，格式变更为：
// [1,2,3]
// [[1,2],[3,4]]
// [[[[1,2],[11,22]],[3,4]]]

// https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/type_test.go
func (ap *AbiParam) parseParam(blob, value string) (interface{}, error) {
	if strings.Count(value, "[") != strings.Count(value, "]") {
		return nil, fmt.Errorf("left block count != right block count")
	}

	// 移除用户填写的空格
	value = strings.ReplaceAll(value, " ", "")

	// value 应该支持科学技术法, eg: 1e18
	_, match := new(big.Int).SetString(value, 10)
	if !match {
		x, _, err := new(big.Float).Parse(value, 10)
		if err == nil {
			// value 为数值类型参数
			val, _ := x.Int(big.NewInt(0))
			//_value = val
			value = val.String()
		}
	}

	typ, err := abi.NewType(blob, "", nil)
	if err != nil {
		return nil, fmt.Errorf("blob to go type error: %s", err)
	}

	switch typ.T {
	case abi.SliceTy:
		return ap.forEachUnpackForString(typ, value)
	case abi.ArrayTy:
		return ap.forEachUnpackForString(typ, value)
	case abi.StringTy:
		return readString(value)
	case abi.IntTy, abi.UintTy:
		return readInteger(typ, value)
	case abi.BoolTy:
		return readBool(value)
	case abi.AddressTy:
		return readAddress(value)
	case abi.HashTy:
		return common.HexToHash(value), nil
	case abi.BytesTy:
		return readBytes(value)
	case abi.FixedBytesTy:
		bytesVal, dErr := hexutil.Decode(value)
		if dErr != nil {
			return nil, err
		}
		return readFixedBytes(typ, bytesVal)
	case abi.FunctionTy:
		bytesVal, dErr := hexutil.Decode(value)
		if dErr != nil {
			return nil, err
		}
		return readFunctionType(typ, bytesVal)
	default:
		return nil, fmt.Errorf("abi: unknown type %v", typ.T)
	}
}

func getType(t abi.Type) string {
	switch t.T {
	case abi.IntTy:
		return fmt.Sprintf("%s%d", "int", t.Size)
	case abi.UintTy:
		return fmt.Sprintf("%s%d", "uint", t.Size)
	case abi.BoolTy:
		return "bool"
	case abi.StringTy:
		return "string"
	case abi.SliceTy:
		//return fmt.Sprintf("%s%s", t.Elem.GetType().String(), "[]")
		return fmt.Sprintf("%s%s", t.Elem.String(), "[]")
	case abi.ArrayTy:
		return fmt.Sprintf("%s[%d]", t.Elem.String(), t.Size)
	case abi.TupleTy:
		return "tuple"
	case abi.AddressTy:
		return "address"
	//case abi.FixedBytesTy, abi.BytesTy:
	case abi.BytesTy:
		return "bytes"
	case abi.FixedBytesTy:
		return "bytes32"
	case abi.HashTy:
		return reflect.ArrayOf(32, reflect.TypeOf(byte(0))).String()
	case abi.FixedPointTy:
		return reflect.ArrayOf(32, reflect.TypeOf(byte(0))).String()
	case abi.FunctionTy:
		return "function"
	default:
		panic("Invalid type")
	}
}

func readInteger(typ abi.Type, value string) (interface{}, error) {
	if typ.T == abi.UintTy {
		switch typ.Size {
		case 8:
			return readUint8(value)
		case 16:
			return readUint16(value)
		case 32:
			return readUint32(value)
		case 64:
			return readUint64(value)
		default:
			return readBigInt(value)
		}
	}

	// int
	switch typ.Size {
	case 8:
		return readInt8(value)
	case 16:
		return readInt16(value)
	case 32:
		return readInt32(value)
	case 64:
		return readInt64(value)
	default:
		ret, _ := new(big.Int).SetString(value, 10)
		if ret.Bit(255) == 1 {
			ret.Add(abi.MaxUint256, new(big.Int).Neg(ret))
			ret.Add(ret, common.Big1)
			ret.Neg(ret)
		}
		return ret, nil
	}
}
