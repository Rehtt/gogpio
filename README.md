# gogpio

对GPIO简单的操作
功能：
1. 设置指定的GPIO口是发射还是接收口（in/out）
2. 设置发射口是高电平还是低电平（1/0）
3. 释放GPIO口（Close()）
4. 软件实现PWM（实验性功能，可能会导致系统不稳定）

安装：

    go get github.com/Rehtt/gogpio	

使用示例：

    func main() {
    //绑定针脚号(BCM)
    pin1, err1 := gogpio.PinBind(20)
    pin2, err2 := gogpio.PinBind(21)
    pin3, err3 := gogpio.PinBind(22)
    pin4, err4 := gogpio.PinBind(23)
    
    //声明针脚为out输出
    if err1 != nil {
    log.Println(err1)
    }
    out := pin1.SetOut()
    out.High() //输出高电平
    out.Low()  //输出低电平
    
    //声明针脚为in输入
    if err2 != nil {
    log.Println(err2)
    }
    in := pin2.SetIn()
    log.Println(in.Read()) //读取输入的数据
    
    //不声明，直接读取数据。时合在其他程序使用此针脚时读取其数据
    if err3 != nil {
    log.Println(err3)
    }
    log.Println(pin3.Read())
    
    //声明为PWM（此功能为实验性功能，这里的PWM由软件生成，所以运行时会占用一定的cpu资源，频率越高cpu占用也越高）
    if err4 != nil {
    log.Println(err4)
    }
    pwm := pin4.SetPWM()
    err = pwm.SetFreq(5) //频率单位为Hz，数值 > 0.0
    err = pwm.SetDC(20) //占空比单位为% ，0.0 < 数值 < 100
    if err != nil {
    log.Println(err)
    }
    pwm.StartPWM() //开启PWM
    pwm.StopPWM() //关闭PWM
    
    //释放
    pin1.Close()
    pin2.Close()
    pin3.Close()
    pin4.Close()
    }