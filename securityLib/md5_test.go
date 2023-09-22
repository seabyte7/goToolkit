package securityLib

import (
	"testing"
)

func TestMd5SumBytes(t *testing.T) {
	data := []byte("123456")
	if !MD5CompareBytes(data, "e10adc3949ba59abbe56e057f20f883e") {
		t.Errorf("MD5CompareBytes sum string:123456(result:%v) upper:false is wrong.", MD5SumBytes(data, false))
	}

	if !MD5CompareBytes(data, "E10ADC3949BA59ABBE56E057F20F883E") {
		t.Errorf("MD5CompareBytes sum string:123456(result:%v) upper:true is wrong.", MD5SumBytes(data, true))
	}
}

func TestMd5SumString(t *testing.T) {
	data := "123456"
	if !MD5CompareString(data, "e10adc3949ba59abbe56e057f20f883e") {
		t.Errorf("MD5CompareString sum string:123456(result:%v) upper:false is wrong.", MD5SumString(data, false))
	}

	if !MD5CompareString(data, "E10ADC3949BA59ABBE56E057F20F883E") {
		t.Errorf("MD5CompareString sum string:123456(result:%v) upper:true is wrong.", MD5SumString(data, true))
	}
}
