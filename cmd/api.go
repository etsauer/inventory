package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/openbar/inventory/pkg/apis"
)

var Items []apis.InventoryItem

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllItems(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Items)
}

func returnSingleItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, item := range Items {
		if item.Metadata.Name == key {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func createNewItem(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Item struct
	// append this to our Item array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var item apis.InventoryItem
	json.Unmarshal(reqBody, &item)

	for _, i := range Items {
		if item.Metadata.Name == i.Metadata.Name {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}
	// update our global Item array to include
	// our new Item
	Items = append(Items, item)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	// once again, we will need to parse the path parameters
	vars := mux.Vars(r)
	// we will need to extract the `id` of the article we
	// wish to delete
	id := vars["id"]

	// we are going to use a downward loop (from end of slice to first)
	for i := len(Items) - 1; i >= 0; i-- {
		item := Items[i]
		// if our id path parameter matches one of our items
		if item.Metadata.Name == id {
			// updates our Items array to remove the item
			w.WriteHeader(http.StatusNoContent)
			Items = append(Items[:i], Items[i+1:]...)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for i, item := range Items {
		if item.Metadata.Name == id {
			reqBody, _ := ioutil.ReadAll(r.Body)
			var item apis.InventoryItem
			json.Unmarshal(reqBody, &item)

			Items[i] = item
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/items", returnAllItems).Methods("GET")
	myRouter.HandleFunc("/item", createNewItem).Methods("POST")
	myRouter.HandleFunc("/item/{id}", returnSingleItem).Methods("GET")
	myRouter.HandleFunc("/item/{id}", deleteItem).Methods("DELETE")
	myRouter.HandleFunc("/item/{id}", updateItem).Methods("PUT")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	Items = []apis.InventoryItem{
		apis.InventoryItem{
			Metadata: apis.ResourceMeta{
				Name:    "appleton12yr",
				Barname: "kitchen",
				Labels: map[string]string{
					"category": "rum",
					"type":     "liquor",
				},
			},
			Spec: apis.InventorySpec{
				FillLevel:   50,
				Maker:       "Appleton",
				Product:     "Estate 12 Year",
				Description: "Funky Rum",
			},
		},
		apis.InventoryItem{
			Metadata: apis.ResourceMeta{
				Name:    "brugal1888",
				Barname: "kitchen",
				Labels: map[string]string{
					"category": "rum",
					"type":     "liquor",
				},
			},
			Spec: apis.InventorySpec{
				FillLevel:   75,
				Maker:       "Brugal",
				Product:     "1888",
				Description: "Smooth, Silky",
			},
		},
	}
	handleRequests()
}
