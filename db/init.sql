CREATE TABLE actors (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(255) NOT NULL,
                        gender VARCHAR(10),
                        birthdate DATE
);

CREATE TABLE movies (
                        id SERIAL PRIMARY KEY,
                        title VARCHAR(150) NOT NULL,
                        description TEXT,
                        release_date DATE,
                        rating FLOAT
);

CREATE TABLE movie_actors (
                              movie_id INT,
                              actor_id INT,
                              PRIMARY KEY (movie_id, actor_id),
                              FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
                              FOREIGN KEY (actor_id) REFERENCES actors(id) ON DELETE CASCADE
);
