package dataModel

type OpenClose struct {
	Status     string
	Date       string `json:"from"`
	Volume     int
	AfterHours float32
	Symbol     string
	Open       float32
	High       float32
	Low        float32
	Close      float32
	PreMarket  float32
}
