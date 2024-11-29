package telegram

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jellydator/ttlcache/v3"
	"github.com/oybek/choguuket/model"
	"github.com/samber/lo"
)

type AptekaPayload struct {
	Name      string   `json:"name"`
	Phone     string   `json:"phone"`
	Address   string   `json:"address"`
	Medicines []string `json:"medicines"`
}

func (lp *LongPoll) GetRequest(w http.ResponseWriter, r *http.Request) {
	lp.requestCache.Set(
		uuid.MustParse("26b5201f-c1a2-453a-bccb-1e268191b3bf"),
		[]lo.Tuple2[model.Apteka, []string]{
			lo.T2(
				model.Apteka{
					Name:    "Неман",
					Phone:   "0559171775",
					Address: "Токтоналиева 61",
				},
				[]string{"Анальгин", "Парацетамол"},
			)},
		ttlcache.DefaultTTL,
	)

	vars := mux.Vars(r)
	requestId, ok := vars["uuid"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rUuid, err := uuid.Parse(requestId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	entry := lp.requestCache.Get(rUuid)
	if entry == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	aptekas := entry.Value()
	aptekaPayload := make([]AptekaPayload, 0, len(aptekas))
	for _, apteka := range aptekas {
		aptekaPayload = append(aptekaPayload, AptekaPayload{
			Name:      apteka.A.Name,
			Phone:     apteka.A.Phone,
			Address:   apteka.A.Address,
			Medicines: apteka.B,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	jsonRaw, err := json.Marshal(aptekaPayload)
	if err != nil {
		log.Printf("json marshall error %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRaw)
}
