package files

import (
	"bytes"
	"encoding/csv"
	"gestione-ordini/pkg/database"
	"strconv"
)

func ExportToList(orders []database.Order) []byte {
	builder := bytes.Buffer{}

	for s, ordersForSupplier := range groupOrdersBySupplier(orders) {
		builder.WriteString("Ordini per il fornitore \"" + s.Name + "\" (" + s.Email + "):\n")

		for _, ordersForProductId := range groupOrdersByProductId(ordersForSupplier) {
			product := ordersForProductId[0].Product
			totalAmount := 0

			for _, o := range ordersForProductId {
				totalAmount += o.Amount
			}

			builder.WriteString("  • ")
			builder.WriteString(product.Description)
			builder.WriteString(" (" + product.Code + ") - ")
			builder.WriteString(strconv.Itoa(totalAmount))
			builder.WriteString(" " + product.UnitOfMeasure.Symbol)
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
		"Richiesto il",
	})
	for _, ordersForSupplier := range groupOrdersBySupplier(orders) {
		for _, order := range ordersForSupplier {
			writer.Write([]string{
				order.Product.Description,
				order.Product.Code,
				strconv.Itoa(order.Amount),
				order.Product.UnitOfMeasure.Symbol,
				order.Product.Supplier.Name,
				order.ExpiresAt.Format("02/01/2006"),
				order.User.Username,
				order.CreatedAt.Format("02/01/2006"),
			})
		}
	}

	writer.Flush()

	return builder.Bytes()
}

func ExportToCSVCollapseProducts(orders []database.Order) []byte {
	builder := bytes.Buffer{}
	writer := csv.NewWriter(&builder)
	writer.Comma = ';'

	writer.Write([]string{
		"Descrizione prodotto",
		"Codice",
		"Quantità",
		"Unità di misura",
		"Fornitore",
	})
	for _, ordersForSupplier := range groupOrdersBySupplier(orders) {
		for _, ordersForProductId := range groupOrdersByProductId(ordersForSupplier) {
			product := ordersForProductId[0].Product
			totalAmount := 0

			for _, o := range ordersForProductId {
				totalAmount += o.Amount
			}

			writer.Write([]string{
				product.Description,
				product.Code,
				strconv.Itoa(totalAmount),
				product.UnitOfMeasure.Symbol,
				product.Supplier.Name,
			})
		}
	}

	writer.Flush()

	return builder.Bytes()
}

func groupOrdersBySupplier(orders []database.Order) map[database.Supplier][]database.Order {
	ordersBySupplier := map[database.Supplier][]database.Order{}

	for _, o := range orders {
		ordersBySupplier[o.Product.Supplier] = append(ordersBySupplier[o.Product.Supplier], o)
	}

	return ordersBySupplier
}

func groupOrdersByProductId(orders []database.Order) map[int][]database.Order {
	ordersByProductId := map[int][]database.Order{}

	for _, o := range orders {
		ordersByProductId[o.ProductID] = append(ordersByProductId[o.ProductID], o)
	}

	return ordersByProductId
}
