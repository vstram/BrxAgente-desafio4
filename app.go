package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// TestExcelReading tests reading an Excel file using the excelize library
func (a *App) TestExcelReading() (string, error) {
	// Using the ADMISSÃO ABRIL.xlsx file as an example
	filePath := filepath.Join("files", "ADMISSÃO ABRIL.xlsx")
	
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening Excel file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()

	// Get all sheets from the workbook
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return "", fmt.Errorf("no sheets found in the Excel file")
	}

	// Read the first sheet
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return "", fmt.Errorf("error reading rows from sheet %s: %w", sheets[0], err)
	}

	// Return information about the file
	return fmt.Sprintf("Successfully read Excel file with %d sheets. First sheet '%s' has %d rows.", 
		len(sheets), sheets[0], len(rows)), nil
}
