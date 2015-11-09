package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/sdvdxl/Go-Geetest/geetest"
	"net/http"
	_ "net/http"
	"encoding/gob"
)

const (
	SessionKey = "session"
)

var sessionStore = sessions.NewCookieStore([]byte("something-very-secret"))

func main() {
	gob.Register(geetest.GeeTestLib{})

	geetest.Config.CaptchaId = "157e7df54d8deb46238cef3c5848a2bf"
	geetest.Config.PrivateKey = "f48a9f88c30f4f01696d96ea0d220f98"



	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/validate", func(writer http.ResponseWriter, req *http.Request) {
		session, err := sessionStore.Get(req, SessionKey)
		if err != nil {
			panic(err)
		}

		result := JsonResult{}
		geetestLib, ok := session.Values["_geetest"].(geetest.GeeTestLib)
		if !ok {
			fmt.Println(err)
			result.Msg = "系统错误"
		} else {
			if geetestLib.Valid(req.PostFormValue("geetest_challenge") , req.PostFormValue("geetest_validate")) {
				result.Msg = "验证成功"
				result.Success = true
			} else {
				result.Msg = "验证失败"
			}
		}

		resultBytes, _ := json.Marshal(result)
		writer.Write(resultBytes)
	})

	http.HandleFunc("/getChallenge", func(writer http.ResponseWriter, req *http.Request) {
		session, err := sessionStore.Get(req, SessionKey)
		if err != nil {
			panic(err)
		}

		geetestLib := geetest.GeeTestLib{}
		fmt.Println(geetestLib.GenerateChallenge())
		session.Values["_geetest"] = geetestLib
		result := JsonResult{}
		if err = session.Save(req, writer); err != nil {
			fmt.Println(err)
			result.Msg = "系统错误"
		} else {
			result.Msg = "成功"
			result.Data = geetestLib.Challenge
			fmt.Println("challenge:", geetestLib.Challenge)
			result.Success = true
		}
		resultBytes, _ := json.Marshal(result)
		writer.Write(resultBytes)
	})

	http.ListenAndServe(":8080", nil)
}

type JsonResult struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Data    interface{} `json:"data"`
}
