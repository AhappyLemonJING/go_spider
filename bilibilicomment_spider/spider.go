package bilibilicommentspider

import (
	"encoding/json"
	"fmt"
	"go_spider/config"
	"io/ioutil"

	"net/http"
)

// https://api.bilibili.com/x/v2/reply/main?next=0&type=1&oid=864489807&mode=3&plat=1
func Spider() {
	// 发送请求
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/v2/reply/main?next=0&type=1&oid=864489807&mode=3&plat=1", nil)
	if err != nil {
		fmt.Println("req err", err)
	}
	// 加请求头伪造成浏览器访问
	req.Header.Set("authority", "api.bilibili.com")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-fetch-dest", "script")
	req.Header.Set("sec-fetch-mode", "no-cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("referer", "https://www.bilibili.com/bangumi/play/ep733750?from_spmid=666.4.banner.2")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败", err)
	}
	defer resp.Body.Close()
	// 读取信息
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("解析失败", err)
	}
	// 获取节点 网页开发者模式进行搜索想要的selector
	var resultList BilibiliComment
	_ = json.Unmarshal(bodyText, &resultList)
	for _, result := range resultList.Data.Replies {
		fmt.Println("一级评论", result.Content.Message)
		data1 := &BiliComment{
			Comments: result.Content.Message,
		}
		err := config.DB.Model(new(BiliComment)).Create(&data1).Error
		if err != nil {
			fmt.Println("插入失败")
		}
		for _, replay := range result.Replies {
			fmt.Println("二级评论", replay.Content.Message)
			data2 := &BiliComment{
				Comments: replay.Content.Message,
			}
			err := config.DB.Model(new(BiliComment)).Create(&data2).Error
			if err != nil {
				fmt.Println("插入失败")
			}
		}
	}

}
