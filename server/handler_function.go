package server

import (
	"net/http"
)

type Item struct {
	Item string
}

func (rt Router) AddItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		id, err := rt.MongoDBClient.CreateItem(Item{Item: "Database Connected"})
		if err != nil {
			println(err.Error())
		}

		w.Write([]byte(id))

		return
	}
}
