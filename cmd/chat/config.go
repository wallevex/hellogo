package main

type Config struct {
	Listen string `json:"listen"`
	Debug  string `json:"debug"`
	Logger Logger `json:"logger"`
}

type Logger struct {
	Dir      string `json:"dir"`
	File     string `json:"file"`
	Count    int32  `json:"count"`
	Size     int64  `json:"size"`
	Unit     string `json:"unit"`
	Level    string `json:"level"`
	Compress int64  `json:"compress"`
}
