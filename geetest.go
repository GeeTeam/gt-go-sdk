package geetest

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//极验配置
type geetestConfig struct {
	CaptchaId       string //应用captcha id 必填
	PrivateKey      string //应用密钥  必填
	IsHttps         bool   //是否是https 可选,如果是https,则设置为true, 暂未实现
	Debug           bool   //是否开启debug模式 暂未实现
	ServerStatusUrl string //服务器状态校验url, 可选
	RegisterUrl     string //注册获取challenge地址 可选
	ServerValidUrl  string //二次验证地址
	ServerValid     bool   // 二次验证,向服务器发起验证,默认为false
}

// 初始化基本配置
func init() {
	Config.ServerValidUrl = "http://api.geetest.com/validate.php"
	Config.ServerStatusUrl = "http://api.geetest.com/check_status.php"
	Config.RegisterUrl = "http://api.geetest.com/register.php?gt="
}

// 极验配置项
var Config geetestConfig

type GeeTestLib struct {
	Challenge string //Challenge
}

// 校验服务器是否正常
// 如果服务器正常,且返回正确状态(ok),则返回true,否则false
func (self GeeTestLib) CheckServerStatus() bool {
	resp, err := http.Get(Config.ServerStatusUrl)

	if err != nil {
		//log.Println(err)
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
func (self *GeeTestLib) GenerateChallenge() (string, error) {
	resp, err := http.Get(Config.RegisterUrl + Config.CaptchaId)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	self.Challenge = getResponseResult(resp)
	return self.Challenge, nil
}

// 校验验证码是否正确
// 由于极验challenge传回服务器会自动在后面加上2位随机字母,所以需要传回后台
// 传回极验服务器的challenge=challenge+2字母,所以是34位,否则校验失败
// 根据极验服务器加密方式,将从极验服务器得到的加密码和后台加密码对比,
// 如果相同则校验成功,返回true, 否则false
// challenge 前端传过来的challenge, 默认是 geetest_challenge 参数
// validateCode 前端传过来的加密后的值,默认是 geetest_validate 参数
// secCode 用于二次验证,如果开启了2次验证请填写
func (self GeeTestLib) Valid(challenge, validateCode string, secCode ...string) (bool, error) {
	return ValidChallenge(challenge, self.Challenge, validateCode, secCode...)
}

// 用于校验验证码, 和Valid方法功能相同,但是允许自行传入之前生成的challenge进行校验
// frontChallenge 前端传过来的34位challenge
// backChallenge 后台生成的challenge
// validateCode 前台传过来的加密校验码
// secCode 用于二次验证,如果开启了2次验证请填写
func ValidChallenge(frontChallenge, backChallenge, validateCode string, secCode ...string) (bool, error) {
	if len(frontChallenge) != 34 || frontChallenge[:32] != backChallenge {
		return false, nil
	}

	if md5sum(Config.PrivateKey+"geetest"+frontChallenge) != validateCode {
		return false, nil
	}

	//二次验证
	if Config.ServerValid {
		if len(secCode) == 0 {
			return false, nil
		}

		params := make(url.Values, 1)
		params.Add("seccode", secCode[0])
		resp, err := http.PostForm(Config.ServerValidUrl, params)
		if err != nil {
			return false, err
		}

		defer resp.Body.Close()

		result := getResponseResult(resp)
		if md5sum(secCode[0]) == result {
			return true, nil
		}
	}

	return false, nil
}

// 计算字符串的md5
func md5sum(text string) string {
	result := md5.Sum([]byte(text))
	return hex.EncodeToString(result[:])
}

// 处理http返回内容
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
