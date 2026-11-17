
-- Tabla de laptops disponibles en el inventario
CREATE TABLE IF NOT EXISTS laptops (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    brand VARCHAR(80) NOT NULL,
    serial_number VARCHAR(100) NOT NULL UNIQUE,
    status ENUM('available', 'loaned', 'maintenance')
        NOT NULL DEFAULT 'available',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla independiente para los préstamos de laptops
CREATE TABLE IF NOT EXISTS laptop_loans (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    laptop_id INT NOT NULL,
    loaned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_at TIMESTAMP NOT NULL,
    returned_at TIMESTAMP NULL,
    status ENUM('active', 'returned') NOT NULL DEFAULT 'active',
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (laptop_id) REFERENCES laptops(id)
);

-- Índices para facilitar las consultas de préstamos
CREATE INDEX idx_laptop_loans_user
    ON laptop_loans(user_id);

CREATE INDEX idx_laptop_loans_laptop
    ON laptop_loans(laptop_id);
