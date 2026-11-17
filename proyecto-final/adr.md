# ADR – Decisión de arquitectura para el módulo de préstamos de laptops

## 1. Contexto

El proyecto Biblioteca API contaba inicialmente con un sistema de gestión de libros y préstamos. Como parte del proyecto final de Arquitectura de Software, se requiere ampliar el sistema para administrar un nuevo recurso: las laptops.

La solución debe permitir registrar laptops, consultar el inventario disponible, realizar préstamos y registrar devoluciones. También debe conservar el funcionamiento del módulo de libros y respetar los permisos de los usuarios.

La aplicación utiliza Go, Gin para la gestión de rutas HTTP, MySQL como base de datos y una arquitectura por capas:

- Handler: recibe las solicitudes HTTP y devuelve las respuestas.
- Service: contiene las reglas del negocio.
- Repository: realiza las operaciones sobre la base de datos.
- Models: representan las entidades y estructuras de datos.

## 2. Problema arquitectónico

Se necesita incorporar el préstamo de laptops sin afectar el sistema de préstamos de libros que ya existe.

La decisión principal consiste en determinar cómo representar las laptops y sus préstamos en la base de datos y cómo organizar su lógica dentro de la arquitectura actual.

## 3. Alternativas evaluadas

### Alternativa 1. Crear un módulo independiente para laptops

Consiste en implementar tablas, modelos, repositorios, servicios y handlers específicos para las laptops.

Se utilizan las tablas `laptops` y `laptop_loans`, independientes de las tablas utilizadas para los libros y sus préstamos.

**Ventajas:**
- No modifica la estructura actual de los préstamos de libros.
- Mantiene separadas las responsabilidades de cada módulo.
- Permite establecer relaciones mediante claves foráneas.
- Facilita realizar cambios futuros en el módulo de laptops.
- Reduce el riesgo de afectar funcionalidades existentes.

**Desventajas:**
- Puede generar código repetido entre los módulos.
- Requiere mantener dos estructuras de préstamos.
- Aumenta la cantidad de archivos y componentes del proyecto.

### Alternativa 2. Unificar los préstamos en una sola tabla

Consiste en utilizar una única tabla de préstamos para libros y laptops, agregando un campo que identifique el tipo de recurso solicitado.

**Ventajas:**
- Reduce la duplicación de algunas operaciones.
- Centraliza la consulta de préstamos.
- Puede facilitar la creación de reportes generales.

**Desventajas:**
- Aumenta la complejidad de las consultas.
- Requiere manejar diferentes tipos de recursos en una misma estructura.
- Puede dificultar la integridad referencial, porque una misma columna tendría que representar recursos de tablas diferentes.
- Implica modificar la lógica existente de préstamos de libros.

### Alternativa 3. Crear una entidad general de recursos

Consiste en diseñar una estructura común para representar libros y laptops como recursos prestables. Los detalles particulares de cada recurso se almacenarían en estructuras relacionadas.

**Ventajas:**
- Permite compartir reglas generales de préstamo.
- Facilita incorporar nuevos tipos de recursos en el futuro.
- Reduce parte de la duplicación de lógica.

**Desventajas:**
- Requiere una modificación más amplia de la arquitectura.
- Aumenta el tiempo de implementación.
- Puede introducir complejidad innecesaria para el alcance del proyecto.
- Exige adaptar el módulo de libros que ya funciona.

## 4. Decisión seleccionada

Se seleccionó la **Alternativa 1: crear un módulo independiente para laptops**.

La decisión se tomó porque permite ampliar el sistema sin modificar directamente la lógica de préstamos de libros. Además, conserva la separación por capas y permite establecer relaciones mediante claves foráneas entre usuarios, laptops y préstamos.

Se implementaron los siguientes componentes:

- `internal/models/laptop.go`
- `internal/repository/laptop_repository.go`
- `internal/service/laptop_service.go`
- `internal/handler/laptop_handler.go`
- `migrations/002_laptops.sql`

También se actualizaron las rutas y la inicialización de dependencias para incorporar el nuevo módulo.

## 5. Reglas de negocio

1. Solamente los usuarios autenticados pueden consultar las laptops disponibles.
2. Solamente los administradores pueden registrar laptops.
3. Solamente los administradores pueden consultar el inventario completo.
4. Una laptop solamente puede prestarse cuando su estado es `available`.
5. Una laptop prestada cambia su estado a `loaned`.
6. Cada préstamo tiene una duración de 10 días.
7. Un usuario no puede tener más de un préstamo activo de la misma laptop.
8. Solamente el usuario propietario del préstamo puede realizar su devolución.
9. Una laptop devuelta cambia nuevamente su estado a `available`.
10. Una laptop en mantenimiento no puede prestarse.

## 6. Consideraciones de integridad y concurrencia

Las operaciones de préstamo y devolución se implementaron utilizando transacciones de MySQL.

En el préstamo, la transacción verifica el estado de la laptop, registra el préstamo y actualiza el inventario. De esta manera, las operaciones se confirman como una unidad.

En la devolución, la transacción verifica que el préstamo exista, que pertenezca al usuario que realiza la solicitud y que todavía esté activo. Después registra la devolución y actualiza el estado de la laptop.

También se utiliza un bloqueo de fila mediante `SELECT ... FOR UPDATE` para evitar que dos solicitudes concurrentes presten la misma laptop al mismo tiempo.

## 7. Consecuencias de la decisión

**Consecuencias positivas:**
- Se conserva el funcionamiento del módulo de libros.
- Las responsabilidades de cada módulo quedan separadas.
- Se facilita el mantenimiento del código.
- Se mantiene la integridad referencial en la base de datos.
- Las operaciones de préstamo y devolución cuentan con control transaccional.

**Consecuencias negativas:**
- Existe cierta duplicación de código entre los módulos de libros y laptops.
- Si en el futuro se agregan muchos tipos de recursos, podría ser conveniente evaluar una arquitectura más general.

## 8. Conclusión

La creación de un módulo independiente para laptops es la alternativa más adecuada para el alcance del proyecto. Permite incorporar nuevas funcionalidades sin realizar una reestructuración innecesaria del sistema existente.

La decisión prioriza la separación de responsabilidades, la integridad de los datos, la seguridad y la facilidad de mantenimiento, conservando la arquitectura que ya utiliza Biblioteca API.