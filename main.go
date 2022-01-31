package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/EZ4BRUCE/athena-server/global"
	"github.com/EZ4BRUCE/athena-server/internal/model"
	"github.com/EZ4BRUCE/athena-server/internal/service"
	"github.com/EZ4BRUCE/athena-server/pkg/setting"
	pb "github.com/EZ4BRUCE/athena-server/proto"
	"github.com/EZ4BRUCE/athena-server/server"
	"google.golang.org/grpc"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
}

func main() {

	svc := service.NewRuleService(context.Background())
	results, _ := svc.SearchAggregators("cpu_rate")
	for _, result := range results {
		fmt.Println(result)
	}
	s := grpc.NewServer()
	pb.RegisterReportServerServer(s, server.NewReportServer())

	lis, err := net.Listen("tcp", ":"+global.RPCSetting.Port)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("server.Serve: %v", err)
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
