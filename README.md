# Baofeng Cloud UGC Site Demo 


安装NSQ
从[NSQ网站](http://nsq.io/) 下载NSQ，并运行在一台公网服务器上。


配置文件conf.json
``` json
{"AccessKey":"",
"SecretKey":"",
"CallbackUrl" : "http://NSQD服务器地址:4151/pub?topic=bfcloud",
"NsqdAddress":"NSQD服务器地址:4150",
"DataPath" : "上传文件的保存路径"
}
```

# 关于

基于 [暴风云视频GO-SDK](https://github.com/baofengcloud/go-sdk) 构建。