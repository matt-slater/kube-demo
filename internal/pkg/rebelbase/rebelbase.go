package rebelbase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/rs/zerolog"
)

type CommandCenter struct {
	WonBattles  int
	LostBattles int
	Commander   string
	Target      string
}
type RebelBase struct {
	lock        sync.Mutex
	WonBattles  int
	LostBattles int
	Commander   string
	Target      string
	Port        string
	Logger      zerolog.Logger
	HTTPClient  http.Client
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
	bytes, err := json.Marshal(CommandCenter{
		WonBattles:  rb.WonBattles,
		LostBattles: rb.LostBattles,
		Commander:   rb.Commander,
		Target:      rb.Target,
	})
	if err != nil {
		rb.Logger.Err(err)
		http.Error(w, "error getting command center", http.StatusInternalServerError)

		return
	}

	w.Write(bytes)
}

func (rb *RebelBase) Attack(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("target")
	if t == "" {
		http.Error(w, "please specify what you want to attack with the target query parameter", http.StatusBadRequest)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/%s", rb.Target, t), nil)
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

	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "target not found", http.StatusNotFound)
		return
	}

	if resp.StatusCode == http.StatusLocked {
		http.Error(w, "too many active defenses to attack this target", http.StatusBadRequest)
	}

	if resp.StatusCode == http.StatusForbidden {
		rb.lock.Lock()

		rb.LostBattles++

		rb.lock.Unlock()

		w.Write([]byte("lost x-wing in battle"))
		return
	}

	if resp.StatusCode == http.StatusOK {
		rb.lock.Lock()

		rb.WonBattles++

		rb.lock.Unlock()
		w.Write([]byte("successful attack"))
	}
}

func (rb *RebelBase) Recon(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/plans", rb.Target), nil)
	if err != nil {
		rb.Logger.Error().Err(err).Msg("error spying")
		http.Error(w, "error spying", http.StatusInternalServerError)

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

	rb.Logger.Info().Msg(fmt.Sprintf("status code: %d", resp.StatusCode))

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			rb.Logger.Error().Err(err).Msg("error reading body")
			http.Error(w, "could not read response from recon mission", http.StatusInternalServerError)
		}

		w.Write(bodyBytes)
	}
}

func (rb *RebelBase) Launch() error {
	http.HandleFunc("/healthz", rb.Health)
	http.HandleFunc("/", rb.CommandCenter)
	http.HandleFunc("/attack", rb.Attack)
	http.HandleFunc("/recon", rb.Recon)
	return http.ListenAndServe(fmt.Sprintf(":%s", rb.Port), nil)
}
