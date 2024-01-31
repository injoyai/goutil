package http

import (
	"github.com/PuerkitoBio/goquery"
)

/*
Document 解析body
使用方法参考
https://blog.csdn.net/qq_38334677/article/details/129225231

	doc, err := r.Document()
	if err!=nil{
		return
	}
	//查找标签: 		doc.Find("body,div,...") 多个用,隔开
	//查找ID: 		doc.Find("#id1")
	//查找class: 	doc.Find(".class1")
	//查找属性: 		doc.Find("div[lang]") doc.Find("div[lang=zh]") doc.Find("div[id][lang=zh]")
	//查找子节点: 	doc.Find("body>div")
	//过滤数据: 		doc.Find("div:contains(xxx)")
	//过滤节点: 		dom.Find("span:has(div)")
	doc.Find("body").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})

选择器					说明
Find(“div[lang]”)		筛选含有lang属性的div元素
Find(“div[lang=zh]”)	筛选lang属性为zh的div元素
Find(“div[lang!=zh]”)	筛选lang属性不等于zh的div元素
Find(“div[lang¦=zh]”)	筛选lang属性为zh或者zh-开头的div元素
Find(“div[lang*=zh]”)	筛选lang属性包含zh这个字符串的div元素
Find(“div[lang~=zh]”)	筛选lang属性包含zh这个单词的div元素，单词以空格分开的
Find(“div[lang$=zh]”)	筛选lang属性以zh结尾的div元素，区分大小写
Find(“div[lang^=zh]”)	筛选lang属性以zh开头的div元素，区分大小写
*/
func (this *Response) Document() (*goquery.Document, error) {
	if this.Error != nil {
		return nil, this.Error
	}
	if this.doc != nil {
		return this.doc, nil
	}
	doc, err := goquery.NewDocumentFromReader(this.Response.Body)
	if err != nil {
		return nil, err
	}
	this.doc = doc
	return doc, nil
}
