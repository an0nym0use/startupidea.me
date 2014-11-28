package main

import (
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"github.com/soveran/redisurl"
)

var redisPool *redis.Pool

func idea(redisConn redis.Conn) string {
	results, err := redis.Values(redisConn.Do("SRANDMEMBER", "organisations", "2"))
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("A cross between %s, and %s.\n", results[0], results[1])
}

func main() {
	redisPool = &redis.Pool{
		MaxIdle:   3,
		MaxActive: 25,
		Dial: func() (redis.Conn, error) {
			redisConn, err := redisurl.ConnectToURL(os.Getenv("REDISTOGO_URL"))
			if err != nil {
				panic(err.Error())
			}
			return redisConn, err
		},
	}

	app := martini.Classic()
	app.Get("/", func() string {
		redisConn := redisPool.Get()
		defer redisConn.Close()
		return idea(redisConn)
	})
	app.Run()
}
