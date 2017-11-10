package main

import (
	"github.com/lestrrat/go-libxml2"
	//"github.com/lestrrat/go-libxml2/parser"
	"encoding/json"
	"github.com/lestrrat/go-libxml2/types"
	"github.com/lestrrat/go-libxml2/xpath"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"github.com/golang/glog"
)

func InitDoc(rc io.Reader) types.Document {
	doc, err := libxml2.ParseHTMLReader(rc)
	if err != nil {
		panic(err.Error())
	}
	return doc
}

func Xpath(p string, doc types.Node) interface{} {
	nodes := xpath.NodeList(doc.Find(p))
	ret := make([]string, 0)
	for _, node := range nodes {
		base := strings.TrimSpace(node.TextContent())
		ret = append(ret, base)
	}
	return strings.Join(ret, "\n")
}

func XpathList(p string, doc types.Node) []interface{} {
	nodes := xpath.NodeList(doc.Find(p))
	ret := make([]interface{}, 0)
	for _, node := range nodes {
		base := strings.TrimSpace(node.TextContent())
		ret = append(ret, base)
	}
	return ret
}
func XpathNodeList(p string, doc types.Node) []types.Node {
	nodes := xpath.NodeList(doc.Find(p))
	ret := make([]types.Node, 0)
	for _, node := range nodes {
		ret = append(ret, node)
	}
	return ret
}

func GetPage(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode != 200 {
		glog.Infof("url: %s ; return a bad code : %s", url, resp.StatusCode)
		return nil
	}
	return resp
}

func GetJson(url string) map[string]interface{} {
	resp, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		glog.Infof("url: %s ; return a bad code : %s\n", url, resp.StatusCode)
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Info(err.Error())
		return nil
	}
	var json_ret map[string]interface{}
	json.Unmarshal([]byte(string(body)), &json_ret)
	return json_ret
}

func GetAbsoluteURL(main_url string, ref_url string) string {
	m, err := url.Parse(main_url)
	if err != nil {
		glog.Info(err.Error())
		return ""
	}
	r, err := url.Parse(ref_url)
	if err != nil {
		glog.Info(err.Error())
		return ""
	}
	return m.ResolveReference(r).String()
}

func quote(str string) string {
	return url.QueryEscape(str)
}
