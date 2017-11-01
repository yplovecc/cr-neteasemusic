package main

import (
	"flag"
	"github.com/golang/glog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	listUrl   = flag.String("list", "http://music.163.com/discover/artist/cat?", "list url")
	detailUrl = flag.String("detail", "http://music.163.com", "detail url")
	interval  = flag.Int("interval", 4*7*24*3600, "crawler interval")
	dsn       = flag.String("dsn", "null", "mysql dsn,eg:[username[:password]@][protocol[(address)]]/dbname")
	table     = flag.String("table", "music_artist_spider", "mysql table name")
)

func stop(sigs chan os.Signal, exitCh chan int) {
	<-sigs
	glog.Info("receive stop signal")
	close(exitCh)
}

func main() {
	flag.Parse()
	defer glog.Flush()
	defer glog.Info("server exit")

    dataSource = *dsn
	initDB()

	cr := NewCrawler(*interval)
	exitCh := make(chan int)

	sigs := make(chan os.Signal)
	var wg sync.WaitGroup
	wg.Add(1)
	glog.Info("server start")
	go cr.Run(&wg, exitCh)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go stop(sigs, exitCh)
	wg.Wait()

}
