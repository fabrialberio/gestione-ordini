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

	// chef
	DestChef       = "/chef/"
	DestOrders     = "/chef/orders/"
	DestOrdersList = "/chef/ordersList"

	// admin and manager
	DestConsole        = "/console/"
	DestAllOrders      = "/console/allOrders/"
	DestProducts       = "/console/products/"
	DestProductsTable  = "/console/productsTable"
	DestSuppliers      = "/console/suppliers/"
	DestSuppliersTable = "/console/suppliersTable"

	// admin
	DestUsers      = "/console/users/"
	DestUsersTable = "/console/usersTable"
)
