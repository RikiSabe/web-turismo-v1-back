package services

var QueryUsuariosResumen = `
	SELECT 
		u.id_usuario AS id,
		u.rol,
		u.nombre,
		u.apellido_paterno,
		u.apellido_materno,
		u.fecha_nacimiento,
		u.ci, 
		u.correo, 
		u.telefono,
		u.foto,
		u.estado
	FROM "GestUsuarios" AS u
	ORDER BY u.id_usuario ASC`

var QueryUsuarioUnique = `
	SELECT
		u.id_usuario AS id,
		u.rol,
		u.nombre,
		u.apellido_paterno,
		u.apellido_materno,
		TO_CHAR(u.fecha_nacimiento, 'YYYY-MM-DD') AS fecha_nacimiento,
		u.ci, 
		u.correo, 
		u.telefono,
		u.contra,
		u.estado,
		u.foto,
		u.id_ubicacion
	FROM "GestUsuarios" AS u
	WHERE u.id_usuario = ?
	LIMIT 1`

var QueryUsuarioMenu = `
	SELECT 
		CONCAT(
			u.nombre, 
			' ', 
			u.apellido_paterno, 
			' ', 
			u.apellido_materno
		) AS nombre_completo,
		u.foto
	FROM "GestUsuarios" as u
	WHERE id_usuario = ?
	LIMIT 1`
