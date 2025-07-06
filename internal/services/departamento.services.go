package services

var QueryDepartamentosTODO = `
	SELECT
		d.id_departamento as id,
		d.nombre
	FROM
		departamentos as d
	WHERE
		d.estado = True`
