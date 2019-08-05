/*
对GPIO简单的操作
功能：
	1.设置指定的GPIO口是发射还是接收口（in/out）
	2.设置发射口是高电平还是低电平（1/0）
	3.释放GPIO口（Close()）

示例：
	package main

	import (
		"github.com/rehtt/gogpio"
		"log"
	)

	func main() {
		out := gogpio.Open(20)
		out.Way("out")
		out.High()
	
		in := gogpio.Open(21)
		in.Way("in")
		s,err := in.Read()
		if err != nil{
			log.Println(err)
		}
		log.Println(string(s))
	}

*/

package gogpio

import (
	"io/ioutil"
	"strconv"
)

type Config struct {
	Port string
	cWay string
}

type Operating interface {
	High()
	Low()
	Read() ([]byte, error)
	Way(way string)
	Close()
}

func (c *Config) High() {
	if c.cWay != "in" && c.cWay != "" {
		ioutil.WriteFile("/sys/class/gpio/gpio"+c.Port+"/value", []byte("1"), 0644)
	} else if c.cWay == "in" {
		panic("Must be \"out\" for available")
	} else {
		panic("Must specify how it works(.Way())")
	}
}

func (c *Config) Low() {
	if c.cWay != "in" && c.cWay != "" {
		ioutil.WriteFile("/sys/class/gpio/gpio"+c.Port+"/value", []byte("0"), 0644)
	} else if c.cWay == "in" {
		panic("Must be \"out\" for available")
	} else {
		panic("Must specify how it works(.Way())")
	}
}

func (c *Config) Read() ([]byte, error) {
	s, err := ioutil.ReadFile("/sys/class/gpio/gpio" + c.Port + "/value")
	if err != nil {
		return nil, err
	}

	return s, err
}

func (c *Config) Way(way string) {
	if way != "in" && way != "out" {
		panic("Must specify how it works(in/out)")
	}
	err := ioutil.WriteFile("/sys/class/gpio/gpio"+c.Port+"/direction", []byte(way), 0644)
	if err != nil {
		panic(err)
	}
	c.cWay = way

}

func (c *Config) Close() {
	ioutil.WriteFile("/sys/class/gpio/unexport", []byte(c.Port), 0644)
}

func Open(port int) Operating {
	Port := strconv.Itoa(port)
	c := &Config{Port: Port, cWay: ""}
	p := "/sys/class/gpio/"
	ioutil.WriteFile(p+"unexport", []byte(c.Port), 0644)
	ioutil.WriteFile(p+"export", []byte(c.Port), 0644)

	return c

}
