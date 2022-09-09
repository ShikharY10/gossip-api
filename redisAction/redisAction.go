package redisAction

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ShikharY10/goAPI/utils"
	"github.com/go-redis/redis"
)

type Redis struct {
	Client *redis.Client
}

func (r *Redis) Init(redisIP string) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisIP + ":6379",
		Password: "",
		DB:       0,
	})
	s := client.Ping()
	fmt.Println(s.String())

	r.Client = client
	fmt.Println("Redis client connected!")
}

func (r *Redis) SetUserData(id int, data map[string]interface{}) {
	key := strconv.Itoa(id) + "data"
	status := r.Client.HMSet(key, data)
	s, e := status.Result()
	if e != nil {
		panic(e)
	}
	fmt.Println("s: ", s)
}

func (r *Redis) GetUserIsOnline(id int) int {
	key := strconv.Itoa(id) + "data"
	value := r.Client.HMGet(key, "IsOnline")
	val, err := value.Result()
	if err != nil {
		panic(err)
	}

	v := val[0].(string)
	i, e := strconv.Atoi(v)
	if e != nil {
		panic(e)
	}
	return i
}

func (r *Redis) GetUserLastSeen(id int) string {
	key := strconv.Itoa(id) + "data"
	value := r.Client.HMGet(key, "LastSeen")
	val, err := value.Result()
	if err != nil {
		panic(err)
	}

	v := val[0].(string)
	return v
}

func (r *Redis) GetUserServerName(id int) int {
	key := strconv.Itoa(id) + "data"
	value := r.Client.HMGet(key, "servername")
	val, err := value.Result()
	if err != nil {
		panic(err)
	}

	v := val[0].(string)
	i, e := strconv.Atoi(v)
	if e != nil {
		panic(e)
	}
	return i
}

// func (r *Redis) GetClientData(id int) utils.UserData {
// 	key := strconv.Itoa(id) + "data"
// 	value := r.Client.HMGet(key, "IsOnline", "LastSeen", "servername")
// 	in := value.Val()
// 	var ud utils.UserData
// 	ud.LastSeen = in[1].(string)
// 	i, e := strconv.Atoi(in[0].(string))
// 	if e != nil {
// 		panic(e)
// 	}
// 	ud.IsOnline = i

// 	ir, er := strconv.Atoi(in[0].(string))
// 	if er != nil {
// 		panic(er)
// 	}
// 	ud.IsOnline = ir

// 	return ud
// }

func (r *Redis) GetUserMsg(id int) (string, error) {
	key := strconv.Itoa(id)
	s := r.Client.LPop(key)
	str, err := s.Result()
	if err != nil {
		panic(err)
	}
	return str, nil
}

func (r *Redis) SetUserMsg(id int, msg string) error {
	key := strconv.Itoa(id)
	s := r.Client.RPush(key, msg)
	_, e := s.Result()
	if e != nil {
		return e
	}
	return nil
}

func (r *Redis) RegisterOTP() (string, string) {
	id64 := utils.GenerateRandomId()
	otp := utils.GenerateOTP(6)
	print("id64: ", id64)
	r.Client.Set(id64, otp, time.Duration(5*time.Minute))
	return id64, otp
}

func (r *Redis) VarifyOTP(id string, otp string) bool {
	res := r.Client.Get(id)
	_otp := res.Val()
	if otp == _otp {
		return true
	} else {
		return false
	}
}

func (r *Redis) GetSecretekey() (string, error) {
	res := r.Client.Get("secretekey")
	if res.Err() != nil {
		return "", res.Err()
	}
	key := res.Val()
	return key, nil
}

func (r *Redis) SetSecretekey(key string) bool {
	res := r.Client.Set("secretekey", key, 0)
	return res.Err() == nil
}

func (r *Redis) GetEngineName() []string {
	fmt.Println("readig engines names")
	ress := r.Client.LRange("engines", 0, -1)
	fmt.Println("read completed: ", ress)
	engines, err := ress.Result()
	if err != nil {
		log.Println(err.Error())
		return []string{}
	}
	return engines
}
