package securityLib

import (
	"fmt"
	"testing"
)

func TestSHA256SumBytes(t *testing.T) {
	data := []byte("123456")
	if !SHA256CompareBytes(data, "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92") {
		t.Errorf("TestSHA256SumBytes sum string:123456(result:%v) upper:false is wrong.", SHA256SumBytes(data, false))
	}

	if !SHA256CompareBytes(data, "8D969EEF6ECAD3C29A3A629280E686CF0C3F5D5A86AFF3CA12020C923ADC6C92") {
		t.Errorf("TestSHA256SumBytes sum string:123456(result:%v) upper:true is wrong.", SHA256SumBytes(data, true))
	}
}

func TestSHA256SumString(t *testing.T) {

	data := "123456"
	if !SHA256CompareString(data, "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92") {
		t.Errorf("TestSHA256SumString sum string:123456(result:%v) upper:false is wrong.", SHA256SumString(data, false))
	}

	if !SHA256CompareString(data, "8D969EEF6ECAD3C29A3A629280E686CF0C3F5D5A86AFF3CA12020C923ADC6C92") {
		t.Errorf("TestSHA256SumString sum string:123456(result:%v) upper:true is wrong.", SHA256SumString(data, true))
	}

	fmt.Println(SHA256CompareFile("./md5.go", "E9A15E4AF8CEAF207BB537CAC3909D2706837D6F260C22F60E4C3CF0C555C0ED"))
}
