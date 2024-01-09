package redisLib

import (
	"context"
	"testing"
)

func TestSimpleDial(t *testing.T) {
	redisClientPtr := SimpleDial("127.0.0.1:6379")
	defer redisClientPtr.Close()

	redisClientPtr.GetClient().Set(context.Background(), "test_key", "test_value", 0)
	val, err := redisClientPtr.GetClient().Get(context.Background(), "test_key").Result()
	if err != nil {
		t.Errorf("TestSimpleDial err:%s", err.Error())
	}
	t.Logf("TestSimpleDial val:%s", val)
}
