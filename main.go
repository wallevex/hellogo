package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
	"time"
)

type Event struct {
	FlowState  int
	MAC        string
	CreateTime time.Time
}

func main() {
	// 创建 Line 折线图
	line := charts.NewLine()
	x := true
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "1000 点时序图"}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
			Name: "时间",
			AxisLabel: &opts.AxisLabel{
				Interval: "10", // 每 10 个点显示一个 X 轴标签
				Rotate:   45,
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "值",
			Min:  0,
			Max:  10,
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: &x}),
	)

	// 模拟 1000 个时间点的数据（每秒一个）
	now := time.Now()
	xAxis := make([]string, 0)
	yAxis := make([]opts.LineData, 0)

	for i := 0; i < 1000; i++ {
		t := now.Add(time.Duration(i) * time.Second).Format("15:04:05")
		xAxis = append(xAxis, t)
		// 示例数据：值在 0-10 范围内波动
		value := float64((i * 7) % 11) // 模拟周期波动
		yAxis = append(yAxis, opts.LineData{Value: value})
	}

	// 添加数据
	line.SetXAxis(xAxis).
		AddSeries("value", yAxis).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: &x}),
		)

	// 渲染图表
	f, _ := os.Create("line_1000points.html")
	_ = line.Render(f)
}
