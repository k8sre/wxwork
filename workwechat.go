package wxwork

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//New 实例化微信企业号应用
func New(cropID string, agentID int64, AgentSecret string) *Client {

	c := new(Client)
	c.CropID = cropID
	c.AgentID = agentID
	c.AgentSecret = AgentSecret
	return c
}

//Send 发送信息
func (c *Client) Send(msg Message) error {

	if err := c.GetAccessToken(); err != nil {
		return err
	}

	msg.AgentID = c.AgentID

	url := "https://qyapi.weixin.qq.com/cgi-bin/linkedcorp/message/send?access_token=" + c.Token.AccessToken

	resultByte, err := JSONPost(url, msg)

	if err != nil {
		err = errors.New("请求微信接口失败: " + err.Error())
		return err
	}

	result := Result{}
	fmt.Println(string(resultByte))

	err = json.Unmarshal(resultByte, &result)
	if err != nil {
		err = errors.New("解析微信接口返回数据失败: " + err.Error())
		return err
	}

	if result.ErrCode != 0 {
		err = errors.New("发送消息失败: " + result.ErrMsg)
	}

	if len(result.InvalidUser) != 0 || len(result.InvalidTag) != 0 || len(result.InvalidParty) != 0 {
		err = fmt.Errorf("消息发送成功, 但是有部分目标无法送达: InvalidUser(%s)InvalidParty(%s)InvalidTag(%s)", result.InvalidUser, result.InvalidParty, result.InvalidTag)
	}

	return err
}

//GetAccessToken 获取回话token
func (c *Client) GetAccessToken() error {
	var err error
	if c.Token.AccessToken == "" || c.Token.ExpiresInTime.Before(time.Now()) {
		c.Token, err = getAccessTokenFromWeixin(c.CropID, c.AgentSecret)
		if err != nil {
			return err
		}
		c.Token.ExpiresInTime = time.Now().Add(time.Duration(c.Token.ExpiresIn-1000) * time.Second)
	}
	return nil
}

//从微信服务器获取token
func getAccessTokenFromWeixin(cropID, secret string) (TokenSession AccessToken, err error) {
	WxAccessTokenURL := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + cropID + "&corpsecret=" + secret

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	result, err := client.Get(WxAccessTokenURL)
	if err != nil {
		return TokenSession, err
	}

	res, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return TokenSession, err
	}

	defer result.Body.Close()

	err = json.Unmarshal(res, &TokenSession)
	if err != nil {
		return TokenSession, err
	}

	if TokenSession.ExpiresIn == 0 || TokenSession.AccessToken == "" {
		err = fmt.Errorf("获取微信错误代码: %v, 错误信息: %v", TokenSession.ErrCode, TokenSession.ErrMsg)
		return TokenSession, err
	}

	return TokenSession, err
}

//JSONPost Post请求json数据
func JSONPost(url string, data interface{}) ([]byte, error) {
	jsonBody, err := encodeJSON(data)

	if err != nil {
		return nil, err
	}
	r, err := http.Post(url, "application/json;charset=utf-8", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}
	return body, err
}

func encodeJSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
