package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Contract struct {
	Maps        []Mapping    `json:"mappings"`
	Credentials *Credentials `json:"credentials"`
}

type Mapping struct {
	Req  *Request    `json:"request"`
	Resp []*Response `json:"responses"`
}

type Request struct {
	Method   string    `json:"method"`
	Url      string    `json:"url"`
	AuthType *AuthType `json:"authType"`
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

type Credentials struct {
	BasicAuth   *BasicAuth `json:"basicAuth"`
	BearerToken string     `json:"bearerToken"`
}

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthType string

const (
	BasicType     AuthType = "basic"
	BasicPrefix   string   = "Basic"
	BearerType    AuthType = "bearer"
	BearerPrefix  string   = "Bearer"
	Authorization string   = "Authorization"
)

func main() {
	fmt.Println("Go Mock Server Started...")

	file := read("contract.json")

	c := contracts(file)

	initRouters(c)

	_ = http.ListenAndServe(":8090", nil)
}

func initRouters(c Contract) {
	maps := c.Maps
	credentials := c.Credentials

	for i := 0; i < len(maps); i++ {
		m := maps[i]

		// TODO add HTTP method
		http.HandleFunc(m.Req.Url, func(w http.ResponseWriter, req *http.Request) {

			if m.Req.AuthType != nil {
				authorizationHeaders := req.Header[Authorization]
				if authorizationHeaders == nil || !isAuthorized(m.Req.AuthType, authorizationHeaders[0], credentials) {
					w.WriteHeader(http.StatusUnauthorized)
					responseBytes := []byte(http.StatusText(http.StatusUnauthorized))
					_, _ = w.Write(responseBytes)
					return
				}
			}

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

func isAuthorized(authType *AuthType, authorizationHeader string, credentials *Credentials) bool {
	if *authType == BasicType {
		if !strings.HasPrefix(authorizationHeader, BasicPrefix) {
			return false
		}

		base64String := strings.SplitAfter(authorizationHeader, BasicPrefix+" ")
		decodedBase64String, _ := b64.StdEncoding.DecodeString(base64String[1])
		sendUserNamePassword := string(decodedBase64String)
		authorizedUserNamePassword := credentials.BasicAuth.Username + ":" + credentials.BasicAuth.Password
		if strings.Compare(sendUserNamePassword, authorizedUserNamePassword) != 0 {
			return false
		}
	} else if *authType == BearerType {
		if !strings.HasPrefix(authorizationHeader, BearerPrefix) {
			return false
		}

		token := strings.SplitAfter(authorizationHeader, BearerPrefix+" ")
		if strings.Compare(token[1], credentials.BearerToken) != 0 {
			return false
		}
	} else {
		return false
	}

	return true
}
