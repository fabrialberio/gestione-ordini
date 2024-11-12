package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"time"
)

func GetNewOrder(w http.ResponseWriter, r *http.Request) {
	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		LogoutError(w, r, err)
		return
	}

	productTypes, err := appContext.Database(r).FindAllProductTypes()
	if err != nil {
		ShowDatabaseQueryError(w, r, err)
		return
	}

	data := struct {
		Sidebar        []components.SidebarDest
		ProductInput   components.ProductInput
		AmountInputURL string
		AmountInput    components.Input
		ExpiresAtInput components.Input
	}{
		Sidebar: currentSidebar(0, user.RoleID == database.RoleIDAdministrator),
		ProductInput: components.ProductInput{
			ProductSelectName: keyOrderProductID,
			ProductSearchURL:  DestProductSearch,
			SearchInputName:   keyProductSearchQuery,
			ProductTypesName:  keyProductSearchProductTypes,
			ProductTypes:      productTypes,
		},
		AmountInputURL: DestOrderAmountInput,
		ExpiresAtInput: components.Input{
			Label:        "Data di consegna richiesta",
			Name:         keyOrderRequestedAt,
			Type:         "date",
			DefaultValue: time.Now().Format(dateFormat),
		},
	}

	appContext.ExecuteTemplate(w, r, "newOrder.html", data)
}
