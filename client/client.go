package main

import (
	"flag"
	"log"
	"time"

	pb "github.com/EZ4BRUCE/athena-server/proto"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "启动端口号")
	flag.Parse()
}

func Report(client pb.ReportServerClient) error {
	resp, _ := client.Report(context.Background(), &pb.ReportReq{
		Timestamp: time.Now().Unix(), Metric: "cpu rate", Dimensions: map[string]string{"computer": "Bruce's desktop"}, Value: 0.45,
	})
	log.Printf("client.Report resp{code: %d, message: %s}", resp.Code, resp.Msg)
	return nil
}

func main() {
	// conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
	conn, _ := grpc.Dial(":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))

	defer conn.Close()

	client := pb.NewReportServerClient(conn)
	_ = Report(client)
}
