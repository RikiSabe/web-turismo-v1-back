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

var QueryUsuarios = `
	SELECT 
		u.id_usuario,
		u.rol,
		u.nombre,
		CONCAT(u.apellido_paterno, ' ', COALESCE(u.apellido_materno, '')) as apellidos,
		u.ci,
		u.correo,
		COALESCE(u.telefono, 'N/A') as telefono,
		TO_CHAR(u.fecha_nacimiento, 'DD/MM/YYYY') as fecha_nacimiento,
		COALESCE(p.nombre, 'Sin ubicación') as provincia,
		CASE WHEN u.estado = true THEN 'Activo' ELSE 'Inactivo' END as estado
	FROM "GestUsuarios" u
	LEFT JOIN provincias p ON u.id_ubicacion = p.id_provincia
	ORDER BY u.rol, u.nombre`

type UsuarioReporte struct {
	ID              uint   `json:"id_usuario"`
	Rol             string `json:"rol"`
	Nombre          string `json:"nombre"`
	Apellidos       string `json:"apellidos"`
	CI              string `json:"ci"`
	Correo          string `json:"correo"`
	Telefono        string `json:"telefono"`
	FechaNacimiento string `json:"fecha_nacimiento"`
	Provincia       string `json:"provincia"`
	Estado          string `json:"estado"`
}

func ReporteUsuarios(w http.ResponseWriter, r *http.Request) {
	m, err := makePDFUsuarios()
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
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_usuarios.pdf\"")

	if _, err := w.Write(doc.GetBytes()); err != nil {
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
		return
	}
}

func makePDFUsuarios() (core.Maroto, error) {
	var usuarios []UsuarioReporte
	err := db.GDB.Raw(QueryUsuarios).Scan(&usuarios).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %v", err)
	}

	if len(usuarios) == 0 {
		return nil, fmt.Errorf("no se encontraron usuarios")
	}

	cfg := config.NewBuilder().
		WithLeftMargin(10).
		WithRightMargin(10).
		WithTopMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	rows, err := getUsuariosRows(usuarios)
	if err != nil {
		return nil, err
	}

	m.AddRows(rows...)

	return m, nil
}

func getUsuariosRows(usuarios []UsuarioReporte) ([]core.Row, error) {
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
				text.New("REPORTE DE USUARIOS", props.Text{
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
				Size:  8,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(2).Add(
			text.New("NOMBRE COMPLETO", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  8,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(1).Add(
			text.New("ROL", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  8,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(1).Add(
			text.New("CI", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  8,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(2).Add(
			text.New("CORREO", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  8,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(2).Add(
			text.New("TELÉFONO", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  8,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
		col.New(1).Add(
			text.New("ESTADO", props.Text{
				Top:   0.5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  8,
			}),
		).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.2}),
	)
	rows = append(rows, headerRow)

	var contentsRow []core.Row
	contadorPorRol := make(map[string]int)

	for _, usuario := range usuarios {
		contadorPorRol[usuario.Rol]++

		nombreCompleto := usuario.Nombre + " " + usuario.Apellidos
		if len(nombreCompleto) > 25 {
			nombreCompleto = nombreCompleto[:22] + "..."
		}

		correo := usuario.Correo
		if len(correo) > 20 {
			correo = correo[:17] + "..."
		}

		r := row.New(6).Add(
			col.New(1).Add(
				text.New(fmt.Sprintf("%d", usuario.ID), props.Text{
					Top:   0.5,
					Size:  7,
					Align: align.Center,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(2).Add(
				text.New(nombreCompleto, props.Text{
					Top:   0.5,
					Size:  7,
					Align: align.Left,
					Left:  0.5,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(1).Add(
				text.New(usuario.Rol, props.Text{
					Top:   0.5,
					Size:  7,
					Align: align.Center,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(1).Add(
				text.New(usuario.CI, props.Text{
					Top:   0.5,
					Size:  7,
					Align: align.Center,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(2).Add(
				text.New(correo, props.Text{
					Top:   0.5,
					Size:  7,
					Align: align.Left,
					Left:  0.5,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(2).Add(
				text.New(usuario.Telefono, props.Text{
					Top:   0.5,
					Size:  7,
					Align: align.Center,
				}),
			).WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.1}),
			col.New(1).Add(
				text.New(usuario.Estado, props.Text{
					Top:   0.5,
					Size:  7,
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
				text.New("RESUMEN POR ROL", props.Text{
					Top:   1.0,
					Align: align.Left,
					Style: fontstyle.Bold,
					Size:  10,
				}),
			),
		),
	)

	for rol, cantidad := range contadorPorRol {
		rows = append(rows,
			row.New(6).Add(
				col.New(12).Add(
					text.New(fmt.Sprintf("  • %s: %d usuario(s)", rol, cantidad), props.Text{
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
				text.New(fmt.Sprintf("TOTAL DE USUARIOS: %d", len(usuarios)), props.Text{
					Top:   1.0,
					Align: align.Left,
					Style: fontstyle.Bold,
					Size:  10,
					Right: 1,
				}),
			),
		),
	)

	return rows, nil
}
