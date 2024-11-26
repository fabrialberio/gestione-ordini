package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

func PostProductSearch(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue(keyProductSearchQuery)

	productTypeIds := make([]int, 0)
	for _, s := range r.Form[keyProductSearchProductTypes] {
		id, err := strconv.Atoi(s)
		if err == nil {
			productTypeIds = append(productTypeIds, id)
		}
	}

	selectedProductId := 0
	orderId, err := strconv.Atoi(r.FormValue(keyOrderID))
	if err != nil {
		if s := r.Form[keyOrderProductID]; len(s) > 0 {
			selectedProductId, _ = strconv.Atoi(s[0])
		}
	} else {
		order, err := appContext.Database(r).FindOrder(orderId)
		if err != nil {
			logError(r, err)
		}

		selectedProductId = order.ProductID
	}

	products, err := appContext.Database(r).FindAllProducts(database.OrderProductByDescription, false, -1)
	if err != nil {
		logError(r, err)
	}

	products = filterProducts(products, query, productTypeIds)

	if len(products) > 0 {
		for _, p := range products {
			appContext.ExecuteTemplate(w, r, "productSearchResult", components.ProductSearchResult{
				ProductSelectName: keyOrderProductID,
				Product:           p,
				IsSelected:        p.ID == selectedProductId,
			})
		}
	} else {
		appContext.ExecuteTemplate(w, r, "productSearchNoResults", nil)
	}
}

func filterProducts(allProducts []database.Product, query string, productTypeIds []int) []database.Product {
	matchesQuery := func(p database.Product, q string) bool {
		if q == "" {
			return true
		}

		return strings.Contains(strings.ToLower(p.Description), strings.ToLower(q))
	}

	matchesProductTypeIds := func(p database.Product, ids []int) bool {
		if len(ids) == 0 {
			return true
		}

		return slices.Contains(ids, p.ProductTypeID)
	}

	products := make([]database.Product, 0)
	for _, p := range allProducts {
		if matchesQuery(p, query) && matchesProductTypeIds(p, productTypeIds) {
			products = append(products, p)
		}
	}

	return products
}
