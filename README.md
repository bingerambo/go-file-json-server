# Go-file-json-server
Simple HTTP server to uploading or downloading files, and fetch the uploading files which is json format. 

**The project is based on existing project** **[go-simple-upload-server](https://github.com/mayth/go-simple-upload-server)**

**The project features：**

- **File upload and download.**
- **Watch the upload files automaticly,  and filter the JSON file to JSON Server**
- **JSON Server provides the api to get json content**



实现了简单的文件和JSON的http服务，特点:

- 简单的上传和下载使用的文件服务器
- 增加了文件状态缓存机制，实时监控上传文件的变化，用于提供可用json数据访问服务
- JSON服务，是个mock接口服务，根据上传的json文件提供json内容访问，方便业务模拟接口测试。
- 场景：业务开发需调用第三方接口，接口数据json格式和内容已定义，但实际上开发环境中无可用和可调试的第三方接口服务，这时可以使用这个工具上传接口数据json文件，然后再调用json服务的mock接口就可以进行桩测试。



# Usage

## Start Server

```
$ mkdir $HOME/tmp
$ ./go-file-json-server -token f9403fc5f537b4ab332d $HOME/tmp
```

(see "Security" section below for `-token` option)

## Uploading

You can upload files with `POST /upload`.
The filename is taken from the original file if available. If not, SHA1 hex digest will be used as the filename.

```
$ echo 'Hello, world!' > sample.txt
$ curl -Ffile=@sample.txt 'http://localhost:25478/upload?token=f9403fc5f537b4ab332d'
{"ok":true,"path":"/files/sample.txt"}
```

```shell script

```

```
$ cat $HOME/tmp/sample.txt
hello, world!
```

**OR**

Use `PUT /files/(filename)`.
In this case, the original file name is ignored, and the name is taken from the URL.

```
$ curl -X PUT -Ffile=@sample.txt "http://localhost:25478/files/another_sample.txt?token=f9403fc5f537b4ab332d"
{"ok":true,"path":"/files/another_sample.txt"}
```

## Downloading

`GET /files/(filename)`.

```
$ curl 'http://localhost:25478/files/sample.txt?token=f9403fc5f537b4ab332d'
hello, world!
```

## Existence Check

`HEAD /files/(filename)`.

```
$ curl -I 'http://localhost:25478/files/foobar.txt?token=f9403fc5f537b4ab332d'
HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 9
Content-Type: text/plain; charset=utf-8
Last-Modified: Sun, 09 Oct 2016 14:35:39 GMT
Date: Sun, 09 Oct 2016 14:35:43 GMT

$ curl 'http://localhost:25478/files/foobar.txt?token=f9403fc5f537b4ab332d'
hello!!!

$ curl -I 'http://localhost:25478/files/unknown?token=f9403fc5f537b4ab332d'
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sun, 09 Oct 2016 14:37:48 GMT
Content-Length: 19
```


## CORS Preflight Request

* `OPTIONS /files/(filename)`
* `OPTIONS /upload`

```
$ curl -I 'http://localhost:25478/files/foo'
HTTP/1.1 204 No Content
Access-Control-Allow-Methods: PUT,GET,HEAD
Access-Control-Allow-Origin: *
Date: Sun, 06 Sep 2020 09:45:20 GMT

$ curl -I -XOPTIONS 'http://localhost:25478/upload'
HTTP/1.1 204 No Content
Access-Control-Allow-Methods: POST
Access-Control-Allow-Origin: *
Date: Sun, 06 Sep 2020 09:45:32 GMT
```

notes:

* Requests using `*` as a path, like as `OPTIONS * HTTP/1.1`, are not supported.
* On sending `OPTIONS` request, `token` parameter is not required.
* For `/files/(filename)` request, server replies "204 No Content" even if the specified file does not exist.



## JSON Server 

users can fetch json data from JSON Server .

1. upload the json file : sample.json

```shell
curl -Ffile=@sample.json 'http://localhost:25478/upload?token=f9403fc5f537b4ab332d'

```

curl 'http://127.0.0.1:25478/mock?name=sample'

2. then vist mock api to get json file: sample.json

eg. http://127.0.0.1:25478/mock?name=xxx

```shell

$ curl 'http://127.0.0.1:25478/mock?name=sample'
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   194  100   194    0     0   189k      0 --:--:-- --:--:-- --:--:--  189k{
 "code": "200",
 "message": "Return Successd!",
 "result": {
  "province": "test",
  "city": "example",
  "areacode": "0571",
  "zip": "310000",
  "company": "example",
  "card": ""
 }
}

```

also，you can open browser  to vist http://127.0.0.1:25478/mock?name=xxx

sample.json

```json
{
 "code": "200",
 "message": "Return Successd!",
 "result": {
  "province": "浙江",
  "city": "杭州",
  "areacode": "0571",
  "zip": "310000",
  "company": "中国移动",
  "card": ""
 }
}

```





# TLS

To enable TLS support, add `-cert` and `-key` options:

```
$ ./simple_upload_server -cert ./cert.pem -key ./key.pem root/
INFO[0000] starting up simple-upload-server
WARN[0000] token generated                               token=28d93c74c8589ab62b5e
INFO[0000] start listening TLS                           cert=./cert.pem key=./key.pem port=25443
INFO[0000] start listening                               ip=0.0.0.0 port=25478 root=root token=28d93c74c8589ab62b5e upload_limit=5242880
...
```

This server listens on `25443/tcp` for TLS connections by default. This can be changed by passing `-tlsport` option.

NOTE: The endpoint using HTTP is still active even if TLS is enabled.


# Security

## Token

There is no Basic/Digest authentication. This app implements dead simple authentication: "security token".

All requests should have `token` parameter (it can be passed as a query string or a form parameter). The server accepts the request only when the token is matched; otherwise, the server rejects the request and respond `401 Unauthorized`.

You can specify the server's token on startup by `-token` option. If you don't so, the server generates the token and writes it to STDOUT at WARN level log, like as:

```
$ ./simple_upload_server root
INFO[0000] starting up simple-upload-server
WARN[0000] token generated                               token=2dd30b90536d688e19f7
INFO[0000] start listening                               ip=0.0.0.0 port=25478 root=root token=2dd30b90536d688e19f7 upload_limit=5242880
```

NOTE: The token is generated from the random number, so it will change every time you start the server.

## CORS

If you enable CORS support using `-cors` option, the server append `Access-Control-Allow-Origin` header to the response. This feature is disabled by default.

# Docker

```
$ docker run -p 25478:25478 -v $HOME/tmp:/var/root mayth/simple-upload-server -token f9403fc5f537b4ab332d /var/root
```

# Example

go-file-json-server program running start

```shell
# windows
./go-file-json-server_windows_amd64 -token f9403fc5f537b4ab332d ./tmp

# linux
./go-file-json-server_linux_amd64 -token f9403fc5f537b4ab332d ./tmp
```


go-file-json-server program running printf

```shell script
$ ./go-file-json-server_windows_amd64 -token f9403fc5f537b4ab332d ./tmp
time="2022-03-30T15:18:43+08:00" level=info msg="starting up file-json-server, for upload file and fetch json"
start to monitor:  github.com\bingerambo\go-file-json-server\tmp
 start syncCacheServer...
 start syncCacheServer...
update files[CREATE] :  error
 start syncCacheServer...
update files[CREATE] :  no_json
 start syncCacheServer...
update files[CREATE] :  sample
time="2022-03-30T15:18:43+08:00" level=info msg="file json server start ok"
time="2022-03-30T15:18:43+08:00" level=info msg="start listening" cors=false ip=0.0.0.0 port=25478 protected_method="[GET POST HEAD PUT]" root=./tmp token=f9403fc5f537b4ab332d upload_limit=5242880
this get request params[name]:  error
write file:  D:\GO_projects\src\github.com\bingerambo\go-file-json-server\tmp\error.json
 start syncCacheServer...
========== get json content start ==========
&{map[code:-5000 message:error message]}
========== get json content end ==========
2022/03/30 15:18:57 http: superfluous response.WriteHeader call from main.JsonServer.handleGet (json_server.go:95)
2022/03/30 15:19:05 httpserver is exiting...
2022/03/30 15:19:05 notify sigs
2022/03/30 15:19:05 httpserver shutdown...
time="2022-03-30T15:19:05+08:00" level=info msg="file json server exit ok"

```