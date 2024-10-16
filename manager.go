package main

import (
	"gestione-ordini/database"
	"log"
	"net/http"
	"strconv"
)

func manager(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, checkRole(r, database.RoleIDManager)
}

func managerProductsTable(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if err := checkPerm(r, database.PermIDViewProducts); err != nil {
		return nil, err
	}

	var err error
	var data struct {
		OrderBy   int
		OrderDesc bool
		Headers   interface{}
		Products  []database.Product
	}

	data.Headers = []struct {
		Index int
		Name  string
	}{}

	data.OrderBy, err = strconv.Atoi(r.URL.Query().Get("orderBy"))
	if err != nil {
		data.OrderBy = database.UserOrderByID
	}
	data.OrderDesc = r.URL.Query().Get("orderDesc") == "true"

	data.Products, err = db.GetProducts()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func managerProduct(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := checkPerm(r, database.PermIDEditProducts)
	if err != nil {
		return nil, err
	}

	var data struct {
		Product        database.Product
		ProductTypes   []database.ProductType
		Suppliers      []database.Supplier
		UnitsOfMeasure []database.UnitOfMeasure
		IsNew          bool
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		data.IsNew = true
		data.Product = database.Product{}
	} else {
		product, err := db.GetProduct(id)
		if err != nil {
			return nil, err
		}

		data.Product = product
	}

	data.ProductTypes, err = db.GetProductTypes()
	if err != nil {
		return nil, err
	}

	data.Suppliers, err = db.GetSuppliers()
	if err != nil {
		return nil, err
	}

	data.UnitsOfMeasure, err = db.GetUnitsOfMeasure()
	log.Println(data.UnitsOfMeasure)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func managerProductEdit(w http.ResponseWriter, r *http.Request) error {
	if err := checkPerm(r, database.PermIDEditProducts); err != nil {
		return err
	}

	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	productTypeId, _ := strconv.Atoi(r.FormValue("productTypeId"))
	supplierId, _ := strconv.Atoi(r.FormValue("supplierId"))
	unitOfMeasureId, _ := strconv.Atoi(r.FormValue("unitOfMeasureId"))
	name := r.FormValue("requestedAt")

	if isNew {
		err := db.CreateProduct(database.Product{
			ProductTypeID:   productTypeId,
			SupplierID:      supplierId,
			UnitOfMeasureID: unitOfMeasureId,
			Name:            name,
		})
		if err != nil {
			return err
		}
	} else {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		if delete {
			err = db.DeleteProduct(id)
			if err != nil {
				return err
			}
		} else {
			err = db.UpdateProduct(database.Product{
				ID:              id,
				ProductTypeID:   productTypeId,
				SupplierID:      supplierId,
				UnitOfMeasureID: unitOfMeasureId,
				Name:            name,
			})
			if err != nil {
				return err
			}
		}
	}

	http.Redirect(w, r, "/manager", http.StatusSeeOther)
	return nil
}
