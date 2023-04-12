package BOKA

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/tidwall/gjson"
)

type HTTP struct {
	Url     string
	Params  map[string]string
	Headers map[string]string
	Body    map[string]interface{}
	RAW     []byte
	JSON    gjson.Result
}

func (R *HTTP) STRUCT(structName any) {
	if err := json.Unmarshal(R.RAW, structName); err != nil {
		log.Println(err)
	}
}

func (R *HTTP) get() error {
	err := R.send("GET")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (R *HTTP) post() error {
	err := R.send("POST")
	if err != nil {
		return err
	}
	return nil
}

func (R *HTTP) send(method string) error {
	//add post body
	var bodyJson []byte
	if len(R.Body) != 0 {
		var err error
		bodyJson, err = json.Marshal(R.Body)
		if err != nil {
			log.Println(err)
		}
	} else {
		bodyJson = nil
	}

	req, err := http.NewRequest(method, R.Url, bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	//add params
	q := req.URL.Query()
	for key, val := range R.Params {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()
	//add headers
	for key, val := range R.Headers {
		req.Header.Add(key, val)
	}
	//http client
	client := &http.Client{}
	log.Printf("Go POST URL : %s \n", req.URL.String())

	//发送请求
	res, err := client.Do(req)

	if err != nil {
		log.Println("request error")
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body) //一定要关闭res.Body
	//读取body
	resBody, err := io.ReadAll(res.Body) //把  body 内容读入字符串
	if err != nil {
		log.Println(err)
		return err
	}
	R.RAW = resBody
	R.JSON = gjson.Parse(string(R.RAW))

	return nil

}
