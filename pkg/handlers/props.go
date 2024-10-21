package handlers

const (
	keyUserID       = "user.id"
	keyUserRoleID   = "user.roleId"
	keyUserUsername = "user.username"
	keyUserPassword = "user.password"
	keyUserName     = "user.name"
	keyUserSurname  = "user.surname"

	keyOrderID          = "order.id"
	keyOrderProductID   = "order.productId"
	keyOrderUserID      = "order.userId"
	keyOrderAmount      = "order.amount"
	keyOrderRequestedAt = "order.requestedAt"

	keyProductID              = "product.id"
	keyProductProductTypeID   = "product.productTypeId"
	keyProductSupplierID      = "product.supplierId"
	keyProductUnitOfMeasureID = "product.unitOfMeasureId"
	keyProductName            = "product.name"

	destCook                 = "/cook"
	destAdminUsers           = "/admin/users"
	destAdminProducts        = "/admin/products"
	destAdminUsersTable      = "/admin/usersTable"
	destManagerAllOrders     = "/manager/allOrders"
	destManagerProducts      = "/manager/products"
	destManagerProductsTable = "/manager/productsTable"
)
