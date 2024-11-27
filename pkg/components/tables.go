package components

import (
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
)

type Table struct {
	TableURL        string
	OrderBy         int
	OrderDesc       bool
	MaxRowCount     int
	NextMaxRowCount int
	Headings        []TableHeading
	Rows            []TableRow
}

type TableHeading struct {
	Index int
	Name  string
}

type TableRow struct {
	EditURL string
	Cells   []string
}

type TableQuery struct {
	TableURL    string
	OrderBy     int
	OrderDesc   bool
	MaxRowCount int
}

func ParseTableQuery(r *http.Request, initialRowCount int) TableQuery {
	orderBy, err := strconv.Atoi(r.URL.Query().Get("orderBy"))
	if err != nil {
		orderBy = database.OrderProductByID
	}
	orderDesc := r.URL.Query().Get("orderDesc") == "true"

	maxRowCount, err := strconv.Atoi(r.URL.Query().Get("maxRowCount"))
	if err != nil {
		maxRowCount = initialRowCount
	}

	return TableQuery{
		TableURL:    r.URL.Path,
		OrderBy:     orderBy,
		OrderDesc:   orderDesc,
		MaxRowCount: maxRowCount,
	}
}

func ComposeTable[T any](
	query TableQuery,
	headings []TableHeading,
	rowFunc func(T) TableRow,
	rowData []T,
) Table {
	rows := make([]TableRow, len(rowData))
	for i, data := range rowData {
		rows[i] = rowFunc(data)
	}

	maxRowCount := min(query.MaxRowCount, len(rows))

	return Table{
		TableURL:        query.TableURL,
		OrderBy:         query.OrderBy,
		OrderDesc:       query.OrderDesc,
		MaxRowCount:     maxRowCount,
		NextMaxRowCount: maxRowCount * 2,
		Headings:        headings,
		Rows:            rows,
	}
}
