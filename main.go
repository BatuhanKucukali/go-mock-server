package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template/parse"
	"time"
)

type Contract struct {
	Maps []Mapping `json:"mappings"`
	Auth *Auth     `json:"auth"`
}

type Mapping struct {
	Req  *Request  `json:"request"`
	Resp *Response `json:"response"`
}

type Request struct {
	Method string `json:"method"`
	Url    string `json:"url"`
}

type Response struct {
	Status       uint16            `json:"status"`
	FixedDelayMs uint64            `json:"fixedDelayMilliseconds"`
	Body         string            `json:"body"`
	JsonBody     parse.Tree        `json:"jsonBody"` // TODO json parse
	Headers      map[string]string `json:"headers"`
}

type Auth struct {
	BasicAuth *BasicAuth `json:"basicAuthCredentials"`
	TokenAuth *TokenAuth `json:"tokenCredentials"`
}

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenAuth struct {
	Token string `json:"token"`
}

func main() {
	fmt.Println("Go Mock Server Started...")

	file := read("contract.json")

	c := contracts(file)

	handlers(c.Maps)

	// TODO add auth (basic & token)

	_ = http.ListenAndServe(":8090", nil)
}

func handlers(maps []Mapping) {
	for i := 0; i < len(maps); i++ {
		m := maps[i]

		// TODO add HTTP method
		http.HandleFunc(m.Req.Url, func(w http.ResponseWriter, req *http.Request) {

			sleep(m.Resp.FixedDelayMs)

			for k, v := range m.Resp.Headers {
				w.Header().Add(k, v)
			}
			w.WriteHeader(int(m.Resp.Status))
			_, _ = w.Write([]byte(m.Resp.Body)) // TODO add body || jsonBody
		})

	}
}

func sleep(millisecond uint64) {
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
}

func contracts(file []byte) Contract {
	var c Contract

	err := json.Unmarshal(file, &c)
	if err != nil {
		log.Fatal("Contract file is not serializable. Error: ", err)
	}
	return c
}

func read(fileName string) []byte {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Contract file does not read. Error:", err)
	}
	return file
}
