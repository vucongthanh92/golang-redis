package database

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

// Connection funct
func Connection() redis.Conn {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, errConn := redis.Dial("tcp", "127.0.0.1:6379")
			if errConn != nil {
				return nil, errConn
			}
			_, errDatabase := conn.Do("select", 0)
			if errDatabase != nil {
				conn.Close()
				return nil, errDatabase
			}
			return conn, nil
		},
	}
	conn := pool.Get()
	_, err := conn.Do("ping")
	if err != nil {
		panic(err.Error())
	} else {
		return conn
	}
}
