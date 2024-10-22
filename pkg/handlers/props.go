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

	keySupplierID = "supplier.id"

	destCook                 = "/cook"
	destAdminUsers           = "/admin/users"
	destAdminUsersTable      = "/admin/usersTable"
	destAdminProducts        = "/admin/products"
	destAdminProductsTable   = "/admin/productsTable"
	destAdminSuppliers       = "/admin/suppliers"
	destAdminSuppliersTable  = "/admin/suppliersTable"
	destManagerAllOrders     = "/manager/allOrders"
	destManagerProducts      = "/manager/products"
	destManagerProductsTable = "/manager/productsTable"
)
