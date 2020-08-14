package tools

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type CommandLine struct {
	Data []struct {
		DoOrSend string `json:"do_or_send"`
		Command  string `json:"command"`
		Params   string `json:"params"`
	} `json:"data"`
}

var conn redis.Conn

func init() {
	var err error
	address := os.Getenv("CACHE_ADDRESS")
	dbNum, _ := strconv.Atoi(os.Getenv("CACHE_DB_NUM"))
	conn, err = redis.Dial("tcp", address, redis.DialDatabase(dbNum))

	if err != nil {
		log.Println(err.Error())
		return
	}
}

func parseParams(params string) []interface{} {
	args := []interface{}{}
	if len(params) != 0 {
		paramSlice := strings.Split(params, ",")
		for _, param := range paramSlice {
			args = append(args, param)
		}
	}
	return args
}

func errorRespond(w http.ResponseWriter, err string){
	log.Println(err)
	http.Error(w, err, http.StatusBadRequest)
}

func ExecCloudRedis(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorRespond(w, err.Error())
		return
	}

	cmd := CommandLine{}
	if err := json.Unmarshal(data, &cmd); err != nil {
		errorRespond(w, err.Error())
		return
	}

	if len(cmd.Data) == 0 {
		errorRespond(w, "data is empty")
		return
	}

	for _, val := range cmd.Data {
		args := parseParams(val.Params)
		switch strings.ToLower(val.DoOrSend) {
		case "send":
			redis.Values(conn.Send(val.Command, args...), nil)
			continue
		case "do":
			reply, err := conn.Do(val.Command, args...)
			result, err := redis.Values(reply, err)
			if err != nil {
				if strings.Contains(err.Error(), "unexpected type for Values") {
					w.Write([]byte(fmt.Sprintf("%s\n", reply)))
				}else {
					errorRespond(w, err.Error())
				}
				return
			}
			var response string
			for _, str := range result {
				tmpStr, _ := redis.String(str, nil)
				response += tmpStr
				response += "\n"
			}
			w.Write([]byte(response))
			return
		default:
			errorRespond(w, "wrong param " + val.DoOrSend)
			return
		}
	}
}
