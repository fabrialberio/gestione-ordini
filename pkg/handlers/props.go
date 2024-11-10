package handlers

const (
	keyUserID       = "user.id"
	keyUserRoleID   = "user.roleId"
	keyUserUsername = "user.username"
	keyUserPassword = "user.password"
	keyUserName     = "user.name"
	keyUserSurname  = "user.surname"

	keyOrderID                   = "order.id"
	keyOrderProductID            = "order.productId"
	keyOrderUserID               = "order.userId"
	keyOrderAmount               = "order.amount"
	keyOrderRequestedAt          = "order.requestedAt"
	keyProductSearchQuery        = "product.searchQuery"
	keyProductSearchProductTypes = "product.searchProductTypes"

	keyProductID              = "product.id"
	keyProductProductTypeID   = "product.productTypeId"
	keyProductSupplierID      = "product.supplierId"
	keyProductUnitOfMeasureID = "product.unitOfMeasureId"
	keyProductDescription     = "product.description"
	keyProductCode            = "product.code"

	keySupplierID    = "supplier.id"
	keySupplierEmail = "supplier.email"
	keySupplierName  = "supplier.name"

	keyOrderSelectionStart      = "orderSelection.start"
	keyOrderSelectionEnd        = "orderSelection.end"
	keyOrderSelectionSupplierID = "orderSelection.supplier"

	// chef
	DestChef             = "/chef/"
	DestChefOrders       = "/chef/orders/"
	DestChefOrdersView   = "/chef/ordersView"
	DestProductSearch    = "/chef/productSearch"
	DestOrderAmountInput = "/chef/orderAmountInput"

	// admin and manager
	DestConsole             = "/console/"
	DestAllOrders           = "/console/allOrders/"
	DestAllOrdersView       = "/console/allOrdersView"
	DestOrderSelection      = "/console/orderSelection/"
	DestOrderSelectionCount = "/console/orderSelectionCount"
	DestProducts            = "/console/products/"
	DestProductsTable       = "/console/productsTable"
	DestSuppliers           = "/console/suppliers/"
	DestSuppliersTable      = "/console/suppliersTable"

	// admin
	DestUsers      = "/console/users/"
	DestUsersTable = "/console/usersTable"
	DestUpload     = "/console/upload/"
)
