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
		pin2 := gogpio.PinBind(21)
		pin3 := gogpio.PinBind(22)
		pin4 := gogpio.PinBind(23)

		//声明针脚为out输出
		out, err := pin1.SetOut()
		if err != nil {
			log.Println(err)
		}
		out.High() //输出高电平
		out.Low()  //输出低电平

		//声明针脚为in输入
		in, err := pin2.SetIn()
		if err != nil {
			log.Println(err)
		}
		log.Println(in.Read()) //读取输入的数据

		//不声明，直接读取数据。时合在其他程序使用此针脚时读取其数据
		log.Println(pin3.Read())

		//声明为PWM（此功能为实验性功能，这里的PWM由软件生成，所以运行时会占用一定的cpu资源，频率越高cpu占用也越高）
		pwm, err := pin4.SetPWM()
		if err != nil {
			log.Println(err)
		}
		err = pwm.SetFreq(5) //频率单位为Hz，数值 > 0.0
		err = pwm.SetDC(20)  //占空比单位为% ，0.0 < 数值 < 100
		if err != nil {
			log.Println(err)
		}
		pwm.StartPWM() //开启PWM
		pwm.StopPWM()  //关闭PWM

		//释放
		pin1.Close()
		pin2.Close()
		pin3.Close()
		pin4.Close()
	}
*/

package gogpio

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

var (
	modePath  = ""
	valuePath = ""
	closePath = ""
	cPath     = ""
)

type _Config struct {
	Port string
	PWM  struct {
		freq float32
		dc   float32
		stop bool
	}
}

type _Bind interface {
	SetOut() (_OutOperating, error)
	SetIn() (_InOperating, error)
	SetPWM() (_PWM, error)
	Read() ([]byte, error)
	Close()
}
type _OutOperating interface {
	High()
	Low()
}
type _InOperating interface {
	Read() ([]byte, error)
}
type _PWM interface {
	SetFreq(float32) error
	SetDC(float32) error
	StartPWM()
	StopPWM()
}

func PinBind(port int) _Bind {
	c := &_Config{
		Port: strconv.Itoa(port),
	}
	modePath = "/sys/class/gpio/gpio" + c.Port + "/direction"
	valuePath = "/sys/class/gpio/gpio" + c.Port + "/value"
	closePath = "/sys/class/gpio/unexport"
	cPath = "/sys/class/gpio/export"

	return c
}

func (c *_Config) SetOut() (_OutOperating, error) {
	err := ioutil.WriteFile(cPath, []byte(c.Port), 0644)
	ioutil.WriteFile(modePath, []byte("out"), 0644)
	return c, err
}

func (c *_Config) SetIn() (_InOperating, error) {
	err := ioutil.WriteFile(cPath, []byte(c.Port), 0644)
	ioutil.WriteFile(modePath, []byte("in"), 0644)
	return c, err
}

//freq	:PWM频率（Hz）	freq > 0.0
//dc	:PWM占空比		0.0<=dc<=100.0
func (c *_Config) SetPWM() (_PWM, error) {
	_, err := c.SetOut()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *_Config) Read() ([]byte, error) {
	s, err := ioutil.ReadFile(valuePath)
	if err != nil {
		return nil, err
	}
	return s, err
}

func (c *_Config) Close() {
	err := ioutil.WriteFile(closePath, []byte(c.Port), 0644)
	if err != nil {
		log.Println(err)
	}
}

func (c *_Config) High() {
	ioutil.WriteFile(valuePath, []byte("1"), 0644)

}
func (c *_Config) Low() {
	ioutil.WriteFile(valuePath, []byte("0"), 0644)

}

func (c *_Config) SetFreq(freq float32) error {
	if freq <= 0.0 {
		return errors.New("freq 需要大于0.0")
	}
	c.PWM.freq = freq
	return nil
}
func (c *_Config) SetDC(dc float32) error {
	if dc >= 100.0 || dc <= 0.0 {
		return errors.New("dc 需要在0.0到100.0之间")
	}
	c.PWM.dc = dc
	return nil
}
func (c *_Config) StartPWM() {
	c.PWM.stop = false

	freq := 1 / c.PWM.freq
	high := time.Duration(freq*c.PWM.dc/100*1e6) * time.Microsecond
	low := time.Duration(freq*1e6)*time.Microsecond - high
	go func() {
		for {
			if !c.PWM.stop {
				break
			}
			c.High()
			time.Sleep(high)
			c.Low()
			time.Sleep(low)
		}
	}()
}
func (c *_Config) StopPWM() {
	c.PWM.stop = true
}
