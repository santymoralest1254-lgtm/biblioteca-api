# Verificación y pruebas – Biblioteca API

## 1. Objetivo

Documentar las pruebas realizadas para comprobar el funcionamiento del módulo de préstamos de laptops y verificar que las funcionalidades existentes de libros continúen operando correctamente.

## 2. Entorno de pruebas

- Lenguaje: Go.
- Framework HTTP: Gin.
- Base de datos: MySQL 8.
- Contenedores: Docker Compose.
- Puerto de la API: 8080.
- Puerto de MySQL: 3307.
- Autenticación: JWT.
- Herramienta de pruebas: PowerShell mediante solicitudes HTTP.

## 3. Verificación de ejecución

| Prueba | Resultado |
|---|---|
| Construcción de la imagen Docker | Exitosa |
| Inicio del contenedor de la API | Exitoso |
| Estado del contenedor MySQL | Healthy |
| Solicitud GET /health | Respondió status: ok |

## 4. Pruebas funcionales del módulo de laptops

### 4.1. Consulta de laptops disponibles

Se realizó una solicitud autenticada al endpoint:

`GET /api/laptops`

**Resultado:** Correcto. El endpoint permitió consultar las laptops disponibles.

### 4.2. Registro de una laptop

Se utilizó una cuenta administradora para registrar una laptop mediante:

`POST /api/laptops`

**Resultado:** Correcto. Se registró una laptop de prueba con identificador 1, marca Lenovo y número de serie TEST-LAP-001.

### 4.3. Préstamo de una laptop

Se utilizó una cuenta de estudiante para solicitar la laptop con identificador 1:

`POST /api/laptop-loans`

**Resultado:** Correcto. Se creó el préstamo con identificador 2, usuario 3 y estado `active`.

La fecha límite de devolución quedó establecida a 10 días después de la fecha del préstamo.

### 4.4. Actualización del inventario

Después de registrar el préstamo, se consultó nuevamente:

`GET /api/laptops`

**Resultado:** Correcto. La laptop prestada dejó de aparecer en el catálogo de disponibles.

### 4.5. Prevención de préstamos duplicados

Se intentó solicitar nuevamente la laptop que ya estaba prestada.

**Resultado:** Correcto. El sistema rechazó la operación y respondió:

`esta laptop no está disponible`

### 4.6. Devolución de laptop

Se ejecutó:

`PUT /api/laptop-loans/2/return`

**Resultado:** Correcto. El sistema registró la devolución del préstamo.

### 4.7. Restauración del inventario

Después de la devolución se consultó nuevamente el catálogo.

**Resultado:** Correcto. La laptop volvió a aparecer como disponible.

## 5. Pruebas de seguridad y autorización

Se verificaron los permisos de los usuarios.

| Prueba | Resultado esperado | Resultado obtenido |
|---|---|---|
| Consultar laptops sin autenticación | HTTP 401 | HTTP 401 |
| Estudiante intenta registrar una laptop | HTTP 403 | HTTP 403 |
| Estudiante intenta consultar inventario administrativo | HTTP 403 | HTTP 403 |

Estas pruebas permiten comprobar que las operaciones administrativas están restringidas a usuarios con rol de administrador.

## 6. Pruebas de regresión del módulo de libros

Se realizaron pruebas para verificar que la incorporación del módulo de laptops no afectara las funcionalidades existentes.

| Prueba | Resultado |
|---|---|
| Consultar catálogo de libros | Correcto |
| Registrar préstamo de libro | Correcto |
| Disminuir unidades disponibles | Correcto |
| Devolver libro | Correcto |
| Restaurar unidades disponibles | Correcto |

El préstamo y la devolución de libros continuaron funcionando después de incorporar el nuevo módulo.

## 7. Verificación de transacciones

Se comprobó el funcionamiento de las operaciones transaccionales de préstamo y devolución de laptops.

En el préstamo se utiliza una transacción para verificar la disponibilidad, insertar el préstamo y actualizar el estado de la laptop.

En la devolución se verifica la existencia del préstamo y su propietario, se actualiza el estado del préstamo y se restablece la disponibilidad de la laptop.

También se utiliza `SELECT ... FOR UPDATE` para bloquear la fila de la laptop durante la operación de préstamo.

## 8. Conclusión

Las pruebas realizadas permiten confirmar el funcionamiento de las operaciones principales del módulo de laptops: consulta, préstamo, control de disponibilidad y devolución.

También se verificaron los permisos por roles y se comprobó que el módulo de libros conserva sus operaciones principales.

La implementación amplía las capacidades de Biblioteca API manteniendo la separación de responsabilidades y el control de integridad de los datos.