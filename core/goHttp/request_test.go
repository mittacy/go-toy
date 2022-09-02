package goHttp

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/mittacy/go-toy/core/log"
	"github.com/mittacy/go-toy/core/response"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"net/http"
	"sync"
	"testing"
	"time"
)

var (
	startServer sync.Once
)

type Student struct {
	Name string
	Age  int
}

func startHttpServerAndLog() {
	log.Init(log.WithPath("./"),
		log.WithTimeFormat("2006-01-02 15:04:05"),
		log.WithLevel(log.DebugLevel),
		log.WithEncoderJSON(true),
		log.WithLogInConsole(true),
		log.WithDistinguish(true))

	r := gin.New()
	r.GET("/info", func(c *gin.Context) {
		data := c.Query("name")
		if data == "" {
			data = "mittacy"
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": data,
		})
	})
	r.GET("/student/info", func(c *gin.Context) {
		student := Student{
			Name: c.Query("name"),
			Age:  18,
		}

		response.Success(c, student)
	})
	r.POST("/student/create", func(c *gin.Context) {
		var student Student
		if err := c.ShouldBindJSON(&student); err != nil {
			response.ValidateErr(c, err)
			return
		}

		response.Success(c, student)
	})

	r.Run(":10110")
}

func TestSimpleGet(t *testing.T) {
	go startServer.Do(startHttpServerAndLog)

	time.Sleep(time.Second)
	client := http.Client{}
	Init(client, "thirdHttp")

	hc := NewApiServer("http://127.0.0.1:10110")
	res := struct {
		Code int
		Msg  string
		Data string
	}{}
	headers := map[string]string{
		"union_id": "sdfasfd",
	}
	name := "yoyo"
	params := map[string]string{
		"name": name,
	}
	ctx := context.WithValue(context.Background(), b3.TraceID, "92934223")

	err := hc.SimpleGet(ctx, "/info", &res, WithAddHeaders(headers), WithFormParams(params))
	if err != nil {
		t.Fatalf("simple get err: %+v\n", err)
	}

	t.Logf("simple get success, data: %+v\n", res)
}

func TestSimplePost(t *testing.T) {
	go startServer.Do(startHttpServerAndLog)

	client := http.Client{}
	Init(client, "thirdHttp")

	hc := NewApiServer("http://127.0.0.1:10110")
	res := struct {
		Code int
		Msg  string
		Data interface{}
	}{}
	headers := map[string]string{
		"union_id": "sdfasfd",
	}
	name := "yoyo"
	params := map[string]string{
		"name": name,
	}
	body := map[string]interface{}{
		"name": name,
		"age":  12,
	}
	ctx := context.WithValue(context.Background(), b3.TraceID, "8abfaj8df")

	err := hc.SimplePost(ctx, "/student/create", &res, WithAddHeaders(headers), WithFormParams(params), WithJSONData(body))
	if err != nil {
		t.Fatalf("simple post err: %+v\n", err)
	}

	t.Logf("simple post success, data: %+v\n", res)
}

func TestGet(t *testing.T) {
	go startServer.Do(startHttpServerAndLog)

	client := http.Client{}
	Init(client, "thirdHttp")

	hc := NewApiServer("http://127.0.0.1:10110")
	headers := map[string]string{
		"union_id": "sdfasfd",
	}
	name := "yoyo"
	params := map[string]string{
		"name": name,
	}
	ctx := context.WithValue(context.Background(), b3.TraceID, "92934223")

	var res Student
	code, err := hc.Get(ctx, "/student/info", &res, WithAddHeaders(headers), WithFormParams(params))
	if err != nil {
		t.Fatalf("get err: %+v\n", err)
	}

	t.Logf("get success, code: %d, data: %+v\n", code, res)
}

func TestPost(t *testing.T) {
	go startServer.Do(startHttpServerAndLog)

	client := http.Client{}
	Init(client, "thirdHttp")

	hc := NewApiServer("http://127.0.0.1:10110")
	headers := map[string]string{
		"union_id": "sdfasfd",
	}
	name := "yoyo"
	params := map[string]string{
		"name": name,
	}
	body := map[string]interface{}{
		"name": name,
		"age":  12,
	}
	ctx := context.WithValue(context.Background(), b3.TraceID, "8abfaj8df")

	var res Student
	code, err := hc.Post(ctx, "/student/create", &res, WithAddHeaders(headers), WithFormParams(params), WithJSONData(body))
	if err != nil {
		t.Fatalf("post err: %+v\n", err)
	}

	t.Logf("post success, code: %d, data: %+v\n", code, res)
}

func TestApiOptions(t *testing.T) {
	go startServer.Do(startHttpServerAndLog)

	client := http.Client{}
	Init(client, "thirdHttp")

	hc := NewApiServer("http://127.0.0.1:10110", WithLogName("customLogName"))
	headers := map[string]string{
		"union_id": "sdfasfd",
	}
	name := "yoyo"
	params := map[string]string{
		"name": name,
	}
	body := map[string]interface{}{
		"name": name,
		"age":  12,
	}
	ctx := context.WithValue(context.Background(), b3.TraceID, "8abfaj8df")

	var res Student
	code, err := hc.Post(ctx, "/student/create", &res, WithAddHeaders(headers), WithFormParams(params), WithJSONData(body))
	if err != nil {
		t.Fatalf("post err: %+v\n", err)
	}

	t.Logf("post success, code: %d, data: %+v\n", code, res)
}

func TestTransport(t *testing.T) {
	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   500, // 单个Host的最大空闲连接数
			MaxConnsPerHost:       500, // 单个Host的最大连接总数，>=MaxIdleConnsPerHost
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: time.Second * 5,
	}
	Init(client, "thirdHttp")

	hc := NewApiServer("http://134.175.50.206:1028")

	startTime := time.Now()
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			if err := hc.SimpleGet(context.Background(), "", nil); err != nil {
				panic(fmt.Sprintf("transport get err: %+v\n", err))
			} else {
				wg.Done()
			}
		}()
	}

	wg.Wait()
	takeTime := time.Now().Sub(startTime)

	t.Logf("transport get success, take time: %+v\n", takeTime)
}
