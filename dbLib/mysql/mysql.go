package mysql

import (
	"github.com/seabyte7/goToolkit/logLib"
	. "github.com/seabyte7/goToolkit/protocol"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLClient struct {
	*gorm.DB
}

func Dial(uri string) (*MySQLClient, Result) {
	dialectorPtr := mysql.Open(uri)
	dbConfigPtr := &gorm.Config{}
	dbPtr, err := gorm.Open(dialectorPtr, dbConfigPtr)
	if err != nil {
		logLib.Sugar().Errorf("NewMySQLClient failed, uri: %s, error: %v", uri, err.Error())
		return nil, err
	}

	clientPtr := &MySQLClient{
		DB: dbPtr,
	}

	return clientPtr, Success
}

func (this *MySQLClient) Close() {
	sqlDBPtr, err := this.DB.DB()
	if err != nil {
		logLib.Sugar().Errorf("Close failed, error: %v", err.Error())
		return
	}

	err = sqlDBPtr.Close()
	if err != nil {
		logLib.Sugar().Errorf("Close failed, error: %v", err.Error())
	}
}
