package manhua

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var client = &http.Client{}

func Home() {
	req, _ := http.NewRequest("GET", "https://www.mhua5.com/", nil)
	req.Header.Set("Connection", "keep-alive") //设置请求头模拟浏览器登录
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Referer", "https://www.mhua5.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Find the review items
	doc.Find(".in-sec-wr .in-comic--type-a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find(".comic__title a").Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}

type ManhuaInfoListResponse struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data []ManhuaInfoModel `json:"data"`
}

type ManhuaInfoModel struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Pic         string   `json:"pic"`
	Picx        string   `json:"picx"`
	LocalPic    string   `json:"localPic"`
	Serialize   string   `json:"serialize"`
	Author      string   `json:"author"`
	Text        string   `json:"text"`
	Content     string   `json:"content"`
	Hits        string   `json:"hits"`
	Ticket      string   `json:"ticket"`
	Score       string   `json:"score"`
	Addtime     string   `json:"addtime"`
	URL         string   `json:"url"`
	ChapterName string   `json:"chapter_name"`
	ChapterURL  string   `json:"chapter_url"`
	ChapterNums string   `json:"chapter_nums"`
	Tags        []string `json:"tags"`
}

func Search(key string, page int) ManhuaInfoListResponse {
	callback := "jQuery22402428632092686871_1667476283510"
	requestUrl := "https://www.mhua5.com/index.php/api/data/comic?callback=" + callback + "&key="
	requestUrl = requestUrl + key + "&page=" + strconv.FormatInt(int64(page), 10) + "&_=1667476283515"
	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("Connection", "keep-alive") //设置请求头模拟浏览器登录
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Referer", "https://www.mhua5.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(body)
	jsonString = strings.Trim(jsonString, ");")
	jsonString = strings.TrimLeft(jsonString, callback+"(")
	result := ManhuaInfoListResponse{}
	json.Unmarshal([]byte(jsonString), &result)
	return result
}

type ChapterModel struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func GetChapterList(manhuaUrl string) []ChapterModel {
	req, _ := http.NewRequest("GET", "https://www.mhua5.com"+manhuaUrl, nil)
	req.Header.Set("Connection", "keep-alive") //设置请求头模拟浏览器登录
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Referer", "https://www.mhua5.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var list []ChapterModel
	doc.Find(".container--response .de-container .de-chapter .chapter__list .chapter__list-box .j-chapter-item").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		capter := s.Find(".j-chapter-link")
		title := capter.Text()
		// 去除首尾空格
		title = strings.TrimSpace(title)
		href, _ := capter.Attr("href")
		list = append(list, ChapterModel{Title: title, Url: href})
	})
	return list
}

func GetChapterDetail(detailUrl string) []string {
	req, _ := http.NewRequest("GET", "https://www.mhua5.com"+detailUrl, nil)
	req.Header.Set("Connection", "keep-alive") //设置请求头模拟浏览器登录
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Referer", "https://www.mhua5.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var list []string
	doc.Find(".read-container .rd-article-wr .rd-article__pic").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		imgUrl, _ := s.Find("img").Attr("data-original")
		list = append(list, imgUrl)
	})
	return list
}
