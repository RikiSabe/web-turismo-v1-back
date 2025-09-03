package controllers

import (
	"encoding/json"
	"net/http"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"

	"github.com/gorilla/mux"
)

func ObtenerSubCategoriasPorCategoria(w http.ResponseWriter, r *http.Request) {
	id_categoria := mux.Vars(r)["id"]
	var subCategorias []models.SubCategoria

	err := db.GDB.
		Table("subcategoria").
		Select("nombre, descripcion, estado").
		Scan(&subCategorias).
		Where("id_categoria = ?", id_categoria).
		Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
