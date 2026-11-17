# Diagramas de arquitectura – Biblioteca API

## 1. Diagrama C4 – Nivel 1: Contexto del sistema

Este diagrama representa el sistema Biblioteca API y las personas que interactúan con él.

```mermaid
flowchart TB
    Admin["Administrador"]
    Student["Estudiante"]
    API["Biblioteca API<br/>Sistema de gestión de libros y laptops"]
    DB[("MySQL<br/>Base de datos")]

    Admin -->|"Gestiona libros y laptops"| API
    Student -->|"Consulta recursos y solicita préstamos"| API
    API -->|"Consulta y modifica información"| DB
```

## 2. Diagrama C4 – Nivel 2: Contenedores

Representa los principales contenedores tecnológicos que componen la aplicación.

```mermaid
flowchart TB
    User["Administrador / Estudiante"]
    API["API REST<br/>Go + Gin"]
    DB[("Base de datos<br/>MySQL 8")]
    
    User -->|"HTTP / JSON"| API
    API -->|"SQL mediante sqlx"| DB
```

## 3. Diagrama C4 – Nivel 3: Componentes

Representa la organización interna del módulo de laptops y su comunicación con las demás capas.

```mermaid
flowchart TB
    Client["Cliente HTTP"]

    subgraph API["Biblioteca API"]
        Router["Routes"]
        Middleware["AuthRequired / AdminRequired"]
        Handler["LaptopHandler"]
        Service["LaptopService"]
        Repository["LaptopRepository"]
        Models["Laptop Models"]
    end

    DB[("MySQL<br/>laptops y laptop_loans")]

    Client --> Router
    Router --> Middleware
    Middleware --> Handler
    Handler --> Service
    Service --> Repository
    Handler -.-> Models
    Repository -.-> Models
    Repository --> DB
```

## 4. Diagrama de secuencia – Préstamo de laptop

Representa el proceso que ocurre cuando un estudiante solicita una laptop.

```mermaid
sequenceDiagram
    actor Estudiante
    participant H as LaptopHandler
    participant S as LaptopService
    participant R as LaptopRepository
    participant DB as MySQL

    Estudiante->>H: POST /api/laptop-loans
    H->>S: Borrow(userID, laptopID)
    S->>R: BorrowAtomically(userID, laptopID, dueAt)
    R->>DB: Iniciar transacción
    R->>DB: Consultar laptop con FOR UPDATE
    DB-->>R: Estado de la laptop

    alt Laptop disponible
        R->>DB: Verificar préstamo activo
        R->>DB: Insertar préstamo
        R->>DB: Actualizar laptop a loaned
        R->>DB: Confirmar transacción
        R-->>S: Préstamo registrado
        S-->>H: Datos del préstamo
        H-->>Estudiante: HTTP 200 / respuesta JSON
    else Laptop no disponible
        R-->>S: Error de disponibilidad
        S-->>H: Error de negocio
        H-->>Estudiante: Respuesta de error
    end
```

## 5. Diagrama de secuencia – Devolución de laptop

Representa el proceso de devolución de una laptop prestada.

```mermaid
sequenceDiagram
    actor Estudiante
    participant H as LaptopHandler
    participant S as LaptopService
    participant R as LaptopRepository
    participant DB as MySQL

    Estudiante->>H: PUT /api/laptop-loans/{id}/return
    H->>S: Return(userID, loanID)
    S->>R: ReturnAtomically(userID, loanID)
    R->>DB: Iniciar transacción
    R->>DB: Consultar préstamo con FOR UPDATE
    DB-->>R: Datos del préstamo

    alt Préstamo válido y pertenece al usuario
        R->>DB: Actualizar préstamo a returned
        R->>DB: Actualizar laptop a available
        R->>DB: Confirmar transacción
        R-->>S: Devolución exitosa
        S-->>H: Resultado
        H-->>Estudiante: HTTP 200 / devolución correcta
    else Préstamo inexistente, ajeno o devuelto
        R-->>S: Error de validación
        S-->>H: Error de negocio
        H-->>Estudiante: Respuesta de error
    end
```

## 6. Explicación general

La arquitectura mantiene una separación por capas. Las solicitudes ingresan por las rutas HTTP y pasan por los middlewares de autenticación y autorización.

Después llegan al handler, que recibe y valida los datos de la solicitud. El service aplica las reglas de negocio y el repository ejecuta las operaciones sobre MySQL.

Las transacciones permiten que el registro del préstamo y la actualización del inventario se realicen de manera consistente. Si una operación falla antes