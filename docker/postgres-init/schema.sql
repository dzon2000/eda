CREATE TABLE outbox_events (
    id              UUID PRIMARY KEY,
    aggregate_type  TEXT NOT NULL,
    aggregate_id    TEXT NOT NULL,
    event_type      TEXT NOT NULL,
    payload         JSONB NOT NULL,
    schema_version  INT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'PENDING',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at    TIMESTAMPTZ,
    error           TEXT
);

CREATE INDEX idx_outbox_status_created
    ON outbox_events (status, created_at);

CREATE TABLE orders (
    id          UUID PRIMARY KEY,
    customer_id UUID NOT NULL,
    amount      NUMERIC(10, 2) NOT NULL,
    status      TEXT NOT NULL,
    discount    NUMERIC(5, 2) DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_customer_id
    ON orders (customer_id);
