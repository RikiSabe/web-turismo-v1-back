package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"

	"github.com/gorilla/mux"
)

func ObtenerSubCategoriasPorCategoria(w http.ResponseWriter, r *http.Request) {
	id_categoria := mux.Vars(r)["id"]
	var subCategorias []models.SubCategoria

	err := db.GDB.
		Table("subcategorias s").
		Select("s.id_subcategoria, s.nombre, s.descripcion, s.estado").
		Where("s.id_categoria = ?", id_categoria).
		Scan(&subCategorias).
		Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if b, err := json.MarshalIndent(subCategorias, " ", "  "); err == nil {
		fmt.Println(string(b))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subCategorias)
}

func ObtenerSubCategoria(w http.ResponseWriter, r *http.Request) {

}

func AgregarSubCategoria(w http.ResponseWriter, r *http.Request) {

}

func ModificarSubCategoria(w http.ResponseWriter, r *http.Request) {

}
