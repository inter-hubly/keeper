CREATE TABLE IF NOT EXISTS client
(
    id              BIGSERIAL PRIMARY KEY,
    name            TEXT        NOT NULL,
    email           TEXT UNIQUE NOT NULL,
    app_id          TEXT        NOT NULL,
    phone_number_id TEXT        NOT NULL,
    business_id     TEXT        NOT NULL,
    access_token    TEXT        NOT NULL
);

CREATE TABLE IF NOT EXISTS "user"
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255)        NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password      VARCHAR(255)        NOT NULL,
    login_attempt SMALLINT DEFAULT 0,
    created_at    DATE                NOT NULL,
    updated_at    DATE                NOT NULL,
    client_id     INT                 NOT NULL,
    CONSTRAINT fk_client FOREIGN KEY (client_id) REFERENCES client (id)
);




INSERT INTO client (name, email, app_id, phone_number_id, business_id,
                    access_token)
VALUES ('associados', 'associados@teste.com', '1121875729525113',
        '515719138282305', '510006955530686',
        'EAAP8VwxXKXkBO9K7ndBh2IminvZCy2WlF4cTjRWWbeKi7CTbf6B60crXVADlO7JGzZCstqoS8WxhKhTwqK8DQxQRz2xe5tbg3hY99SIdzsJslGZBmnhuhRrlZAvJDlO7JGzZCstqoS8WxhKhTwqK8DQxQRz2xe5tbg3hY99SIdzsJslGZBmnhuhRrlZAvJDlO5ZAftf0e'
        )
ON CONFLICT (email) DO NOTHING;

INSERT INTO client (name, email, app_id, phone_number_id, business_id,
                    access_token)
VALUES ('Guimarães & Staziaki', 'guimaraes@staziaki.com', '1121875729525113',
        '559153210606318', '525894273941090',
        'EAAP8VwxXKXkBO3WJuZAX2ZBsW2mVzBAHqG7gb38x2fIgd3ydReFv6VVNRZBGnZBsmhOt8CF7R9msJgPSprnjS7dM7g6XQdfvo50dguTVp6NxZAzC6JX5SvDqqBv3CpHX39O1uEAkkCdBrbwE4NKbCBHKVNOC7dtc2OSFQWSnIx2HUbvCTFFbg3lXxFt3j4PQZA'
        )
ON CONFLICT (email) DO NOTHING;


INSERT INTO "user" (name, email, password, login_attempt, created_at,
                    updated_at, client_id)
VALUES ('userTest', 'user@test.com', '12345', 0, CURRENT_DATE, CURRENT_DATE,(SELECT id FROM client WHERE email = 'guimaraes@staziaki.com'))
ON CONFLICT (email) DO NOTHING;