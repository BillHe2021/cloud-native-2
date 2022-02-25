package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

//自定义返回
type JsonRes struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	TimeStamp int64       `json:"timestmap"`
}

func apiResult(w http.ResponseWriter, code int, data interface{}, msg string) {
	body, _ := json.Marshal(JsonRes{
		Code: code,
		Data: data,
		Msg:  msg,
		// 获取时间戳
		TimeStamp: time.Now().Unix(),
	})
	w.Write(body)
	println(w.Write(body))
}

func main() {
	srv := http.Server{
		Addr:    ":8080",
		Handler: http.TimeoutHandler(http.HandlerFunc(defaultHttp), 2*time.Second, "Timeout!!!"),
	}
	srv.ListenAndServe()
}

// 默认http处理
func defaultHttp(w http.ResponseWriter, r *http.Request) {
	path, httpMethod := r.URL.Path, r.Method
	ipaddress := r.RemoteAddr
	fmt.Printf("ipaddress is:%s\n", ipaddress)

	if path == "/healthz" {
		w.Write([]byte("This is my first web server,Welcome"))
		w.Header().Set("Go version", runtime.Version())

		for k, v := range w.Header() {
			fmt.Printf("%s=%s\n", k, v)

		}
		//println(runtime.Version())
		fmt.Printf("response health code is :%d\n", http.StatusOK)
		loggerv1(ipaddress, http.StatusOK)
		return
	}

	if path == "/welcome" && httpMethod == "POST" {
		sayHello(w, r)
		return
	}

	if path == "/sleep" {
		// 模拟一下业务处理超时
		time.Sleep(4 * time.Second)
		return
	}

	if path == "/path" {
		w.Write([]byte("path:" + path + ", method:" + httpMethod))
		return
	}

	// 自定义404
	http.Error(w, "page not exits", http.StatusNotFound)
}

// 处理hello，并接收参数输出json
func sayHello(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// 第一种方式，但是没有name参数会报错
	// name := query["name"][0]

	// 第二种方式
	name := query.Get("name")
	println(name)

	apiResult(w, 0, name+" say "+r.PostFormValue("some"), "success")
}

type logcontent struct {
	ipadd      string
	httpstatus int
}

/*
func logger(ip string, status int ...) {
	logfile, err := os.OpenFile("/Users/bill/go/src/practice/httpserver/test.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	defer logfile.Close()
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	logger.Println(ip)
	logger.Println(status)
}
*/
//重写logger

func loggerv1(args ...interface{}) {
	logfile, err := os.OpenFile("/Users/bill/go/src/practice/httpserver/test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	logger := log.New(logfile, "\n", log.Ldate|log.Ltime|log.Llongfile)
	defer logfile.Close()
	for _, arg := range args {
		switch arg.(type) {
		case int:
			logger.Println(arg)
		case string:
			logger.Println(arg)
		case int64:
			logger.Println(arg)
		default:
			logger.Println(arg)
		}

	}
}
