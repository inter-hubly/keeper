CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS client
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(255)        NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    app_id          VARCHAR(255)        NOT NULL,
    phone_number_id VARCHAR(255)    UNIQUE    NOT NULL,
    business_id     VARCHAR(255)        NOT NULL,
    access_token    VARCHAR(255)        NOT NULL,
    created_at      TIMESTAMP        NOT NULL,
    updated_at      TIMESTAMP        NOT NULL,
    removed          BOOL DEFAULT false
);

CREATE TABLE IF NOT EXISTS "user"
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(255)        NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password      VARCHAR(255)        NOT NULL,
    login_attempt SMALLINT DEFAULT 0,
    created_at    TIMESTAMP                NOT NULL,
    updated_at    TIMESTAMP                NOT NULL,
    removed        BOOL     DEFAULT false,
    tenant_id   VARCHAR(255)  ,
    CONSTRAINT fk_client FOREIGN KEY (tenant_id) REFERENCES client (phone_number_id)
);

INSERT INTO client (name, email, app_id, phone_number_id, business_id,
                    access_token, created_at, updated_at)
VALUES ('Hubly', 'hubly@teste.com', '1121875729525113',
        '459185417288378', '510006955530686',
        'EAAP8VwxXKXkBO0cdthInNO4einT5cE0UZAibrl6pdZColBjFoJydbIirutBe4rPX4yHEZAEdwEIPZAubugT6EXWCJLqZCgSORn9AWn7V1fFGp3VsDn9AUfVLP7wKcHtYyqW4YIW9ZAQZAJ3lZCgpvzqDZBZAvR0ZCBsEh4zbPoN8C2md7z0MgizAI1nk1sPRBBTkRlMHASCN9Towf1kZAdQe1DwN2d2LeFXfp69VSk9m9PVh',CURRENT_DATE, CURRENT_DATE)
ON CONFLICT (email) DO NOTHING;

INSERT INTO client (name, email, app_id, phone_number_id, business_id,
                    access_token, created_at, updated_at)
VALUES ('Guimarães & Staziaki', 'guimaraes@staziaki.com', '1121875729525113',
        '559153210606318', '525894273941090',
        'EAAP8VwxXKXkBO0cdthInNO4einT5cE0UZAibrl6pdZColBjFoJydbIirutBe4rPX4yHEZAEdwEIPZAubugT6EXWCJLqZCgSORn9AWn7V1fFGp3VsDn9AUfVLP7wKcHtYyqW4YIW9ZAQZAJ3lZCgpvzqDZBZAvR0ZCBsEh4zbPoN8C2md7z0MgizAI1nk1sPRBBTkRlMHASCN9Towf1kZAdQe1DwN2d2LeFXfp69VSk9m9PVh',CURRENT_DATE, CURRENT_DATE)
ON CONFLICT (email) DO NOTHING;


INSERT INTO "user" (name, email, password, login_attempt, created_at,
                    updated_at, tenant_id)
VALUES ('Fabiane', 'fabiane@test.com', '12345', 0, CURRENT_DATE, CURRENT_DATE,
        (SELECT tenant_id FROM client WHERE email = 'guimaraes@staziaki.com'))
ON CONFLICT (email) DO NOTHING;

INSERT INTO "user" (name, email, password, login_attempt, created_at,
                    updated_at, tenant_id)
VALUES ('Saimon', 'saimon@test.com', '12345', 0, CURRENT_DATE, CURRENT_DATE,
        (SELECT tenant_id FROM client WHERE email = 'hubly@teste.com'))
ON CONFLICT (email) DO NOTHING;
