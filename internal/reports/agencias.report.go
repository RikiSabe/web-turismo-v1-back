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

var QueryAgencias = `
	SELECT 
		a.id_agencia,
		a.nombre,
		a.descripcion,
		a.direccion,
		COALESCE(a.telefono, 'N/A') as telefono,
		COALESCE(a.correo, 'N/A') as correo,
		COALESCE(d.nombre, 'Sin departamento') as departamento,
		COALESCE(u.nombre, 'Sin encargado') as encargado,
		COALESCE(CONCAT(u.apellido_paterno, ' ', COALESCE(u.apellido_materno, '')), '') as apellido_encargado,
		CASE WHEN a.estado = true THEN 'Activo' ELSE 'Inactivo' END as estado
	FROM "GestAgencias" a
	LEFT JOIN departamentos d ON a.id_departamento = d.id_departamento
	LEFT JOIN "GestUsuarios" u ON a.id_encargado = u.id_usuario
	ORDER BY d.nombre, a.nombre`

type AgenciaReporte struct {
	ID                uint   `json:"id_agencia"`
	Nombre            string `json:"nombre"`
	Descripcion       string `json:"descripcion"`
	Direccion         string `json:"direccion"`
	Telefono          string `json:"telefono"`
	Correo            string `json:"correo"`
	Departamento      string `json:"departamento"`
	Encargado         string `json:"encargado"`
	ApellidoEncargado string `json:"apellido_encargado"`
	Estado            string `json:"estado"`
}

func ReporteAgencias(w http.ResponseWriter, r *http.Request) {
	m, err := makePDFAgencias()
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
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_agencias.pdf\"")

	if _, err := w.Write(doc.GetBytes()); err != nil {
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
		return
	}
}

func makePDFAgencias() (core.Maroto, error) {
	var agencias []AgenciaReporte
	err := db.GDB.Raw(QueryAgencias).Scan(&agencias).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener agencias: %v", err)
	}

	if len(agencias) == 0 {
		return nil, fmt.Errorf("no se encontraron agencias")
	}

	cfg := config.NewBuilder().
		WithLeftMargin(10).
		WithRightMargin(10).
		WithTopMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	rows, err := getAgenciasRows(agencias)
	if err != nil {
		return nil, err
	}

	m.AddRows(rows...)

	return m, nil
}

func getAgenciasRows(agencias []AgenciaReporte) ([]core.Row, error) {
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
				text.New("REPORTE DE AGENCIAS TURÍSTICAS", props.Text{
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
		col.New(2).Add(
			text.New("NOMBRE", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  7,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(2).Add(
			text.New("DIRECCIÓN", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  7,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(2).Add(
			text.New("CORREO", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  7,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(2).Add(
			text.New("DEPARTAMENTO", props.Text{
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
	contadorPorDepartamento := make(map[string]int)

	for _, agencia := range agencias {
		contadorPorDepartamento[agencia.Departamento]++

		nombre := agencia.Nombre
		if len(nombre) > 25 {
			nombre = nombre[:22] + "..."
		}

		direccion := agencia.Direccion
		if len(direccion) > 25 {
			direccion = direccion[:22] + "..."
		}

		correo := agencia.Correo
		if len(correo) > 25 {
			correo = correo[:22] + "..."
		}

		departamento := agencia.Departamento
		if len(departamento) > 20 {
			departamento = departamento[:17] + "..."
		}

		encargadoCompleto := agencia.Encargado
		if agencia.ApellidoEncargado != "" {
			encargadoCompleto = agencia.Encargado + " " + agencia.ApellidoEncargado
		}
		if len(encargadoCompleto) > 25 {
			encargadoCompleto = encargadoCompleto[:22] + "..."
		}

		r := row.New(6).Add(
			col.New(1).Add(
				text.New(fmt.Sprintf("%d", agencia.ID), props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Center,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(2).Add(
				text.New(nombre, props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Left,
					Left:  0.3,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(2).Add(
				text.New(direccion, props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Left,
					Left:  0.3,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(2).Add(
				text.New(correo, props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Left,
					Left:  0.3,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(2).Add(
				text.New(departamento, props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Center,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(2).Add(
				text.New(encargadoCompleto, props.Text{
					Top:   0.5,
					Size:  6,
					Align: align.Center,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(1).Add(
				text.New(agencia.Estado, props.Text{
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
				text.New("RESUMEN POR DEPARTAMENTO", props.Text{
					Top:   1.0,
					Align: align.Left,
					Style: fontstyle.Bold,
					Size:  10,
				}),
			),
		),
	)

	for departamento, cantidad := range contadorPorDepartamento {
		rows = append(rows,
			row.New(6).Add(
				col.New(12).Add(
					text.New(fmt.Sprintf("  • %s: %d agencia(s)", departamento, cantidad), props.Text{
						Top:   0.5,
						Align: align.Left,
						Size:  9,
					}),
				),
			),
		)
	}

	rows = append(rows, row.New(3))

	rows = append(rows,
		row.New(7).Add(
			col.New(12).Add(
				text.New(fmt.Sprintf("TOTAL DE AGENCIAS: %d", len(agencias)), props.Text{
					Top:   1.0,
					Align: align.Right,
					Style: fontstyle.Bold,
					Size:  10,
					Right: 1,
				}),
			),
		),
	)

	return rows, nil
}
