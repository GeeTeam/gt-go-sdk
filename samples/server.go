package main

import (
	"encoding/json"
	"fmt"
	"github.com/sdvdxl/geetest"
	"net/http"
	_ "net/http"
)

var (
	cacheMap   = map[string]interface{}{}
	geetestKey = "_geetest"
)

func main() {
	fmt.Println("server is starting")
	geetest.Config.CaptchaId = "157e7df54d8deb46238cef3c5848a2bf"
	geetest.Config.PrivateKey = "f48a9f88c30f4f01696d96ea0d220f98"

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/validate", func(writer http.ResponseWriter, req *http.Request) {

		result := JsonResult{}
		geetestLib, ok := cacheMap["_geetest"].(geetest.GeeTestLib)
		if !ok {
			result.Msg = "系统错误"
		} else {
			if geetestLib.Valid(req.PostFormValue("geetest_challenge"), req.PostFormValue("geetest_validate")) {
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
		geetestLib := geetest.GeeTestLib{}
		fmt.Println(geetestLib.GenerateChallenge())
		cacheMap[geetestKey] = geetestLib
		result := JsonResult{}
		result.Msg = "成功"
		result.Data = geetestLib.Challenge
		fmt.Println("challenge:", geetestLib.Challenge)
		result.Success = true
		resultBytes, _ := json.Marshal(result)
		writer.Write(resultBytes)
	})

	fmt.Println("server starged and listen on 8080")
	http.ListenAndServe(":8080", nil)
}

type JsonResult struct {
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
