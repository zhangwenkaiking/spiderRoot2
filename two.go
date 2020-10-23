package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	_ "math/rand"
	"net/http"
	_ "net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var flag int

func main() {
	// 目标:http://xiaohua.zol.com.cn/new/3.html
	var (
		start int
		end   int
	)

	fmt.Println("***********欢迎***********")
	for m := 0; m <= 2; m++ {
		fmt.Printf("***                          \n")
	}
	fmt.Println("***请输入起始页(>=1):           ")
	fmt.Scan(&start)
	fmt.Println("***请输入end页:               ")
	fmt.Scan(&end)
	for m := 0; m <= 8; m++ {
		fmt.Printf("***                          \n")
	}
	for m := 0; m <= 4; m++ {
		fmt.Printf("***************************\n")
	}
	DoWork(start, end)
}
func DoWork(start int, end int) {
	fmt.Printf("正在爬取第%d到%d的范围\n", start, end)

	page := make(chan int)
	//明确目标
	for i := start + 6; i <= end+6; i++ {
		go SpiderPage(i, page)
	}
	for i := start; i <= 10000; i++ {
		fmt.Printf("第%d个已经完成爬取:", <-page)
	}
}

//, page chan int
func SpiderPage(i int, page chan int) {
	//// 目标:http://xiaohua.zol.com.cn/new/3.html
	var url string
	url = "http://xiaohua.zol.com.cn/new/" +
		strconv.Itoa(i+1) + ".html"
	//// 目标:http://xiaohua.zol.com.cn/new/3.html

	fmt.Println("正在爬取到网页：", i)
	fmt.Println("正在爬取的网页地址是：", url)

	//爬取全部内容
	result, err := HttpGet(url) //----------------->HttpGet(url)
	if err != nil {
		fmt.Println("err=:", err)
		return
	}

	//筛选信息 //       /detail60/59366.html
	re := regexp.MustCompile(`\/detail60\/[0-9]{5}.html`)
	if re == nil {
		fmt.Println("regexp.MustCompile err:=", err)
		return
	}

	//取关键信息
	//http://xiaohua.zol.com.cn/detail60/59366.html

	newUrls := make(map[int]string)
	// newjoyURLs := make(map[string]string)

	//获取第二个网页的地址
	//goquery.NewDocument("")

	joyURLs := re.FindAllStringSubmatch(result, -1)
	//fmt.Println("********************")
	//fmt.Println("joyURLs:\n", joyURLs)
	//fmt.Println("joyURLs(len):\n", len(joyURLs))
	//fmt.Println("********************")

	for k := 0; k < len(joyURLs); k += 3 {

		newUrls[k] = "http://xiaohua.zol.com.cn" + joyURLs[k][0]

		//fmt.Println(" len(joyURLs)：---->", len(joyURLs))
		fmt.Println("第二页到网址是：---->", newUrls[k])

		title, content, err := SpiderJoy(newUrls[k])
		//fmt.Println("content:------------>", content)
		if err != nil {
			fmt.Println("进入了continue之前的错误中：")
			fmt.Println("title, content, err:=", err)
			return
		}

		//*******筛选信息**********

		//***把内容输入到文件
		flag++
		filename := strconv.Itoa(flag) + ".txt"
		f, err1 := os.Create(filename)
		if err1 != nil {
			fmt.Println("os.Create:", err1)
			return //
		}
		f.WriteString(content)
		f.WriteString(title)
		// f.WriteString(content)
		f.Close()
		page <- flag
	}
}

//爬取第二个网页的内容
func SpiderJoy(url string) (title, content string, err error) {
	result, err1 := HttpGet(url) //----------------->HttpGet(url)
	if err1 != nil {
		err = err1
		return
	}

	//取关键信息
	content = Dosome(result)
	title = "dont tell you title"
	return title, content, err

}

//爬取方法
// func HttpGet(url string) (result string, err error) {
// 	resp, err1 := http.Get(url) //发送get请求
// 	if err1 != nil {
// 		err = err1
// 		return
// 	}
// 	defer resp.Body.Close() //关闭

// 	//读取网页Body内容
// 	buf := make([]byte, 1024*4)

// 	for {
// 		n, err := resp.Body.Read(buf)
// 		if n == 0 { //可能读取结束或者出问题
// 			fmt.Println("resp.Body.Read:---->\n", err)
// 			break
// 		}
// 		result += string(buf[:n]) //累加读取的内容
// 	}
// 	return
// }

//爬取方法
func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	//读取网页Body内容
	buf := make([]byte, 1024*4)
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 { //可能读取结束或者出问题
			fmt.Println("resp.Body.Read", err)
			break
		}

		d, err3 := GbkToUtf8(buf[:n])
		if err3 != nil {
			err = err3
			continue
		}
		//result += string(buf[:n])
		result += string(d)
		// fmt.Println(result)
	}
	return result, err
}
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 获取text内容
func Dosome(url string) (text string) {
	var result string
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(url))
	if err != nil {
		log.Fatalln(err)
	}

	dom.Find(".article-text").Each(func(i int, selection *goquery.Selection) {
		//fmt.Println("********************")
		//fmt.Println(selection.Text())
		result = selection.Text()
		//log.Println(selection.Text())
	})
	return result
}
