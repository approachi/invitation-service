CREATE TABLE invitations (
    id         BIGSERIAL    PRIMARY KEY,
    code       VARCHAR(255) NOT NULL UNIQUE,
    max_uses   INT          NOT NULL,
    used_count INT          NOT NULL
);

CREATE TABLE registrations (
    id            BIGSERIAL    PRIMARY KEY,
    invitation_id BIGINT       NOT NULL REFERENCES invitations(id),
    email         VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    UNIQUE (invitation_id, email)
);
