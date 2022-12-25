CREATE TABLE IF NOT EXISTS players (
    id varchar(64) CONSTRAINT playerid PRIMARY KEY,
    name varchar(128),
    rating numeric(10, 6),
    num_games integer
);

CREATE TABLE IF NOT EXISTS games (
    id varchar(64) CONSTRAINT gameid PRIMARY KEY,
    blue_won boolean,
    blue_win_expectation numeric(8, 6)
);

CREATE TABLE IF NOT EXISTS players_by_game (
    game_id varchar(64),
    player_id varchar(64),
    player_rating numeric(10, 6),
    player_num_games integer,
    player_role varchar(16),
    PRIMARY KEY (game_id, player_id)
);
