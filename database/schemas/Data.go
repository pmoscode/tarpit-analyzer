package schemas

import (
	"endlessh-analyzer/importData/structs"
	"gorm.io/gorm"
)

type Data struct {
	gorm.Model
	structs.ImportItem
}
