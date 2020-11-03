// @Title  	gogpio.go
// @Author  Rehtt  2019/8/4 下午 4:30
// @Update  Rehtt  2020/11/1 下午 2：00

// gogpio 是对GPIO简单的操作的库
//	功能：
//		1.设置指定的GPIO口是发射还是接收口（in/out）
//		2.设置发射口是高电平还是低电平（1/0）
//		3.释放GPIO口（Close()）
//		4.软件实现PWM（实验性，频率越高越消耗cpu资源）
//
//	示例：
//		package main
//
//		import (
//			"github.com/rehtt/gogpio"
//			"log"
//		)
//
//		func main() {
//			//绑定针脚号(BCM)
//			pin1, err1 := gogpio.PinBind(20)
//			pin2, err2 := gogpio.PinBind(21)
//			pin3, err3 := gogpio.PinBind(22)
//			pin4, err4 := gogpio.PinBind(23)
//
//			//声明针脚为out输出
//			if err1 != nil {
//				log.Println(err1)
//			}
//			out := pin1.SetOut()
//			out.High() //输出高电平
//			out.Low()  //输出低电平
//
//			//声明针脚为in输入
//			if err2 != nil {
//				log.Println(err2)
//			}
//			in := pin2.SetIn()
//			log.Println(in.Read()) //读取输入的数据
//
//			//不声明，直接读取数据。时合在其他程序使用此针脚时读取其数据
//			if err3 != nil {
//				log.Println(err3)
//			}
//			log.Println(pin3.Read())
//
//			//声明为PWM（此功能为实验性功能，这里的PWM由软件生成，所以运行时会占用一定的cpu资源，频率越高cpu占用也越高）
//			if err4 != nil {
//				log.Println(err4)
//			}
//			pwm := pin4.SetPWM()
//			err = pwm.SetFreq(5) //频率单位为Hz，数值 > 0.0
//			err = pwm.SetDC(20)  //占空比单位为% ，0.0 < 数值 < 100
//			if err != nil {
//				log.Println(err)
//			}
//			pwm.StartPWM() //开启PWM
//			pwm.StopPWM()  //关闭PWM
//
//			//释放
//			pin1.Close()
//			pin2.Close()
//			pin3.Close()
//			pin4.Close()
//		}
package gogpio

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	modePath  = ""
	valuePath = ""
	closePath = ""
	cPath     = ""
)

/** _Config
 * @Description: 针脚配置信息
 */
type _Config struct {
	Port string
	PWM  struct {
		freq float32
		dc   float32
		stop bool
	}
}

/** _Bind
 * @Description: 针脚模式
 */
type _Bind interface {
	SetOut() _OutOperating
	SetIn() _InOperating
	SetPWM() _PWM
	Read() ([]byte, error)
	Close()
}

/** _OutOperating
 * @Description: out(输出)模式的操作
 */
type _OutOperating interface {
	High()
	Low()
}

/** _InOperating
 * @Description: in(输入)模式的操作
 */
type _InOperating interface {
	Read() ([]byte, error)
}

/** _PWM
 * @Description: PWM模式的操作
 */
type _PWM interface {
	SetFreq(float32) error
	SetDC(float32) error
	StartPWM() error
	StopPWM()
}

/** PinBind
 * @Description: 输入针脚编号启动针脚
 * @param port
 * @return _Bind
 * @return error
 */
func PinBind(port int) (_Bind, error) {
	c := &_Config{
		Port: strconv.Itoa(port),
	}
	modePath = "/sys/class/gpio/gpio" + c.Port + "/direction"
	valuePath = "/sys/class/gpio/gpio" + c.Port + "/value"
	closePath = "/sys/class/gpio/unexport"
	cPath = "/sys/class/gpio/export"

	_, err := os.Lstat(valuePath)
	if err != nil {
		err = nil
		err = ioutil.WriteFile(cPath, []byte(c.Port), 0644)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

/** SetOut
 * @Description: 将目标针脚设为out（输出）模式
 * @receiver c
 * @return _OutOperating
 */
func (c *_Config) SetOut() _OutOperating {
	ioutil.WriteFile(modePath, []byte("out"), 0644)
	return c
}

/** SetIn
 * @Description: 将目标针脚设为in(输入)模式
 * @receiver c
 * @return _InOperating
 */
func (c *_Config) SetIn() _InOperating {
	ioutil.WriteFile(modePath, []byte("in"), 0644)
	return c
}

/** SetPWM
 * @Description: 将目标针脚设为PWM模式
 * @receiver c
 * @return _PWM
 */
func (c *_Config) SetPWM() _PWM {
	c.SetOut()
	return c
}

/** Read
 * @Description: in(输入)模式下读取针脚信息
 * @receiver c
 * @return []byte
 * @return error
 */
func (c *_Config) Read() ([]byte, error) {
	s, err := ioutil.ReadFile(valuePath)
	if err != nil {
		return nil, err
	}
	return s, err
}

/** Close
 * @Description: 释放针脚
 * @receiver c
 */
func (c *_Config) Close() {
	err := ioutil.WriteFile(closePath, []byte(c.Port), 0644)
	if err != nil {
		log.Println(err)
	}
}

/** High
 * @Description: 设置为高电平
 * @receiver c
 */
func (c *_Config) High() {
	ioutil.WriteFile(valuePath, []byte("1"), 0644)
}

/** Low
 * @Description: 设置为低电平
 * @receiver c
 */
func (c *_Config) Low() {
	ioutil.WriteFile(valuePath, []byte("0"), 0644)
}

/** SetFreq
 * @Description: 设置频率，单位为Hz，freq > 0.0
 * @receiver c
 * @param freq
 * @return error
 */
func (c *_Config) SetFreq(freq float32) error {
	if freq <= 0.0 {
		return errors.New("freq 需要大于0.0")
	}
	c.PWM.freq = freq
	return nil
}

/** SetDC
 * @Description: 设置空占比，单位为%，0.0 < dc <100.0
 * @receiver c
 * @param dc
 * @return error
 */
func (c *_Config) SetDC(dc float32) error {
	if dc >= 100.0 || dc <= 0.0 {
		return errors.New("dc 需要在0.0到100.0之间")
	}
	c.PWM.dc = dc
	return nil
}

/** StartPWM
 * @Description: 开启PWM
 * @receiver c
 * @return error
 */
func (c *_Config) StartPWM() error {
	c.PWM.stop = false
	if c.PWM.freq == 0 || c.PWM.dc == 0 {
		return errors.New("需要设置freq（频率）和dc（空占比）")
	}
	go func() {
		freq := 1 / c.PWM.freq
		high := time.Duration(freq*c.PWM.dc/100*1e6) * time.Microsecond
		low := time.Duration(freq*1e6)*time.Microsecond - high
		for {
			if c.PWM.stop {
				return
			}

			c.High()
			time.Sleep(high)
			c.Low()
			time.Sleep(low)
		}
	}()
	return nil
}

/** StopPWM
 * @Description: 关闭PWM
 * @receiver c
 */
func (c *_Config) StopPWM() {
	c.PWM.stop = true
}
