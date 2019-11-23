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
		//绑定针脚号(BCM)
		pin1 := gogpio.PinBind(20)
		pin2:=gogpio.PinBind(21)
		pin3:=gogpio.PinBind(22)

		//声明针脚为out输出
		out, err := pin1.SetOut()
		if err!=nil {
			log.Println(err)
		}
		out.High()								//输出高电平
		out.Low()								//输出低电平

		//声明针脚为in输入
		in,err:=pin2.SetIn()
		if err!=nil {
			log.Println(err)
		}
		log.Println(in.Read())					//读取输入的数据

		//不声明，直接读取数据。时合在其他程序使用此针脚时读取其数据
		log.Println(pin3.Read())

		//释放
		pin1.Close()
		pin2.Close()
		pin3.Close()
	}
*/

package gogpio

import (
	"io/ioutil"
	"log"
	"strconv"
)

var (
	modePath  = ""
	valuePath = ""
	closePath = ""
	cPath     = ""
)

type Config struct {
	Port string
}

type Bind interface {
	SetOut() (OutOperating, error)
	SetIn() (InOperating, error)
	Read() ([]byte, error)
	Close()
}
type OutOperating interface {
	High()
	Low()
}
type InOperating interface {
	Read() ([]byte, error)
}

func PinBind(port int) (Bind) {
	c := &Config{
		Port: strconv.Itoa(port),
	}
	modePath = "/sys/class/gpio/gpio" + c.Port + "/direction"
	valuePath = "/sys/class/gpio/gpio" + c.Port + "/value"
	closePath = "/sys/class/gpio/unexport"
	cPath = "/sys/class/gpio/export"

	return c
}



func (c *Config) SetOut() (OutOperating, error) {
	err := ioutil.WriteFile(cPath, []byte(c.Port), 0644)
	ioutil.WriteFile(modePath, []byte("out"), 0644)
	return c, err
}

func (c *Config) SetIn() (InOperating, error) {
	err := ioutil.WriteFile(cPath, []byte(c.Port), 0644)
	ioutil.WriteFile(modePath, []byte("in"), 0644)
	return c, err
}
func (c *Config) Read() ([]byte, error) {
	s, err := ioutil.ReadFile(valuePath)
	if err != nil {
		return nil, err
	}
	return s, err
}

func (c *Config) Close() {
	err := ioutil.WriteFile(closePath, []byte(c.Port), 0644)
	if err != nil {
		log.Println(err)
	}
}

func (c *Config) High() {
	ioutil.WriteFile(valuePath, []byte("1"), 0644)

}
func (c *Config) Low() {
	ioutil.WriteFile(valuePath, []byte("0"), 0644)

}
