package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func RegisterIdeas() {
	registry.IdeaTypes = append(registry.IdeaTypes, func() string {
		redisConn := redisPool.Get()
		defer redisConn.Close()
		results, err := redis.Values(redisConn.Do("SRANDMEMBER", "organisations", "2"))
		if err != nil {
			panic(err.Error())
		}
		return fmt.Sprintf("A cross between %s, and %s.\n", results[0], results[1])
	})
}
