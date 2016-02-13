package main

import (
	"container/list"
	"log"
	"reflect"

	"github.com/tealeg/xlsx"
)

// NAXIMUM ALLOWED ROWS BY EXCEL.
const ExcelMaxRows = 1048576

func generateReport(fname string, connections *list.List) {
	doc := xlsx.NewFile()
	report, _ := doc.AddSheet("Report")
	generateHeaders(connections.Front().Value, report.AddRow())
	for e := connections.Front(); e != nil; e = e.Next() {
		generateRow(e, report.AddRow())
		if report.MaxRow >= ExcelMaxRows {
			log.Println("Too many entries for report, stopping at ", report.MaxRow)
			break
		}
	}

	doc.Save(fname)
}

func generateHeaders(i interface{}, row *xlsx.Row) {
	fields := reflect.ValueOf(i)
	for i := 0; i < fields.NumField(); i++ {
		cell := row.AddCell()
		cell.SetString(fields.Type().Field(i).Name)
	}
}

func generateRow(i interface{}, row *xlsx.Row) {
	fields := reflect.ValueOf(i)
	for i := 0; i < fields.NumField(); i++ {
		cell := row.AddCell()
		cell.SetValue(fields.Interface())
	}
}
