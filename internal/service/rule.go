package service

import (
	"fmt"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/pkg/email"
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
	mailer := email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	})
	subject := fmt.Sprintf("[%s告警] %s 状态异常 timestamp: %v", a.Rule.Level, r.GetMetric(), r.GetTimestamp())
	body := fmt.Sprintf("异常主机Uid: %s\t 主机描述 %s\t Timestamp: %v\t告警等级: %s\n异常指标: %s\t指标维度: %v\t函数类型: %s\t描述: %s\t异常值: %v\n",
		r.GetUId(), agentDescription, r.GetTimestamp(), a.Rule.Level, r.GetMetric(), r.GetDimensions(), a.Function.Type, a.Function.Description, result)
	err := mailer.SendMail(global.EmailSetting.To, subject, body)
	if err != nil {
		global.Logger.Errorf("[邮件错误] 邮件发送失败！mailer.SendMail err:%s", err)
	}
	global.Logger.Infof("<EMAIL操作> 邮件已发送")
}

func doCall() {
	global.Logger.Infof("已进行电话通知")
}

func sendMessage() {
	global.Logger.Infof("已进行短信通知")
}
