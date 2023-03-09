package main

import (
	"fmt"
	bilibilicommentspider "go_spider/bilibilicomment_spider"
	doubanspider "go_spider/douban_spider"

	"strconv"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Printf("正在爬取第 %d 页数据\n", i)
		doubanspider.Spider(strconv.Itoa(i * 25))
	}
	bilibilicommentspider.Spider()
	fmt.Println("爬取完毕")
}
