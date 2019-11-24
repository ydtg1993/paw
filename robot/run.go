package robot

import (
	"bytes"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"main/google"
	"net/http"
	"time"
)

func Run(username,password,secret string)  {
	done := make(chan int)
	const (
		seleniumPath = `D:\chromedriver.exe`
		port            = 9515
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
	login(username,password,secret,w_b1)

	<-done
}

func login(username,password,secret string,w_b1 selenium.WebDriver)  {
	err := w_b1.Get("http://bibi.cnluyao.cn/bi")
	if err != nil {
		fmt.Println("get page faild", err.Error())
		return
	}
	wes,err := w_b1. FindElements(selenium.ByClassName,"form-control")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var input selenium.WebElement
	for _,input = range wes {
		idName,_ := input.GetAttribute("id")
		if idName == "username" {
			input.SendKeys(username)
		}else {
			input.SendKeys(password)
		}
	}

	we,err := w_b1.FindElement(selenium.ByTagName,"button")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for i:=0;i<5;i++ {
		we.Click()
	}
	time.Sleep(1 * time.Second)

	we,err = w_b1.FindElement(selenium.ByID,"googleCodeNum")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	code := google.Index(secret)
	fmt.Println(code)
	we.SendKeys(code)
	time.Sleep(1 * time.Second)
	we,err = w_b1.FindElement(selenium.ByTagName,"button")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for i:=0;i<5;i++ {
		we.Click()
	}
}

func getUserList()  {
	
}

func request(url string,body []byte) (response *http.Response,err error)  {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	req.Header.Set("Accept","application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Encoding","gzip, deflate")
	req.Header.Set("Accept-Language","zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Content-Length","84")
	req.Header.Set("Content-Type","application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host","bibi.cnluyao.cn")
	req.Header.Set("Origin","http://bibi.cnluyao.cn")
	req.Header.Set("Referer","http://bibi.cnluyao.cn/bi/login_public.html")
	req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	req.Header.Set("X-Requested-With","XMLHttpRequest")

	defer req.Body.Close()
	return client.Do(req)
}
