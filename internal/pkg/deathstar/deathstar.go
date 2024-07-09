package deathstar

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"

	"github.com/rs/zerolog"
)

type DeathStar struct {
	Destroyed    bool
	TIEFighters  int
	LaserCannons int
	Port         string
	Logger       zerolog.Logger
	lock         sync.Mutex
}

func New(port string, logger zerolog.Logger) *DeathStar {
	return &DeathStar{
		Destroyed:    false,
		TIEFighters:  10,
		LaserCannons: 10,
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

	if ds.TIEFighters > 0 {
		http.Error(w, fmt.Sprintf("too many defenses. %d TIE fighters remain", ds.TIEFighters), http.StatusLocked)

		return
	}

	ds.lock.Lock()

	ds.Destroyed = true

	ds.lock.Unlock()

	ds.Logger.Info().Msg(fmt.Sprintf("rebel commander %s destroyed death star!", r.Header.Get("commander")))

	w.Write([]byte(fmt.Sprintf("rebel commander %s destroyed death star!", r.Header.Get("commander"))))
}

func (ds *DeathStar) Battle(w http.ResponseWriter, r *http.Request) {
	if ds.TIEFighters > 0 {
		battleProb := rand.Intn(100)

		if battleProb >= 50 { //succeed

			ds.lock.Lock()

			ds.TIEFighters--

			ds.lock.Unlock()

			ds.Logger.Info().Msg(fmt.Sprintf("rebel commander %s destroyed 1 TIE fighter!", r.Header.Get("commander")))
			ds.Logger.Info().Msg(fmt.Sprintf("there are %d TIE fighters remaining", ds.TIEFighters))
			w.Write([]byte("destroyed 1 TIE Fighter"))

			return
		}

		ds.Logger.Info().Msg(fmt.Sprintf("destroyed rebel x-wing commanded by %s", r.Header.Get("commander")))
		http.Error(w, "attack unsuccessful", http.StatusForbidden)

		return
	}
	
	http.Error(w, "no more tie fighters left", http.StatusNotFound)
}

func (ds *DeathStar) Plans(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("weakness: thermalexhaustport"))
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
	http.HandleFunc("/thermalexhaustport", ds.Destroy)
	http.HandleFunc("/commandcenter", ds.CommandCenter)
	http.HandleFunc("/tiefighter", ds.Battle)
	http.HandleFunc("/plans", ds.Plans)
	return http.ListenAndServe(fmt.Sprintf(":%s", ds.Port), nil)
}
