package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/samuelastech/products-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.getProducts(w, r)
	case http.MethodPost:
		p.addProduct(w, r)
	case http.MethodPut:
		idRegex := regexp.MustCompile(`([0-9]+)`)
		group := idRegex.FindAllStringSubmatch(r.URL.Path, -1)

		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Recovered from panic:", err)
			}
		}()

		if len(group) != 1 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(group[0]) != 2 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		id, _ := strconv.Atoi(group[0][1])
		p.updateProduct(id, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	product := &data.Product{}
	err := product.FromJSON(r.Body)

	if err != nil {
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
	}

	data.AddProduct(product)
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()
	err := products.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to write JSON", http.StatusInternalServerError)
		return
	}
}

func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	product := &data.Product{}
	err := product.FromJSON(r.Body)

	if err != nil {
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, product)

	if err == data.ErrorProductNotFound {
		http.Error(w, data.ErrorProductNotFound.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
