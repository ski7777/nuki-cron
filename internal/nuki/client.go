package nuki

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type NukiClient struct {
	apikey string
	client *http.Client
}

func (n *NukiClient) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+n.apikey)
	return http.DefaultTransport.RoundTrip(req)
}

func (n *NukiClient) do(req *http.Request) (res []byte, err error) {
	hres, err := n.client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = hres.Body.Close()
	}()
	if hres.StatusCode < 200 || hres.StatusCode >= 300 {
		err = errors.New(fmt.Sprintf("unexpected status code %d", hres.StatusCode))
		return
	}
	res, err = io.ReadAll(hres.Body)
	return
}

type GetAuthReponseAuth struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetAuthsResponse []GetAuthReponseAuth

func (n *NukiClient) GetAuths() (res GetAuthsResponse, err error) {
	req, err := http.NewRequest("GET", "https://api.nuki.io/smartlock/auth", nil)
	if err != nil {
		return
	}
	rres, err := n.do(req)
	if err != nil {
		return
	}
	res = GetAuthsResponse{}
	err = json.Unmarshal(rres, &res)
	return
}

type SmartLockAuthMultiUpdate struct {
	GetAuthReponseAuth
	AllowedFromDate  string `json:"allowedFromDate"`
	AllowedUntilDate string `json:"allowedUntilDate"`
	AllowedFromTime  int    `json:"allowedFromTime"`
	AllowedUntilTime int    `json:"allowedUntilTime"`
	AllowedWeekDays  int    `json:"allowedWeekDays"`
	Enabled          bool   `json:"enabled"`
}

type SmartLockAuthMultiUpdateRequest []SmartLockAuthMultiUpdate

func (n *NukiClient) UpdateAuths(reqBody SmartLockAuthMultiUpdateRequest) (err error) {
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", "https://api.nuki.io/smartlock/auth", bytes.NewReader(reqBodyBytes))
	if err != nil {
		return
	}
	_, err = n.do(req)
	return
}

func NewNukiClient(apikey string) (client *NukiClient) {
	client = &NukiClient{
		apikey: apikey,
	}
	client.client = &http.Client{Transport: client}
	return
}
