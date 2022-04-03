package main

import (
	"encoding/json"
	"fmt"
	"github.com/bingerambo/go-file-json-server/utils"
	"net/http"
	"path/filepath"
)

// Server represents a simple-upload server.
type JsonServer struct {
	DocumentRoot     string
	EnableCORS       bool
	ProtectedMethods []string
	Sjp              *utils.SimpleJsonParser
	cache            *Cache
}

// NewServer creates a new simple-upload server.
func NewJsonServer(documentRoot string, enableCORS bool) JsonServer {
	return JsonServer{
		DocumentRoot: documentRoot,
		EnableCORS:   enableCORS,
		Sjp:          utils.NewSimpleJsonParse(),
		cache:        NewCache(documentRoot),
	}
}

func (s JsonServer) Start() {
	s.cache.Boot()
}

func (s JsonServer) HandleGet(w http.ResponseWriter, r *http.Request) {
	s.handleGet(w, r)
}

func (s JsonServer) handleGet(w http.ResponseWriter, r *http.Request) {

	if s.EnableCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	q := r.URL.Query()

	//if len(q) != 0 {
	//	// 方式1：通过字典下标取值
	//	fmt.Println("id1:", q["id"][0])
	//}
	// 方式2：使用Get方法，如果没有值会返回空字符串
	fmt.Println("this get request params[name]: ", q.Get("name"))
	if "" == q.Get("name") {
		failedResp(w, 400, "GET request paramers is empty, failed")
		return
	}

	//Content-Type 用于描述本次请求的body的内容是json格式，且编码为UTF-8
	//Accept 用于描述客户端希望返回的结果以json来组织，且UTF-8
	// Accept:application/json;charset=UTF-8
	//Content-Type 用于描述request,而Accept用于描述reponse
	w.Header().Set("content-type", "application/json;charset=UTF-8")
	w.Header().Set("accept", "application/json;charset=UTF-8")
	// nosuffix
	file_name := q.Get("name")

	// not in cache， failed
	if !s.cache.data.Set().Contains(file_name) {
		err_msg := "request json not exist, failed"
		logger.Error(err_msg)
		failedResp(w, 420, err_msg).WriteHeader(http.StatusInternalServerError)
		return
	}

	//file_path := "github.com/bingerambo/go-file-json-server/tmp/sample.json"
	file_path := filepath.Join(s.cache.watchRoot, file_name+JSON_FILESUFFIX)

	datas, err := s.Sjp.Load(file_path)
	if err != nil {
		logger.Error(err)
		failedResp(w, 440, "load json, failed").WriteHeader(http.StatusInternalServerError)
		return
	}

	js, err := s.Sjp.Parse(datas)
	if err != nil {
		logger.Error(err)
		failedResp(w, 450, "parse json, failed").WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("========== get json content start ==========")
	fmt.Println(js)
	fmt.Println("========== get json content end ==========")

	w.Write(datas)
	w.WriteHeader(http.StatusOK)

}

func (s JsonServer) handlePost(w http.ResponseWriter, r *http.Request) {

}

func (s JsonServer) handlePut(w http.ResponseWriter, r *http.Request) {

}

func (s JsonServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet, http.MethodHead:
		s.handleGet(w, r)
	case http.MethodPost:
		//s.handlePost(w, r)
	case http.MethodPut:
		//s.handlePut(w, r)
	case http.MethodOptions:
		//s.handleOptions(w, r)
	default:
		//w.Header().Add("Allow", "GET,HEAD,POST,PUT")
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.WriteError(w, fmt.Errorf("method \"%s\" is not allowed", r.Method))
	}
}

type JsonResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func failedResp(w http.ResponseWriter, code int, msg_content string) http.ResponseWriter {
	msg, _ := json.Marshal(JsonResult{Code: code, Msg: msg_content})
	w.Header().Set("content-type", "text/json")
	w.Write(msg)
	return w
}
