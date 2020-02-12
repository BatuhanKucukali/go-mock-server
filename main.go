package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Contract struct {
	Req  *Request  `json:"request"`
	Res  *Response `json:"response"`
	Auth *Auth     `json:"auth"`
}

type Request struct {
	Method string `json:"method"`
	Url    string `json:"url"`
}

type Response struct {
	Status       uint16            `json:"status"`
	FixedDelayMs uint64            `json:"fixedDelayMilliseconds"`
	Body         string            `json:"body"`
	JsonBody     string            `json:"jsonBody"`
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

	handlers(c)

	_ = http.ListenAndServe(":8090", nil)
}

func handlers(c []Contract) {
	for i := 0; i < len(c); i++ {
		r := c[i]

		http.HandleFunc(r.Path, func(w http.ResponseWriter, req *http.Request) {

			code := r.SuccessStatusCode
			body := r.SuccessBody
			if err(req) {
				code = r.ErrorStatusCode
				body = r.ErrorBody
			}

			w.Header().Add("Content-Type", r.ContentType)
			w.WriteHeader(int(code))
			_, _ = w.Write([]byte(body))
		})

	}
}

func err(req *http.Request) bool {
	key, ok := req.URL.Query()["success"]

	if !ok || len(key[0]) < 1 {
		log.Println("Url Param 'success' is missing")
		return false
	}

	return key[0] == "false"
}

func contracts(file []byte) []Contract {
	var c []Contract

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
