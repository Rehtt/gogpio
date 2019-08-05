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
		out := gogpio.Open(20,gogpio.OUT)
		out.High()

		in := gogpio.Open(21,gogpio.IN)
		s,err := in.Read()
		if err != nil{
			log.Println(err)
		}
		log.Println(string(s))
	}

*/

package gogpio

import (
	"errors"
	"io/ioutil"
	"strconv"
)

var OUT string = "out"
var IN string = "in"

type Config struct {
	Port string
	cWay string
}

type Operating interface {
	High()
	Low()
	Read() ([]byte, error)
	Close()
}

func (c *Config) High() {
	if c.cWay != OUT {
		ioutil.WriteFile("/sys/class/gpio/gpio"+c.Port+"/value", []byte("1"), 0644)
	} else {
		panic("Must be \"gogpio.OUT\" for available")
	}
}

func (c *Config) Low() {
	if c.cWay != IN {
		ioutil.WriteFile("/sys/class/gpio/gpio"+c.Port+"/value", []byte("0"), 0644)
	} else {
		panic("Must be \"gogpio.OUT\" for available")
	}
}

func (c *Config) Read() ([]byte, error) {
	s, err := ioutil.ReadFile("/sys/class/gpio/gpio" + c.Port + "/value")
	if err != nil {
		return nil, err
	}

	return s, err
}

func (c *Config) Close() {
	ioutil.WriteFile("/sys/class/gpio/unexport", []byte(c.Port), 0644)
}

func Open(port int, way string) (Operating, error) {
	Port := strconv.Itoa(port)
	c := &Config{Port: Port, cWay: ""}
	if way != IN && way != OUT {
		c.Close()
		return nil, errors.New("Must specify how it works(gogpio.OUT/gogpio.IN)")
	}
	err := ioutil.WriteFile("/sys/class/gpio/gpio"+c.Port+"/direction", []byte(way), 0644)
	if err != nil {
		c.Close()
		return nil, err
	}
	c.cWay = way
	c.Close()
	ioutil.WriteFile("/sys/class/gpio/export", []byte(c.Port), 0644)

	return c, nil

}
