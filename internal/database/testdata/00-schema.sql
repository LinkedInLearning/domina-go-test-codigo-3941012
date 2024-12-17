CREATE TABLE pokemons (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    hp INT NOT NULL,
    attack INT NOT NULL,
    defense INT NOT NULL
);

CREATE TABLE battles (
    id SERIAL PRIMARY KEY,
    pokemon1_id INT NOT NULL,
    pokemon2_id INT NOT NULL,
    winner_id INT NOT NULL,
    turns INT NOT NULL,
    FOREIGN KEY (pokemon1_id) REFERENCES pokemons (id),
    FOREIGN KEY (pokemon2_id) REFERENCES pokemons (id),
    FOREIGN KEY (winner_id) REFERENCES pokemons (id)
);
