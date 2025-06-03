package services

var QueryUsuariosTODO = `
	SELECT 
		u.id_usuario as id,
		u.rol,
		u.nombre,
		u.apellido, 
		u.correo, 
		u.telefono,
		u.direccion,
		u.contra,
		u.estado,
		u.foto
	FROM "GestUsuarios" as u
	ORDER BY id_usuario ASC`

var QueryUsuarioUnique = `
	SELECT
		u.id_usuario as id,
		u.rol,
		u.nombre,
		u.apellido, 
		u.correo, 
		u.telefono,
		u.direccion,
		u.contra,
		u.estado,
		u.foto
	FROM "GestUsuarios" as u
	WHERE id_usuario = ?
	LIMIT 1`
