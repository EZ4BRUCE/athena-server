package service

import (
	"fmt"
	"time"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/global"
)

// service层方法，接收请求结构体或特定参数执行特定操作

// 执行告警规则函数
func (svc *RuleService) ExecuteRule(r *pb.ReportReq, a *Aggregator, result float64, agentDescription string) {
	switch a.Rule.Action {
	case global.EMAILAction:
		sendWarningEmail(r, a, result, agentDescription)
	case global.PHONEAction:
		doCall()
	case global.MESSAGEAction:
		sendMessage()
	default:
		global.Logger.Errorf("[参数错误] Rule：找不到 %s 告警动作\n", a.Rule.Action)
		return
	}
}

// 定义告警行为的逻辑实现
func sendWarningEmail(r *pb.ReportReq, a *Aggregator, result float64, agentDescription string) {
	mailer := global.GetMailer()
	subject := fmt.Sprintf("[%s告警] %s 状态异常 timestamp: %v", a.Rule.Level, r.GetMetric(), r.GetTimestamp())
	body := fmt.Sprintf("异常主机Uid: %s\t 主机描述 %s\t Timestamp: %v\t告警等级: %s\n异常指标: %s\t指标维度: %v\t函数类型: %s\t描述: %s\t异常值: %v\n",
		r.GetUId(), agentDescription, r.GetTimestamp(), a.Rule.Level, r.GetMetric(), r.GetDimensions(), a.Function.Type, a.Function.Description, result)
	for i := 0; i < 3; i++ {
		err := mailer.SendMail(global.EmailSetting.To, subject, body)
		// 发送一次sleep一秒，防止一秒内大量发送被禁止
		time.Sleep(time.Second)
		if err != nil {
			global.Logger.Errorf("[邮件错误] 第 %d 次发送邮件失败！mailer.SendMail err:%s", i+1, err)
		} else {
			global.Logger.Infof("<EMAIL操作成功> 邮件已发送")
			break
		}
	}
	global.PutMailer(mailer)
}

func doCall() {
	global.Logger.Infof("已进行电话通知")
}

func sendMessage() {
	global.Logger.Infof("已进行短信通知")
}
