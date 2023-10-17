-- Membuat table inventori
CREATE TABLE inventories (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    name TEXT NOT NULL,
    item_code TEXT NOT NULL UNIQUE,
    stock INTEGER NOT NULL,
    description TEXT,
    status TEXT CHECK (status IN ('active', 'broken'))
);

-- Menambahkan data inventori
INSERT INTO inventories (name, item_code, stock, description, status) VALUES
    ('Item 1', 'ITEM001', 100, 'Deskripsi Item 1', 'active'),
    ('Item 2', 'ITEM002', 50, 'Deskripsi Item 2', 'active'),
    ('Item 3', 'ITEM003', 75, 'Deskripsi Item 3', 'broken');
