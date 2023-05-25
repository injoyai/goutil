package influx

import (
	"log"
	"testing"
)

func TestNewHTTPClient(t *testing.T) {
	c := NewHTTPClient(&HTTPOption{})
	if err := c.Err(); err != nil {
		t.Error(err)
		return
	}
	c.Ping()

	{
		result, err := c.Exec(`select * from "system_data_log" order by time desc limit 10`)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(result)
	}
	result, err := c.Exec("select * from test")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range result.GMap1() {
		t.Log(v)
	}
}

// TestHTTPWrite 测试写入速度 (32.15s|29.33s)/1000条
func TestHTTPWrite(t *testing.T) {
	c := NewHTTPClient(&HTTPOption{})
	if err := c.Err(); err != nil {
		t.Error(err)
		return
	}
	for i := 100000; i < 10000000; i++ {
		if err := c.Write("test", map[string]string{
			"key":  "123456",
			"key2": "key2",
		}, map[string]interface{}{
			"f1":  i,
			"f2":  i,
			"f3":  i,
			"f4":  i,
			"f5":  i,
			"f6":  i,
			"f7":  i,
			"f8":  i,
			"f9":  i,
			"f10": i,
		}); err != nil {
			t.Error(err)
			return
		}
	}
}

func TestHTTPWrite2(t *testing.T) {
	c := NewHTTPClient(&HTTPOption{})
	if err := c.Err(); err != nil {
		t.Error(err)
		return
	}

	//if err := c.Write("push_fail_log", nil, map[string]interface{}{
	//	"pushKey":  "",
	//	"pushMark": "PUSH_MQTT",
	//	"msg":      "钱测试",
	//}); err != nil {
	//	t.Error(err)
	//	return
	//}

	result, err := c.Exec(`select * from "push_fail_log"`)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(result.GMap1())
}

func TestHTTPRead(t *testing.T) {
	c := NewHTTPClient(&HTTPOption{})
	if err := c.Err(); err != nil {
		t.Error(err)
		return
	}
	result, err := c.Exec(`select count(*) from "test"`)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(result.GMap1())
}
