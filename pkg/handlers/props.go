package handlers

const (
	keyUserID       = "id"
	keyUserRoleID   = "roleId"
	keyUserUsername = "username"
	keyUserPassword = "password"
	keyUserName     = "name"
	keyUserSurname  = "surname"

	keyOrderID          = "id"
	keyOrderProductID   = "productId"
	keyOrderUserID      = "userId"
	keyOrderAmount      = "amount"
	keyOrderRequestedAt = "requestedAt"

	keyProductID              = "id"
	keyProductProductTypeID   = "productTypeId"
	keyProductSupplierID      = "supplierId"
	keyProductUnitOfMeasureID = "unitOfMeasureId"
	keyProductName            = "name"

	destCook             = "/cook"
	destAdminUsers       = "/admin/users"
	destAdminProducts    = "/admin/products"
	destManagerAllOrders = "/manager/allOrders"
	destManagerProducts  = "/manager/products"
)
