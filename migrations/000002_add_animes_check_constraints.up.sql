ALTER TABLE animes ADD CONSTRAINT animes_runtime_check CHECK (runtime >= 0);
ALTER TABLE animes ADD CONSTRAINT animes_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
ALTER TABLE animes ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);
