package services

var QueryAtraccionesTuristicasTODO = `
	SELECT 
  at.id_atraccion AS id,
  at.categoria,
  at.nombre,
  at.direccion,
  at.descripcion,
  at.horario_apertura,
  at.horario_cierre,
  at.precio,
  at.estado,
  jsonb_build_object(
    'id', u.id_usuario,
    'nombre_completo', concat(u.nombre, ' ', u.apellido_paterno, ' ', u.apellido_materno)
  ) AS encargado,
  jsonb_build_object(
    'provincia', jsonb_build_object(
      'id', pro.id_provincia,
      'nombre', pro.nombre
    ),
    'departamento', jsonb_build_object(
      'id', d.id_departamento,
      'nombre', d.nombre
    )
  ) AS ubicacion
  FROM "GestAtraccionesTuristicas" AS at
  LEFT JOIN "GestUsuarios" AS u ON at.id_encargado = u.id_usuario 
  LEFT JOIN provincias AS pro ON at.id_ubicacion = pro.id_provincia  
  LEFT JOIN departamentos AS d ON d.id_departamento = pro.id_departamento
  GROUP BY
    at.id_atraccion, at.categoria, at.nombre, at.direccion, at.descripcion,
    at.horario_apertura, at.horario_cierre, at.precio, at.estado,
    u.id_usuario, u.nombre, u.apellido_paterno, u.apellido_materno,
    pro.id_provincia, pro.nombre,
    d.id_departamento, d.nombre
  ORDER BY at.id_atraccion asc;`

var QueryAtraccionesTuristicaUnique = `
	SELECT 
  at.id_atraccion AS id,
  at.categoria,
  at.nombre,
  at.direccion,
  at.descripcion,
  at.horario_apertura,
  at.horario_cierre,
  at.precio,
  at.estado,
  jsonb_build_object(
    'id', u.id_usuario,
    'nombre_completo', concat(u.nombre, ' ', u.apellido_paterno, ' ', u.apellido_materno)
  ) AS encargado,
  jsonb_build_object(
    'provincia', jsonb_build_object(
      'id', pro.id_provincia,
      'nombre', pro.nombre
    ),
    'departamento', jsonb_build_object(
      'id', d.id_departamento,
      'nombre', d.nombre
    )
  ) AS ubicacion,
  CASE
    WHEN COUNT(fa.id_foto) = 0 THEN to_jsonb('N/A'::text)
    ELSE jsonb_agg(
      DISTINCT jsonb_build_object(
        'id', fa.id_foto,
        'foto', fa.foto,
        'orden', fa.orden
      )
    )
  END AS fotos
	FROM "GestAtraccionesTuristicas" AS at
	LEFT JOIN "GestUsuarios" AS u ON at.id_encargado = u.id_usuario 
	LEFT JOIN provincias AS pro ON at.id_ubicacion = pro.id_provincia  
	LEFT JOIN departamentos AS d ON d.id_departamento = pro.id_departamento
	LEFT JOIN fotos_atracciones AS fa ON fa.id_atraccion = at.id_atraccion
	WHERE at.id_atraccion = ?
	GROUP BY 
		at.id_atraccion, at.categoria, at.nombre, at.direccion, at.descripcion,
		at.horario_apertura, at.horario_cierre, at.precio, at.estado,
		u.id_usuario, u.nombre, u.apellido_paterno, u.apellido_materno,
		pro.id_provincia, pro.nombre,
		d.id_departamento, d.nombre
	LIMIT 1`

var QueryAtraccionesTuristicasEnlazadas = `
	SELECT
		at.id_atracciones as id, 
		at.tipo, 
		at.nombre, 
		at.ubicacion, 
		at.descripcion, 
		at.horarios, 
		at.precio, 
		at.estado
	FROM "GestAtraccionesTuristicas" as at, "PaquetesAtracciones" as pa
	WHERE at.id_atracciones = pa.id_atraccion
	AND pa.id_paquete = ?`

var QueryEncargadoAtraccionTuristicas = `
  SELECT
    concat(u.nombre, ' ', u.apellido_paterno, ' ', u.apellido_materno) as nombre_completo,
    u.foto as src,
    u.ci,
    u.correo, 
    u.telefono
  FROM "GestUsuarios" as u
  WHERE u.id_usuario = ?
  limit 1;`
