package influx

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"time"
)

type Result []client.Result

func (this Result) GMap1() []map[string]any {
	list := []map[string]any{}
	if len(this) > 0 {
		for _, k := range this[0].Series {
			for _, x := range k.Values {
				m := map[string]any{}
				for i, y := range k.Columns {
					if len(x) > i {
						m[y] = x[i]
					}
				}
				list = append(list, m)
			}
		}
	}
	return list
}

func (this Result) GMaps() [][]map[string]any {
	lists := [][]map[string]any{}
	for _, v := range this {
		list := []map[string]any{}
		for _, k := range v.Series {
			for _, x := range k.Values {
				m := map[string]any{}
				for i, y := range k.Columns {
					if len(x) > i {
						m[y] = x[i]
					}
				}
				list = append(list, m)
			}
		}
		lists = append(lists, list)
	}
	return lists
}

type option struct {
	// Database 数据库
	Database string

	// Precision 精度s
	Precision string
}

type Client struct {
	client    client.Client
	option    *option
	err       error
	newClient func() *Client
}

func (this *Client) Reconnect() {
	if this.newClient != nil {
		this.client.Close()
		*this = *this.newClient()
	}
}

func (this *Client) Ping() {
	_, _, err := this.client.Ping(time.Second)
	if err != nil {
		log.Println("[错误]", err.Error())
	} else {
		log.Println("InfluxDB连接成功...")
	}
}

func (this *Client) Err() error {
	return this.err
}

func (this *Client) Close() error {
	this.option = nil
	if this.client != nil {
		this.client.Close()
	}
	return nil
}

// Write 把数据写入influxdb
// @tableName 表名
// @tags 索引
// @fields 字段
// @t 时间
// 时间一直会覆盖数据
func (this *Client) Write(tableName string, tags map[string]string, fields map[string]any, t ...time.Time) error {
	if this.err != nil {
		this.Reconnect()
		return this.err
	}
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  this.option.Database,  //数据库名
		Precision: this.option.Precision, //时间精度秒
	})
	if err != nil {
		return err
	}
	T := time.Now()
	if len(t) > 0 {
		T = t[0]
	}
	//将创建的表,以及内容字段放入pt
	pt, err := client.NewPoint(tableName, tags, fields, T)
	if err != nil {
		return err
	}
	//把表放入创建的point中
	bp.AddPoint(pt)
	//写入数据
	err = this.client.Write(bp)
	if err != nil {
		this.Reconnect()
	}
	return err
}

// Exec 执行
func (this *Client) Exec(sql string) (Result, error) {
	if this.err != nil {
		this.Reconnect()
		return nil, this.err
	}
	q := client.NewQuery(sql, this.option.Database, this.option.Precision)
	response, err := this.client.Query(q)
	if err != nil {
		this.Reconnect()
		return nil, err
	}
	if err := response.Error(); err != nil {
		return nil, err
	}
	return response.Results, nil
}
