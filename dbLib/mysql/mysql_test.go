package mysql

import (
	. "github.com/seabyte7/goToolkit/protocol"
	"testing"
)

func TestDial(t *testing.T) {
	connectUri := "root:123456@tcp(127.0.0.1:3306)/login?charset=utf8mb4"
	mysqlClientPtr, result := Dial(connectUri)
	if result != Success {
		t.Errorf("TestMysql connect failed, result:%s", result.Error())
		return
	}

	defer mysqlClientPtr.Close()
}
