package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/garyburd/redigo/redis"
)

type Org struct {
	Name string
}

func LoadData(filename string, redisConn redis.Conn) error {
	defer redisConn.Close()

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var orgs []Org
	if err = json.Unmarshal(file, &orgs); err != nil {
		return err
	}

	if err = PushDataToRedis(orgs, redisConn); err != nil {
		return err
	}
	return nil
}

func PushDataToRedis(orgs []Org, redisConn redis.Conn) error {
	for orgIndex := range orgs {
		_, err := redisConn.Do("SADD", "organisations", orgs[orgIndex].Name)
		if err != nil {
			return err
		}
	}
	return nil
}
