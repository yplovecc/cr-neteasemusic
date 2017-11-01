package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/crawlerclub/dl"
	"github.com/golang/glog"
	"strings"
	"sync"
	"time"
)

type Crawler struct {
	Url      string
	Interval int
	Docs     *[]Artist
}

type Artist struct {
	Url      string `json:"url"`
	Alias    string `json:"alias"`
	Gender   int    `json:"gender"`
	Status   int    `json:"status"`
	La       int    `json:"la"`
	Lo       string `json:"lo"`
	Alphabet string `json:"alphabet"`
	Origin   string `json:"origin"`
	Alias1   string `json:"alias1"`
	Alias2   string `json:"alias2"`
	Alias3   string `json:"alias3"`
	Alias4   string `json:"alias4"`
	Alias5   string `json:"alias5"`
}

func NewCrawler(interval int) (crawler *Crawler) {
	crawler = &Crawler{Interval: interval}
	return
}

func (self *Crawler) Run(wg *sync.WaitGroup, exitCh chan int) {
	defer wg.Done()
	glog.Info("start crawler crontab")
	defer glog.Info("crawler crontab exit")
	for {
		select {
		case <-exitCh:
			return
		default:
			glog.Info("execute crontab")
			self.Process()
			select {
			case <-exitCh:
				return
			case <-time.After(time.Duration(self.Interval) * time.Second):
				break
			}
		}
	}
}

func (self *Crawler) Process() {
	glog.Info("process start")
	for _, crawleinfo := range Crawlinfos {
		for ialph, salph := range Alphabets {
			listurl := fmt.Sprintf("%sid=%d&initial=%d", *listUrl, crawleinfo.catid, ialph)
			glog.Infof("crawl list : %s", listurl)
			req := &dl.HttpRequest{Url: listurl, Method: "GET", UseProxy: false, Platform: "pc"}
			res := dl.Download(req)
			if res.Error != nil {
				glog.Error(res.Error)
				continue
			}
			doc := InitDoc(bytes.NewReader(res.Content))
			elements := XpathNodeList("//a[@class='nm nm-icn f-thide s-fc0']", doc)
			for _, e := range elements {
				var artist Artist
				artist.Url = *detailUrl + Xpath("@href", e).(string)
				artist.Gender = Gender[crawleinfo.gen]
				artist.La = La[crawleinfo.lo]
				artist.Lo = crawleinfo.lo
				artist.Alphabet = salph
				artist.Alias = e.TextContent()
				artist.Status = 4
				artist.Origin = "wangyi"
				if self.IsArtistExist(artist.Url, artist.La) {
					glog.Infof("%s is exist", artist.Alias)
					continue
				}
				req := &dl.HttpRequest{Url: artist.Url,
					Method: "GET", UseProxy: false, Platform: "pc"}
				res := dl.Download(req)
				if res.Error != nil {
					glog.Error(res.Error)
					continue
				}
				doc1 := InitDoc(bytes.NewReader(res.Content))
				es := XpathNodeList("//h3[@id='artist-alias']", doc1)
				if len(es) > 0 {
					title := Xpath("@title", es[0]).(string)
					if titlesplit := strings.Split(title, " - "); len(titlesplit) == 2 {
						artist.Alias = artist.Alias + ";" + titlesplit[1]
					}
				}
				aliases := strings.Split(artist.Alias, ";")
				for i, v := range aliases {
					switch i {
					case 0:
						artist.Alias1 = v
					case 1:
						artist.Alias2 = v
					case 2:
						artist.Alias3 = v
					case 3:
						artist.Alias4 = v
					case 4:
						artist.Alias5 = v
					}
				}
				doc1.Free()

				var mdata map[string]interface{}
				jsondata, _ := json.Marshal(artist)
				json.Unmarshal(jsondata, &mdata)
				writeDB(*table, mdata)
				glog.Infof("insert %s", artist)
			}
			doc.Free()
		}
	}
	glog.Info("process end")
}

func (self *Crawler) IsArtistExist(url string, la int) bool {
	sql := fmt.Sprintf("select * from %s where url=\"%s\" and la=%d", *table, url, la)
	if len(selectDB(sql)) > 0 {
		return true
	}
	return false
}
