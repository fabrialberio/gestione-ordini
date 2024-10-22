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

	DestCook                 = "/cook/"
	DestCookOrders           = "/cook/orders/"
	DestCookOrdersList       = "/cook/ordersList"
	DestAdmin                = "/admin/"
	DestAdminUsers           = "/admin/users/"
	DestAdminUsersTable      = "/admin/usersTable"
	DestAdminProducts        = "/admin/products/"
	DestAdminProductsTable   = "/admin/productsTable"
	DestAdminSuppliers       = "/admin/suppliers/"
	DestAdminSuppliersTable  = "/admin/suppliersTable"
	DestManager              = "/manager/"
	DestManagerAllOrders     = "/manager/allOrders/"
	DestManagerProducts      = "/manager/products/"
	DestManagerProductsTable = "/manager/productsTable"
)
