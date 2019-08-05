# gogpio

对GPIO简单的操作
<br>
功能：<br>
	1.设置指定的GPIO口是发射还是接收口（in/out）<br>
	2.设置发射口是高电平还是低电平（1/0）<br>
	3.释放GPIO口（Close()）<br>

示例：<br>
```
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

```
