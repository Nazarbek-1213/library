CREATE TABLE books (
    id               BIGSERIAL    PRIMARY KEY,
    title            VARCHAR(255) NOT NULL,
    author           VARCHAR(150) NOT NULL,
    total_copies     INT          NOT NULL CHECK (total_copies >= 0),
    available_copies INT          NOT NULL CHECK (available_copies >= 0),
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT books_available_le_total CHECK (available_copies <= total_copies)
);

CREATE INDEX idx_books_title  ON books (lower(title));
CREATE INDEX idx_books_author ON books (lower(author));