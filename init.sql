CREATE TABLE IF NOT EXISTS clients (
                                       id BIGSERIAL PRIMARY KEY,
                                       name TEXT NOT NULL,
                                       email TEXT UNIQUE NOT NULL,
                                       app_id TEXT NOT NULL,
                                       phone_number_id TEXT NOT NULL,
                                       business_id TEXT NOT NULL,
                                       access_token TEXT NOT NULL
);

INSERT INTO clients (name, email, app_id, phone_number_id, business_id, access_token)
VALUES (
           'associados',
           'associados@teste.com',
           '1121875729525113',
           '515719138282305',
           '510006955530686',
           'EAAP8VwxXKXkBO9K7ndBh2IminvZCy2WlF4cTjRWWbeKi7CTbf6B60crXVADlO7JGzZCstqoS8WxhKhTwqK8DQxQRz2xe5tbg3hY99SIdzsJslGZBmnhuhRrlZAvJD5ZCErgVfAUZAH4OtBn4SZBcUbskcB25sx2qZAlebJuZAERzkSSZAgcN3qDChsQ1ZAfyajkfwcruRNJBIxFXDrggMdNNmbl3oO2'
       )
ON CONFLICT (email) DO NOTHING;

