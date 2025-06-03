package services

var QueryRolesTODO = `
	SELECT 
		r.id_rol as id,
		r.nombre,
		r.estado
	FROM "GestRoles" as r
	ORDER BY id_rol ASC`

var QueryRolUnique = `
	SELECT 
		r.id_rol as id,
		r.nombre,
		r.estado
	FROM "GestRoles" as r
	WHERE id_rol = ?
	LIMIT 1`
