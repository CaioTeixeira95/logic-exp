-- +migrate Up

CREATE TABLE public.expressions (
    id SERIAL PRIMARY KEY,
    expression text
);

-- +migrate Down

DROP TABLE public.expressions;
