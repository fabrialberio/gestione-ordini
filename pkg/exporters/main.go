package exporters

import (
	"bytes"
	"encoding/csv"
	"gestione-ordini/pkg/database"
	"strconv"
)

func ExportToList(orders []database.Order) []byte {
	builder := bytes.Buffer{}

	for s, supplierOrders := range splitOrdersBySupplier(orders) {
		builder.WriteString("Ordini per il fornitore \"" + s.Name + "\" (" + s.Email + "):\n")

		for _, o := range supplierOrders {
			builder.WriteString("  • ")
			builder.WriteString(o.Product.Description)
			builder.WriteString(" (" + o.Product.Code + ") - ")
			builder.WriteString(o.AmountString)
			builder.WriteString("\n")
		}

		builder.WriteString("\n")
	}

	return builder.Bytes()
}

func ExportToCSV(orders []database.Order) []byte {
	builder := bytes.Buffer{}
	writer := csv.NewWriter(&builder)
	writer.Comma = ';'

	writer.Write([]string{
		"Descrizione prodotto",
		"Codice",
		"Quantità",
		"Unità di misura",
		"Fornitore",
		"Data di consegna richiesta",
		"Richiesto da",
	})
	for _, supplierOrders := range splitOrdersBySupplier(orders) {
		for _, order := range supplierOrders {
			writer.Write([]string{
				order.Product.Description,
				order.Product.Code,
				strconv.Itoa(order.Amount),
				order.Product.UnitOfMeasure.Symbol,
				order.Product.Supplier.Name,
				order.ExpiresAt.Format("02/01/2006"),
				order.User.Username,
			})
		}
	}

	writer.Flush()

	return builder.Bytes()
}

func splitOrdersBySupplier(orders []database.Order) map[database.Supplier][]database.Order {
	ordersBySupplier := map[database.Supplier][]database.Order{}

	for _, o := range orders {
		ordersBySupplier[o.Product.Supplier] = append(ordersBySupplier[o.Product.Supplier], o)
	}

	return ordersBySupplier
}
