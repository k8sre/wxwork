

使用示例

```go
package main

import (
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/k8sre/wxwork"
)

var agentid int64

func main() {
	SendWorkWechat("xxx/zhangsan,lisi,tom", "", "", "hello world!", "")
}

// SendWorkWechat 发送微信企业应用消息
func SendWorkWechat(touser, toparty, totag, msg, logsign string) string {
	cropid := "xxxxxx"
	agentid = 1000000
	agentsecret := "xxxxxx"

	workwxapi := wxwork.Client{
		CropID:      cropid,
		AgentID:     agentid,
		AgentSecret: agentsecret,
	}

	slicetouser := strings.Split(touser, ",")
	slicetoparty := strings.Split(toparty, ",")
	slicetotag := strings.Split(totag, ",")

	workwxmsg := wxwork.Message{
		ToUser:   slicetouser,
		ToParty:  slicetoparty,
		ToTag:    slicetotag,
		MsgType:  "markdown",
		Markdown: wxwork.Content{Content: msg},
	}
	if err := workwxapi.Send(workwxmsg); err != nil {
		logs.Error(logsign, "[workwechat]", err.Error())
	}

	logs.Info(logsign, "[workwechat]", "workwechat send ok.")
	return "workwechat send ok"
}
```

