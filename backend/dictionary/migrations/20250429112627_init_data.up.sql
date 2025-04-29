-- Вставка записей в dictionary для английских слов (lang = 'en')
WITH inserted_en AS (
    INSERT INTO dictionary (lang, name, type, created_at)
        VALUES
            ('en', 'accomplish', 1, '2022-06-23 09:07:54'),
            ('en', 'accuse', 1, '2022-06-23 09:14:44'),
            ('en', 'across', 1, '2022-06-23 09:13:33'),
            ('en', 'admit', 1, '2022-06-23 09:12:46'),
            ('en', 'affect', 1, '2022-06-23 09:08:13'),
            ('en', 'afford', 1, '2022-06-23 09:10:54'),
            ('en', 'against', 1, '2022-06-23 09:09:12'),
            ('en', 'ahead', 1, '2022-06-23 09:09:47'),
            ('en', 'almost', 1, '2022-06-23 09:10:43'),
            ('en', 'along', 1, '2022-06-23 09:12:27'),
            ('en', 'although', 1, '2022-06-23 09:14:54'),
            ('en', 'ancient', 1, '2022-06-23 09:07:46'),
            ('en', 'anymore', 1, '2022-06-23 09:08:31'),
            ('en', 'appearance', 1, '2022-06-23 09:08:02'),
            ('en', 'appoint', 1, '2022-06-23 09:09:39'),
            ('en', 'approach', 1, '2022-06-23 09:09:14'),
            ('en', 'approval', 1, '2022-06-23 09:09:28'),
            ('en', 'approve', 1, '2022-06-23 09:11:09'),
            ('en', 'attend', 1, '2022-06-23 09:10:27'),
            ('en', 'attractive', 1, '2022-06-23 09:10:31'),
            ('en', 'attractively', 1, '2022-06-23 09:12:54'),
            ('en', 'aunt', 1, '2022-06-23 09:11:24'),
            ('en', 'back away', 2, '2022-06-23 09:13:04'),
            ('en', 'back off', 2, '2022-06-23 09:11:39'),
            ('en', 'be off', 2, '2022-06-23 09:07:43'),
            ('en', 'be over', 2, '2022-06-23 09:10:33')
        ON CONFLICT (name, lang) WHERE deleted_at IS NULL DO NOTHING
        RETURNING id, name
),
-- Вставка русских переводов в dictionary (lang = 'ru')
     inserted_ru AS (
         INSERT INTO dictionary (lang, name, type)
             VALUES
                 ('ru', 'выполнить', 1),
                 ('ru', 'совершать', 1),
                 ('ru', 'обвинять', 1),
                 ('ru', 'через', 1),
                 ('ru', 'поперек', 1),
                 ('ru', 'сквозь', 1),
                 ('ru', 'признавать', 1),
                 ('ru', 'допускать', 1),
                 ('ru', 'влиять', 1),
                 ('ru', 'воздействовать', 1),
                 ('ru', 'предоставлять', 1),
                 ('ru', 'против', 1),
                 ('ru', 'с', 1),
                 ('ru', 'вперед', 1),
                 ('ru', 'предстоящий', 1),
                 ('ru', 'почти', 1),
                 ('ru', 'по', 1),
                 ('ru', 'вдоль', 1),
                 ('ru', 'несмотря на то что', 1),
                 ('ru', 'древний', 1),
                 ('ru', 'больше', 1),
                 ('ru', 'внешность', 1),
                 ('ru', 'появление', 1),
                 ('ru', 'назначать', 1),
                 ('ru', 'подход', 1),
                 ('ru', 'одобрение', 1),
                 ('ru', 'одобрить', 1),
                 ('ru', 'посещать', 1),
                 ('ru', 'присутствовать', 1),
                 ('ru', 'привлекательный', 1),
                 ('ru', 'привлекательно', 1),
                 ('ru', 'тетя', 1),
                 ('ru', 'отойти', 1),
                 ('ru', 'отступить', 1),
                 ('ru', 'уйти', 1),
                 ('ru', 'быть свободным', 1),
                 ('ru', 'проходить', 1),
                 ('ru', 'заканчиваться', 1)
             ON CONFLICT (name, lang) WHERE deleted_at IS NULL DO NOTHING
             RETURNING id, name
     ),
-- Вставка предложений в sentence (только для 'across')
     inserted_sentences AS (
         INSERT INTO sentence (text_en, text_ru)
             VALUES
                 ('There is a forest across the river', 'Через реку есть лес'),
                 ('The boy ran across the street', 'Мальчик побежал через улицу')
             RETURNING id
     ),
-- Связывание dictionary и sentence через dictionary_sentence
     inserted_dict_sentences AS (
         INSERT INTO dictionary_sentence (dictionary_id, sentence_id)
             SELECT
                 (SELECT id FROM inserted_en WHERE name = 'across') AS dictionary_id,
                 id AS sentence_id
             FROM inserted_sentences
             RETURNING dictionary_id, sentence_id
     )
-- Вставка переводов в translation
INSERT INTO translation (dictionary_id, translation_id)
SELECT
    en.id AS dictionary_id,
    ru.id AS translation_id
FROM inserted_en en
         JOIN inserted_ru ru ON ru.name IN (
    CASE en.name
        WHEN 'accomplish' THEN 'выполнить'
        WHEN 'accomplish' THEN 'совершать'
        WHEN 'accuse' THEN 'обвинять'
        WHEN 'across' THEN 'через'
        WHEN 'across' THEN 'поперек'
        WHEN 'across' THEN 'сквозь'
        WHEN 'admit' THEN 'признавать'
        WHEN 'admit' THEN 'допускать'
        WHEN 'affect' THEN 'влиять'
        WHEN 'affect' THEN 'воздействовать'
        WHEN 'afford' THEN 'предоставлять'
        WHEN 'against' THEN 'против'
        WHEN 'against' THEN 'с'
        WHEN 'ahead' THEN 'вперед'
        WHEN 'ahead' THEN 'предстоящий'
        WHEN 'almost' THEN 'почти'
        WHEN 'along' THEN 'по'
        WHEN 'along' THEN 'вдоль'
        WHEN 'although' THEN 'несмотря на то что'
        WHEN 'ancient' THEN 'древний'
        WHEN 'anymore' THEN 'больше'
        WHEN 'appearance' THEN 'внешность'
        WHEN 'appearance' THEN 'появление'
        WHEN 'appoint' THEN 'назначать'
        WHEN 'approach' THEN 'подход'
        WHEN 'approval' THEN 'одобрение'
        WHEN 'approve' THEN 'одобрить'
        WHEN 'attend' THEN 'посещать'
        WHEN 'attend' THEN 'присутствовать'
        WHEN 'attractive' THEN 'привлекательный'
        WHEN 'attractively' THEN 'привлекательно'
        WHEN 'aunt' THEN 'тетя'
        WHEN 'back away' THEN 'отойти'
        WHEN 'back away' THEN 'отступить'
        WHEN 'back off' THEN 'отступить'
        WHEN 'be off' THEN 'уйти'
        WHEN 'be off' THEN 'быть свободным'
        WHEN 'be over' THEN 'проходить'
        WHEN 'be over' THEN 'заканчиваться'
        END
    )
WHERE ru.name IN (
    SELECT unnest(ARRAY['выполнить', 'совершать']) WHERE en.name = 'accomplish'
    UNION SELECT unnest(ARRAY['обвинять']) WHERE en.name = 'accuse'
    UNION SELECT unnest(ARRAY['через', 'поперек', 'сквозь']) WHERE en.name = 'across'
    UNION SELECT unnest(ARRAY['признавать', 'допускать']) WHERE en.name = 'admit'
    UNION SELECT unnest(ARRAY['влиять', 'воздействовать']) WHERE en.name = 'affect'
    UNION SELECT unnest(ARRAY['предоставлять']) WHERE en.name = 'afford'
    UNION SELECT unnest(ARRAY['против', 'с']) WHERE en.name = 'against'
    UNION SELECT unnest(ARRAY['вперед', 'предстоящий']) WHERE en.name = 'ahead'
    UNION SELECT unnest(ARRAY['почти']) WHERE en.name = 'almost'
    UNION SELECT unnest(ARRAY['по', 'вдоль']) WHERE en.name = 'along'
    UNION SELECT unnest(ARRAY['несмотря на то что']) WHERE en.name = 'although'
    UNION SELECT unnest(ARRAY['древний']) WHERE en.name = 'ancient'
    UNION SELECT unnest(ARRAY['больше']) WHERE en.name = 'anymore'
    UNION SELECT unnest(ARRAY['внешность', 'появление']) WHERE en.name = 'appearance'
    UNION SELECT unnest(ARRAY['назначать']) WHERE en.name = 'appoint'
    UNION SELECT unnest(ARRAY['подход']) WHERE en.name = 'approach'
    UNION SELECT unnest(ARRAY['одобрение']) WHERE en.name = 'approval'
    UNION SELECT unnest(ARRAY['одобрить']) WHERE en.name = 'approve'
    UNION SELECT unnest(ARRAY['посещать', 'присутствовать']) WHERE en.name = 'attend'
    UNION SELECT unnest(ARRAY['привлекательный']) WHERE en.name = 'attractive'
    UNION SELECT unnest(ARRAY['привлекательно']) WHERE en.name = 'attractively'
    UNION SELECT unnest(ARRAY['тетя']) WHERE en.name = 'aunt'
    UNION SELECT unnest(ARRAY['отойти', 'отступить']) WHERE en.name = 'back away'
    UNION SELECT unnest(ARRAY['отступить']) WHERE en.name = 'back off'
    UNION SELECT unnest(ARRAY['уйти', 'быть свободным']) WHERE en.name = 'be off'
    UNION SELECT unnest(ARRAY['проходить', 'заканчиваться']) WHERE en.name = 'be over'
)
ON CONFLICT (dictionary_id, translation_id) WHERE deleted_at IS NULL DO NOTHING;
