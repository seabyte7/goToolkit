package encodingLib

import "testing"

type Account struct {
	Id   int32
	Name string
}

func TestMarshalGob(t *testing.T) {
	accountObj := Account{
		Id:   1,
		Name: "Account",
	}

	_, err := MarshalGob(accountObj)
	if err != nil {
		t.Errorf("TestMarshalGob MarshalGob object:%+v err:%s", accountObj, err.Error())
	}
}

func TestUnmarshalGob(t *testing.T) {
	accountObj1 := Account{
		Id:   1,
		Name: "accountObj1",
	}

	data, err := MarshalGob(accountObj1)
	if err != nil {
		t.Errorf("TestUnmarshalGob MarshalGob object:%+v err:%s", accountObj1, err.Error())
	}

	var accountObj2 Account
	err = UnmarshalGob(data, &accountObj2)
	if err != nil {
		t.Errorf("TestUnmarshalGob UnmarshalGob object:%+v err:%s", accountObj1, err.Error())
	}

	if accountObj1.Id != accountObj2.Id {
		t.Errorf("TestUnmarshalGob UnmarshalGob object:%+v != UnmarshalGob(%+v)", accountObj1, accountObj2)
	}
}
