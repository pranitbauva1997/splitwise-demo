package pkg

import (
	"crypto/tls"
	"github.com/pranitbauva1997/splitwise-demo/pkg/store"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func initLoggers() Logger {
	return Logger{
		err:  log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
		info: log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
	}
}

func Init() (*Application, error) {
	const PORT int = 8000
	addr := ":" + strconv.Itoa(PORT)

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	}

	server := &http.Server{
		Addr:           addr,
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 10,
		IdleTimeout:    time.Second * 120,
		MaxHeaderBytes: 524288,
		TLSConfig:      tlsConfig,
	}

	pgConfig := store.PGConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "",
		DBName:   "loan_module_test",
	}

	app := &Application{
		Server: server,
		Log:    initLoggers(),
	}

	var err error
	app.StorageClient, err = store.Init(pgConfig)
	if err != nil {
		app.Log.err.Println("couldn't initialize the db object:", err)
		return app, err
	}

	server.Handler = initRoutes(app)
	app.Log.info.Println("Initialized all the routes")
	return app, nil
}
