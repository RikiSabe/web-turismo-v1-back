package reports

import (
	"fmt"
	"net/http"
	"time"
	"web-turismo-v1/internal/db"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

var QueryAtracciones = `
	SELECT 
		a.id_atraccion,
		a.nombre,
		a.direccion,
		a.descripcion,
		a.horario_apertura,
		a.horario_cierre,
		a.precio,
		COALESCE(p.nombre, 'Sin ubicación') as provincia,
		COALESCE(u.nombre, 'Sin encargado') as encargado,
		CASE WHEN a.estado = true THEN 'Activo' ELSE 'Inactivo' END as estado
	FROM "GestAtraccionesTuristicas" a
	LEFT JOIN provincias p ON a.id_ubicacion = p.id_provincia
	LEFT JOIN "GestUsuarios" u ON a.id_encargado = u.id_usuario
	ORDER BY p.nombre, a.nombre`

type AtraccionReporte struct {
	ID              uint    `json:"id_atraccion"`
	Nombre          string  `json:"nombre"`
	Direccion       string  `json:"direccion"`
	Descripcion     string  `json:"descripcion"`
	HorarioApertura string  `json:"horario_apertura"`
	HorarioCierre   string  `json:"horario_cierre"`
	Precio          float64 `json:"precio"`
	Provincia       string  `json:"provincia"`
	Encargado       string  `json:"encargado"`
	Estado          string  `json:"estado"`
}

func ReporteAtracciones(w http.ResponseWriter, r *http.Request) {
	m, err := makePDFAtracciones()
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	doc, err := m.Generate()
	if err != nil {
		http.Error(w, "Error al generar pdf: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_atracciones.pdf\"")

	if _, err := w.Write(doc.GetBytes()); err != nil {
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
		return
	}
}

func makePDFAtracciones() (core.Maroto, error) {
	var atracciones []AtraccionReporte
	err := db.GDB.Raw(QueryAtracciones).Scan(&atracciones).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener atracciones: %v", err)
	}

	if len(atracciones) == 0 {
		return nil, fmt.Errorf("no se encontraron atracciones turísticas")
	}

	cfg := config.NewBuilder().
		WithLeftMargin(10).
		WithRightMargin(10).
		WithTopMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	rows, err := getAtraccionesRows(atracciones)
	if err != nil {
		return nil, err
	}

	m.AddRows(rows...)

	return m, nil
}

func getAtraccionesRows(atracciones []AtraccionReporte) ([]core.Row, error) {
	var rows []core.Row

	rows = append(rows,
		row.New(12).Add(
			col.New(12).Add(
				text.New("SISTEMA DE GESTIÓN TURÍSTICA", props.Text{
					Top:   3,
					Style: fontstyle.Bold,
					Align: align.Center,
					Size:  16,
				}),
			),
		),
		row.New(8).Add(
			col.New(12).Add(
				text.New("REPORTE DE ATRACCIONES TURÍSTICAS", props.Text{
					Top:   1.0,
					Align: align.Center,
					Style: fontstyle.Bold,
					Size:  14,
				}),
			),
		),
		row.New(8).Add(
			col.New(12).Add(
				text.New("Fecha de generación: "+time.Now().Format("02/01/2006 15:04:05"), props.Text{
					Top:   1.0,
					Align: align.Center,
					Size:  10,
				}),
			),
		),
		row.New(5),
	)

	headerRow := row.New(7).Add(
		col.New(1).Add(
			text.New("ID", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  7,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(3).Add(
			text.New("NOMBRE", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  7,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(2).Add(
			text.New("HORARIO", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  7,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(1).Add(
			text.New("PRECIO", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  7,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(2).Add(
			text.New("PROVINCIA", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  7,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(2).Add(
			text.New("ENCARGADO", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  7,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(1).Add(
			text.New("ESTADO", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  7,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
	)
	rows = append(rows, headerRow)

	var contentsRow []core.Row
	contadorPorProvincia := make(map[string]int)
	var precioTotal float64

	for _, atraccion := range atracciones {
		contadorPorProvincia[atraccion.Provincia]++
		precioTotal += atraccion.Precio

		nombre := atraccion.Nombre
		if len(nombre) > 30 {
			nombre = nombre[:27] + "..."
		}

		horario := "Cerrado"
		if atraccion.HorarioApertura != "" && atraccion.HorarioCierre != "" {
			horario = fmt.Sprintf("%s-%s", atraccion.HorarioApertura, atraccion.HorarioCierre)
			if len(horario) > 15 {
				horario = horario[:12] + "..."
			}
		}

		provincia := atraccion.Provincia
		if len(provincia) > 20 {
			provincia = provincia[:17] + "..."
		}

		encargado := atraccion.Encargado
		if len(encargado) > 20 {
			encargado = encargado[:17] + "..."
		}

		r := row.New(6).Add(
			col.New(1).Add(
				text.New(fmt.Sprintf("%d", atraccion.ID), props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Center,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(3).Add(
				text.New(nombre, props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Left,
					Left:  0.3,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(2).Add(
				text.New(horario, props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Center,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(1).Add(
				text.New(fmt.Sprintf("%.2f", atraccion.Precio), props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Right,
					Right: 0.3,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(2).Add(
				text.New(provincia, props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Center,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(2).Add(
				text.New(encargado, props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Center,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(1).Add(
				text.New(atraccion.Estado, props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Center,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
		)
		contentsRow = append(contentsRow, r)
	}

	rows = append(rows, contentsRow...)

	rows = append(rows, row.New(5))

	rows = append(rows,
		row.New(7).Add(
			col.New(12).Add(
				text.New("RESUMEN POR PROVINCIA", props.Text{
					Top:   1.0,
					Align: align.Left,
					Style: fontstyle.Bold,
					Size:  10,
				}),
			),
		),
	)

	for provincia, cantidad := range contadorPorProvincia {
		rows = append(rows,
			row.New(6).Add(
				col.New(12).Add(
					text.New(fmt.Sprintf("  • %s: %d atracción(es)", provincia, cantidad), props.Text{
						Top:   0.5,
						Align: align.Left,
						Size:  9,
					}),
				),
			),
		)
	}

	rows = append(rows, row.New(3))

	return rows, nil
}
