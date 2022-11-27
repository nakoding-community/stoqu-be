package constant

// db
const DB_DEFAULT_SYSTEM = "system"

// code
const (
	CODE_LENGTH = 5

	CODE_ROLE_PREFIX            = "R-"
	CODE_UNIT_PREFIX            = "U-"
	CODE_PACKET_PREFIX          = "PKT-"
	CODE_CONVERTION_UNIT_PREFIX = "CVU-"
	CODE_CURRENCY_PREFIX        = "CUR-"
	CODE_REMINDER_STOCK_PREFIX  = "RMS-"
	CODE_RACK_PREFIX            = "RCK-"

	CODE_BRAND_PREFIX        = "B-"
	CODE_VARIANT_PREFIX      = "V-"
	CODE_PRODUCT_PREFIX      = "P-"
	CODE_STOCK_LOOKUP_PREFIX = "STL-"

	CODE_ORDER_TRX_PREFIX = "OTX-"
	CODE_STOCK_TRX_PREFIX = "STX-"
)

// trx
const (
	TRX_TYPE_IN      = "in"
	TRX_TYPE_OUT     = "out"
	TRX_TYPE_CONVERT = "convert"
)

// status
const (
	STATUS_SUCCESS = "SUCCESS"
	STATUS_FAILED  = "FAILED"
)

// convert
const (
	CONVERT_TYPE_ORIGIN      = "origin"
	CONVERT_TYPE_DESTINATION = "destination"
)

// action
const (
	ACTION_INSERT = "insert"
	ACTION_UPDATE = "update"
	ACTION_DELETE = "delete"
)

// reminder stock
const (
	REMINDER_STOCK_SECOND = "second"
	REMINDER_STOCK_DAILY  = "daily"
)

// firebase
const (
	FIRESTORE_MAX_DATA                   = 50
	FIRESTORE_COLLECTION_DASHBOARD_ORDER = "dashboard-order"
	FIRESTORE_COLLECTION_TOTAL_ORDER     = "total-order"
)

// report order product
const (
	GROUP_BY_BRAND   = "brand"
	GROUP_BY_VARIANT = "variant"
	GROUP_BY_PACKET  = "packet"
)
