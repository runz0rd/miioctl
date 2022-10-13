package miio

import (
	"encoding/json"
	"fmt"
	"sync"

	miiogo "github.com/ofen/miio-go"
	"github.com/ofen/miio-go/proto"
)

type Client struct {
	*miiogo.Client
	requestID int
}

func New(addr, token string) *Client {
	conn, err := proto.Dial(fmt.Sprintf("%v:54321", addr), nil)
	if err != nil {
		panic(err)
	}
	c := &miiogo.Client{Mutex: sync.Mutex{}, Conn: conn}
	c.SetToken(token)
	return &Client{c, 1}
}

// Send sends request to device.
func (c *Client) Send(method string, params interface{}) ([]byte, error) {
	req := struct {
		RequestID int         `json:"id"`
		Method    string      `json:"method"`
		Params    interface{} `json:"params"`
	}{
		RequestID: c.requestID,
		Method:    method,
		Params:    params,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	if _, err := c.Write(payload); err != nil {
		return nil, err
	}

	resp := make([]byte, 4096)
	n, err := c.Read(resp)
	if err != nil {
		return nil, err
	}

	if err == nil {
		c.Lock()
		c.requestID++
		c.Unlock()
	}

	return resp[:n], nil

	// // trim non-printable characters
	// return bytes.TrimFunc(resp[:n], func(r rune) bool {
	// 	return !unicode.IsGraphic(r)
	// }), err
}

type Device interface {
	DeviceId() string
	Query() error
	ToString(field string) string
}

type Result struct {
	Did   string      `json:"did,omitempty"`
	Siid  int         `json:"siid,omitempty"`
	Piid  int         `json:"piid,omitempty"`
	Code  int         `json:"code,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type Response struct {
	Id            int      `json:"id,omitempty"`
	Results       []Result `json:"result,omitempty"`
	ExecutionTime int      `json:"exe_time,omitempty"`
}

func (c Client) GetProperties(payload interface{}) (Response, error) {
	var response Response
	resp, err := c.Send("get_properties", payload)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return response, err
	}
	return response, err
}

func (c Client) SetProperties(payload interface{}) error {
	_, err := c.Send("set_properties", payload)
	return err
}
