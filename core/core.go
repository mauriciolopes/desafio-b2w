package core

import (
	"desafio-b2w/db"
	"log"
	"os"
	"os/signal"

	"github.com/paked/configure"
)

// Facade represents exported API.
type Facade struct {
	conf     *configure.Configure
	dbURI    string
	dbName   string
	env      map[string]*string
	shutdown func()
	db       *db.DB
}

const (
	defaultDBURI  = "mongodb://localhost:27017"
	defaultDBName = "starwars"
)

// Init initializes the facade API.
func Init() *Facade {
	f := Facade{}
	f.env = make(map[string]*string)
	f.conf = configure.New(configure.NewEnvironment(), configure.NewFlag())
	f.conf.StringVar(&f.dbURI, "db-uri", defaultDBURI, "URI de conexão com o banco de dados")
	f.conf.StringVar(&f.dbName, "db-name", defaultDBName, "Nome do banco de dados")
	return &f
}

// Configure sets a configuration variable.
func (f *Facade) Configure(name, defValue, description string) {
	f.env[name] = f.conf.String(name, defValue, description)
}

// Get gets a value from name.
func (f *Facade) Get(name string) string {
	if val, ok := f.env[name]; ok {
		return *val
	}
	return ""
}

// RegisterShutdown registers shutdown function.
func (f *Facade) RegisterShutdown(shutdownFunc func()) {
	f.shutdown = shutdownFunc
}

// Run executes the application.
func (f *Facade) Run(app func(*Facade)) {
	f.conf.Parse()
	f.db = db.New()
	if err := f.db.Open(f.dbURI, f.dbName); err != nil {
		log.Fatal(err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go app(f)
	<-stop
	if f.shutdown != nil {
		f.shutdown()
	}
	if err := f.db.Close(); err != nil {
		log.Fatal("Não foi possível fechar a conexão com o banco de dados:", err)
	}
}

// Planets returns the planet service.
func (f *Facade) Planets() PlanetService {
	return PlanetService{f}
}
