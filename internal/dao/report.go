package dao

import "github.com/EZ4BRUCE/athena-server/internal/model"

// dao的方法，dao作为接收者
func (d ReportDao) CreateReport(uId string, timestamp int64, metric string, dimensions map[string]string, value float64) error {
	Report := model.Report{UId: uId, Timestamp: timestamp, Metric: metric, Dimensions: dimensions, Value: value}
	return Report.Create(d.engine)
}
