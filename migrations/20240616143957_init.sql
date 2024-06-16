-- +goose Up
-- +goose StatementBegin
CREATE TABLE images (
    entity_id VARCHAR(25),
    entity_type VARCHAR(25),
    image_url VARCHAR(255),
    content_type VARCHAR(50),
    filename VARCHAR(255),
    format integer
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE images 
    ADD CONSTRAINT uq_entity_image UNIQUE (entity_id, entity_type, content_type, format);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE images;
-- +goose StatementEnd
