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
	keyOrderSelectionExportMode = "orderSelection.exportMode"

	// common
	DestFirstLogin       = "/firstLogin"
	DestApi              = "/api/"
	DestProductSearch    = "/api/productSearch"
	DestOrderAmountInput = "/api/orderAmountInput"
	DestOwnOrdersView    = "/api/ownOrdersView"

	// chef
	DestChef       = "/chef/"
	DestChefOrders = "/chef/orders/"

	// admin and manager
	DestConsole             = "/console/"
	DestNewOrder            = "/console/newOrder/"
	DestOrders              = "/console/orders/"
	DestAllOrders           = "/console/allOrders/"
	DestAllOrdersView       = "/console/allOrdersView"
	DestOrderSelection      = "/console/orderSelection/"
	DestOrderSelectionCount = "/console/orderSelectionCount"
	DestProducts            = "/console/products/"
	DestProductsTable       = "/console/productsTable"
	DestProductsTableSearch = "/console/productsTableSearch"
	DestSuppliers           = "/console/suppliers/"
	DestSuppliersTable      = "/console/suppliersTable"

	// admin
	DestUsers         = "/console/users/"
	DestUsersTable    = "/console/usersTable"
	DestUpload        = "/console/upload/"
	DestUploadPreview = "/console/uploadPreview"
)
