package constant

// db
const DB_DEFAULT_SYSTEM = "system"

// code
const (
	LENGTH_CODE = 5
)

// master
const (
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
