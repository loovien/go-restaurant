# Golang 实现现实的餐馆

## 描述

学习golang tcp编程的一个例子. 无任何目的, 纯粹的学习. 玩. 😄

## 依赖

1. golang (1.8)
2. godep `go get github.com/tools/godep` (可选)


## 安装

- 使用已经编译好的二进制文件 `go-restaurant.exe`(windows) `go-restaurant`(unix-like)
- 自己编译使用

    1. `cd /path/to/project`
    1. `godep restore`
    2. `go build`

## 启动服务

- 默认启动 `go-restaurent.exe`
- 自定义绑定参数 `go-restaurant.exe --addr=0.0.0.0:1018` 服务监听1018端口


## 使用说明

1. 功能说明: 以我请人吃饭为🌰. 服务端服务(餐馆), 客户端A(我), 客户端B(我朋友), 服务端启动服务, 我连上务服务后询问是否有空座(2人座, 4人座, ..), 有座后, 我坐下来. 等待我朋友(客户端B)连上服务, 找到对应我坐下的桌号坐下. 此时开始点餐. 开始吃饭.

1. 项目目录

    - **conf** 项目配置文件目录
    - **Godeps** 工具`godep`生成的依赖文件
    - **logs** 程序日志
    - ***.go**  程序文件

1. 功能(协议列表) _协议为简单的使用前4个字节使用小端序表示包体长度_

    - 请求是否还有空座位. (协议内容)

        ```json
            {
                "cmd": "emptySeat", // 标识
                "acceptUnion": "1" // 是否接受拼桌
            }
        ```

    - 就坐. (协议内容)

        ```json
            {
                "cmd": "sitDown", // 标识
                "tableNo": "A01", // 桌号
                "token": "luowen" // 唯一标识, 表示您做下了这个桌号, 下次有人坐这个位置, 必须报口令
            }
        ```

    - 获取餐馆菜单. (协议内容)

        ```json
            {
                "cmd": "orderMeat", // 标识
                "tableNo": "A01", // 桌号
                "token": "luowen", // 桌子口令
                "menu": ["gbjd", "yxrs"] // 菜单
            }
        ```

## 任务呆完成

1. 确保用户支付在做菜
1. 菜的余量处理
2. 拼桌处理没完成

