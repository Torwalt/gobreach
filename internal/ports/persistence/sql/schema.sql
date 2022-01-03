CREATE TABLE Breach (
	email_hash BYTEA NOT NULL,
	domain VARCHAR(255) NOT NULL,
	breached_info VARCHAR(255),
	breach_date TIMESTAMPTZ,
	breach_source VARCHAR(255),

	CONSTRAINT unique_email UNIQUE(email_hash)
);

