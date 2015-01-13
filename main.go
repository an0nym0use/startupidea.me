package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"github.com/soveran/redisurl"
)

var redisPool *redis.Pool

type idea func() string

type IdeaRegistry struct {
	IdeaTypes []idea
}

var registry IdeaRegistry

func (registry IdeaRegistry) GetIdea() string {
	rand.Seed(time.Now().UnixNano())
	chosenType := registry.IdeaTypes[rand.Intn(len(registry.IdeaTypes))]
	return chosenType()
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
		LoadData(dataFile)
	}

	RegisterIdeas()

	app := martini.Classic()
	app.Get("/", func() string {
		return registry.GetIdea()
	})
	app.Run()
}
