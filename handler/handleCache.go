package handler

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/ShikharY10/gbAPI/logger"
	"github.com/ShikharY10/gbAPI/utils"
	"github.com/go-redis/redis"
)

type CacheHandler struct {
	Client *redis.Client
	logger *logger.Logger
}

func CreateCacheHandler(client *redis.Client, logger *logger.Logger) *CacheHandler {
	return &CacheHandler{
		Client: client,
		logger: logger,
	}
}

func (r *CacheHandler) GetEngineChannel() (string, error) {
	names := r.GetEngineName()
	if len(names) == 0 {
		return "", errors.New("no engine found")
	}
	randomIndex := rand.Intn(len(names))
	pick := names[randomIndex]
	return pick, nil
}

func (r *CacheHandler) SetUserData(id int, data map[string]interface{}) {
	key := strconv.Itoa(id) + "data"
	status := r.Client.HMSet(key, data)
	s, e := status.Result()
	if e != nil {
		panic(e)
	}
	fmt.Println("s: ", s)
}

func (r *CacheHandler) GetUserIsOnline(id int) int {
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

func (r *CacheHandler) GetUserLastSeen(id int) string {
	key := strconv.Itoa(id) + "data"
	value := r.Client.HMGet(key, "LastSeen")
	val, err := value.Result()
	if err != nil {
		panic(err)
	}

	v := val[0].(string)
	return v
}

func (r *CacheHandler) GetUserServerName(id int) int {
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

// func (r *CacheHandler) GetUserMsg(id int) (string, error) {
// 	key := strconv.Itoa(id)
// 	s :
// 	var r CacheHandler
// 	r.Client = client
// 	return r= r.Client.LPop(key)
// 	str, err := s.Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return str, nil
// }

func (r *CacheHandler) SetUserMsg(id int, msg string) error {
	key := strconv.Itoa(id)
	s := r.Client.RPush(key, msg)
	_, e := s.Result()
	if e != nil {
		return e
	}
	return nil
}

func (r *CacheHandler) RegisterOTP() (string, string) {
	id64 := utils.GenerateRandomId()
	otp := utils.GenerateOTP(6)
	r.Client.Set(id64, otp, time.Duration(5*time.Minute))
	return id64, otp
}

func (r *CacheHandler) VarifyOTP(id string, otp string) bool {
	res := r.Client.Get(id)
	_otp := res.Val()
	if otp == _otp {
		return true
	} else {
		return false
	}
}

func (r *CacheHandler) GetSecretekey() (string, error) {
	res := r.Client.Get("secretekey")
	if res.Err() != nil {
		return "", res.Err()
	}
	key := res.Val()
	return key, nil
}

func (r *CacheHandler) SetSecretekey(key string) bool {
	res := r.Client.Set("secretekey", key, 0)
	return res.Err() == nil
}

func (r *CacheHandler) GetEngineName() []string {
	fmt.Println("reading engines names")
	ress := r.Client.LRange("engines", 0, -1)
	fmt.Println("read completed: ", ress)
	engines, err := ress.Result()
	if err != nil {
		log.Println(err.Error())
		return []string{}
	}
	return engines
}
