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
	productTypeIds := r.Form[keyProductSearchProductTypes]

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

	allProducts, err := appContext.Database(r).FindAllProducts(database.OrderProductByDescription, false, -1)
	if err != nil {
		logError(r, err)
	}

	matchesQuery := func(p database.Product, q string) bool {
		if q == "" {
			return true
		}

		return strings.Contains(strings.ToLower(p.Description), strings.ToLower(q))
	}

	matchesProductTypeIds := func(p database.Product, ids []string) bool {
		if len(ids) == 0 {
			return true
		}

		return slices.Contains(ids, strconv.Itoa(p.ProductTypeID))
	}

	products := make([]database.Product, 0)
	for _, p := range allProducts {
		if matchesQuery(p, query) && matchesProductTypeIds(p, productTypeIds) {
			products = append(products, p)
		}
	}

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
