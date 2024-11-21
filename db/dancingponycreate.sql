-- restaurants
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
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Last update timestamp
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE
);

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

CREATE TABLE dish_images (
    id SERIAL PRIMARY KEY,
    dish_id INTEGER NOT NULL,
    filename VARCHAR(255) NOT NULL,
    content BYTEA NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dish_id) REFERENCES restaurant_dishes(id) ON DELETE CASCADE
);

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- users
CREATE TABLE users (
    id SERIAL PRIMARY KEY, -- Unique identifier for each user
    name VARCHAR(255) NOT NULL, -- User's full name
    email VARCHAR(255) NOT NULL UNIQUE, -- User's email, must be unique
    role VARCHAR(50) NOT NULL, -- User's role (e.g., admin, manager, chef, etc.)
    restaurant_id INTEGER NOT NULL, -- Restaurant ID the user is associated with (0 for system-wide roles)
    hashed_password TEXT NOT NULL, -- Hashed version of the user's password
    salt TEXT NOT NULL, -- Salt used for password hashing
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of account creation
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Timestamp of the last update
);

INSERT INTO users (name, email, role, restaurant_id, hashed_password, salt)
VALUES 
    ('Frans Slabber', 'byeboer@gmail.com', 'admin', NULL, crypt('frans', gen_salt('bf')), gen_salt('bf')),
    ('Frodo Slabber', 'frodo@gmail.com', 'admin', 1, crypt('frodo', gen_salt('bf')), gen_salt('bf'));

-- Customer reviews
CREATE TABLE restaurant_reviews (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    review TEXT NOT NULL,
    rating NUMERIC(3, 2) CHECK (rating >= 1.00 AND rating <= 5.00), -- Rating from 1.00 to 5.00
    sentiment_score NUMERIC,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE    
);

INSERT INTO restaurant_reviews (restaurant_id, user_id, review, rating, sentiment_score) 
    VALUES
        ((select id from restaurants where path_name = 'OrcShack'), 2, 'Situated inconspicuously above Van Rensburgs - every carnivore favourite meat emporium - Kafe Serefe has access to the best quality meat in town.  Their creamy beef stroganoff is a huge favorite  and their generous helpings allow for a doggie bag evening snack.  We do miss dear Erica Arries, the warm and engaging manager who passed away too soon.', 3.30, 0.00 );
        ((select id from restaurants where path_name = 'OrcShack'), 2, 'Went on a Monday night thinking we wouldnt have to book - it was full! But they found us a table in 5 mins...the food was great. Massive portions at a VERY good price. Woukd go again and recommend!', 3.30, 0.00 );

-- User rating for dish
CREATE TABLE user_dish_ratings (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER NOT NULL,
    dish_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    rating NUMERIC(3, 2) CHECK (rating >= 1.00 AND rating <= 5.00), -- Rating from 1.00 to 5.00
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
    FOREIGN KEY (dish_id) REFERENCES restaurant_dishes(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
