package main

import (
	"context"
	"log"
	"net"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/internal/model"
	"github.com/EZ4BRUCE/athena-server/internal/service"
	"github.com/EZ4BRUCE/athena-server/pkg/email"
	"github.com/EZ4BRUCE/athena-server/pkg/logger"
	"github.com/EZ4BRUCE/athena-server/pkg/setting"
	"github.com/EZ4BRUCE/athena-server/server"
	"google.golang.org/grpc"
)

// 项目配置初始化，仅在程序开始时执行一次
func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
		panic(err)
	}
	setupLogger()
	setupServer()
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
		panic(err)
	}
	err = setupEmail()
	if err != nil {
		log.Fatalf("init.setupEmail err: %v", err)
		panic(err)
	}
	err = setupSMTP()
	if err != nil {
		log.Fatalf("init.setupSMTP err: %v", err)
		panic(err)
	}

}

func main() {
	s := grpc.NewServer()
	pb.RegisterReportServerServer(s, server.NewReportServer())
	global.Logger.Infof("告警系统正在运行")
	lis, err := net.Listen("tcp", ":"+global.RPCSetting.Port)
	if err != nil {
		global.Logger.Fatalf("net.Listen error: %v", err)
	}
	err = s.Serve(lis)
	if err != nil {
		global.Logger.Fatalf("server.Serve error: %v", err)
	}

}

// 从configs中载入global配置
func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("RPC", &global.RPCSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("LOG", &global.LOGSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("ReportDB", &global.ReportDBSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("RuleDB", &global.RuleDBSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}

	return nil
}

// 根据global的设置初始化数据库
func setupDBEngine() error {
	var err error
	global.ReportDBEngine, err = model.NewReportDBEngine(global.ReportDBSetting)
	if err != nil {
		return err
	}
	global.RuleDBEngine, err = model.NewRuleDBEngine(global.RuleDBSetting)
	if err != nil {
		return err
	}
	return nil
}

// 初始化zap日志与lumberjack
func setupLogger() {
	global.Logger = logger.NewLogger(global.LOGSetting)
	if global.Logger == nil {
		panic("logger is not initialized")
	}
}

// 初始化 gRPC 服务需要的变量
func setupServer() {
	server.RegisterMap = make(map[string]*server.Agent, global.RPCSetting.MaxConn)
	if server.RegisterMap == nil {
		panic("server.RegisterMap is not initialized")
	}
}

// 初始化告警用户邮箱（只会载入一次）
func setupEmail() error {
	var err error
	svc := service.NewRuleService(context.Background())
	global.EmailSetting.To, err = svc.GetAllEmails()
	if err != nil {
		return err
	}
	if len(global.EmailSetting.To) == 0 {
		setting, err := setting.NewSetting()
		if err != nil {
			return err
		}
		err = setting.ReadSection("Email", &global.EmailSetting)
		if err != nil {
			return err
		}
	}
	return nil
}

// 初始化STMP服务邮箱（只会载入一次）
func setupSMTP() error {
	svc := service.NewRuleService(context.Background())
	smtps, err := svc.GetAllSmtps()
	if err != nil {
		return err
	}
	if len(smtps) == 0 {
		global.MailerPool = make(chan *email.Mailer, 1)
		log.Printf("未配置STMP！已载入默认STMP服务邮箱：%s \n", global.EmailSetting.From)
		global.PutMailer(&email.Mailer{
			Host:     global.EmailSetting.Host,
			Port:     global.EmailSetting.Port,
			IsSSL:    global.EmailSetting.IsSSL,
			UserName: global.EmailSetting.UserName,
			Password: global.EmailSetting.Password,
			From:     global.EmailSetting.From,
		})
	} else {
		global.MailerPool = make(chan *email.Mailer, len(smtps))
		for _, smtp := range smtps {
			log.Printf("已载入STMP服务邮箱：%s \n", smtp.From)
			global.PutMailer(&email.Mailer{
				Host:     smtp.Host,
				Port:     smtp.Port,
				IsSSL:    smtp.IsSSL,
				UserName: smtp.UserName,
				Password: smtp.Password,
				From:     smtp.From,
			})
		}
	}

	return nil
}
