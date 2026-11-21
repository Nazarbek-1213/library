CREATE TABLE loans (
    id          BIGSERIAL   PRIMARY KEY,

    member_id   BIGINT      NOT NULL REFERENCES members (id) ON DELETE RESTRICT,
    book_id     BIGINT      NOT NULL REFERENCES books (id)   ON DELETE RESTRICT,
    taken_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    due_at      TIMESTAMPTZ NOT NULL,
    returned_at TIMESTAMPTZ,                              -- NULL = kitob hali qaytarilmagan

    CONSTRAINT loans_due_after_taken      CHECK (due_at > taken_at),
    CONSTRAINT loans_returned_after_taken CHECK (returned_at IS NULL OR returned_at >= taken_at)
);

CREATE UNIQUE INDEX uq_loans_active_member_book
    ON loans (member_id, book_id)
    WHERE returned_at IS NULL;

CREATE INDEX idx_loans_member_active ON loans (member_id) WHERE returned_at IS NULL;

CREATE INDEX idx_loans_due_active    ON loans (due_at)    WHERE returned_at IS NULL;

CREATE INDEX idx_loans_book          ON loans (book_id);