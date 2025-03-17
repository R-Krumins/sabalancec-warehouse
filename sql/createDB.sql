-- Enable foreign key support
PRAGMA foreign_keys = ON;

-- Create Users Table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    address TEXT,
    phone TEXT,
    role TEXT
);

-- Create Products Table
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    price REAL NOT NULL,
    category TEXT,
    image TEXT,
    amount_sold INTEGER DEFAULT 0,
    amount_in_stock INTEGER NOT NULL,
    has_allergens INTEGER DEFAULT 0, -- Boolean replacement (0 = false, 1 = true)
    rating REAL DEFAULT 0
);

-- Create Allergens Table
CREATE TABLE allergens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    image TEXT,
	info TEXT
);

-- Create Product-Allergens Relationship Table
CREATE TABLE product_allergens (
    product_id INTEGER,
    allergen_id INTEGER,
    PRIMARY KEY (product_id, allergen_id),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (allergen_id) REFERENCES allergens(id) ON DELETE CASCADE
);

-- Create Cards Table
CREATE TABLE cards (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    card_number TEXT NOT NULL,
    expiry_date TEXT NOT NULL,
    cvv TEXT NOT NULL,
    card_type TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create Carts Table
CREATE TABLE carts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    product_id INTEGER,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Create Payment Methods Table
CREATE TABLE payment_methods (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

-- Create Transactions Table
CREATE TABLE transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    cart_id INTEGER,
    payment_id INTEGER,
    total_cost REAL NOT NULL CHECK (total_cost >= 0),
    status TEXT NOT NULL,
    payment_method_id INTEGER NOT NULL,
    transaction_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    delivery_status TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE,
    FOREIGN KEY (payment_id) REFERENCES cards(id) ON DELETE SET NULL,
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);

