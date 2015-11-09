# geetest
极验golang sdk

# 使用方法
- 主要配置是geetest.Config对象.首先我们需要添加自己的公钥和密钥,公钥请配置在geetest.Config.CaptchaId,密钥配置在geetest.Config.PrivateKey.
- geetest.GeetestLib 是主要的操作方法集合,CheckServerStatus 校验服务器状态,GenerateChallenge 用于生成challenge,Valid 校验验证码是否正确.您可能需要将GeetestLib对象保存的session,在校验验证码是否正确的时候将其取出,调用Valid的方法.如果您不愿意保存该对象,那么可以使用ValidChallenge函数自行校验验证码是否正确.

# 示例
在samples文件夹中有相关使用示例,`go run server.go`,然后在浏览器访问[localhost:8080](http://localhost:8080)即可访问.

# TODO
[] https 支持
[] debug信息打印
