package server

import (
	"fmt"
	"log"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/pkg/email"
)

// 发送新增主机提醒邮件
func sendLoginEmail(r *pb.RegisterReq, uId string) {
	mailer := email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	})
	subject := fmt.Sprintf("[新增主机] 收到新主机注册 Timestamp %v", r.GetTimestamp())
	body := fmt.Sprintf("新主机注册，分配其UID为: %s\t\n主机详情：注册时间戳：%v\t 监控指标：%v\t 描述：%v", uId, r.GetTimestamp(), r.GetMetrics(), r.GetDescription())
	err := mailer.SendMail(global.EmailSetting.To, subject, body)
	if err != nil {
		log.Printf("[邮件错误] 邮件发送失败！mailer.SendMail err:%s", err)
	}
	log.Printf("<新增主机> 邮件通知已发送")
}

// 发送链接异常主机告警邮件
func sendOfflineEmail(agent *Agent) {
	mailer := email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	})
	subject := fmt.Sprintf("[主机异常] %s 主机连接异常", agent.UId)
	body := fmt.Sprintf("异常主机Uid: %s\t，已%d未收到其上报信息", agent.UId, agent.CheckAliveTime)
	err := mailer.SendMail(global.EmailSetting.To, subject, body)
	if err != nil {
		log.Printf("[邮件错误] 邮件发送失败！mailer.SendMail err:%s", err)
	}
	log.Printf("<主机异常> 邮件已发送")
}
