-- Terminate existing connections to the database
DO
$$
BEGIN
    IF EXISTS (
        SELECT FROM pg_database WHERE datname = 'dancingpony'
    ) THEN
        -- Terminate connections to the database
        PERFORM pg_terminate_backend(pg_stat_activity.pid)
        FROM pg_stat_activity
        WHERE pg_stat_activity.datname = 'dancingpony'
          AND pid <> pg_backend_pid();
        
        -- Drop the database
        EXECUTE 'DROP DATABASE dancingpony';
    END IF;
END
$$;

-- Create the new database
CREATE DATABASE dancingpony;

\c dancingpony

-- tenants/restaurants
CREATE TABLE restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    path_name VARCHAR(255) NOT NULL,
    location VARCHAR(255)
);

INSERT INTO restaurants (name,path_name, location)
VALUES 
    ('The Orc Shack', 'OrcShack', 'The Docks'),
    ('Meats Back On The Menu', 'MeatyMenu', 'Orc City');

-- restaurant dishes
CREATE TABLE restaurant_dishes (
    id SERIAL PRIMARY KEY, -- Auto-incrementing integer ID
    name VARCHAR(255) NOT NULL, -- Dish name
    description TEXT, -- Detailed description of the dish
    price NUMERIC(10, 2) NOT NULL, -- Price with two decimal places
    category VARCHAR(100), -- Category like Appetizer, Dessert
    is_vegetarian BOOLEAN DEFAULT FALSE, -- Indicates if the dish is vegetarian
    is_available BOOLEAN DEFAULT TRUE, -- Indicates availability
    rating NUMERIC(3, 2) CHECK (rating >= 1.00 AND rating <= 5.00), -- Rating from 1.00 to 5.00
    restaurant_id INTEGER NOT NULL, -- ID of the restaurant the dish belongs to
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Last update timestamp
);

ALTER TABLE restaurant_dishes
ADD CONSTRAINT fk_restaurant
FOREIGN KEY (restaurant_id)
REFERENCES restaurants (id)
ON DELETE CASCADE;

-- CREATE OR REPLACE FUNCTION update_updated_at_column()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     NEW.updated_at = CURRENT_TIMESTAMP;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER set_updated_at
-- BEFORE UPDATE ON restaurant_dishes
-- FOR EACH ROW
-- EXECUTE FUNCTION update_updated_at_column();

-- Dishes for Restaurant 1
INSERT INTO restaurant_dishes (name, description, price, category, is_vegetarian, is_available, rating, restaurant_id)
VALUES 
    ('Spaghetti Carbonara', 'Creamy pasta with pancetta and Parmesan.', 12.99, 'Main Course', FALSE, TRUE, 4.80, 1),
    ('Garlic Bread', 'Crispy bread with garlic butter.', 3.99, 'Appetizer', TRUE, TRUE, 4.50, 1),
    ('Tiramisu', 'Classic Italian dessert with coffee and mascarpone.', 5.99, 'Dessert', TRUE, TRUE, 4.90, 1),
    ('Fettuccine Alfredo', 'Fettuccine in a creamy white sauce.', 11.99, 'Main Course', TRUE, TRUE, 4.70, 1),
    ('Minestrone Soup', 'Hearty Italian vegetable soup.', 6.50, 'Appetizer', TRUE, FALSE, 4.60, 1);

-- Dishes for Restaurant 2
INSERT INTO restaurant_dishes (name, description, price, category, is_vegetarian, is_available, rating, restaurant_id)
VALUES 
    ('Pepperoni Pizza', 'Classic pizza with pepperoni and mozzarella.', 9.99, 'Main Course', FALSE, TRUE, 4.85, 2),
    ('Veggie Supreme Pizza', 'Loaded with vegetables and cheese.', 10.99, 'Main Course', TRUE, TRUE, 4.65, 2),
    ('Garlic Knots', 'Soft dough knots with garlic and herbs.', 4.50, 'Appetizer', TRUE, TRUE, 4.70, 2),
    ('Cheesecake', 'Rich and creamy cheesecake.', 6.99, 'Dessert', TRUE, TRUE, 4.95, 2),
    ('BBQ Chicken Pizza', 'Pizza topped with BBQ sauce and chicken.', 11.50, 'Main Course', FALSE, TRUE, 4.75, 2);


CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- users
CREATE TABLE users (
    id SERIAL PRIMARY KEY, -- Unique identifier for each user
    name VARCHAR(255) NOT NULL, -- User's full name
    email VARCHAR(255) NOT NULL UNIQUE, -- User's email, must be unique
    role VARCHAR(50) NOT NULL, -- User's role (e.g., admin, manager, chef, etc.)
    restaurant_id INTEGER, -- Restaurant ID the user is associated with (nullable for system-wide roles)
    hashed_password TEXT NOT NULL, -- Hashed version of the user's password
    salt TEXT NOT NULL, -- Salt used for password hashing
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of account creation
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Timestamp of the last update
);

INSERT INTO users (name, email, role, restaurant_id, hashed_password, salt)
VALUES 
    ('Frans Slabber', 'frans@byeboer@gmail.com', 'admin', NULL, crypt('secure_password_1', gen_salt('bf')), gen_salt('bf'));
    -- ('Bob Chef', 'bob@pizzapalace.com', 'chef', 1, crypt('secure_password_2', gen_salt('bf')), gen_salt('bf')),
    -- ('Charlie Admin', 'charlie@admin.com', 'admin', NULL, crypt('secure_password_3', gen_salt('bf')), gen_salt('bf')),
    -- ('Diana Staff', 'diana@pastaheaven.com', 'staff', 2, crypt('secure_password_4', gen_salt('bf')), gen_salt('bf'));


-- SELECT * 
-- FROM users
-- WHERE email = 'alice@pizzapalace.com'
--   AND hashed_password = crypt('secure_password_1', hashed_password);
