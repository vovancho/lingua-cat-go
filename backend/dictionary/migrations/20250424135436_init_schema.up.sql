CREATE TABLE public.dictionary
(
    id         bigserial
        CONSTRAINT dictionary_pk PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP,
    lang       VARCHAR(2)              NOT NULL,
    name       VARCHAR(255)            NOT NULL,
    type       SMALLINT                NOT NULL
);

CREATE TABLE public.sentence
(
    id         bigserial
        CONSTRAINT sentence_pk PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP,
    text_ru    VARCHAR(255)            NOT NULL,
    text_en    VARCHAR(255)            NOT NULL
);

CREATE TABLE public.dictionary_sentence
(
    id            bigserial,
    dictionary_id bigint NOT NULL
        CONSTRAINT dictionary_sentence_dictionary_id_fk REFERENCES public.dictionary,
    sentence_id   bigint NOT NULL
        CONSTRAINT dictionary_sentence_sentence_id_fk REFERENCES public.sentence,
    CONSTRAINT dictionary_sentence_pk UNIQUE (sentence_id, dictionary_id)
);

CREATE TABLE public.translation
(
    id             bigserial
        CONSTRAINT translation_pk PRIMARY KEY,
    created_at     TIMESTAMP DEFAULT NOW() NOT NULL,
    deleted_at     TIMESTAMP,
    dictionary_id  bigint                  NOT NULL
        CONSTRAINT translation_dictionary_id_fk REFERENCES public.dictionary,
    translation_id bigint                  NOT NULL
        CONSTRAINT translation_dictionary_id_fk_2 REFERENCES public.dictionary,
    CONSTRAINT circular_reference_check CHECK (dictionary_id <> translation_id)
);

CREATE UNIQUE INDEX translation_id_dictionary_id_uniq ON public.translation (translation_id, dictionary_id) WHERE (deleted_at IS NULL);
