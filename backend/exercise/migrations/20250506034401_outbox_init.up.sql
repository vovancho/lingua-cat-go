CREATE TABLE public.watermill_lcg_exercise_completed
(
    "offset"       BIGSERIAL,
    uuid           VARCHAR(36)                         NOT NULL,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    payload        JSON,
    metadata       JSON,
    transaction_id XID8                                NOT NULL,
    PRIMARY KEY (transaction_id, "offset")
);

CREATE TABLE public.watermill_offsets_lcg_exercise_completed
(
    consumer_group                VARCHAR(255) NOT NULL PRIMARY KEY,
    offset_acked                  BIGINT,
    last_processed_transaction_id XID8         NOT NULL
);
