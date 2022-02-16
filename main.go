package main

import (
	"context"
	"log"
	"net"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/internal/model"
	"github.com/EZ4BRUCE/athena-server/internal/service"
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
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
		panic(err)
	}
	server.RegisterMap = make(map[string]*server.Agent, global.RPCSetting.MaxConn)
	svc := service.NewRuleService(context.Background())
	global.EmailSetting.To, err = svc.GetAllEmails()
	if err != nil {
		log.Fatalf("svc.GetAllEmails err: %v", err)
		panic(err)
	}
	setupLogger()

}

func main() {
	s := grpc.NewServer()
	pb.RegisterReportServerServer(s, server.NewReportServer())
	global.Logger.Infof("告警系统正在运行")
	lis, err := net.Listen("tcp", ":"+global.RPCSetting.Port)
	if err != nil {
		log.Fatalf("net.Listen error: %v", err)
	}
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("server.Serve error: %v", err)
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

func setupLogger() {
	global.Logger = logger.NewLogger(global.LOGSetting)
}
