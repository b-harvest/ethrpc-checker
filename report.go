package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/xuri/excelize/v2"
)

// ReportResults prints or saves the RPC results based on the verbosity flag and output format
func ReportResults(results []*RpcResult, verbose bool, outputExcel bool) {
	if outputExcel {
		f := excelize.NewFile()
		name := fmt.Sprintf("geth%s", GethVersion)
		if err := f.SetSheetName("Sheet1", name); err != nil {
			log.Fatalf("Failed to set sheet name: %v", err)
		}

		// set header
		header := []string{"Method", "Status", "Value", "Warnings", "ErrMsg"}
		for col, h := range header {
			cell := fmt.Sprintf("%s1", string(rune('A'+col)))
			if err := f.SetCellValue(name, cell, h); err != nil {
				log.Fatalf("Failed to set cell value: %v", err)
			}
		}

		// set columns width
		if err := f.SetColWidth(name, "A", "A", 30); err != nil {
			log.Fatalf("Failed to set col width: %v", err)
		}
		if err := f.SetColWidth(name, "C", "C", 40); err != nil {
			log.Fatalf("Failed to set col width: %v", err)
		}
		if err := f.SetColWidth(name, "E", "E", 40); err != nil {
			log.Fatalf("Failed to set col width: %v", err)
		}

		// set style for method column
		methodColStyle, err := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{Vertical: "center"},
		})
		if err != nil {
			log.Fatalf("Failed to create style: %v", err)
		}
		if err = f.SetColStyle(name, "A", methodColStyle); err != nil {
			log.Fatalf("Failed to set col style: %v", err)
		}

		// set style for value column
		valueColStyle, err := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				WrapText:   false,
				Horizontal: "left",
			},
		})
		if err != nil {
			log.Fatalf("Failed to create style: %v", err)
		}
		if err = f.SetColStyle(name, "C", valueColStyle); err != nil {
			log.Fatalf("Failed to set col style: %v", err)
		}

		fontStyle := &excelize.Style{Font: &excelize.Font{Bold: true}}
		for i, result := range results {
			row := i + 2
			warnings, _ := json.Marshal(result.Warnings)
			methodCell := fmt.Sprintf("A%d", row)
			if err = f.SetCellValue(name, methodCell, result.Method); err != nil {
				log.Fatalf("Failed to set cell value: %v", err)
			}
			statusCell := fmt.Sprintf("B%d", row)
			if err = f.SetCellValue(name, statusCell, result.Status); err != nil {
				log.Fatalf("Failed to set cell value: %v", err)
			}
			valueCell := fmt.Sprintf("C%d", row)
			if err = f.SetCellValue(name, valueCell, result.Value); err != nil {
				log.Fatalf("Failed to set cell value: %v", err)
			}
			warningsCell := fmt.Sprintf("D%d", row)
			if err = f.SetCellValue(name, warningsCell, string(warnings)); err != nil {
				log.Fatalf("Failed to set cell value: %v", err)
			}
			errCell := fmt.Sprintf("E%d", row)
			if err = f.SetCellValue(name, errCell, result.ErrMsg); err != nil {
				log.Fatalf("Failed to set cell value: %v", err)
			}

			// SET STYLES
			// set status column style based on status
			switch result.Status {
			case Ok:
				fontStyle.Font.Color = GREEN
				s, err := f.NewStyle(fontStyle)
				if err != nil {
					log.Fatalf("Failed to create style: %v", err)
				}
				if err = f.SetCellStyle(name, statusCell, statusCell, s); err != nil {
					log.Fatalf("Failed to set cell style: %v", err)
				}
			case Warning:
				fontStyle.Font.Color = YELLOW
				s, err := f.NewStyle(fontStyle)
				if err != nil {
					log.Fatalf("Failed to create style: %v", err)
				}
				if err = f.SetCellStyle(name, statusCell, statusCell, s); err != nil {
					log.Fatalf("Failed to set cell style: %v", err)
				}
			case Error:
				fontStyle.Font.Color = RED
				s, err := f.NewStyle(fontStyle)
				if err != nil {
					log.Fatalf("Failed to create style: %v", err)
				}
				if err = f.SetCellStyle(name, statusCell, statusCell, s); err != nil {
					log.Fatalf("Failed to set cell style: %v", err)
				}
			}

			if err = f.SetRowHeight(name, row, 20); err != nil {
				log.Fatalf("Failed to set row height: %v", err)
			}
		}
		// Set header style at last to avoid override by other styles
		headerStyle, err := f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#D3D3D3"}},
			Font: &excelize.Font{Bold: true},
		})
		if err != nil {
			log.Fatalf("Failed to create style: %v", err)
		}
		if err = f.SetRowStyle(name, 1, 1, headerStyle); err != nil {
			log.Fatalf("Failed to set row style: %v", err)
		}

		fileName := fmt.Sprintf("rpc_results_%s.xlsx", time.Now().Format("15:04:05"))
		if err := f.SaveAs(fileName); err != nil {
			log.Fatalf("Failed to save Excel file: %v", err)
		}
		fmt.Println("Results saved to " + fileName)
	}

	fmt.Println(`
██████╗ ██████╗  ██████╗    ██████╗ ███████╗███████╗██╗   ██╗██╗  ████████╗███████╗
██╔══██╗██╔══██╗██╔════╝    ██╔══██╗██╔════╝██╔════╝██║   ██║██║  ╚══██╔══╝██╔════╝
██████╔╝██████╔╝██║         ██████╔╝█████╗  ███████╗██║   ██║██║     ██║   ███████╗
██╔══██╗██╔═══╝ ██║         ██╔══██╗██╔══╝  ╚════██║██║   ██║██║     ██║   ╚════██║
██║  ██║██║     ╚██████╗    ██║  ██║███████╗███████║╚██████╔╝███████╗██║   ███████║
╚═╝  ╚═╝╚═╝      ╚═════╝    ╚═╝  ╚═╝╚══════╝╚══════╝ ╚═════╝ ╚══════╝╚═╝   ╚══════╝
------------------------------------------------------------------------------------
                                                                                   `)
	for _, result := range results {
		ColorPrint(result, verbose)
	}
}

func ColorPrint(result *RpcResult, verbose bool) {
	method := result.Method
	status := result.Status
	switch status {
	case Ok:
		value := result.Value
		if !verbose {
			value = ""
		}
		color.Green("%-40s: %s (value: %v)", method, status, value)
	case Warning:
		color.Yellow("%-40s: %s (%v)", method, status, result.Warnings)
	case Error:
		color.Red("%-40s: %s (%v)", method, status, result.ErrMsg)
	}
}
