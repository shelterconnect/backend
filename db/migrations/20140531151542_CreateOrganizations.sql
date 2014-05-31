
-- +goose Up
CREATE TABLE organizations (
  id serial not null primary key,
  created timestamp not null,
  updated timestamp not null,
  name text not null,
  email text not null unique,
  type integer not null,
  address text not null,
  location point not null,
  password text not null
);


-- +goose Down
DROP TABLE organizations;

