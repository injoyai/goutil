package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

func newResponse(r *http.Response, startTime time.Time, tryNum int) *Response {
	if r == nil {
		return nil
	}
	return &Response{
		Response: r,
		spend:    time.Since(startTime),
		tryNum:   tryNum,
	}
}

type Response struct {
	*http.Response
	spend  time.Duration //花费时间
	tryNum int           //
}

// Spend 获取花费时间
func (this *Response) Spend() time.Duration {
	return this.spend
}

func (this *Response) TryNum() int {
	return this.tryNum
}

func (this *Response) String() string {
	bs, _ := httputil.DumpResponse(this.Response, true)
	return string(bs)
}

// ReadBody 一次性读取全部字节(适用于小数据)
func (this *Response) ReadBody() ([]byte, error) {
	if this.Response == nil {
		return nil, nil
	}
	defer this.Response.Body.Close()
	return io.ReadAll(this.Response.Body)
}

func (this *Response) Bind(ptr any) error {
	switch v := ptr.(type) {
	case nil:
		return nil
	case io.Writer:
		defer this.Response.Body.Close()
		_, err := io.Copy(v, this.Body)
		return err
	default:
		body, err := this.ReadBody()
		if err != nil {
			return err
		}
		switch val := ptr.(type) {
		case *string:
			*val = string(body)
		case *[]byte:
			//val不为nil,this.body为nil可以赋值成功
			*val = body
		case io.Writer:
			_, err := io.Copy(val, this.Body)
			return err
		default:
			//尝试解析,错误不处理,不然会返回错误,看不到正常请求的数据
			return json.Unmarshal(body, ptr)
		}
	}
	return nil
}
