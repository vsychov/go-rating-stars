package storage

// KeyValue represent key-value storage in ORM
type KeyValue struct {
	ID    string  `gorm:"primaryKey,type:varchar(64)"`
	Value float64 `gorm:"type:decimal(40,15);"`
	TTL   int64   `gorm:"default:0"`
}
