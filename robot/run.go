package robot

import (
	"compress/gzip"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type MyCookieData struct {
	Session string
	Batman  string
}

func Run(username, password, secret string) {
	done := make(chan int)
	const (
		seleniumPath = `D:\chromedriver.exe`
		port         = 9515
	)
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(seleniumPath, port, opts...)
	if nil != err {
		fmt.Println("start a chromedriver service falid", err.Error())
		return
	}
	defer service.Stop()
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			//"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			//"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36",
		},
	}

	caps.AddChrome(chromeCaps)
	w_b1, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		fmt.Println("connect to the webDriver faild", err.Error())
		return
	}
	//但是不会导致seleniumServer关闭
	defer w_b1.Quit()
	flag := Login(username, password, secret, w_b1)
	w_b1.Refresh()
	if !flag {
		return
	}
	getUserList(w_b1)

	<-done
}

func getUserList(wb selenium.WebDriver) {
	data := url.Values{}
	data.Set("page", "1")
	data.Set("limit", "500")
	data.Set("start", "2012-11-24 00:00:00")
	data.Set("end", "2019-11-24 23:00:00")
	data.Set("pid", "")
	data.Set("rid", "")

	var ck *MyCookieData
	ck = new(MyCookieData)
	cookies, _ := wb.GetCookies()
	var cookie selenium.Cookie
	for _, cookie = range cookies {
		if cookie.Name == "batmanCok" {
			ck.Batman = cookie.Value
		} else if cookie.Name == "JSESSIONID" {
			ck.Session = cookie.Value
		}
	}

	for i := 0; i < 5; i++ {
		res := request("http://bibi.cnluyao.cn/bi/extension/extensions", ck, strings.NewReader(data.Encode()))
		fmt.Println(string(res))
	}
}

func request(url string, cookie *MyCookieData, body io.Reader) []byte {
	payload := strings.NewReader("page=1&limit=500&start=2012-11-24%2000%3A00%3A00&end=2019-11-24%2023%3A00%3A00&pid=&rid=")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	ck := fmt.Sprintf("JSESSIONID=%s;batmanCok=%s", cookie.Session, cookie.Batman)
	fmt.Println(ck)
	req.Header.Add("Cookie", ck)
	req.Header.Add("Referer", "http://bibi.cnluyao.cn/bi/gameMaster/Newgeneralizedetails.html")
	req.Header.Add("Origin", "http://bibi.cnluyao.cn")
	req.Header.Add("Host", "bibi.cnluyao.cn")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Length", "83")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var response []byte
	if err != nil {
		fmt.Println(err.Error())
		return response
	}

	reader, error := gzip.NewReader(res.Body)
	if error != nil {
		fmt.Println(error.Error())
		return response
	}
	response, _ = ioutil.ReadAll(reader)
	cookies := res.Cookies()
	for _, c := range cookies {
		if c.Name == "batmanCok" {
			cookie.Batman = c.Value
		}
	}

	return response
}
