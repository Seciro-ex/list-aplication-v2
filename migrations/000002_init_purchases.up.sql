CREATE TABLE purchases (
	id SERIAL PRIMARY KEY,
	image_url TEXT,
	cost DECIMAL(20,2) NOT NULL,
	item_url TEXT,
	item_name TEXT NOT NULL,
	item_description TEXT,
	user_id INT REFERENCES users(id) ON DELETE CASCADE
);