package controllers

import (
	"encoding/json"
	"net/http"
	"web-turismo-v1/internal/db"
	"web-turismo-v1/internal/models"
	"web-turismo-v1/internal/types"

	"github.com/gorilla/mux"
)

func ObtenerCategorias(w http.ResponseWriter, r *http.Request) {
	var categorias []models.Categoria

	err := db.GDB.
		Table("categorias").
		Select("id_categoria, nombre, descripcion, estado").
		Scan(&categorias).
		Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categorias)
}

func ObtenerCategoria(w http.ResponseWriter, r *http.Request) {
	id_categoria := mux.Vars(r)["id"]
	var categoria types.Categoria

	err := db.GDB.
		Table("categorias").
		Select("nombre, descripcion, estado").
		Where("id_categoria = ?", id_categoria).
		First(&categoria).
		Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoria)
}

func AgregarCategoria(w http.ResponseWriter, r *http.Request) {
	var nuevaCategoria models.Categoria

	if err := json.NewDecoder(r.Body).Decode(&nuevaCategoria); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevaCategoria).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al crear la categoria: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nuevaCategoria)
}

func ModificarCategoria(w http.ResponseWriter, r *http.Request) {
	id_categoria := mux.Vars(r)["id"]
	var categoriaExistente models.Categoria

	err := db.GDB.
		Where("id_categoria = ?", id_categoria).
		First(&categoriaExistente).
		Error

	if err != nil {
		http.Error(w, "Categoria no encontrada", http.StatusNotFound)
		return
	}

	var nuevaCategoria types.Categoria

	if err := json.NewDecoder(r.Body).Decode(&nuevaCategoria); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	categoriaExistente.Nombre = nuevaCategoria.Nombre
	categoriaExistente.Descripcion = nuevaCategoria.Descripcion
	categoriaExistente.Estado = nuevaCategoria.Estado

	if err := db.GDB.Save(&categoriaExistente).Error; err != nil {
		http.Error(w, "Error al modificar categoria", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoriaExistente)
}
