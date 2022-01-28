package model

type Report struct {
	Timestamp  int64             `bson:"timestamp"`
	Metric     string            `bson:"metric"`
	Dimensions map[string]string `bson:"dimensions"`
	Value      float64           `bson:"value"`
}
