package main

import (
	"fmt"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

var DIM_COLUMN = []string{"isp", "loc"}
var DATA_COLUMN = []string{"numbers", "flux"}

func RestoreColName(df *dataframe.DataFrame) {
	for _, n := range df.Names() {
		if !strings.Contains(n, "_SUM") {
			continue
		}
		new_name := strings.Replace(n, "_SUM", "", 1)
		*df = df.Rename(new_name, n)
	}
}

func test1() {
	fmt.Println("isp loc 数据")
	df := dataframe.New(
		series.New([]string{"bj", "bj", "sh", "sh", "bj", "sh"}, series.String, "loc"),
		series.New([]string{"cmnet", "ct", "cmnet", "ct", "cmnet", "ct"}, series.String, "isp"),
		series.New([]int{6, 5, 4, 2, 3, 1}, series.Int, "numbers"),
		series.New([]int{1, 2, 3, 4, 5, 6}, series.Float, "flux"),
	)
	fmt.Println(df)

	grouped := df.GroupBy(DIM_COLUMN...)
	// for k, v := range grouped.GetGroups() {
	// 	fmt.Println(k)
	// 	fmt.Println(v)
	// }

	df = grouped.Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM, dataframe.Aggregation_SUM}, DATA_COLUMN)
	RestoreColName(&df)
	// df.SetNames(append(DIM_COLUMN, DATA_COLUMN...)...) // 顺序会变，不能这么用
	fmt.Println("按照isp loc合并")
	fmt.Println(df)

	AddNewData(df)
}

func AddNewData(df dataframe.DataFrame) {
	df_new := dataframe.New(
		series.New([]string{"bj", "sh"}, series.String, "loc"),
		series.New([]string{"cmnet", "ct"}, series.String, "isp"),
		series.New([]int{1, 7}, series.Int, "numbers"),
		series.New([]int{3, 2}, series.Float, "flux"),
	)
	fmt.Println("新数据")
	fmt.Println(df_new)

	// df = df.Concat(df_new)
	df = df.RBind(df_new)
	fmt.Println(df)

	grouped := df.GroupBy(DIM_COLUMN...)
	df = grouped.Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM, dataframe.Aggregation_SUM}, DATA_COLUMN)
	// df = df.Rename("numbers", "numbers_SUM")
	RestoreColName(&df)
	fmt.Println("合并新数据")
	fmt.Println(df)

	OneDim(df)
}

func OneDim(df dataframe.DataFrame) {
	grouped := df.GroupBy("isp")
	df_isp := grouped.Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM, dataframe.Aggregation_SUM}, DATA_COLUMN)
	fmt.Println("按照isp合并")
	fmt.Println(df_isp)

	grouped = df.GroupBy("loc")
	df_loc := grouped.Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM, dataframe.Aggregation_SUM}, DATA_COLUMN)
	fmt.Println("按照loc合并")
	fmt.Println(df_loc)
}

func main() {
	test1()
}
