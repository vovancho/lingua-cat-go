CREATE TABLE exercise_complete
(
    user_id               UUID,
    user_name             String,
    exercise_id           UInt64,
    exercise_lang         String,
    spent_time            UInt64,
    words_count           UInt16,
    words_corrected_count UInt16,
    event_time            DateTime DEFAULT now()
)
    ENGINE = MergeTree PARTITION BY toYYYYMM(event_time) ORDER BY (user_id, event_time, exercise_id);
