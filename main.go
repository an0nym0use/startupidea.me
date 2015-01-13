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
	defer redisConn.Close()
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
		Wait: true,
	}

	if dataFile := os.Getenv("DATA_FILE"); len(dataFile) > 0 {
		LoadData(dataFile, redisPool.Get())
	}

	app := martini.Classic()
	app.Get("/", func() string {
		return idea(redisPool.Get())
	})
	app.Run()
}
