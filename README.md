# gogpio

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
		//指定GPIO口及类型（in/out）
		c := &gogpio.Config{
			Port: "21",
			Way:  "in",
		}
		c2 := &gogpio.Config{
			Port: "20",
			Way:  "out",
		}

		//开启GPIO口
		in := gogpio.Open(c2)
		out := gogpio.Open(c)

		//发射口输出高电平（1）
		err := out.Write([]byte("1"))
		if err != nil {
			log.Fatal(err)
		}

		//接受口读取
		res, err := in.Read()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(res))

		//释放
		in.Close()
		out.Close()
	}
