## `Bilibili` 直播间弹幕处理库

`golang`的b站信息流处理, 弹幕/礼物/大航海/SC等一切你想要的数据

### 介绍

`go-bili-chat` 是一个用于处理 Bilibili 直播间信息流库，可以用于开发自己的 Bilibili 直播间信息流处理程序。

b站直播间信息流以Websocket传输并加密, 含有几十个不同的命令, 本项目对其进行了解析, 并提供了一些简单的处理方法, 以便于开发者快速开发自己的程序。

### 快速使用

信息流处理流程为: 客户端收到信息 -> 解析后由处理器进行分发

因此, 你需要先将命令处理函数绑定到处理器, 再开启直播间进行处理,

#### 1. 一个简单的使用示例:

```go
package main

import (
	"fmt"
	bili "github.com/FishZe/go-bili-chat"
	handle "github.com/FishZe/go-bili-chat/handler"
)

func main() {
	// 新建一个命令处理器
	h := bili.GetNewHandler()
	// 注册一个处理，将该直播间的弹幕消息绑定到这个函数
	h.AddOption(handle.CmdDanmuMsg, 26097368, func(event handle.MsgEvent) {
		// 打印出弹幕消息
		fmt.Printf("[%v] %v: %v\n", event.RoomId, event.DanMuMsg.Data.Sender.Name, event.DanMuMsg.Data.Content)
	})
	// 连接到直播间
	_ = h.AddRoom(26097368)
	// 启动处理器
	h.Run()
}

```
特殊地, 绑定函数的直播间号为0时，绑定所有房间

#### 2. 也可以先运行命令处理器，再添加房间：

```go
package main

import (
	"fmt"
	bili "github.com/FishZe/go-bili-chat"
	handle "github.com/FishZe/go-bili-chat/handler"
	"time"
)

func main() {
	h := bili.GetNewHandler()
	// 运行处理器
	go h.Run()
	h.AddOption(handle.CmdDanmuMsg, 26097368, func(event handle.MsgEvent) {
		fmt.Printf("[%v] %v: %v\n", event.RoomId, event.DanMuMsg.Data.Sender.Name, event.DanMuMsg.Data.Content)
	})
	_ = h.AddRoom(26097368)
	for {
		time.Sleep(time.Second)
	}
}
```
#### 3. 当然了，也可以删除房间：

```go
package main

import (
	"fmt"
	bili "github.com/FishZe/go-bili-chat"
	handle "github.com/FishZe/go-bili-chat/handler"
)

func main() {
	h := bili.GetNewHandler()
	h.AddOption(handle.CmdDanmuMsg, 26097368, func(event handle.MsgEvent) {
		fmt.Printf("[%v] %v: %v\n", event.RoomId, event.DanMuMsg.Data.Sender.Name, event.DanMuMsg.Data.Content)
	})
	_ = h.AddRoom(26097368)
	_ = h.DelRoom(26097368)
	h.Run()
}
```

#### 4. 对于连接房间过多的情况, 建议先`go h.Run()`再添加房间, 否则会阻塞通道造成内存大量占用

```go
go h.Run()
_ = h.AddRoom(26097368)
_ = h.AddRoom(26097369)
_ = h.AddRoom(26097370)
// sth...
```

#### 5. 如果你有删除房间或debug的需求, 可以绑定函数时传入函数的名称
```go
h.AddOption(handle.CmdDanmuMsg, 26097368, func(event handle.MsgEvent) {
    fmt.Printf("[%v] %v: %v\n", event.RoomId, event.DanMuMsg.Data.Sender.Name, event.DanMuMsg.Data.Content)
}, "弹幕")
h.DelOption("弹幕")
```

**关于为什么在处理绑定函数时, 多一个直播间号的参数, 因为考虑到可能会有根据不同的直播间分发处理消息的需求**

#### 6. 也可以获得当前的连接房间个数
```go
sum := h.CountRoom()
fmt.Println("当前连接房间数: ", sum)
```

#### 7. 对于接入的服务器: 程序提供了3个模式:

```text
DefaultPriority: 默认模式, 按照b站提供的顺序进行连接
DelayPriority: 低延迟模式, 将会在连接前计算到所有提供的服务器的延迟, 并选择最低的(连接会慢一些)
NoCDNPriority: 不使用CDN模式, 适合大量连接的情况
```

默认使用的模式为`DefaultPriority`, 有需要的可以使用这种方式修改:

```go
bili.SetClientPriorityMode(bili.DelayPriority)
```

#### 8. 默认使用的Json解析器为`sonic`, 有需要自定义`json`解析器的可以使用这种方式:

```go
type Json struct{}

func (j *Json) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j *Json) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

bili.SetJsonCoder(&Json{})
```

具体的示例代码可以查看`example/main.gp`

#### 9. 对于需要查看`DEBUG`日志的情况, 可以修改日志等级
```go
bili.ChangeLogLevel(log.DebugLevel)
```

开启`DEBUG`日志后, 可以看到所有的数据包和分发过程, 如果想要直观的看到分发到的函数名称, 可以在`AddOption`中写明该函数的备注名
```go
h.AddOption(handle.CmdDanmuMsg, 21545805, func(event handle.MsgEvent) {
    fmt.Printf("[%v][弹幕] %v (%v): %v\n", event.RoomId, event.DanMuMsg.Data.Sender.Name, event.DanMuMsg.Data.Medal.MedalName, event.DanMuMsg.Data.Content)
}, "弹幕处理")
````

控制台会显示类似如下的日志

```text
[bili-live][02-21 16:42:55][DEBUG]: distribute DANMU_MSG cmd to 弹幕处理
```

#### 所有的消息类型:
这些常量请填入`go-bili-chat.GetNewHandler().AddOption()`的第一个参数
```text
常量名                            原始命令
CmdDanmuMsg                     "DANMU_MSG"
CmdSuperChatMessage             "SUPER_CHAT_MESSAGE"
CmdSuperChatMessageJpn          "SUPER_CHAT_MESSAGE_JPN"
CmdWatchedChange                "WATCHED_CHANGE"
CmdSendGift                     "SEND_GIFT"
CmdOnlineRankCount              "ONLINE_RANK_COUNT"
CmdOnlineRankV2                 "ONLINE_RANK_V2"
CmdOnlineRankTop3               "ONLINE_RANK_TOP3"
CmdLikeInfoV3Click              "LIKE_INFO_V3_CLICK"
CmdInteractWord                 "INTERACT_WORD"
CmdStopLiveRoomList             "STOP_LIVE_ROOM_LIST"
CmdLikeInfoV3Update             "LIKE_INFO_V3_UPDATE"
CmdHotRankChange                "HOT_RANK_CHANGED"
CmdNoticeMsg                    "NOTICE_MSG"
CmdRoomRealTimeMessageUpdate    "ROOM_REAL_TIME_MESSAGE_UPDATE"
CmdWidgetBanner                 "WIDGET_BANNER"
CmdHotRankChangedV2             "HOT_RANK_CHANGED_V2"
CmdGuardHonorThousand           "GUARD_HONOR_THOUSAND"
CmdLive                         "LIVE"
CmdRoomChange                   "ROOM_CHANGE"
CmdRoomBlockMsg                 "ROOM_BLOCK_MSG"
CmdFullScreenSpecialEffect      "FULL_SCREEN_SPECIAL_EFFECT"
CmdCommonNoticeDanmaku          "COMMON_NOTICE_DANMAKU"
CmdTradingScore                 "TRADING_SCORE"
CmdPreparing                    "PREPARING"
CmdGuardBuy                     "GUARD_BUY"
CmdGiftStarProcess              "GIFT_STAR_PROCESS"
CmdRoomSkinMsg                  "ROOM_SKIN_MSG"
CmdEnterEffect                  "ENTER_EFFECT"
CmdUserToastMsg                 "USER_TOAST_MSG"
CmdHeartBeatReply               "HEARTBEAT_REPLY"
CmdPopularityRedPocketNew       "POPULARITY_RED_POCKET_NEW"
CmdAreaRankChanged              "AREA_RANK_CHANGED"
CmdSuperChatEntrance            "SUPER_CHAT_ENTRANCE"
CmdPlayTogether                 "PLAY_TOGETHER"
CmdComboSend                    "COMBO_SEND"
CmdPopularityRedPocketStart     "POPULARITY_RED_POCKET_START"
```

由于我也不是很明白b站的命令, 所以这里只是列出了我知道的命令, 如果有人知道更多的命令, 请在issue中提出, 我会及时更新。

相信可以通过直译看懂这些命令...


### 消息处理函数

消息处理函数的原型为:

```go
func someFunc(event MsgEvent)
```

其中`MsgEvent`的定义为:

```go
type MsgEvent struct {
	//原始命令 
	Cmd    string
	//直播间号 
	RoomId int
    // 以下为不同的消息类型 
	DanMuMsg *DanMuMsg
	SuperChatMessage *SuperChatMessage
    // 下同, 可参考上方的消息类型, 取消Cmd即为结构体名称...
	...
}
```

### 消息结构体
这是弹幕消息的例子
```go
type DanMuMsg struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Sender struct {
			Uid    int64
			Name   string
			RoomId int64
		}
		Medal                    FansMedal
		Content                  string
		SendTimeStamp            int
		SendMillionTimeStamp     int64
		SenderEnterRoomTimeStamp int
	}
}
```
结构体内容, 请参考`go-bili-chat/Handler/template.go`中的定义
因为30个实在是太多了, 所以我就不一一列出来了...

## **在弹幕消息中, 我还没有完全明白每个参数的含义, 强烈欢迎各位PR完善**

### 性能

具体可以参考[Issue #24](https://github.com/FishZe/go-bili-chat/issues/24)

感谢[@XiaoMiku01](https://github.com/XiaoMiku01)提供的压测：

上海腾讯云 4C8G 10M Linux

同时连接5.8w个直播中的房间，内存占用 4.5G 左右, CPU 213% (平均 42%)


### 抽象画

```text
+============================+   +=================+
|        LiveRoom 1          |   |                 |
| goroutine heart beat       |   |                 → → goroutine your function
| goroutine receive msg → → ↓|   |                 |
| goroutine msg handler ← ← ←|   |                 → → goroutine your function
+================|↓|=========+   _                 |
                  → → → → → → → →    goroutine     → → goroutine your function
+============================+   _                 |
|        LiveRoom 1          |   |  msg distribute → → goroutine your function
| goroutine heart beat       |   |                 |
| goroutine receive msg → → ↓|   |                 → → goroutine your function
| goroutine msg handler ← ← ←|   |                 |
+================|↓|=========+   _                 → → goroutine your function
                  → → → → → → → →                  |
                                 _                 → → goroutine your function
                                 +======|↑|========+
                                         ↑
                                         ↑
                                +=======|↑|========+
                                |                  | 
                                |  main goroutine  |
                                |                  |
                                +==================+
```
