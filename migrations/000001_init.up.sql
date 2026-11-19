-- 000001_init: skema awal LaporKampus.
-- Semua waktu memakai TIMESTAMPTZ agar perbandingan SLA tidak bergantung zona.

CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(30)  NOT NULL UNIQUE,
    email       VARCHAR(120) NOT NULL UNIQUE,
    password    TEXT         NOT NULL,
    role        VARCHAR(20)  NOT NULL DEFAULT 'mahasiswa'
                CHECK (role IN ('mahasiswa', 'agen', 'supervisor', 'admin')),
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          BIGSERIAL PRIMARY KEY,
    user_id     INTEGER     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash  TEXT        NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS refresh_tokens_user_id_idx
    ON refresh_tokens (user_id);

CREATE TABLE IF NOT EXISTS categories (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL DEFAULT '',
    -- Target penyelesaian dalam jam; dasar perhitungan SLA.
    sla_hours   INTEGER      NOT NULL DEFAULT 24 CHECK (sla_hours > 0),
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tickets (
    id           SERIAL PRIMARY KEY,
    code         VARCHAR(20)  NOT NULL UNIQUE,
    category_id  INTEGER      NOT NULL REFERENCES categories (id),
    requester_id INTEGER      NOT NULL REFERENCES users (id),
    assignee_id  INTEGER      REFERENCES users (id),
    priority     VARCHAR(10)  NOT NULL DEFAULT 'NORMAL'
                 CHECK (priority IN ('LOW', 'NORMAL', 'HIGH', 'URGENT')),
    status       VARCHAR(20)  NOT NULL DEFAULT 'OPEN'
                 CHECK (status IN ('OPEN', 'ASSIGNED', 'IN_PROGRESS',
                                   'WAITING_USER', 'RESOLVED', 'CLOSED')),
    subject      VARCHAR(150) NOT NULL,
    description  TEXT         NOT NULL,
    due_at       TIMESTAMPTZ  NOT NULL,
    resolved_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS tickets_status_due_idx ON tickets (status, due_at);
CREATE INDEX IF NOT EXISTS tickets_assignee_idx   ON tickets (assignee_id);
CREATE INDEX IF NOT EXISTS tickets_requester_idx  ON tickets (requester_id);

-- Outbox: ditulis dalam transaksi yang sama dengan perubahan tiket, lalu
-- dikirim ke Redis oleh relay. Menjamin tidak ada notifikasi yang hilang.
CREATE TABLE IF NOT EXISTS ticket_outbox (
    id           BIGSERIAL PRIMARY KEY,
    ticket_id    INTEGER     NOT NULL REFERENCES tickets (id) ON DELETE CASCADE,
    event        VARCHAR(50) NOT NULL,
    payload      JSONB       NOT NULL DEFAULT '{}'::jsonb,
    attempts     INTEGER     NOT NULL DEFAULT 0,
    available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ticket_outbox_pending_idx
    ON ticket_outbox (available_at)
    WHERE processed_at IS NULL;
