package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var port string
var metric string

func init() {
	flag.StringVar(&port, "p", "8888", "启动端口号")
	// 还可以选memory_used
	flag.StringVar(&metric, "m", "cpu_rate", "启动端口号")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
}

func Report(client pb.ReportServerClient, uid string) error {
	resp, _ := client.Report(context.Background(), &pb.ReportReq{
		UId:        uid,
		Timestamp:  time.Now().Unix(),
		Metric:     metric,
		Dimensions: map[string]string{"computer": "Bruce's desktop"},
		Value:      rand.Float64(),
	})
	log.Printf("client.Report resp{code: %d, message: %s}", resp.Code, resp.Msg)
	return nil
}

func Register(client pb.ReportServerClient) (string, error) {
	resp, _ := client.Register(context.Background(), &pb.RegisterReq{
		Timestamp:   time.Now().Unix(),
		Description: "bruce",
	})
	log.Printf("client.Report resp{code: %d, Uid:%s, message: %s}", resp.Code, resp.UId, resp.Msg)
	return resp.UId, nil
}

func main() {
	// conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
	conn, _ := grpc.Dial(":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	client := pb.NewReportServerClient(conn)
	uId, _ := Register(client)
	for i := 0; i < 100; i++ {
		_ = Report(client, uId)
		time.Sleep(time.Second * 2)
	}

}
