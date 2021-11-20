package schemas

import (
	"database/sql"
	"endlessh-analyzer/importData/structs"
	"time"
)

type Data struct {
	ID        string `gorm:"primaryKey;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
	structs.ImportItem
}
