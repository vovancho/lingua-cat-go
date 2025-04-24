CREATE TABLE public.dictionary
(
    id         bigserial
        CONSTRAINT dictionary_pk PRIMARY KEY,
    name       VARCHAR(255)            NOT NULL,
    type       SMALLINT                NOT NULL,
    lang       VARCHAR(2)              NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE public.sentence
(
    id            bigserial
        CONSTRAINT sentence_pk PRIMARY KEY,
    dictionary_id bigint                  NOT NULL
        CONSTRAINT sentence_dictionary_id_fk REFERENCES public.dictionary,
    text          VARCHAR(255)            NOT NULL,
    translation   VARCHAR(255)            NOT NULL,
    created_at    TIMESTAMP DEFAULT NOW() NOT NULL,
    deleted_at    TIMESTAMP
);

CREATE TABLE public.translation
(
    id             bigserial
        CONSTRAINT translation_pk PRIMARY KEY,
    dictionary_id  bigint                  NOT NULL
        CONSTRAINT translation_dictionary_id_fk REFERENCES public.dictionary,
    translation_id bigint                  NOT NULL
        CONSTRAINT translation_dictionary_id_fk_2 REFERENCES public.dictionary,
    created_at     TIMESTAMP DEFAULT NOW() NOT NULL,
    deleted_at     TIMESTAMP,
    CONSTRAINT circular_reference_check CHECK (dictionary_id <> translation_id)
);

CREATE UNIQUE INDEX translation_id_dictionary_id_uniq ON public.translation (translation_id, dictionary_id) WHERE (deleted_at IS NULL);
