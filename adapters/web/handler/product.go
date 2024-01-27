package handler

import (
	"encoding/json"
	"go/src/application"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func MakeProductHandler(r *mux.Router, n *negroni.Negroni, s application.ProductServiceInterface) {
	r.Handle("/product/{id}", n.With(
		negroni.Wrap(getProduct(s)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/product", n.With(
		negroni.Wrap(saveProduct(s)),
	)).Methods("POST")
}

func getProduct(service application.ProductServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		id := vars["id"]
		product, err := service.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})
}

func saveProduct(service application.ProductServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		var productData struct {
			Price float64 `json:"price"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&productData); err != nil {
			http.Error(w, "Erro ao decodificar o corpo da solicitação", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		price := productData.Price
		newProduct, err := service.Create(name, price)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(newProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
