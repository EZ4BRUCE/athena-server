package service

import (
	"fmt"
	"log"

	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/pkg/email"
	pb "github.com/EZ4BRUCE/athena-server/proto"
)

func (svc *RuleService) ExecuteRule(r *pb.ReportReq, a *Aggregator, result float64) {
	switch a.Rule.Action {
	case "EMAIL":
		sendWarningEmail(r, a, result)
	case "PHONE":
		doCall()
	case "MESSAGE":
		sendMessage()
	default:
		log.Printf("Rule：找不到 %s 告警动作\n", a.Rule.Action)
		return
	}
}

func sendWarningEmail(r *pb.ReportReq, a *Aggregator, result float64) {
	mailer := email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	})
	subject := fmt.Sprintf("[%s告警] %s 状态异常 timestamp: %v", a.Rule.Level, r.GetMetric(), r.GetTimestamp())
	body := fmt.Sprintf("异常主机Uid: %s\tTimestamp: %v\t告警等级: %s\n异常指标: %s\t指标维度: %v\t函数类型: %s\t描述: %s\t异常值: %v\n",
		r.GetUId(), r.GetTimestamp(), a.Rule.Level, r.GetMetric(), r.GetDimensions(), a.Function.Type, a.Function.Description, result)
	err := mailer.SendMail(global.EmailSetting.To, subject, body)
	if err != nil {
		log.Printf("邮件发送失败！mailer.SendMail err:%s", err)
	}
	log.Printf("邮件发送成功！")
}

func doCall() {
	log.Printf("已进行电话通知")
}

func sendMessage() {
	log.Printf("已进行短信通知")
}