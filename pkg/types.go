// Package pkg contains the shared data structures and functions
// which the application requires to run.
package pkg

import (
	"log"
	"net/http"

	"github.com/pranitbauva1997/splitwise-demo/pkg/store"
)

type Logger struct {
	err  *log.Logger
	info *log.Logger
}

type Application struct {
	Server        *http.Server
	StorageClient *store.PgStorage
	Log           Logger
}
