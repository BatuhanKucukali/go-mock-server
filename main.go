package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Contract struct {
	Maps []Mapping `json:"mappings"`
	Auth *Auth     `json:"auth"`
}

type Mapping struct {
	Req  *Request     `json:"request"`
	Resp [] *Response `json:"responses"`
}

type Request struct {
	Method string `json:"method"`
	Url    string `json:"url"`
}

type Response struct {
	Condition    *Condition        `json:"condition"`
	Status       uint16            `json:"status"`
	FixedDelayMs uint64            `json:"fixedDelayMilliseconds"`
	Body         string            `json:"body"`
	JsonBody     map[string]string `json:"jsonBody"` // TODO json parse
	Headers      map[string]string `json:"headers"`
}

type Condition struct {
	Body map[string]interface{} `json:"body"`
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

			var response *Response

			var request map[string]interface{}
			decoder := json.NewDecoder(req.Body)

			_ = decoder.Decode(&request)

			for i := 0; i < len(m.Resp); i++ {
				r := m.Resp[i]

				if response == nil && r.Condition == nil {
					response = r
					continue
				}

				if r.Condition.Body != nil {

					allParametersSatisfied := true
					for k, v := range r.Condition.Body {
						parameter := request[k]

						if parameter != v {
							allParametersSatisfied = false
							break
						}
					}

					if allParametersSatisfied {
						response = r
					}
				}

			}

			sleep(response.FixedDelayMs)

			for k, v := range response.Headers {
				w.Header().Add(k, v)
			}
			w.WriteHeader(int(response.Status))
			var responseBytes []byte

			if response.JsonBody != nil {
				responseBytes, _ = json.Marshal(response.JsonBody)
			} else if response.Body != "" {
				responseBytes = []byte(response.Body)
			}

			_, _ = w.Write(responseBytes)

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
