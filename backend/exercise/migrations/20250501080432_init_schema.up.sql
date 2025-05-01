CREATE TABLE public.exercise
(
    id                BIGSERIAL
        CONSTRAINT exercise_pk PRIMARY KEY,
    created_at        TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at        TIMESTAMP DEFAULT NOW() NOT NULL,
    user_id           UUID                    NOT NULL,
    lang              VARCHAR(2)              NOT NULL,
    task_amount       SMALLINT                NOT NULL,
    processed_counter SMALLINT  DEFAULT 0     NOT NULL,
    selected_counter  SMALLINT  DEFAULT 0     NOT NULL,
    corrected_counter SMALLINT  DEFAULT 0     NOT NULL
);

CREATE TABLE public.task
(
    id            BIGSERIAL
        CONSTRAINT task_pk PRIMARY KEY,
    exercise_id   BIGINT   NOT NULL
        CONSTRAINT task_exercise_id_fk REFERENCES public.exercise,
    words         BIGINT[] NOT NULL,
    word_correct  BIGINT   NOT NULL,
    word_selected BIGINT
);
