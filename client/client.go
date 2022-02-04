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

func init() {
	flag.StringVar(&port, "p", "8888", "启动端口号")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
}

func Report(client pb.ReportServerClient, uid string, metric string) error {
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
	resp, err := client.Register(context.Background(), &pb.RegisterReq{
		Timestamp:   time.Now().Unix(),
		Metrics:     []string{"cpu_rate", "memory_used"},
		Description: "just a test",
	})
	if err != nil {
		return "", err
	}
	log.Printf("client.Report resp{code: %d, Uid:%s, message: %s}", resp.Code, resp.UId, resp.Msg)
	return resp.UId, nil
}

func main() {
	// conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
	conn, _ := grpc.Dial(":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	client := pb.NewReportServerClient(conn)
	uId, err := Register(client)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		_ = Report(client, uId, "cpu_rate")
		_ = Report(client, uId, "memory_used")
		time.Sleep(time.Second * 2)
	}

}
