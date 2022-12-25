CREATE TABLE IF NOT EXISTS players (
    id varchar(64) CONSTRAINT primary_id_players_id PRIMARY KEY,
    name varchar(128),
    rating numeric(10, 6),
    num_games integer
);

CREATE TABLE IF NOT EXISTS games (
    id varchar(64) CONSTRAINT primary_id_games_id PRIMARY KEY,
    blue_won boolean,
);

CREATE TABLE IF NOT EXISTS players_by_game (
    game_id varchar(64),
    player_id varchar(64),
    player_name varchar(128),
    player_rating numeric(10, 6),
    player_num_games integer,
    player_role varchar(16),
    CONSTRAINT primary_id_players_by_game_game_id_player_id PRIMARY KEY (game_id, player_id)
);
