package core

import (
	"os"
	"testing"
	"time"
)

const defaultTestDBName = "starwars_test"

func TestInit(t *testing.T) {
	f := Init()
	f.conf.Parse()

	if f.dbURI != defaultDBURI {
		t.Errorf("default db-uri got %#v; want %#v", f.dbURI, defaultDBURI)
	}
	if f.dbName != defaultDBName {
		t.Errorf("default db-name got %#v; want %#v", f.dbURI, defaultDBName)
	}
}

func TestConfigure(t *testing.T) {
	f := Init()
	want := "hello!"
	f.Configure("test", want, "test var")
	f.conf.Parse()

	got := f.Get("test")
	if got != want {
		t.Errorf("configure var \"test\" got %#v; want %#v", got, want)
	}
}

func TestRunAndShutdown(t *testing.T) {
	f := initForTest()

	var shutdown bool
	f.RegisterShutdown(func() {
		shutdown = true
	})

	f.Run(func(f *Facade) {
		if f.dbName != defaultTestDBName {
			t.Errorf("default test db name got %#v; want %#v", f.dbName, defaultTestDBName)
		}
		interrupt(t)
	})
	if !shutdown {
		t.Error("shutdown was not called")
	}
}

func initForTest() *Facade {
	os.Setenv("DB_NAME", defaultTestDBName)
	return Init()
}

func interrupt(t *testing.T) {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	p.Signal(os.Interrupt)
	time.Sleep(1 * time.Second)
}
