package go_abi_param

import (
	"errors"
	"github.com/sirupsen/logrus"
)

var (
	errBadBlob  = errors.New("param: blob is null")
	errBadValue = errors.New("param: value is null")
	errBadBool  = errors.New("param: improperly encoded boolean value")
)

type AbiParam struct {
	blob   string
	value  string
	logger *logrus.Logger
}

func NewAbiParam(blob string, value string) (*AbiParam, error) {
	ap := &AbiParam{blob: blob, value: value, logger: logrus.New()}
	return ap, ap.check()
}

func (ap *AbiParam) check() error {
	if ap.blob == "" {
		return errBadBlob
	}

	if ap.value == "" {
		return errBadValue
	}
	return nil
}

func (ap *AbiParam) Parse() (interface{}, error) {
	return ap.parseParam(ap.blob, ap.value)
}
