package go_abi_param

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/magiconair/properties/assert"
	"math/big"
	"reflect"
	"testing"
	"unsafe"
)

var (
	biVal     *big.Int
	bytesVal  []byte
	byte32Val [32]byte
)

func byte32(s []byte) (a *[32]byte) {
	if len(a) <= len(s) {
		a = (*[len(a)]byte)(unsafe.Pointer(&s[0]))
	}
	return a
}

func init() {
	biVal, _ = new(big.Int).SetString("1000", 10)
	bytesVal, _ = hexutil.Decode("0x9e99847ecf80af04f0808e017172bc71b71a5d1bb7b82ab1ce4b2ec666f009425419ad6e1d42c27f0d8408976e20276e5fd2411c6dc42d06d885b4c25d71fbb31c")
	_byte32Val, _ := hexutil.Decode("0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000")
	byte32Val = *byte32(_byte32Val)
}

func TestAbiParam_Parse(t *testing.T) {
	tests := []struct {
		name       string
		blob       string
		value      string
		goArgument string
		want       interface{}
	}{
		{
			name:       "error: abi unsupported int",
			blob:       "int",
			goArgument: "error",
			value:      "1000",
		},
		{
			name:       "normal: int8",
			blob:       "int8",
			value:      "3",
			goArgument: "int8",
			want:       int8(3),
		},
		{
			name:       "normal: int16",
			blob:       "int16",
			value:      "11",
			goArgument: "int16",
			want:       int16(11),
		},
		{
			name:       "normal: int32",
			blob:       "int32",
			value:      "1122",
			goArgument: "int32",
			want:       int32(1122),
		},
		{
			name:       "normal: int64",
			blob:       "int64",
			value:      "111111",
			goArgument: "int64",
			want:       int64(111111),
		},
		{
			name:       "normal: int128",
			blob:       "int128",
			value:      "1000",
			goArgument: "*big.Int",
			want:       biVal,
		},
		{
			name:       "normal: int256",
			blob:       "int256",
			value:      "1000",
			goArgument: "*big.Int",
			want:       biVal,
		},
		{
			name:       "normal: uint8",
			blob:       "uint8",
			value:      "3",
			goArgument: "uint8",
			want:       uint8(3),
		},
		{
			name:       "normal: uint16",
			blob:       "uint16",
			value:      "20",
			goArgument: "uint16",
			want:       uint16(20),
		},
		{
			name:       "normal: uint32",
			blob:       "uint32",
			value:      "100",
			goArgument: "uint32",
			want:       uint32(100),
		},
		{
			name:       "normal: uint64",
			blob:       "uint64",
			value:      "100",
			goArgument: "uint64",
			want:       uint64(100),
		},
		{
			name:       "normal: uint128",
			blob:       "uint128",
			value:      "1000",
			goArgument: "*big.Int",
			want:       biVal,
		},
		{
			name:       "normal: uint256",
			blob:       "uint256",
			value:      "1000",
			goArgument: "*big.Int",
			want:       biVal,
		},
		{
			name:       "error: uint, abi unSupported",
			blob:       "uint",
			value:      "1000",
			goArgument: "uint",
		},
		{
			name:       "normal: address",
			blob:       "address",
			value:      "0x00000000006c3852cbef3e08e8df289169ede581",
			goArgument: "common.Address",
			want:       common.HexToAddress("0x00000000006c3852cbef3e08e8df289169ede581"),
		},
		{
			name:       "normal: address[]",
			blob:       "address[]",
			value:      `["0x00000000006c3852cbef3e08e8df289169ede581","0x1b2667862b2a4f46DfD6C53f561C58a8B0EED0D6"]`,
			goArgument: "[]common.Address",
			want: []common.Address{common.HexToAddress("0x00000000006c3852cbef3e08e8df289169ede581"),
				common.HexToAddress("0x1b2667862b2a4f46DfD6C53f561C58a8B0EED0D6")},
		},
		{
			name:       "error: unsupported int[] type",
			blob:       "int[]",
			value:      "[1,2]",
			goArgument: "int[]",
		},
		{
			name:       "normal: int8[]",
			blob:       "int8[]",
			value:      "[1,3]",
			goArgument: "[]int8",
			want:       []int8{1, 3},
		},
		{
			name:       "normal: int16[]",
			blob:       "int16[]",
			value:      "[3,233]",
			goArgument: "[]int16",
			want:       []int16{3, 233},
		},
		{
			name:       "normal: int32[]",
			blob:       "int32[]",
			value:      "3,344",
			goArgument: "[]int32",
			want:       []int32{3, 344},
		},
		{
			name:       "normal: int64[]",
			blob:       "int64[]",
			value:      "1000,10001",
			goArgument: "[]int64",
			want:       []int64{1000, 10001},
		},
		{
			name:       "normal: int128[]",
			blob:       "int128[]",
			value:      "[1000]",
			goArgument: "[]*big.Int",
			want:       []*big.Int{biVal},
		},
		{
			name:       "normal: int256[]",
			blob:       "int256[]",
			value:      "[1000]",
			goArgument: "[]*big.Int",
			want:       []*big.Int{biVal},
		},
		{
			name:       "error: unSupported uint[] type",
			blob:       "uint[]",
			value:      "",
			goArgument: "uint64",
		},
		{
			name:       "normal: uint8[]",
			blob:       "uint8[]",
			value:      "[1,2,3]",
			goArgument: "[]uint8",
			want:       []uint8{1, 2, 3},
		},
		{
			name:       "normal: uint16[]",
			blob:       "uint16[]",
			value:      "[4,5,6]",
			goArgument: "[]uint16",
			want:       []uint16{4, 5, 6},
		},
		{
			name:       "normal: uint32[]",
			blob:       "uint32[]",
			value:      "[100,200,400]",
			goArgument: "[]uint32",
			want:       []uint32{100, 200, 400},
		},
		{
			name:       "normal: uint64[]",
			blob:       "uint64[]",
			value:      "[1000,2000]",
			goArgument: "[]uint64",
			want:       []uint64{1000, 2000},
		},
		{
			name:       "normal: uint128[]",
			blob:       "uint128[]",
			value:      "1000",
			goArgument: "[]*big.Int",
			want:       []*big.Int{biVal},
		},
		{
			name:       "normal: uint256[]",
			blob:       "uint256[]",
			value:      "1000",
			goArgument: "[]*big.Int",
			want:       []*big.Int{biVal},
		},
		{
			name:       "normal: bytes",
			blob:       "bytes",
			value:      "0x9e99847ecf80af04f0808e017172bc71b71a5d1bb7b82ab1ce4b2ec666f009425419ad6e1d42c27f0d8408976e20276e5fd2411c6dc42d06d885b4c25d71fbb31c",
			goArgument: "[]uint8",
			want:       bytesVal,
		},
		{
			name:       "normal: slice bytes",
			blob:       "bytes[]",
			value:      "[0x9e99847ecf80af04f0808e017172bc71b71a5d1bb7b82ab1ce4b2ec666f009425419ad6e1d42c27f0d8408976e20276e5fd2411c6dc42d06d885b4c25d71fbb31c,0x9e99847ecf80af04f0808e017172bc71b71a5d1bb7b82ab1ce4b2ec666f009425419ad6e1d42c27f0d8408976e20276e5fd2411c6dc42d06d885b4c25d71fbb31c]",
			goArgument: "[][]uint8",
			want:       [][]uint8{bytesVal, bytesVal},
		},
		{
			name:       "normal: bytes32",
			blob:       "bytes32",
			value:      "0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000",
			goArgument: "[32]uint8",
			want:       byte32Val,
		},
		{
			name:       "normal: bytes32[]",
			blob:       "bytes32[]",
			value:      "[0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000,0x0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000]",
			goArgument: "[][32]uint8",
			want:       [][32]byte{byte32Val, byte32Val},
		},
		{
			name:       "normal: string",
			blob:       "string",
			value:      "abcd434d32",
			goArgument: "string",
			want:       "abcd434d32",
		},
		{
			name:       "normal: bool",
			blob:       "bool",
			value:      "1",
			goArgument: "bool",
			want:       true,
		},
		{
			name:       "normal: bool",
			blob:       "bool",
			value:      "true",
			goArgument: "bool",
			want:       true,
		},
		{
			name:       "normal: array",
			blob:       "address[3]",
			value:      `0x543a5aed5abc902553a92547701ac38f73a70785,0x7d2768de32b0b80b7a3454c06bdac94a69ddc7a9,0x028171bCA77440897B824Ca71D1c56caC55b68A3`,
			goArgument: "[3]common.Address",
			want: [3]common.Address{common.HexToAddress("0x543a5aed5abc902553a92547701ac38f73a70785"),
				common.HexToAddress("0x7d2768de32b0b80b7a3454c06bdac94a69ddc7a9"),
				common.HexToAddress("0x028171bCA77440897B824Ca71D1c56caC55b68A3")},
		},
		{
			name:       "normal: new fmt array",
			blob:       "address[3]",
			value:      `[0x543a5aed5abc902553a92547701ac38f73a70785,0x7d2768de32b0b80b7a3454c06bdac94a69ddc7a9,0x028171bCA77440897B824Ca71D1c56caC55b68A3]`,
			goArgument: "[3]common.Address",
			want: [3]common.Address{common.HexToAddress("0x543a5aed5abc902553a92547701ac38f73a70785"),
				common.HexToAddress("0x7d2768de32b0b80b7a3454c06bdac94a69ddc7a9"),
				common.HexToAddress("0x028171bCA77440897B824Ca71D1c56caC55b68A3")},
		},
		{
			name:       "normal: slice && array with bool",
			blob:       "bool[][2]",
			value:      "[[1,0,1],[0,1]]",
			goArgument: "[2][]bool",
			want:       [2][]bool{{true, false, true}, {false, true}},
		},
		{
			name:       "normal: slice && array with bool",
			blob:       "bool[2][2]",
			value:      "[[true,false],[false,true]]",
			goArgument: "[2][2]bool",
			want:       [2][2]bool{{true, false}, {false, true}},
		},
		{
			name:       "normal: slice && array with string",
			blob:       "string[][3]",
			value:      `[[aaa,vvv,bbb],[w4f,6s%#],[14c14,c423,f34e&*^,fjhvfw]]`,
			goArgument: "[3][]string",
			want:       [3][]string{{"aaa", "vvv", "bbb"}, {"w4f", "6s%#"}, {"14c14", "c423", "f34e&*^", "fjhvfw"}},
		},
		{
			name:  "error: slice && array, out of index",
			blob:  "bool[][2]",
			value: "[[1,0],[0,0,1]]",
			want:  false,
		},
		{
			name:  "error: slice && array, arguments",
			blob:  "bool[][2]",
			value: "[[1,5],[2,4]]]",
			want:  false,
		},
		{
			name:       "nested array",
			blob:       "int8[2][2][2]",
			value:      "[[[1,2],[3,4]],[[5,6],[7,8]]]",
			goArgument: "[2][2][2]int8",
			want:       [2][2][2]int8{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param, err := NewAbiParam(tt.blob, tt.value)
			if err != nil {
				t.Errorf("new abi param error: %s", err)
				return
			}
			parsedData, err := param.Parse()
			if err != nil {
				t.Errorf("parsed abi params error: %s", err)
				return
			}
			t.Logf("res: %v", parsedData)
			assert.Equal(t, reflect.TypeOf(parsedData).String(), tt.goArgument, "argument: %s", tt.goArgument)
			assert.Equal(t, parsedData, tt.want, fmt.Sprintf("parsed data: %v", tt.want))
		})
	}
}

func TestParseUnpackString(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "normal:not contains block",
			value: "1,2,3",
			want:  false,
		},
		{
			name:  "normal:not contains block",
			value: "[1,2,3]",
			want:  true,
		},
		{
			name:  "normal:not contains block",
			value: `["1","2","3"]`,
			want:  true,
		},
		{
			name:  "normal:not contains block",
			value: `0x543a5aed5abc902553a92547701ac38f73a70785,0x7d2768de32b0b80b7a3454c06bdac94a69ddc7a9,0x028171bCA77440897B824Ca71D1c56caC55b68A3`,
			want:  true,
		},
		{
			name:  "normal:not contains block",
			value: "0x543a5aed5abc902553a92547701ac38f73a70785,0x7d2768de32b0b80b7a3454c06bdac94a69ddc7a9,0x028171bCA77440897B824Ca71D1c56caC55b68A3",
			want:  true,
		},
		{
			name:  "normal:not contains block",
			value: `["0x543a5aed5abc902553a92547701ac38f73a70785","0x7d2768de32b0b80b7a3454c06bdac94a69ddc7a9","0x028171bCA77440897B824Ca71D1c56caC55b68A3"]`,
			want:  true,
		},
		{
			name:  "normal: paired block",
			value: "[1]",
			want:  true,
		},
		{
			name:  "normal: paired block",
			value: `["1]`,
			want:  true,
		},
		{
			name:  "normal: paired block",
			value: "[1,2]",
			want:  true,
		},
		{
			name:  "normal: paired block",
			value: "[[20],[21,23]]",
			want:  true,
		},
		{
			name:  "normal: paired block",
			value: `[[20"],["21,23"]]`,
			want:  true,
		},
		{
			name:  "normal: nest array",
			value: "[[[1,2],[3,4]],[[5,6],[7,8]]]",
			want:  true,
		},
		{
			name:  "normal: nest array",
			value: "[[[[1]]]]",
			want:  true,
		},
		{
			name:  "normal: nest array",
			value: "[[[[1,2],[1,2]],[[1,2],[3,4]],[[5,6],[7,8]]]]",
			want:  true,
		},
		{
			name:  "error: quote repetition",
			value: `[[20"],["21,23""]]`,
			want:  false,
		},
		{
			name:  "error: unpaired block",
			value: "[[[1,2]]",
			want:  false,
		},
		{
			name:  "error: unpaired block",
			value: "[[1,2]]]",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unpackedData, err := parseUnpackString(tt.value)
			if err != nil {
				t.Errorf("err: %s", err)
				return
			}
			t.Logf("unpackedData: %v", unpackedData)
		})
	}
}
