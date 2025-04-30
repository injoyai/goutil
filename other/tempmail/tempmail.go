package tempmail

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"time"
)

// New 例:username@mailto.plus,填写@前面的
func New(username string) *Client {
	return &Client{Username: username}
}

type Client struct {
	Username string
}

// List 最新邮件列表
func (this *Client) List(limits ...int) ([]*Mail, error) {
	limit := conv.Default[int](0, limits...)
	u := fmt.Sprintf("https://tempmail.plus/api/mails?email=%s%%40mailto.plus&first_id=0&epin=&limit=%d", this.Username, limit)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ll := &struct {
		Mails []*Mail `json:"mail_list"`
	}{}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bs, ll); err != nil {
		return nil, err
	}
	for _, v := range ll.Mails {
		if err := v.resp(this); err != nil {
			return nil, err
		}
	}

	return ll.Mails, nil
}

// Listen 监听邮件,number循环次数,wait等待时间,f筛选邮件
func (this *Client) Listen(number int, wait time.Duration, f func(m *Mail) bool) {
	for num := 0; num < number; num++ {
		<-time.After(wait)
		ls, err := this.List(10)
		if err != nil {
			continue
		}
		for i := len(ls) - 1; i >= 0; i-- {
			if !f(ls[i]) {
				return
			}
		}
	}
}

func (this *Client) Details(id int64) (*Details, error) {
	u := fmt.Sprintf("https://tempmail.plus/api/mails/%d?email=%s%%40mailto.plus&epin=", id, this.Username)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := new(Details)
	if err := json.Unmarshal(bs, data); err != nil {
		return nil, err
	}
	if !data.Result {
		return nil, errors.New("失败")
	}
	return data, data.resp()
}

func (this *Client) Delete(id int64) error {
	u := fmt.Sprintf("https://tempmail.plus/api/mails/%d", id)
	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	m := conv.NewMap(resp.Body)
	if !m.GetBool("result") {
		return errors.New("失败")
	}
	return nil
}

/*
	{
	    "attachment_count": 0,
	    "first_attachment_name": "",
	    "from_mail": "SRS0=q9Uv=kv=qq.com=1113655791@injoy.ink",
	    "from_name": "",
	    "is_new": true,
	    "mail_id": 2969893429,
	    "subject": "titless",
	    "time": "2025-02-19 10:31:25"
	}
*/
type Mail struct {
	ID      int64  `json:"mail_id"`
	Unread  bool   `json:"is_new"`
	Subject string `json:"subject"`
	TimeStr string `json:"time"`
	Time    time.Time
	read    func() (*Details, error)
}

func (this *Mail) resp(c *Client) error {
	var err error
	this.Time, err = time.ParseInLocation(time.DateTime, this.TimeStr, time.Local)
	if err != nil {
		return err
	}
	this.read = func() (*Details, error) {
		return c.Details(this.ID)
	}
	return nil
}

func (this *Mail) Read() (*Details, error) {
	return this.Details()
}

func (this *Mail) Details() (*Details, error) {
	if this.read == nil {
		return nil, errors.New("不可用")
	}
	return this.read()
}

/*
	{
	    "attachments": [],
	    "date": "Wed, 19 Feb 2025 16:25:17 +0800",
	    "from": "1113655791@qq.com",
	    "from_is_local": false,
	    "from_mail": "1113655791@qq.com",
	    "from_name": "",
	    "html": "你好呀\n\n\n",
	    "is_tls": true,
	    "mail_id": 2969880549,
	    "message_id": "<tencent_DC5E1BC65D06D40F14942B65AC3DD839BF08@qq.com>",
	    "result": true,
	    "subject": "titless",
	    "text": "你好呀",
	    "to": "test@injoy.ink"
	}
*/
type Details struct {
	ID          int64     `json:"mail_id"`
	Attachments []string  `json:"attachments"`
	TimeStr     string    `json:"date"`
	Time        time.Time `json:"time"`
	From        string    `json:"from"`
	FromIsLocal bool      `json:"from_is_local"`
	FromMail    string    `json:"from_mail"`
	FromName    string    `json:"from_name"`
	Html        string    `json:"html"`
	Tls         bool      `json:"is_tls"`
	MessageID   string    `json:"message_id"`
	Result      bool      `json:"result"`
	Subject     string    `json:"subject"`
	Text        string    `json:"text"`
	To          string    `json:"to"`
}

func (this *Details) resp() (err error) {
	this.Time, err = time.Parse(time.RFC1123Z, this.TimeStr)
	return
}
