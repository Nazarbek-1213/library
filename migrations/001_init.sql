CREATE TABLE members (
    id            BIGSERIAL    PRIMARY KEY,
    full_name     VARCHAR(100) NOT NULL,
    phone         VARCHAR(20)  NOT NULL UNIQUE,        
    password_hash TEXT         NOT NULL,                 
    role          VARCHAR(10)  NOT NULL DEFAULT 'member'
                  CHECK (role IN ('admin', 'member')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);