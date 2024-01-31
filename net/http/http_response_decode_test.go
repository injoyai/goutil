package http

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"testing"
)

func TestResponse_Document(t *testing.T) {
	t.Log(Url("https://github.com/PuerkitoBio/goquery").Do().
		Execute(func(r *Response) error {
			doc, err := r.Document()
			if err != nil {
				return err
			}
			doc.Find("head>link").Each(func(i int, selection *goquery.Selection) {
				fmt.Println(selection.Attr("href"))
				fmt.Println(selection.Text())
			})

			return nil
		}),
	)
}
