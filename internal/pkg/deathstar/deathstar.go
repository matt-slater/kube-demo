package deathstar

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type DeathStar struct {
	Destroyed    bool
	TIEFighters  int
	LaserCannons int
	Port         string
	Logger       zerolog.Logger
}

func New(port string, logger zerolog.Logger) *DeathStar {
	return &DeathStar{
		Destroyed:    false,
		TIEFighters:  100,
		LaserCannons: 100,
		Port:         port,
		Logger:       logger,
	}
}

func (ds *DeathStar) Health(w http.ResponseWriter, r *http.Request) {
	if ds.Destroyed {
		http.Error(w, "death star has been destroyed", http.StatusGone)

		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("healthy"))
}

func (ds *DeathStar) Destroy(w http.ResponseWriter, r *http.Request) {
	if ds.Destroyed {
		http.Error(w, "the death star has been destroyed!", http.StatusGone)

		return
	}

	ds.Destroyed = true

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("destroyed death star!"))
}

func (ds *DeathStar) Index(w http.ResponseWriter, r *http.Request) {
	if ds.TIEFighters > 0 {
		ds.TIEFighters--
		ds.Logger.Info().Msg(fmt.Sprintf("rebel commander %s destroyed 1 TIE fighter!", r.Header.Get("commander")))
		ds.Logger.Info().Msg(fmt.Sprintf("there are %d TIE fighters remaining", ds.TIEFighters))
		w.Write([]byte("destroyed 1 TIE Fighter"))
		return
	}
	w.Write([]byte("all TIE fighters destroyed"))
}

func (ds *DeathStar) CommandCenter(w http.ResponseWriter, _ *http.Request) {
	bytes, err := json.Marshal(ds)
	if err != nil {
		http.Error(w, "error getting command center", http.StatusInternalServerError)

		return
	}

	w.Write(bytes)
}

func (ds *DeathStar) Launch() error {
	http.HandleFunc("/healthz", ds.Health)
	http.HandleFunc("/destroy", ds.Destroy)
	http.HandleFunc("/commandcenter", ds.CommandCenter)
	http.HandleFunc("/", ds.Index)
	return http.ListenAndServe(fmt.Sprintf(":%s", ds.Port), nil)
}
