package server

import (
	"fmt"
	"time"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/global"
)

// 发送新增主机提醒邮件
func sendLoginEmail(r *pb.RegisterReq, uId string) {
	mailer := global.GetMailer()
	subject := fmt.Sprintf("[新增主机] 收到新主机注册 Timestamp %v", r.GetTimestamp())
	body := fmt.Sprintf("新主机注册，分配其UID为: %s\t\n主机详情：注册时间戳：%v\t 监控指标：%v\t 描述：%v", uId, r.GetTimestamp(), r.GetMetrics(), r.GetDescription())
	// 邮件尝试发送3次，不行就不发了
	for i := 0; i < 3; i++ {
		err := mailer.SendMail(global.EmailSetting.To, subject, body)
		// 发送一次sleep一秒，防止一秒内大量发送被禁止
		time.Sleep(time.Second)
		if err != nil {
			global.Logger.Errorf("[邮件错误] 第 %d 次发送邮件失败！mailer.SendMail err:%s", i+1, err)
		} else {
			global.Logger.Infof("<EMAIL操作成功> 新注册 Agent:%s 邮件已发送", uId)
			break
		}

	}
	global.PutMailer(mailer)

}

// 发送链接异常主机告警邮件
func sendOfflineEmail(agent *Agent) {
	mailer := global.GetMailer()
	subject := fmt.Sprintf("[主机异常] %s 主机连接异常", agent.UId)
	body := fmt.Sprintf("异常主机Uid: %s\t，已%d秒未收到其上报信息", agent.UId, agent.CheckAliveTime*2)
	for i := 0; i < 3; i++ {
		err := mailer.SendMail(global.EmailSetting.To, subject, body)
		// 发送一次sleep一秒，防止一秒内大量发送被禁止
		time.Sleep(time.Second)
		if err != nil {
			global.Logger.Errorf("[邮件错误] 第 %d 次发送邮件失败！mailer.SendMail err:%s", i+1, err)
		} else {
			global.Logger.Infof("<EMAIL操作成功> Agent: %s 连接异常通知邮件已发送", agent.UId)
			break
		}
	}
	global.PutMailer(mailer)
}
