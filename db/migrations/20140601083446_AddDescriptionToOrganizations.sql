
-- +goose Up
ALTER TABLE organizations ADD COLUMN description text not null DEFAULT '';


-- +goose Down
ALTER TABLE organizations DROP COLUMN description;

