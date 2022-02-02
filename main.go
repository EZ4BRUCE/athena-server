package main

import (
	"log"
	"net"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/internal/model"
	"github.com/EZ4BRUCE/athena-server/pkg/setting"
	"github.com/EZ4BRUCE/athena-server/server"
	"google.golang.org/grpc"
)

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
	global.RegisterMap = make(map[string]struct{}, 100)
	global.ReportMap = make(map[string]chan *pb.ReportReq, global.RPCSetting.AggregationTime*2)
}

func main() {
	s := grpc.NewServer()
	pb.RegisterReportServerServer(s, server.NewReportServer())
	lis, err := net.Listen("tcp", ":"+global.RPCSetting.Port)
	if err != nil {
		log.Fatalf("net.Listen error: %v", err)
	}
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("server.Serve error: %v", err)
	}

}

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("RPC", &global.RPCSetting)
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
