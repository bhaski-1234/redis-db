package constant

// Package constant provides constants used in the RESP protocol.
const (
	ErrInvalidRESP = "invalid RESP format"
)

const (
	DataFileExtension = ".rdb"
	Header            = "REDIS-DB"
	EOF               = 0xFF
	HeaderLength      = 8
	TypeString        = 0x00
	TypeInteger       = 0x01
	TypeList          = 0x02
	TypeTTL           = 0x03
)
