package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nynicg/cake/lib/log"
)

const ApnicUrl = "http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest"

const (
	domainFile = "domain.toml"
	ipFile = "apnic-latest.txt"
)

// deprecated
var cniplist map[string]struct{}
var domainlist domainRule
var domainCac domainCache

func init(){
	domainCac = domainCache{
		cache: map[string]int{},
	}
	cniplist = make(map[string]struct{})
}

func loadPassrule(){
	defer func() {
		go func() {
			domainlist = loadDomainList()
		}()
	}()
	info ,e := os.Stat(ipFile)
	// TODO check sum
	if e == nil && !info.IsDir(){
		log.Info("Find " ,ipFile)
		return
	}

	log.Info("Loading ip list from " ,ApnicUrl)
	log.Info("This process may take several minutes")
	f ,e := os.OpenFile(ipFile ,os.O_CREATE | os.O_RDWR | os.O_APPEND ,0755)
	if e != nil {
		panic(e)
	}
	defer f.Close()
	resp ,e := http.Get(ApnicUrl)
	if e != nil{
		panic(e)
	}
	f.WriteString("# Chinese ipv4(*.*.*.0) got from http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest.")
	rd := bufio.NewReader(resp.Body)
	for {
		line ,e := rd.ReadString('\n')
		if e != nil{
			break
		}
		if strings.Contains(line ,"ipv4") && strings.Contains(line ,"CN") {
			f.WriteString(getIP(line))
			f.WriteString("\n")
		}
	}
}

// apnic|CN|ipv4|39.0.8.0|2048|20110412|allocated
func getIP(line string) string{
	nopre := line[14:]
	i := strings.Index(nopre ,"|")
	return nopre[:i-2]
}

// Deprecated
func loadCNIPList() error{
	f ,e := os.OpenFile(ipFile ,os.O_RDONLY ,0755)
	if e != nil {
		panic(e)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	if line ,e := rd.ReadString('\n');e != nil{
		panic(e)
	}else{
		rmip := line[:len(line) - 2]
		cniplist[rmip] = struct{}{}
	}
	return nil
}


type domainRule struct {
	Bypass	map[string]uint8	`toml:"bypass"`
	Ads		map[string]uint8	`toml:"ads"`
}

func loadDomainList() domainRule{
	do := &domainRule{}
	if _ ,e := toml.DecodeFile(domainFile ,do);e != nil{
		panic(e)
	}
	return *do
}

const (
	BypassProxy = iota
	BypassTrue
	BypassDiscard
)
// Bypass 0:proxy ,1:bypass ,2:discard
func Bypass(dm string) int{

	for k := range domainlist.Bypass{
		if strings.HasSuffix(dm ,k) {
			return BypassTrue
		}
	}

	for k := range domainlist.Ads{
		if strings.HasSuffix(dm ,k){
			return BypassDiscard
		}
	}
	return BypassProxy
}

type domainCache struct {
	cache map[string]int
}

func PutDomainCache(domain string ,rule int){
	domainCac.cache[domain] = rule
}

func GetDomainCache(domain string) (int ,bool) {
	i ,r := domainCac.cache[domain]
	log.Debug("get domain bypass rule from cache " ,domain , " -> " ,i ,r)
	return i, r
}

