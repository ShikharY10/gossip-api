package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ShikharY10/goAPI/gbp"
	"github.com/ShikharY10/goAPI/mongoAction"
	"github.com/ShikharY10/goAPI/redisAction"
	"github.com/ShikharY10/goAPI/rmq"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
)

type API_V2 struct {
	Mongo *mongoAction.Mongo
	Redis *redisAction.Redis
	RMQ   *rmq.RMQ
}

func (a *API_V2) SendOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("service", "Gossip API")
	params := mux.Vars(r)
	number := params["number"]

	fmt.Println("number: ", number)

	id, otp := a.Redis.RegisterOTP()
	var otpData map[string]string = map[string]string{
		"otp":    otp,
		"number": number,
	}
	b, _ := json.Marshal(otpData)
	a.RMQ.Produce("OTPd3hdzl8", b)

	var response gbp.Response
	response.Status = true
	response.Disc = "otp send"
	response.Data = id
	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		log.Println("[marshal error]", err.Error())
	}
	w.Write(responseBytes)
}
