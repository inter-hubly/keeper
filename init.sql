CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS client
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            TEXT        NOT NULL,
    email           TEXT UNIQUE NOT NULL,
    app_id          TEXT        NOT NULL,
    phone_number_id TEXT        NOT NULL,
    business_id     TEXT        NOT NULL,
    access_token    TEXT        NOT NULL,
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
    client_id     INT,
    CONSTRAINT fk_client FOREIGN KEY (client_id) REFERENCES client (id)
);

CREATE TABLE IF NOT EXISTS contact
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(255)        NOT NULL,
    phone         VARCHAR(255) UNIQUE NOT NULL,
    created_at    TIMESTAMP                NOT NULL,
    updated_at    TIMESTAMP                NOT NULL,
    removed        BOOL     DEFAULT false
);

INSERT INTO contact (name, phone, created_at, updated_at) VALUES ('Saimon Ribeiro', '48991784586', CURRENT_DATE, CURRENT_DATE) ON CONFLICT (phone) DO NOTHING;
INSERT INTO contact (name, phone, created_at, updated_at) VALUES ('Fabiane Staziaki', '48988356622', CURRENT_DATE, CURRENT_DATE) ON CONFLICT (phone) DO NOTHING;

INSERT INTO client (name, email, app_id, phone_number_id, business_id,
                    access_token, created_at, updated_at)
VALUES ('Hubly', 'hubly@teste.com', '1121875729525113',
        '459185417288378', '510006955530686',
        'EAAP8VwxXKXkBO9K7ndBh2IminvZCy2WlF4cTjRWWbeKi7CTbf6B60crXVADlO7JGzZCstqoS8WxhKhTwqK8DQxQRz2xe5tbg3hY99SIdzsJslGZBmnhuhRrlZAvJDlO7JGzZCstqoS8WxhKhTwqK8DQxQRz2xe5tbg3hY99SIdzsJslGZBmnhuhRrlZAvJDlO5ZAftf0e',CURRENT_DATE, CURRENT_DATE)
ON CONFLICT (email) DO NOTHING;

INSERT INTO client (name, email, app_id, phone_number_id, business_id,
                    access_token, created_at, updated_at)
VALUES ('Guimarães & Staziaki', 'guimaraes@staziaki.com', '1121875729525113',
        '559153210606318', '525894273941090',
        'EAAP8VwxXKXkBO3WJuZAX2ZBsW2mVzBAHqG7gb38x2fIgd3ydReFv6VVNRZBGnZBsmhOt8CF7R9msJgPSprnjS7dM7g6XQdfvo50dguTVp6NxZAzC6JX5SvDqqBv3CpHX39O1uEAkkCdBrbwE4NKbCBHKVNOC7dtc2OSFQWSnIx2HUbvCTFFbg3lXxFt3j4PQZA',CURRENT_DATE, CURRENT_DATE)
ON CONFLICT (email) DO NOTHING;


INSERT INTO "user" (name, email, password, login_attempt, created_at,
                    updated_at, client_id)
VALUES ('Fabiane', 'fabiane@test.com', '12345', 0, CURRENT_DATE, CURRENT_DATE,
        (SELECT id FROM client WHERE email = 'guimaraes@staziaki.com'))
ON CONFLICT (email) DO NOTHING;

INSERT INTO "user" (name, email, password, login_attempt, created_at,
                    updated_at, client_id)
VALUES ('Saimon', 'saimon@test.com', '12345', 0, CURRENT_DATE, CURRENT_DATE,
        (SELECT id FROM client WHERE email = 'hubly@teste.com'))
ON CONFLICT (email) DO NOTHING;

