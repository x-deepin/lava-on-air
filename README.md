# Lava@Air

开发中~~ 

https://github.com/x-deepin/fpanel_controller 的替代者

主要特性

1. 所有信息通过2.4G无线通讯
2. 支持串口转发
3. 支持Fpanel控制(PowerSwitch)
4. 支持Fpanel状态查询(PowerLed)

模块结构

- manager 管理机客户端
  1. 下发Fpanel指令给特定测试机
  2. 查询特定测试机的Fpanel状态
  3. 维护RF24Network的树形网络拓扑结构与lava测试机间的地址转换
  
- driver 管理机驱动
  1. 为每个slaver建立虚拟的串口设备(pts)
  2. 将firmware\_master传递过来的消息解包分发到对应的pts上
  3. 将manager下发的Serial消息以及Fpanel控制指令封包后传递给firmware\_master

- firmware\_master  与管理机硬件连接
  1. 封包slaver的数据，并上传给driver
  2. 解包driver的数据，并分发到对应的slaver
  
- firmware\_slaver  与测试机硬件连接. 
  1. 执行manager下发的指令
  2. 定期汇报测试机当前状态
  3. 与测试机的真实COM口建立无线串口通路
  
- hardware 硬件设计
  1. master
  2. slaver

- include 公用协议的头文件
  1. protocol.h
