# Fase 1. Análisis del problema

## 1.1. Situación actual

Actualmente, el sistema biblioteca-api cuenta con un módulo de préstamos desarrollado bajo una arquitectura por capas: handler, service y repository.

Al revisar la migración inicial `001_init.sql`, se encontró que la tabla `loans` contiene una columna llamada `book_id`, definida como obligatoria mediante `NOT NULL` y relacionada mediante una llave foránea con la tabla `books`.

Esto demuestra que el modelo actual de préstamos está directamente acoplado al concepto de libro.

## 1.2. Problema identificado

Si la biblioteca desea incorporar préstamos de laptops, proyectores, tablets o salas, el diseño actual presenta una limitación, porque la tabla de préstamos solamente permite relacionar cada préstamo con un libro.

Agregar nuevos tipos de recursos utilizando el diseño actual podría requerir modificaciones en la base de datos, los modelos, las solicitudes, los servicios y los repositorios.

Esto aumentaría el acoplamiento y el costo de mantenimiento del sistema.

## 1.3. Puntos de acoplamiento identificados

- Base de datos: la columna `book_id` es obligatoria y tiene una llave foránea hacia `books`.
- Modelo: el préstamo utiliza el identificador de un libro.
- Solicitud: el registro de un préstamo necesita identificar un libro.
- Servicio: las reglas de negocio trabajan con la disponibilidad de ejemplares de libros.
- Repositorio: las consultas de préstamos utilizan la tabla de libros.
- Rutas: el módulo actual está organizado alrededor de los préstamos de libros.

## 1.4. Conclusión inicial

El sistema cumple correctamente su función actual de prestar libros, pero presenta limitaciones de extensibilidad cuando se intenta incorporar nuevos tipos de recursos.

Por esta razón, es necesario evaluar diferentes alternativas arquitectónicas que permitan agregar laptops sin afectar el funcionamiento del módulo de libros.