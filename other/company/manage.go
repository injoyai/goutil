/*****************************************************************************
*名称:	公司管理工具包
*功能:	通过内存快速处理公司结构
*作者:	钱纯净
******************************************************************************/

package company

import (
	"log"
	"sort"
	"sync"
	"time"

	"xorm.io/xorm"
)

type Interface interface {
	GetList() (map[int64]*ICompany, error) //获取最新的公司数据
	//刷新触发事件
}

type ICompany struct {
	ID       int64
	ParentID int64
	Name     string
	Tree     []*ICompany
}

func NewManage(data Interface) IManage {
	//log.Println("开始加载单位")
	l := &CompanyManage{data: data}
	go func() {
		for {
			l.F5()
			//log.Println("加载单位成功")
			time.Sleep(time.Hour * 24)
		}
	}()
	return l
}

//////////////////////////////////////////////////////

type IManage interface {
	F5()                                                //刷新缓存数据
	Look()                                              //查看数据,测试用
	IsRoot(int64) bool                                  //是否是顶级公司
	Exist(int64) bool                                   //是否存在
	GetFather(int64) int64                              //获取上级部门
	GetFathers(int64) []int64                           //获取所有上级部门(不包括本级)
	GetChildrens(int64) []int64                         //获取所有下级部门(不包括本级)
	GetParentAndChildrens(int64) []int64                //获取本级和下级部门
	DbFilter(string, int64, *xorm.Engine) *xorm.Session //添加搜索条件
	GetCompanyTree(int64) *ICompany                     //获取公司树
}

// CompanyManager 管理模块
type CompanyManage struct {
	companyMux sync.Mutex
	companys   map[int64]*ICompany
	data       Interface
}

// Refresh 刷新
func (c *CompanyManage) F5() {
	c.companyMux.Lock()
	defer c.companyMux.Unlock()
	list, err := c.data.GetList()
	if err != nil {
		log.Println("刷新公司错误:", err)
	} else {
		c.companys = list
	}
}

//查看,测试用
func (c *CompanyManage) Look() {
	log.Println("长度:", len(c.companys))
	for _, v := range c.companys {
		log.Println(v)
	}
}

// IsRoot 是否根部门
func (c *CompanyManage) IsRoot(companyID int64) bool {
	c.companyMux.Lock()
	defer c.companyMux.Unlock()
	v, ok := c.companys[companyID]
	return ok && v.ParentID == 0
}

// Exist 是否存在
func (c *CompanyManage) Exist(companyID int64) bool {
	c.companyMux.Lock()
	defer c.companyMux.Unlock()
	_, ok := c.companys[companyID]
	return ok
}

// GetFather 获取上级
func (c *CompanyManage) GetFather(companyId int64) int64 {
	c.companyMux.Lock()
	defer c.companyMux.Unlock()
	if _, ok := c.companys[companyId]; ok {
		return c.companys[companyId].ParentID
	}
	return 0
}

// getFahters 查看上级部门的id
func (c *CompanyManage) getFathers(companyID int64) []int64 {
	var result []int64
	for _, v := range c.companys {
		if v.ID == companyID {
			if v.ParentID != 0 {
				fathers := c.getFathers(v.ParentID)
				result = append(result, fathers...)
			}
			result = append(result, v.ID)
		}
	}
	return result
}

// GetFathers 获取本级上级所有部门编号
func (c *CompanyManage) GetFathers(companyID int64) []int64 {
	if c.Exist(companyID) {
		c.companyMux.Lock()
		defer c.companyMux.Unlock()
		return c.getFathers(companyID)
	}
	return []int64{}
}

// getChildren 查看子部门的id
func (c *CompanyManage) getChildrens(parentID int64) []int64 {
	var result []int64
	for _, v := range c.companys {
		if v.ParentID == parentID {
			childIds := c.getChildrens(v.ID)
			mid := append(childIds, v.ID)
			result = append(result, mid...)
		}
	}
	return result
}

// GetParentAndChildren 获取本级和下级的所有部门编号
func (c *CompanyManage) GetChildrens(parentID int64) []int64 {
	if c.Exist(parentID) {
		c.companyMux.Lock()
		defer c.companyMux.Unlock()
		return c.getChildrens(parentID)
	}
	return []int64{}
}

// GetParentAndChildren 获取本级和下级的所有部门编号
func (c *CompanyManage) GetParentAndChildrens(parentID int64) []int64 {
	if c.Exist(parentID) {
		c.companyMux.Lock()
		defer c.companyMux.Unlock()
		childrenIds := c.getChildrens(parentID)
		result := append(childrenIds, parentID)
		return result
	}
	return []int64{}
}

//数据取反
func (c *CompanyManage) reverse(co []int64) []int64 {
	c.companyMux.Lock()
	defer c.companyMux.Unlock()
	data := []int64{}
	for i, _ := range c.companys {
		data = append(data, i)
		for _, v := range co {
			if i == v && len(data) != 0 {
				data = data[:len(data)-1]
				continue
			}
		}
	}
	return data
}

// DbFilter2 构建session查询条件,优化下,当数量上去用notIn
func (c *CompanyManage) DbFilter(colName string, companyID int64, engine *xorm.Engine) *xorm.Session {
	co := c.GetParentAndChildrens(companyID)
	//判断长度
	n := len(c.companys)
	if len(co) > n/2+1 {
		return engine.NotIn(colName, c.reverse(co))
	}
	return engine.In(colName, co)
}

func (c *CompanyManage) GetCompanyTree(companyId int64) *ICompany {
	data := new(ICompany)
	if val, ok := c.companys[companyId]; ok {
		data = val
		l := append(c.getChildrensInfo(companyId), data)
		m := make(map[int64]*ICompany)
		for _, v := range l {
			v.Tree = nil
			m[v.ID] = v
		}
		for i, v := range m {
			if _, ok := m[v.ParentID]; ok {
				if m[v.ParentID].Tree == nil {
					m[v.ParentID].Tree = []*ICompany{m[i]}
				} else {
					m[v.ParentID].Tree = append(m[v.ParentID].Tree, m[i])
				}
			}
		}
		for i, _ := range m {
			if len(m[i].Tree) > 1 {
				sort.Sort(Sort(m[i].Tree))
			}
		}
	}
	return data
}

// getChildren 查看子部门的信息
func (c *CompanyManage) getChildrensInfo(parentID int64) []*ICompany {
	var result []*ICompany
	for _, v := range c.companys {
		if v.ParentID == parentID {
			result = append(result, append(c.getChildrensInfo(v.ID), v)...)
		}
	}
	return result
}

//=======================================

type Sort []*ICompany

func (this Sort) Len() int {
	return len(this)
}

func (this Sort) Less(i, j int) bool {
	return this[i].ID < this[j].ID
}

func (this Sort) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
