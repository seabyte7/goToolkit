package redisLib

import "github.com/redis/go-redis/v9"

// wrap https://github.com/redis/go-redis
type RedisClient struct {
	client *redis.Client
}

func SimpleDial(connectUri string) *RedisClient {
	redisPtr := redis.NewClient(&redis.Options{
		Addr: connectUri,
	})
	clientPtr := &RedisClient{client: redisPtr}

	return clientPtr
}

func Dial(redisOptPtr *redis.Options) *RedisClient {
	clientPtr := &RedisClient{client: redis.NewClient(redisOptPtr)}

	return clientPtr
}

func (this *RedisClient) GetClient() *redis.Client {
	return this.client
}

func (this *RedisClient) Close() {
	this.client.Close()
}
