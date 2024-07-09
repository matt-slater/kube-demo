package rebelbase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type RebelBase struct {
	Commander  string
	Target     string
	Port       string
	Logger     zerolog.Logger
	HTTPClient http.Client
}

func New(commander, target, port string, logger zerolog.Logger) *RebelBase {
	return &RebelBase{
		Commander:  commander,
		Port:       port,
		Target:     target,
		Logger:     logger,
		HTTPClient: *http.DefaultClient,
	}
}

func (rb *RebelBase) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("healthy"))
}

func (rb *RebelBase) CommandCenter(w http.ResponseWriter, _ *http.Request) {
	bytes, err := json.Marshal(rb)
	if err != nil {
		http.Error(w, "error getting command center", http.StatusInternalServerError)

		return
	}

	w.Write(bytes)
}

func (rb *RebelBase) Attack(w http.ResponseWriter, _ *http.Request) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/", rb.Target), nil)
	if err != nil {
		rb.Logger.Error().Err(err).Msg("error attacking")
		http.Error(w, "error attacking", http.StatusInternalServerError)

		return
	}

	req.Header.Set("commander", rb.Commander)

	resp, err := rb.HTTPClient.Do(req)
	if err != nil {
		rb.Logger.Error().Err(err).Msg("error attacking")
		http.Error(w, "error attacking", http.StatusInternalServerError)

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		w.Write([]byte("successful attack"))
	}
}

func (rb *RebelBase) Launch() error {
	http.HandleFunc("/healthz", rb.Health)
	http.HandleFunc("/commandcenter", rb.CommandCenter)
	http.HandleFunc("/attack", rb.Attack)
	return http.ListenAndServe(fmt.Sprintf(":%s", rb.Port), nil)
}
