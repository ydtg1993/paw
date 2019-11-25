package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"main/robot"
)

func main() {
	done := make(chan bool)
	robots, err := getRobots()
	if err != nil {
		fmt.Println(err)
	}

	var rb robotInfo
	for _, rb = range robots {
		robot.Run(rb.Username, rb.Password, rb.Secret)
	}

	<-done
}

type robotInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Secret   string `json:"secret"`
}

func getRobots() (robots []robotInfo, err error) {
	content, err := ioutil.ReadFile("robots.json")
	if err != nil {
		return nil, errors.New("can't open the config file")
	}

	var robotsInfo []robotInfo
	json.Unmarshal(content, &robotsInfo)
	if err != nil {
		return nil, errors.New("explain json error")
	}

	return robotsInfo, nil
}
