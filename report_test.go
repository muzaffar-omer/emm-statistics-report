package main

import (
	"fmt"
	"github.com/kniren/gota/dataframe"
	"testing"
)

func TestReport_ExtractData(t *testing.T) {
	df := dataframe.LoadRecords([][]string{
		[]string{"TIME", "INPUT FILES", "INPUT CDRS", "INPUT BYTES", "OUTPUT FILES", "OUTPUT CDRS", "OUTPUT BYTES"},
		[]string{"20190325", "67682", "1474082862", "772594852110", "195133", "2130813264", "2206548945208"},
	})

	fmt.Println(df.String())
}
