package robot

import (
	"fmt"
	"github.com/tebeka/selenium"
	"main/google"
	"time"
)

func Login(username, password, secret string, w_b1 selenium.WebDriver) bool {
	err := w_b1.Get("http://bibi.cnluyao.cn/bi")
	if err != nil {
		fmt.Println("get page faild", err.Error())
		return false
	}
	wes, err := w_b1.FindElements(selenium.ByClassName, "form-control")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	var input selenium.WebElement
	for _, input = range wes {
		idName, _ := input.GetAttribute("id")
		if idName == "username" {
			input.SendKeys(username)
		} else {
			input.SendKeys(password)
		}
	}

	we, err := w_b1.FindElement(selenium.ByTagName, "button")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	flag := doubleCheckCurrent(w_b1, we, "http://bibi.cnluyao.cn/bi/gooleCode.html")
	if !flag {
		//return false
	}

	/*google auth*/
	we, err = w_b1.FindElement(selenium.ByID, "googleCodeNum")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for i := 0; i < 4; i++ {
		flag = auth(w_b1, we, secret)
		if flag {
			break
		}
	}

	return true
}

func auth(web selenium.WebDriver, we selenium.WebElement, secret string) bool {
	code := google.Index(secret)
	we.SendKeys(code)
	we, err := web.FindElement(selenium.ByTagName, "button")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	flag := doubleCheckCurrent(web, we, "http://bibi.cnluyao.cn/bi/index.html")
	if !flag {
		we.Clear()
	}
	return flag
}

func doubleCheckCurrent(web selenium.WebDriver, we selenium.WebElement, url string) bool {
	i := 0
	currentUrl, _ := web.CurrentURL()
	for ; currentUrl != url; i++ {
		fmt.Println(currentUrl)
		we.Click()
		we.Click()
		we.Click()
		if i > 3 {
			return false
		}
		time.Sleep(1 * time.Second)
		currentUrl, _ = web.CurrentURL()
	}

	return true
}
