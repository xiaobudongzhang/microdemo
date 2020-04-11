package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	user "user-web/proto/user"

	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/v2/client"
)

var (
	serviceClient user.UserService
)

type Error struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func Init() {
	serviceClient = user.NewUserService("mu.micro.book.service.user", client.DefaultClient)
}

func UserCall(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	userClient := user.NewUserService("mu.micro.book.service.user", client.DefaultClient)
	rsp, err := userClient.Call(context.TODO(), &user.Request{
		Name: request["name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Logf("非法请求")
		http.Error(w, "非法请求", 400)
		return
	}

	r.ParseForm()

	rsp, err := serviceClient.QueryUserByName(context.TODO(), &user.UserRequest{
		UserName: r.Form.Get("userName"),
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	response := map[string]interface{}{
		"ref": time.Now().UnixNano(),
	}

	if rsp.User.Pwd == r.Form.Get("pwd") {
		response["sucees"] = true

		rsp.User.Pwd = ""
		response["data"] = rsp.User
	} else {
		response["success"] = false
		response["error"] = &Error{
			Detail: "密码错误",
		}
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
