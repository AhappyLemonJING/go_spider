package doubanspider

import (
	"fmt"
	"go_spider/config"

	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func Spider(page string) {
	// 发送请求
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250?start="+page+"&filter=", nil)
	if err != nil {
		fmt.Println("req err", err)
	}
	// 加请求头伪造成浏览器访问
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Referer", "https://movie.douban.com/chart")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败", err)
	}
	defer resp.Body.Close()
	// 解析网页
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析失败", err)
	}
	// 获取节点 网页开发者模式进行搜索想要的selector
	docDetail.Find("#content > div > div.article > ol > li").
		Each(func(i int, s *goquery.Selection) {

			title := s.Find("div > div.info > div.hd > a > span:nth-child(1)").Text()
			img := s.Find("div > div.pic > a > img")
			imgTmp, ok := img.Attr("src") // 获取img标签下的src属性内容
			info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text()
			score := s.Find("div > div.info > div.bd > div > span.rating_num").Text()
			quote := s.Find("div > div.info > div.bd > p.quote > span").Text()
			if ok {
				director, actor, year := InfoSpilt(info)
				data := &MovieData{
					Picture:  imgTmp,
					Director: director,
					Actor:    actor,
					Year:     year,
					Title:    title,
					Quote:    quote,
					Score:    score,
				}
				err := config.DB.Model(new(MovieData)).Create(&data).Error

				if err != nil {
					fmt.Println("插入失败")
				}
			}

		})

}

func InfoSpilt(info string) (director, actor, year string) {
	directorRe, _ := regexp.Compile(`导演:(.*)主演`)
	director = string(directorRe.Find([]byte(info)))
	actorRe, _ := regexp.Compile(`主演:(.*)`)
	actor = string(actorRe.Find([]byte(info)))
	yearRe, _ := regexp.Compile(`(\d+)`)
	year = string(yearRe.Find([]byte(info)))
	return
}
