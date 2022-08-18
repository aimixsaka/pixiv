package pixiv

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

const (
	proxy = "http://127.0.0.1:10809"
	imgPath = "C:/ELOI/pixiv/tag/"
)

var (
	proxyURL, _ = url.Parse(proxy)
	trans = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client = &http.Client{
		Transport: trans,
	}
)

func getWorkId(tag string, pageUp int) (ids []string) {
	var waitGroup sync.WaitGroup
	for ;pageUp > 0; pageUp-- {
		waitGroup.Add(1)
		go func(pageUp int) {
			req, err := http.NewRequest("GET", fmt.Sprintf("https://www.pixiv.net/ajax/search/artworks/%s", tag), nil)
			if err != nil {
				log.Fatalln("创建请求失败", err)
			}
			
			q := req.URL.Query()
			q.Add("word", tag)
			q.Add("order", "date_d")
			q.Add("mode", "all")
			q.Add("p", strconv.Itoa(pageUp))
			q.Add("s_mode", "s_tag_full")
			q.Add("type", "all")
			q.Add("lang", "zh")
			req.URL.RawQuery = q.Encode()
			
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Accept-Encoding", "gzip,deflate,br")
			req.Header.Set("Connection", "keep-alive")
			req.Header.Set("Referer", "https://www.pixiv.net/")
			req.Header.Add("Accept-Charset","utf-8")

			res, err := client.Do(req)
			if err != nil {
				log.Fatalln("客户端执行请求失败", err)
			}
			
			reader, _ := gzip.NewReader(res.Body)
			content, err := ioutil.ReadAll(reader)
			defer res.Body.Close()
			if err != nil {
				log.Fatalln("响应体读取失败", err)
			}
			
			idNum := jsoniter.Get(content, "body").Get("illustManga").Get("data").Size()
			if idNum == 0 {
				log.Fatalln("获取data错误", err)
			}

			for ;idNum > 0; idNum-- {
				ids = append(ids, jsoniter.Get(content, "body").Get("illustManga").Get("data", idNum-1).Get("id").ToString())
			}
			waitGroup.Done()

		}(pageUp)	
		waitGroup.Wait()
	}
	log.Printf("作品总数为%d", len(ids))
	return 
}

func getImgUrls(ids []string) (imgUrls []string) {
	var wp sync.WaitGroup
	for index, id := range(ids) {
		log.Printf("尝试获取id为%s的作品", id)
		wp.Add(1)
		go func(id string) {
			req, err := http.NewRequest("GET", fmt.Sprintf("https://www.pixiv.net/artworks/%s", id), nil)
			if err != nil {
				log.Fatalln("请求创建失败", err)
			}

			res, err := client.Do(req)
			if err != nil {
				log.Fatalln("无法得到响应", err)
			}
			if res.StatusCode != 200 {
				log.Fatalf("请求错误， code=%d", res.StatusCode)
			}	
			defer res.Body.Close()

			htmlByte, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalln("读取响应失败")
			}
			html := string(htmlByte)
			
			reg := regexp.MustCompile(`(?s)"original":"(.*?)"}`)
			url := reg.FindStringSubmatch(html)[1]
			if url == "" {
				log.Fatalf("获取id为%s的imgurl失败", id)
			}

			log.Printf("====已获取到第%d张作品", index)
			imgUrls = append(imgUrls, url)
			wp.Done()
		}(id)
		wp.Wait()
	}	
	
	log.Printf("图片url总数为%d", len(imgUrls))
	return
}


func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

func main() {
	ids := getWorkId("R-18", 2)
	downLoadImg(getImgUrls(ids), "R-18")
}