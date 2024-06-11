package configs

import (
	"github.com/dhanekom/cruddy/internal/storage"
)

type App struct {
	DB storage.DBRepo
	// Config          *koanf.Koanf
	//TemplatesFolder fs.FS
}
