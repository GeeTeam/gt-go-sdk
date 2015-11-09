package geetest
import (
	"net/http"
	"io/ioutil"
	"strings"
	"log"
	"crypto/md5"
	"encoding/hex"
)

//极验配置
type geetestConfig struct {
	CaptchaId       string //应用id
	PrivateKey      string //应用密钥
	IsHttps         bool   //是否是https
	Debug           bool   //是否开启debug模式
	ServerStatusUrl string
	RegisterUrl     string
}

func init() {
	Config.ServerStatusUrl = "http://api.geetest.com/check_status.php"
	Config.RegisterUrl = "http://api.geetest.com/register.php?gt="

}

var Config geetestConfig

type GeeTestLib struct {
	Challenge string //Challenge
}


// 校验服务器是否正常
func (self GeeTestLib)CheckServerStatus() bool {
	resp, err := http.Get(Config.ServerStatusUrl)


	if err != nil {
		log.Println(err)
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false
	}

	result := getResponseResult(resp)

	if strings.EqualFold(result, "ok") {
		return true
	}

	return false
}

// 生成challenge
func (self *GeeTestLib)GenerateChallenge() (string, error) {
	resp, err := http.Get(Config.RegisterUrl + Config.CaptchaId)
	if err != nil {
		return "", err
	}

	self.Challenge = getResponseResult(resp)
	return self.Challenge, nil
}

//校验
func (self GeeTestLib)Valid(challenge, secCode string) bool {
	if len(challenge) != 34 || challenge[:32] != self.Challenge {
		return false
	}

	return strings.EqualFold(md5Encode(Config.PrivateKey + "geetest" + challenge), secCode)
}

func md5Encode(text string) string {
	result := md5.Sum([]byte(text))
	return hex.EncodeToString(result[:])
}


func getResponseResult(resp *http.Response) string {
	if resp == nil {
		return ""
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(bodyBytes)

}