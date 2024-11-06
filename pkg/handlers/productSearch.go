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
	if s := r.Form[keyOrderProductID]; len(s) > 0 {
		selectedProductId, _ = strconv.Atoi(s[0])
	}

	allProducts, err := appContext.Database(r).FindAllProducts(database.OrderProductByID, false)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	matchesQuery := func(p database.Product, q string) bool {
		if q == "" {
			return true
		}

		return strings.Contains(strings.ToLower(p.Name), strings.ToLower(q))
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

	for _, p := range products {
		appContext.ExecuteTemplate(w, r, "productSearchResult", components.ProductSearchResult{
			ProductSelectName: keyOrderProductID,
			Product:           p,
			IsSelected:        p.ID == selectedProductId,
		})
	}
}

func GetUnitOfMeasure(w http.ResponseWriter, r *http.Request) {
	selectedProductId := 0
	if s := r.Form[keyOrderProductID]; len(s) > 0 {
		selectedProductId, _ = strconv.Atoi(s[0])
	}

	product, err := appContext.Database(r).FindProduct(selectedProductId)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	w.Write([]byte(product.UnitOfMeasure.Symbol))
}
