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

	keySupplierID    = "supplier.id"
	keySupplierEmail = "supplier.email"
	keySupplierName  = "supplier.name"

	// chef
	DestChef           = "/chef/"
	DestChefOrders     = "/chef/orders/"
	DestChefOrdersList = "/chef/ordersList"

	// admin and manager
	DestConsole        = "/console/"
	DestAllOrders      = "/console/allOrders/"
	DestAllOrdersTable = "/console/allOrdersTable"
	DestProducts       = "/console/products/"
	DestProductsTable  = "/console/productsTable"
	DestSuppliers      = "/console/suppliers/"
	DestSuppliersTable = "/console/suppliersTable"

	// admin
	DestUsers      = "/console/users/"
	DestUsersTable = "/console/usersTable"
)
