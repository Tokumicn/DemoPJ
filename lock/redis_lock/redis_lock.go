package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"os/exec"
	"sync"
	"time"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var (
	counter int64
	wg      sync.WaitGroup
	lockKey = "myrslock"
)

func main() {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			lock(incr)
		}()
	}
	wg.Wait()
	fmt.Println("final counter is %d\n", counter)
}

func incr() {
	counter++
	fmt.Printf("after incr is %d\n", counter)
}

func lock(myFunc func()) {
	defer wg.Done()

	uuid := getUuid()
	lockSuccess, err := client.SetNX(lockKey, uuid, time.Second*5).Result()
	if err != nil || !lockSuccess {
		fmt.Println("get lock fail")
		return
	} else {
		fmt.Println("get lock")
	}

	myFunc()

	// unlock
	//_, err = client.Del(lockKey).Result()
	//if err != nil {
	//	fmt.Println("unlock fail")
	//} else {
	//	fmt.Println("unlock")
	//}

	var luaScript = redis.NewScript(`
      	if redis.call("get", KEYS[1]) == ARGV[1]
			then
					return redis.call("del", KEYS[1])
			else
					return 0
			end
	`)

	res, _ := luaScript.Run(client, []string{lockKey}, uuid).Result()
	if res == 0 {
		fmt.Println("unlock fail")
	} else {
		fmt.Println("unlock")
	}
}

func getUuid() string {
	output, err := exec.Command("uuidgen").Output()
	if err != nil {
		panic(err)
	}

	return string(output)
}
