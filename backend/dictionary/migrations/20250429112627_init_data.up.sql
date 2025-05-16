-- Создание временной таблицы
CREATE TEMP TABLE temp_en_words (
    name VARCHAR(255),
    transcription VARCHAR(255),
    translations JSONB,
    type SMALLINT,
    sentences JSONB
);

-- Вставка данных из JSON
INSERT INTO temp_en_words (name, transcription, translations, type, sentences)
SELECT 
    (elem->>'name')::VARCHAR,
    (elem->>'transcription')::VARCHAR,
    (elem->'translation')::JSONB,
    (elem->>'type')::SMALLINT,
    (elem->'sentences')::JSONB
FROM jsonb_array_elements('[
  {
    "name": "accomplish",
    "transcription": "əˈkɑmplɪʃ",
    "translation": ["выполнить", "совершать"],
    "createdAt": "2022-06-23 09:07:54",
    "type": 1
  },
  {
    "name": "accuse",
    "transcription": "əˈkjuz",
    "translation": ["обвинять"],
    "createdAt": "2022-06-23 09:14:44",
    "type": 1
  },
  {
    "name": "across",
    "transcription": "əˈkrɔs",
    "translation": ["через", "поперек", "сквозь"],
    "createdAt": "2022-06-23 09:13:33",
    "type": 1,
    "sentences": [
      {
        "sentence": "There is a forest across the river",
        "translation": "Через реку есть лес"
      },
      {
        "sentence": "The boy ran across the street",
        "translation": "Мальчик побежал через улицу"
      }
    ],
    "groups": ["Игра престолов S01E01", "Игра престолов S01E02"]
  },
  {
    "name": "admit",
    "transcription": "ədˈmɪt",
    "translation": ["признавать", "допускать"],
    "createdAt": "2022-06-23 09:12:46",
    "type": 1
  },
  {
    "name": "affect",
    "transcription": "əˈfɛkt",
    "translation": ["влиять", "воздействовать"],
    "createdAt": "2022-06-23 09:08:13",
    "type": 1
  },
  {
    "name": "afford",
    "transcription": "əˈfɔrd",
    "translation": ["предоставлять"],
    "createdAt": "2022-06-23 09:10:54",
    "type": 1
  },
  {
    "name": "against",
    "transcription": "əˈgɛnst",
    "translation": ["против", "с"],
    "createdAt": "2022-06-23 09:09:12",
    "type": 1,
    "groups": ["Игра престолов S01E01"]
  },
  {
    "name": "ahead",
    "transcription": "əˈhɛd",
    "translation": ["вперед", "предстоящий"],
    "createdAt": "2022-06-23 09:09:47",
    "type": 1,
    "groups": ["Игра престолов S01E01"]
  },
  {
    "name": "almost",
    "transcription": "ˈɔlˌmoʊst",
    "translation": ["почти"],
    "createdAt": "2022-06-23 09:10:43",
    "type": 1,
    "groups": ["Игра престолов S01E01"]
  },
  {
    "name": "along",
    "transcription": "əˈlɔŋ",
    "translation": ["по", "вдоль"],
    "createdAt": "2022-06-23 09:12:27",
    "type": 1,
    "groups": ["Игра престолов S01E01"]
  },
  {
    "name": "although",
    "transcription": "ˌɔlˈðoʊ",
    "translation": ["несмотря на то что"],
    "createdAt": "2022-06-23 09:14:54",
    "type": 1
  },
  {
    "name": "ancient",
    "transcription": "ˈeɪnʧənt",
    "translation": ["древний"],
    "createdAt": "2022-06-23 09:07:46",
    "type": 1
  },
  {
    "name": "anymore",
    "transcription": "ˌɛniˈmɔr",
    "translation": ["больше"],
    "createdAt": "2022-06-23 09:08:31",
    "type": 1
  },
  {
    "name": "appearance",
    "transcription": "əˈpɪrəns",
    "translation": ["внешность", "появление"],
    "createdAt": "2022-06-23 09:08:02",
    "type": 1
  },
  {
    "name": "appoint",
    "transcription": "əˈpɔɪnt",
    "translation": ["назначать"],
    "createdAt": "2022-06-23 09:09:39",
    "type": 1
  },
  {
    "name": "approach",
    "transcription": "əˈproʊʧ",
    "translation": ["подход"],
    "createdAt": "2022-06-23 09:09:14",
    "type": 1,
    "groups": ["Игра престолов S01E01"]
  },
  {
    "name": "approval",
    "transcription": "əˈpruvəl",
    "translation": ["одобрение"],
    "createdAt": "2022-06-23 09:09:28",
    "type": 1
  },
  {
    "name": "approve",
    "transcription": "əˈpruv",
    "translation": ["одобрить"],
    "createdAt": "2022-06-23 09:11:09",
    "type": 1,
    "phrases": ["approval"]
  },
  {
    "name": "attend",
    "transcription": "əˈtɛnd",
    "translation": ["посещать", "присутствовать"],
    "createdAt": "2022-06-23 09:10:27",
    "type": 1
  },
  {
    "name": "attractive",
    "transcription": "əˈtræktɪv",
    "translation": ["привлекательный"],
    "createdAt": "2022-06-23 09:10:31",
    "type": 1,
    "groups": ["Игра престолов S01E01"],
    "phrases": ["attractively"]
  },
  {
    "name": "attractively",
    "transcription": "əˈtræktɪvli",
    "translation": ["привлекательно"],
    "createdAt": "2022-06-23 09:12:54",
    "type": 1
  },
  {
    "name": "aunt",
    "transcription": "ænt",
    "translation": ["тетя"],
    "createdAt": "2022-06-23 09:11:24",
    "type": 1,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "back away",
    "transcription": "bæk əˈweɪ",
    "translation": ["отойти", "отступить"],
    "createdAt": "2022-06-23 09:13:04",
    "type": 2,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "back off",
    "transcription": "bæk ɔf",
    "translation": ["отступить"],
    "createdAt": "2022-06-23 09:11:39",
    "type": 2,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "be off",
    "transcription": "bi ɔf",
    "translation": ["уйти", "быть свободным"],
    "createdAt": "2022-06-23 09:07:43",
    "type": 2,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "be over",
    "transcription": "bi ˈoʊvər",
    "translation": ["проходить", "заканчиваться"],
    "createdAt": "2022-06-23 09:10:33",
    "type": 2,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "be up to",
    "transcription": "bi ʌp tu",
    "translation": ["быть способным", "зависеть от", "задумать"],
    "createdAt": "2022-06-23 09:08:26",
    "type": 2,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "befool",
    "transcription": "bɪˈfuːl",
    "translation": ["одурачить"],
    "createdAt": "2022-06-23 09:14:30",
    "type": 1
  },
  {
    "name": "behind",
    "transcription": "bɪˈhaɪnd",
    "translation": ["сзади", "за"],
    "createdAt": "2022-06-23 09:08:22",
    "type": 1,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "belief",
    "transcription": "bɪˈlif",
    "translation": ["вера", "убеждение"],
    "createdAt": "2022-06-23 09:12:08",
    "type": 1,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "beside",
    "transcription": "bɪˈsaɪd",
    "translation": ["рядом", "вне"],
    "createdAt": "2022-06-23 09:08:59",
    "type": 1,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "bit",
    "transcription": "bɪt",
    "translation": ["немного"],
    "createdAt": "2022-06-23 09:11:44",
    "type": 1,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "blow up",
    "transcription": "bloʊ ʌp",
    "translation": ["вспылить", "взорвать"],
    "createdAt": "2022-06-23 09:11:03",
    "type": 2,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "book",
    "transcription": "bʊk",
    "translation": ["книга", "бронировать"],
    "createdAt": "2022-06-23 09:13:19",
    "type": 1,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "both",
    "transcription": "boʊθ",
    "translation": ["оба"],
    "createdAt": "2022-06-23 09:12:20",
    "type": 1,
    "groups": ["Игра престолов S01E02"]
  },
  {
    "name": "break down",
    "transcription": "breɪk daʊn",
    "translation": ["разрушить", "сдаться", "разломать"],
    "createdAt": "2022-06-23 09:07:45",
    "type": 2
  },
  {
    "name": "break in",
    "transcription": "breɪk ɪn",
    "translation": ["вломиться"],
    "createdAt": "2022-06-23 09:09:16",
    "type": 2
  },
  {
    "name": "break into",
    "transcription": "breɪk ˈɪntu",
    "translation": ["вступить", "вломиться", "пуститься"],
    "createdAt": "2022-06-23 09:08:05",
    "type": 2
  },
  {
    "name": "break off",
    "transcription": "breɪk ɔf",
    "translation": ["прекратить", "разорвать"],
    "createdAt": "2022-06-23 09:08:37",
    "type": 2
  },
  {
    "name": "break out",
    "transcription": "breɪk aʊt",
    "translation": ["вырваться", "выбраться"],
    "createdAt": "2022-06-23 09:10:48",
    "type": 2
  },
  {
    "name": "break up",
    "transcription": "breɪk ʌp",
    "translation": ["разогнать", "прекратить", "прервать"],
    "createdAt": "2022-06-23 09:13:37",
    "type": 2
  },
  {
    "name": "call back",
    "transcription": "kɔl bæk",
    "translation": ["перезвонить"],
    "createdAt": "2022-06-23 09:14:27",
    "type": 1
  },
  {
    "name": "call off",
    "transcription": "kɔl ɔf",
    "translation": ["отозвать", "перенести (свадьбу)"],
    "createdAt": "2022-06-23 09:12:35",
    "type": 2
  },
  {
    "name": "calm",
    "transcription": "kɑm",
    "translation": ["спокойствие"],
    "createdAt": "2022-06-23 09:09:21",
    "type": 1,
    "phrases": ["calm down"]
  },
  {
    "name": "calm down",
    "transcription": "kɑm daʊn",
    "translation": ["успокоиться"],
    "createdAt": "2022-06-23 09:13:03",
    "type": 2
  },
  {
    "name": "carry",
    "transcription": "ˈkæri",
    "translation": ["нести", "везти"],
    "createdAt": "2022-06-23 09:14:41",
    "type": 1,
    "phrases": ["carry on", "carry out"]
  },
  {
    "name": "carry on",
    "transcription": "ˈkæri ɑn",
    "translation": ["продолжать"],
    "createdAt": "2022-06-23 09:07:38",
    "type": 2
  },
  {
    "name": "carry out",
    "transcription": "ˈkæri aʊt",
    "translation": ["выполнять", "исполнить"],
    "createdAt": "2022-06-23 09:11:28",
    "type": 2
  },
  {
    "name": "casual",
    "transcription": "ˈkæʒəwəl",
    "translation": ["повседневный"],
    "createdAt": "2022-06-23 09:08:32",
    "type": 1
  },
  {
    "name": "catch up",
    "transcription": "kæʧ ʌp",
    "translation": ["догнать", "не отставать", "присоединиться"],
    "createdAt": "2022-06-23 09:11:27",
    "type": 2
  },
  {
    "name": "certain",
    "transcription": "ˈsɜrtən",
    "translation": ["определенный"],
    "createdAt": "2022-06-23 09:09:10",
    "type": 1
  },
  {
    "name": "cheat",
    "transcription": "ʧit",
    "translation": ["обманывать"],
    "createdAt": "2022-06-23 09:08:10",
    "type": 1
  },
  {
    "name": "check in",
    "transcription": "ʧɛk ɪn",
    "translation": ["зарегистрироваться"],
    "createdAt": "2022-06-23 09:11:18",
    "type": 2
  },
  {
    "name": "check out",
    "transcription": "ʧɛk aʊt",
    "translation": ["проверить", "освободить", "выехать"],
    "createdAt": "2022-06-23 09:11:16",
    "type": 2
  },
  {
    "name": "check-in",
    "transcription": "ʧɛk-ɪn",
    "translation": ["время заселения", "регистрация"],
    "createdAt": "2022-06-23 09:11:04",
    "type": 2
  },
  {
    "name": "cheer",
    "transcription": "ʧɪr",
    "translation": ["радость", "настроение"],
    "createdAt": "2022-06-23 09:09:09",
    "type": 1
  },
  {
    "name": "chill",
    "transcription": "ʧɪl",
    "translation": ["холод", "холодный"],
    "createdAt": "2022-06-23 09:09:55",
    "type": 1
  },
  {
    "name": "choose",
    "transcription": "ʧuz",
    "translation": ["выбрать"],
    "createdAt": "2022-06-23 09:11:25",
    "type": 3,
    "v2": {
      "name": "chose",
      "transcription": "ʧoʊz"
    },
    "v3": {
      "name": "chosen",
      "transcription": "ˈʧoʊzən"
    }
  },
  {
    "name": "citizen",
    "transcription": "ˈsɪtəzən",
    "translation": ["гражданин"],
    "createdAt": "2022-06-23 09:13:01",
    "type": 1,
    "phrases": ["citizenship"]
  },
  {
    "name": "citizenship",
    "transcription": "ˈsɪtɪzənˌʃɪp",
    "translation": ["гражданство"],
    "createdAt": "2022-06-23 09:07:48",
    "type": 1
  },
  {
    "name": "claim",
    "transcription": "kleɪm",
    "translation": ["требовать"],
    "createdAt": "2022-06-23 09:07:35",
    "type": 1
  },
  {
    "name": "climb",
    "transcription": "klaɪm",
    "translation": ["забираться"],
    "createdAt": "2022-06-23 09:11:43",
    "type": 1
  },
  {
    "name": "close",
    "transcription": "kloʊs",
    "translation": ["закрывать", "рядом", "близко", "поблизости"],
    "createdAt": "2022-06-23 09:12:21",
    "type": 1
  },
  {
    "name": "closet",
    "transcription": "ˈklɑzət",
    "translation": ["шкаф"],
    "createdAt": "2022-06-23 09:11:46",
    "type": 1
  },
  {
    "name": "clothes",
    "transcription": "kloʊðz",
    "translation": ["одежда"],
    "createdAt": "2022-06-23 09:08:01",
    "type": 1
  },
  {
    "name": "colleague",
    "transcription": "ˈkɑlig",
    "translation": ["коллега"],
    "createdAt": "2022-06-23 09:09:42",
    "type": 1
  },
  {
    "name": "colorful",
    "transcription": "ˈkʌlərfəl",
    "translation": ["красочный"],
    "createdAt": "2022-06-23 09:08:08",
    "type": 1
  },
  {
    "name": "come across",
    "transcription": "kʌm əˈkrɔs",
    "translation": ["столкнуться", "попасться случайно"],
    "createdAt": "2022-06-23 09:11:41",
    "type": 2
  },
  {
    "name": "come along",
    "transcription": "kʌm əˈlɔŋ",
    "translation": ["улучшиться", "выпадать", "продвигаться"],
    "createdAt": "2022-06-23 09:08:09",
    "type": 2
  },
  {
    "name": "come back",
    "transcription": "kʌm bæk",
    "translation": ["вернуться"],
    "createdAt": "2022-06-23 09:08:15",
    "type": 2
  },
  {
    "name": "come by",
    "transcription": "kʌm baɪ",
    "translation": ["подойти", "заглянуть"],
    "createdAt": "2022-06-23 09:08:36",
    "type": 2
  },
  {
    "name": "come down",
    "transcription": "kʌm daʊn",
    "translation": ["прийти", "спуститься"],
    "createdAt": "2022-06-23 09:12:32",
    "type": 2
  },
  {
    "name": "come forward",
    "transcription": "kʌm ˈfɔrwərd",
    "translation": ["выйти вперед", "выступить", "явиться"],
    "createdAt": "2022-06-23 09:08:33",
    "type": 2
  },
  {
    "name": "come in",
    "transcription": "kʌm ɪn",
    "translation": ["входить", "вступать"],
    "createdAt": "2022-06-23 09:10:53",
    "type": 2
  },
  {
    "name": "come off",
    "transcription": "kʌm ɔf",
    "translation": ["выйти вперед", "слететь"],
    "createdAt": "2022-06-23 09:11:48",
    "type": 2
  },
  {
    "name": "come on",
    "transcription": "kʌm ɑn",
    "translation": ["давай", "участвовать", "запуститься"],
    "createdAt": "2022-06-23 09:13:54",
    "type": 2
  },
  {
    "name": "come out",
    "transcription": "kʌm aʊt",
    "translation": ["выйти"],
    "createdAt": "2022-06-23 09:10:36",
    "type": 2
  },
  {
    "name": "come over",
    "transcription": "kʌm ˈoʊvər",
    "translation": ["приходить", "происходить"],
    "createdAt": "2022-06-23 09:08:51",
    "type": 2
  },
  {
    "name": "come up",
    "transcription": "kʌm ʌp",
    "translation": ["пойти", "подойти", "подняться"],
    "createdAt": "2022-06-23 09:07:47",
    "type": 2
  },
  {
    "name": "compete",
    "transcription": "kəmˈpit",
    "translation": ["конкурировать", "соревноваться"],
    "createdAt": "2022-06-23 09:13:41",
    "type": 1
  },
  {
    "name": "concern",
    "transcription": "kənˈsɜrn",
    "translation": ["беспокойство"],
    "createdAt": "2022-06-23 09:07:57",
    "type": 1
  },
  {
    "name": "confuse",
    "transcription": "kənˈfjuz",
    "translation": ["путать", "перепутать"],
    "createdAt": "2022-06-23 09:13:50",
    "type": 1
  },
  {
    "name": "consider",
    "transcription": "kənˈsɪdər",
    "translation": ["рассматривать", "учитывать"],
    "createdAt": "2022-06-23 09:10:26",
    "type": 1
  },
  {
    "name": "contest",
    "transcription": "ˈkɑntɛst",
    "translation": ["конкурс"],
    "createdAt": "2022-06-23 09:10:47",
    "type": 1
  },
  {
    "name": "costs",
    "transcription": "kɑsts",
    "translation": ["расходы"],
    "createdAt": "2022-06-23 09:11:10",
    "type": 1
  },
  {
    "name": "count on",
    "transcription": "kaʊnt ɑn",
    "translation": ["рассчитывать (на кого-то)"],
    "createdAt": "2022-06-23 09:13:29",
    "type": 2
  },
  {
    "name": "country",
    "transcription": "ˈkʌntri",
    "translation": ["страна", "загородный"],
    "createdAt": "2022-06-23 09:10:56",
    "type": 1
  },
  {
    "name": "couple",
    "transcription": "ˈkʌpəl",
    "translation": ["пара"],
    "createdAt": "2022-06-23 09:14:12",
    "type": 1
  },
  {
    "name": "creature",
    "transcription": "ˈkriʧər",
    "translation": ["существо"],
    "createdAt": "2022-06-23 09:12:12",
    "type": 1
  },
  {
    "name": "crew",
    "transcription": "kru",
    "translation": ["экипаж"],
    "createdAt": "2022-06-23 09:10:09",
    "type": 1
  },
  {
    "name": "cross",
    "transcription": "krɔs",
    "translation": ["пересекать"],
    "createdAt": "2022-06-23 09:09:40",
    "type": 1
  },
  {
    "name": "crowd",
    "transcription": "kraʊd",
    "translation": ["толпа"],
    "createdAt": "2022-06-23 09:11:33",
    "type": 1
  },
  {
    "name": "customs",
    "transcription": "ˈkʌstəmz",
    "translation": ["таможня", "таможенный"],
    "createdAt": "2022-06-23 09:10:08",
    "type": 1
  },
  {
    "name": "cut off",
    "transcription": "kʌt ɔf",
    "translation": ["отрезать", "отрубить"],
    "createdAt": "2022-06-23 09:10:37",
    "type": 2
  },
  {
    "name": "cut out",
    "transcription": "kʌt aʊt",
    "translation": ["быть созданным", "подходить"],
    "createdAt": "2022-06-23 09:14:26",
    "type": 2
  },
  {
    "name": "decide",
    "transcription": "ˌdɪˈsaɪd",
    "translation": ["решать"],
    "createdAt": "2022-06-23 09:12:47",
    "type": 1
  },
  {
    "name": "declare",
    "transcription": "dɪˈklɛr",
    "translation": ["объявить"],
    "createdAt": "2022-06-23 09:10:20",
    "type": 1
  },
  {
    "name": "defend",
    "transcription": "dɪˈfɛnd",
    "translation": ["защищать"],
    "createdAt": "2022-06-23 09:14:37",
    "type": 1
  },
  {
    "name": "definitely",
    "transcription": "ˈdɛfənətli",
    "translation": ["определенно"],
    "createdAt": "2022-06-23 09:14:47",
    "type": 1
  },
  {
    "name": "descriptive",
    "transcription": "dɪˈskrɪptɪv",
    "translation": ["описательный", "наглядный"],
    "createdAt": "2022-06-23 09:09:59",
    "type": 1
  },
  {
    "name": "deserve",
    "transcription": "dɪˈzɜrv",
    "translation": ["заслуживать"],
    "createdAt": "2022-06-23 09:11:52",
    "type": 1
  },
  {
    "name": "discover",
    "transcription": "dɪˈskʌvər",
    "translation": ["обнаружить"],
    "createdAt": "2022-06-23 09:08:57",
    "type": 1
  },
  {
    "name": "distant",
    "transcription": "ˈdɪstənt",
    "translation": ["далеко", "далекий", "отдаленный"],
    "createdAt": "2022-06-23 09:12:22",
    "type": 1
  },
  {
    "name": "during",
    "transcription": "ˈdʊrɪŋ",
    "translation": ["в течение", "во время"],
    "createdAt": "2022-06-23 09:13:25",
    "type": 1
  },
  {
    "name": "edge",
    "transcription": "ɛʤ",
    "translation": ["лезвие", "грань", "край"],
    "createdAt": "2022-06-23 09:08:48",
    "type": 1
  },
  {
    "name": "effort",
    "transcription": "ˈɛfərt",
    "translation": ["усилие", "попытка"],
    "createdAt": "2022-06-23 09:12:28",
    "type": 1
  },
  {
    "name": "either",
    "transcription": "ˈiðər",
    "translation": ["также", "один из двух", "каждый"],
    "createdAt": "2022-06-23 09:14:25",
    "type": 1
  },
  {
    "name": "either way",
    "transcription": "ˈiðər weɪ",
    "translation": ["так или иначе", "в любом случае"],
    "createdAt": "2022-06-23 09:11:15",
    "type": 4
  },
  {
    "name": "embarrass",
    "transcription": "ɪmˈbɛrəs",
    "translation": ["смущать", "стесняться"],
    "createdAt": "2022-06-23 09:08:29",
    "type": 1
  },
  {
    "name": "end up",
    "transcription": "ɛnd ʌp",
    "translation": ["очутиться", "добиться", "закончить"],
    "createdAt": "2022-06-23 09:08:11",
    "type": 2
  },
  {
    "name": "enemy",
    "transcription": "ˈɛnəmi",
    "translation": ["враг"],
    "createdAt": "2022-06-23 09:13:59",
    "type": 1
  },
  {
    "name": "enough",
    "transcription": "ɪˈnʌf",
    "translation": ["достаточно"],
    "createdAt": "2022-06-23 09:14:17",
    "type": 1
  },
  {
    "name": "entrance",
    "transcription": "ˈɛntrəns",
    "translation": ["вход"],
    "createdAt": "2022-06-23 09:12:18",
    "type": 1
  },
  {
    "name": "entrepreneur",
    "transcription": "ˌɑntrəprəˈnɜr",
    "translation": ["предприниматель"],
    "createdAt": "2022-06-23 09:12:56",
    "type": 1
  },
  {
    "name": "equally",
    "transcription": "ˈikwəli",
    "translation": ["одинаково"],
    "createdAt": "2022-06-23 09:08:44",
    "type": 1
  },
  {
    "name": "especially",
    "transcription": "əˈspɛʃli",
    "translation": ["особенно"],
    "createdAt": "2022-06-23 09:07:32",
    "type": 1
  },
  {
    "name": "even",
    "transcription": "ˈivɪn",
    "translation": ["даже"],
    "createdAt": "2022-06-23 09:14:49",
    "type": 1
  },
  {
    "name": "ever",
    "transcription": "ˈɛvər",
    "translation": ["когда-нибудь"],
    "createdAt": "2022-06-23 09:10:35",
    "type": 1
  },
  {
    "name": "exaggeration",
    "transcription": "ɪgˌzæʤəˈreɪʃən",
    "translation": ["преувеличение"],
    "createdAt": "2022-06-23 09:09:52",
    "type": 1
  },
  {
    "name": "except",
    "transcription": "ɪkˈsɛpt",
    "translation": ["исключая", "кроме", "вот только"],
    "createdAt": "2022-06-23 09:12:36",
    "type": 1
  },
  {
    "name": "excite",
    "transcription": "ɪkˈsaɪt",
    "translation": ["возбуждать"],
    "createdAt": "2022-06-23 09:09:20",
    "type": 1
  },
  {
    "name": "exhaust",
    "transcription": "ɪgˈzɑst",
    "translation": ["выхлоп", "исчерпывать"],
    "createdAt": "2022-06-23 09:12:23",
    "type": 1
  },
  {
    "name": "facial",
    "transcription": "ˈfeɪʃəl",
    "translation": ["лицевой"],
    "createdAt": "2022-06-23 09:12:45",
    "type": 1
  },
  {
    "name": "faith",
    "transcription": "feɪθ",
    "translation": ["вера"],
    "createdAt": "2022-06-23 09:09:30",
    "type": 1
  },
  {
    "name": "fall down",
    "transcription": "fɔl daʊn",
    "translation": ["рухнуть", "упасть"],
    "createdAt": "2022-06-23 09:13:20",
    "type": 2
  },
  {
    "name": "fall in love",
    "transcription": "fɔl ɪn lʌv",
    "translation": ["влюбляться"],
    "createdAt": "2022-06-23 09:09:54",
    "type": 4
  },
  {
    "name": "fall off",
    "transcription": "fɔl ɔf",
    "translation": ["упасть", "уснуть"],
    "createdAt": "2022-06-23 09:08:03",
    "type": 2
  },
  {
    "name": "fascinate",
    "transcription": "ˈfæsəˌneɪt",
    "translation": ["очаровывать"],
    "createdAt": "2022-06-23 09:14:34",
    "type": 1
  },
  {
    "name": "fictional",
    "transcription": "ˈfɪkʃənəl",
    "translation": ["вымышленный"],
    "createdAt": "2022-06-23 09:09:45",
    "type": 1
  },
  {
    "name": "figure out",
    "transcription": "ˈfɪgjər aʊt",
    "translation": ["выяснить", "понять"],
    "createdAt": "2022-06-23 09:12:53",
    "type": 2
  },
  {
    "name": "find out",
    "transcription": "faɪnd aʊt",
    "translation": ["выяснить", "понять"],
    "createdAt": "2022-06-23 09:10:14",
    "type": 2
  },
  {
    "name": "firm",
    "transcription": "fɜrm",
    "translation": ["твердый", "фирма"],
    "createdAt": "2022-06-23 09:08:55",
    "type": 1
  },
  {
    "name": "fit",
    "transcription": "fɪt",
    "translation": ["поместиться", "соответствовать"],
    "createdAt": "2022-06-23 09:07:49",
    "type": 1
  },
  {
    "name": "forbid",
    "transcription": "fərˈbɪd",
    "translation": ["запрещать"],
    "createdAt": "2022-06-23 09:08:20",
    "type": 1
  },
  {
    "name": "force",
    "transcription": "fɔrs",
    "translation": ["сила", "заставлять"],
    "createdAt": "2022-06-23 09:12:29",
    "type": 1
  },
  {
    "name": "fortunately",
    "transcription": "ˈfɔrʧənətli",
    "translation": ["к счастью"],
    "createdAt": "2022-06-23 09:12:33",
    "type": 1
  },
  {
    "name": "forward",
    "transcription": "ˈfɔrwərd",
    "translation": ["вперед", "переслать"],
    "createdAt": "2022-06-23 09:13:34",
    "type": 1
  },
  {
    "name": "founder",
    "transcription": "ˈfaʊndər",
    "translation": ["основатель"],
    "createdAt": "2022-06-23 09:10:13",
    "type": 1
  },
  {
    "name": "freshen up",
    "transcription": "ˈfrɛʃən ʌp",
    "translation": ["освежиться"],
    "createdAt": "2022-06-23 09:14:29",
    "type": 2
  },
  {
    "name": "fridge",
    "transcription": "frɪʤ",
    "translation": ["холодильник"],
    "createdAt": "2022-06-23 09:11:54",
    "type": 1
  },
  {
    "name": "fried",
    "transcription": "fraɪd",
    "translation": ["жареный"],
    "createdAt": "2022-06-23 09:12:13",
    "type": 1
  },
  {
    "name": "gain",
    "transcription": "geɪn",
    "translation": ["прирост", "усиление"],
    "createdAt": "2022-06-23 09:13:42",
    "type": 1
  },
  {
    "name": "gamble",
    "transcription": "ˈgæmbəl",
    "translation": ["азартная игра"],
    "createdAt": "2022-06-23 09:14:02",
    "type": 1
  },
  {
    "name": "garbage",
    "transcription": "ˈgɑrbɪʤ",
    "translation": ["мусор"],
    "createdAt": "2022-06-23 09:08:23",
    "type": 1
  },
  {
    "name": "gather",
    "transcription": "ˈgæðər",
    "translation": ["собирать"],
    "createdAt": "2022-06-23 09:11:47",
    "type": 1
  },
  {
    "name": "genial",
    "transcription": "ˈʤinjəl",
    "translation": ["прекрасный", "приветливый", "добродушный"],
    "createdAt": "2022-06-23 09:11:08",
    "type": 1
  },
  {
    "name": "gentle",
    "transcription": "ˈʤɛntəl",
    "translation": ["нежный", "вежливый"],
    "createdAt": "2022-06-23 09:10:10",
    "type": 1
  },
  {
    "name": "get along",
    "transcription": "gɛt əˈlɔŋ",
    "translation": ["выживать", "уживаться"],
    "createdAt": "2022-06-23 09:13:53",
    "type": 2
  },
  {
    "name": "get around",
    "transcription": "gɛt əˈraʊnd",
    "translation": ["передвигаться", "обмануть"],
    "createdAt": "2022-06-23 09:10:34",
    "type": 2
  },
  {
    "name": "get away",
    "transcription": "gɛt əˈweɪ",
    "translation": ["уйти", "убрать", "ускользнуть"],
    "createdAt": "2022-06-23 09:12:42",
    "type": 2
  },
  {
    "name": "get back",
    "transcription": "gɛt bæk",
    "translation": ["вернуться"],
    "createdAt": "2022-06-23 09:08:21",
    "type": 2
  },
  {
    "name": "get down",
    "transcription": "gɛt daʊn",
    "translation": ["пригнуться", "спуститься", "сосредоточиться"],
    "createdAt": "2022-06-23 09:12:09",
    "type": 2
  },
  {
    "name": "get in",
    "transcription": "gɛt ɪn",
    "translation": ["втянуть", "оказаться", "забраться"],
    "createdAt": "2022-06-23 09:14:36",
    "type": 2
  },
  {
    "name": "get off",
    "transcription": "gɛt ɔf",
    "translation": ["выходить", "спасти", "уходить"],
    "createdAt": "2022-06-23 09:09:37",
    "type": 2
  },
  {
    "name": "get on",
    "transcription": "gɛt ɑn",
    "translation": ["продолжить", "действовать", "добраться", "сесть"],
    "createdAt": "2022-06-23 09:10:46",
    "type": 2
  },
  {
    "name": "get out",
    "transcription": "gɛt aʊt",
    "translation": ["выйти", "вылезти", "выбраться"],
    "createdAt": "2022-06-23 09:13:05",
    "type": 2
  },
  {
    "name": "get over",
    "transcription": "gɛt ˈoʊvər",
    "translation": ["справиться", "отправиться"],
    "createdAt": "2022-06-23 09:10:57",
    "type": 2
  },
  {
    "name": "get through",
    "transcription": "gɛt θru",
    "translation": ["пройти", "провести", "попасть (куда-то через)"],
    "createdAt": "2022-06-23 09:12:51",
    "type": 2
  },
  {
    "name": "get up",
    "transcription": "gɛt ʌp",
    "translation": ["встать"],
    "createdAt": "2022-06-23 09:08:54",
    "type": 2
  },
  {
    "name": "give up",
    "transcription": "gɪv ʌp",
    "translation": ["сдаваться", "отказаться", "пожертвовать"],
    "createdAt": "2022-06-23 09:11:20",
    "type": 2
  },
  {
    "name": "glitter",
    "transcription": "ˈglɪtər",
    "translation": ["блеск"],
    "createdAt": "2022-06-23 09:10:16",
    "type": 1
  },
  {
    "name": "go alone",
    "transcription": "goʊ əˈloʊn",
    "translation": ["сопровождать", "согласиться"],
    "createdAt": "2022-06-23 09:09:05",
    "type": 2
  },
  {
    "name": "go away",
    "transcription": "goʊ əˈweɪ",
    "translation": ["уйти", "проходить", "отправляться"],
    "createdAt": "2022-06-23 09:08:49",
    "type": 2
  },
  {
    "name": "go back",
    "transcription": "goʊ bæk",
    "translation": ["вернуться"],
    "createdAt": "2022-06-23 09:14:10",
    "type": 2
  },
  {
    "name": "go by",
    "transcription": "goʊ baɪ",
    "translation": ["отправиться на", "проезжать в", "проходить мимо"],
    "createdAt": "2022-06-23 09:10:18",
    "type": 2
  },
  {
    "name": "go down",
    "transcription": "goʊ daʊn",
    "translation": ["спуститься"],
    "createdAt": "2022-06-23 09:11:29",
    "type": 2
  },
  {
    "name": "go in",
    "transcription": "goʊ ɪn",
    "translation": ["вступать", "войти", "зайти за тучи"],
    "createdAt": "2022-06-23 09:11:59",
    "type": 2
  },
  {
    "name": "go off",
    "transcription": "goʊ ɔf",
    "translation": ["потерять сознание", "уходить", "пойти"],
    "createdAt": "2022-06-23 09:11:02",
    "type": 2
  },
  {
    "name": "go on",
    "transcription": "goʊ ɑn",
    "translation": ["продолжать", "включаться", "развивать"],
    "createdAt": "2022-06-23 09:08:45",
    "type": 2
  },
  {
    "name": "go out",
    "transcription": "goʊ aʊt",
    "translation": ["пойти", "проводить время вне дома", "выйти из моды"],
    "createdAt": "2022-06-23 09:08:53",
    "type": 2
  },
  {
    "name": "go over",
    "transcription": "goʊ ˈoʊvər",
    "translation": ["перейти", "повторять"],
    "createdAt": "2022-06-23 09:08:39",
    "type": 2
  },
  {
    "name": "go through",
    "transcription": "goʊ θru",
    "translation": ["разобраться", "разбирать", "жить"],
    "createdAt": "2022-06-23 09:13:16",
    "type": 2
  },
  {
    "name": "go up",
    "transcription": "goʊ ʌp",
    "translation": ["подняться"],
    "createdAt": "2022-06-23 09:09:32",
    "type": 2
  },
  {
    "name": "gonna (going to)",
    "transcription": "ˈgɑnə (ˈgoʊɪŋ tu)",
    "translation": ["собираться"],
    "createdAt": "2022-06-23 09:11:49",
    "type": 4
  },
  {
    "name": "grab",
    "transcription": "græb",
    "translation": ["хватать"],
    "createdAt": "2022-06-23 09:14:14",
    "type": 1
  },
  {
    "name": "greeting",
    "transcription": "ˈgritɪŋ",
    "translation": ["приветствие"],
    "createdAt": "2022-06-23 09:13:58",
    "type": 1
  },
  {
    "name": "grow up",
    "transcription": "groʊ ʌp",
    "translation": ["вырасти", "расти"],
    "createdAt": "2022-06-23 09:12:34",
    "type": 2
  },
  {
    "name": "guerrilla",
    "transcription": "gəˈrɪlə",
    "translation": ["партизан"],
    "createdAt": "2022-06-23 09:12:43",
    "type": 1
  },
  {
    "name": "guess",
    "transcription": "gɛs",
    "translation": ["предполагать", "догадываться"],
    "createdAt": "2022-06-23 09:07:41",
    "type": 1
  },
  {
    "name": "haircut",
    "transcription": "ˈhɛrˌkʌt",
    "translation": ["стрижка"],
    "createdAt": "2022-06-23 09:13:49",
    "type": 1
  },
  {
    "name": "handle",
    "transcription": "ˈhændəl",
    "translation": ["управлять"],
    "createdAt": "2022-06-23 09:11:56",
    "type": 1
  },
  {
    "name": "handsome",
    "transcription": "ˈhænsəm",
    "translation": ["красивый (о мужчине)"],
    "createdAt": "2022-06-23 09:13:02",
    "type": 1
  },
  {
    "name": "hang",
    "transcription": "hæŋ",
    "translation": ["висеть", "вешать"],
    "createdAt": "2022-06-23 09:11:51",
    "type": 3,
    "v2": {
      "name": "hung",
      "transcription": "hʌŋ",
      "sentences": [
        {
          "sentence": "Your calendar hung in our kitchen for years",
          "translation": "Твой календарь висел у нас на кухне много лет"
        }
      ]
    },
    "v3": {
      "name": "hung",
      "transcription": "hʌŋ",
      "sentences": [
        {
          "sentence": "The picture had hung above the fireplace for years",
          "translation": "Картина висела над камином годами"
        }
      ]
    },
    "phrases": ["hang around", "hang on", "hang up"],
    "sentences": [
      {
        "sentence": "She hangs her clothes in the wardrobe",
        "translation": "Она вешает свою одежду в гардероб"
      },
      {
        "sentence": "The lamp hangs above the bed",
        "translation": "Лампа висит над кроватью"
      }
    ]
  },
  {
    "name": "hang around",
    "transcription": "hæŋ əˈraʊnd",
    "translation": ["ошиваться", "шастать"],
    "createdAt": "2022-06-23 09:11:35",
    "type": 2,
    "sentences": [
      {
        "sentence": "He would hang around the kitchen going bump in the night",
        "translation": "Он ошивался на кухне, грохоча в ночи"
      }
    ]
  },
  {
    "name": "hang on",
    "transcription": "hæŋ ɑn",
    "translation": ["подожди", "продержаться", "повиснуть", "цепляться"],
    "createdAt": "2022-06-23 09:09:11",
    "type": 2
  },
  {
    "name": "hang up",
    "transcription": "hæŋ ʌp",
    "translation": ["бросать трубку"],
    "createdAt": "2022-06-23 09:09:29",
    "type": 2
  },
  {
    "name": "headphones",
    "transcription": "ˈhɛdˌfoʊnz",
    "translation": ["наушники"],
    "createdAt": "2022-06-23 09:10:45",
    "type": 1
  },
  {
    "name": "heaven",
    "transcription": "ˈhɛvən",
    "translation": ["небеса"],
    "createdAt": "2022-06-23 09:08:18",
    "type": 1
  },
  {
    "name": "help out",
    "transcription": "hɛlp aʊt",
    "translation": ["выручать", "помочь"],
    "createdAt": "2022-06-23 09:07:52",
    "type": 2
  },
  {
    "name": "herbs",
    "transcription": "ɜrbz",
    "translation": ["трава", "травы"],
    "createdAt": "2022-06-23 09:11:42",
    "type": 1
  },
  {
    "name": "herself",
    "transcription": "hərˈsɛlf",
    "translation": ["сама"],
    "createdAt": "2022-06-23 09:10:39",
    "type": 1
  },
  {
    "name": "himself",
    "transcription": "hɪmˈsɛlf",
    "translation": ["сам"],
    "createdAt": "2022-06-23 09:14:13",
    "type": 1
  },
  {
    "name": "historically",
    "transcription": "hɪˈstɔrɪkəli",
    "translation": ["исторически"],
    "createdAt": "2022-06-23 09:14:48",
    "type": 1
  },
  {
    "name": "hold on",
    "transcription": "hoʊld ɑn",
    "translation": ["подождать", "продолжать"],
    "createdAt": "2022-06-23 09:12:44",
    "type": 2
  },
  {
    "name": "hold out",
    "transcription": "hoʊld aʊt",
    "translation": ["продержаться", "выдержать"],
    "createdAt": "2022-06-23 09:12:31",
    "type": 2
  },
  {
    "name": "hold up",
    "transcription": "hoʊld ʌp",
    "translation": ["остановиться", "задержать"],
    "createdAt": "2022-06-23 09:09:27",
    "type": 2
  },
  {
    "name": "hopelessness",
    "transcription": "ˈhoʊpləsnəs",
    "translation": ["безнадежность"],
    "createdAt": "2022-06-23 09:12:40",
    "type": 1
  },
  {
    "name": "horrible",
    "transcription": "ˈhɔrəbəl",
    "translation": ["ужасный"],
    "createdAt": "2022-06-23 09:09:18",
    "type": 1
  },
  {
    "name": "hostile",
    "transcription": "ˈhɑstəl",
    "translation": ["враждебный"],
    "createdAt": "2022-06-23 09:10:28",
    "type": 1
  },
  {
    "name": "humid",
    "transcription": "ˈhjuməd",
    "translation": ["влажный", "сырой (климат)"],
    "createdAt": "2022-06-23 09:11:26",
    "type": 1
  },
  {
    "name": "humidity",
    "transcription": "hjuˈmɪdəti",
    "translation": ["влажность"],
    "createdAt": "2022-06-23 09:14:38",
    "type": 1
  },
  {
    "name": "icy",
    "transcription": "ˈaɪsi",
    "translation": ["ледяной"],
    "createdAt": "2022-06-23 09:10:32",
    "type": 1
  },
  {
    "name": "impress",
    "transcription": "ˈɪmˌprɛs",
    "translation": ["впечатлить"],
    "createdAt": "2022-06-23 09:07:40",
    "type": 1,
    "phrases": ["impressed"]
  },
  {
    "name": "impressed",
    "transcription": "ɪmˈprɛst",
    "translation": ["впечатленный"],
    "createdAt": "2022-06-23 09:09:31",
    "type": 1
  },
  {
    "name": "in front of",
    "transcription": "ɪn frʌnt ʌv",
    "translation": ["перед"],
    "createdAt": "2022-06-23 09:07:51",
    "type": 4
  },
  {
    "name": "income",
    "transcription": "ˈɪnˌkʌm",
    "translation": ["доход", "прибыль"],
    "createdAt": "2022-06-23 09:12:03",
    "type": 1,
    "phrases": ["income tax"]
  },
  {
    "name": "income tax",
    "transcription": "ˈɪnˌkʌm tæks",
    "translation": ["подоходный налог"],
    "createdAt": "2022-06-23 09:09:56",
    "type": 4
  },
  {
    "name": "independence",
    "transcription": "ˌɪndɪˈpɛndəns",
    "translation": ["независимость"],
    "createdAt": "2022-06-23 09:09:01",
    "type": 1
  },
  {
    "name": "independent",
    "transcription": "ˌɪndɪˈpɛndənt",
    "translation": ["независимый"],
    "createdAt": "2022-06-23 09:13:00",
    "type": 1
  },
  {
    "name": "intelligent",
    "transcription": "ɪnˈtɛləʤənt",
    "translation": ["умный"],
    "createdAt": "2022-06-23 09:12:37",
    "type": 1
  },
  {
    "name": "interest rate",
    "transcription": "ˈɪntrəst reɪt",
    "translation": ["процентная ставка"],
    "createdAt": "2022-06-23 09:07:30",
    "type": 4
  },
  {
    "name": "irritate",
    "transcription": "ˈɪrɪˌteɪt",
    "translation": ["раздражать"],
    "createdAt": "2022-06-23 09:14:21",
    "type": 1
  },
  {
    "name": "jogging",
    "transcription": "ˈʤɑgɪŋ",
    "translation": ["пробежка"],
    "createdAt": "2022-06-23 09:12:02",
    "type": 1
  },
  {
    "name": "journey",
    "transcription": "ˈʤɜrni",
    "translation": ["путешествие"],
    "createdAt": "2022-06-23 09:13:10",
    "type": 1
  },
  {
    "name": "keep on",
    "transcription": "kip ɑn",
    "translation": ["держать", "продолжать"],
    "createdAt": "2022-06-23 09:14:01",
    "type": 2
  },
  {
    "name": "keep up",
    "transcription": "kip ʌp",
    "translation": ["поддерживать", "не отставать", "продолжать", "поспевать"],
    "createdAt": "2022-06-23 09:10:15",
    "type": 2
  },
  {
    "name": "kettle",
    "transcription": "ˈkɛtəl",
    "translation": ["чайник"],
    "createdAt": "2022-06-23 09:12:26",
    "type": 1
  },
  {
    "name": "kind",
    "transcription": "kaɪnd",
    "translation": ["добрый", "вид", "тип"],
    "createdAt": "2022-06-23 09:14:43",
    "type": 1
  },
  {
    "name": "knock",
    "transcription": "nɑk",
    "translation": ["стучать", "ударять"],
    "createdAt": "2022-06-23 09:09:23",
    "type": 1,
    "phrases": ["knock down", "knock off", "knock out"]
  },
  {
    "name": "knock down",
    "transcription": "nɑk daʊn",
    "translation": ["уронить", "продавить", "разрушить"],
    "createdAt": "2022-06-23 09:10:07",
    "type": 2
  },
  {
    "name": "knock off",
    "transcription": "nɑk ɔf",
    "translation": ["сбить", "быстро сделать", "уйти с работы", "копировать"],
    "createdAt": "2022-06-23 09:10:49",
    "type": 2
  },
  {
    "name": "knock out",
    "transcription": "nɑk aʊt",
    "translation": ["выбить", "ошеломить", "разрушить", "оглушить"],
    "createdAt": "2022-06-23 09:12:49",
    "type": 2
  },
  {
    "name": "laugh",
    "transcription": "læf",
    "translation": ["смеяться"],
    "createdAt": "2022-06-23 09:11:30",
    "type": 1
  },
  {
    "name": "lay",
    "transcription": "leɪ",
    "translation": ["лежать", "положить"],
    "createdAt": "2022-06-23 09:08:19",
    "type": 3,
    "v2": {
      "name": "laid",
      "transcription": "leɪd"
    },
    "v3": {
      "name": "laid",
      "transcription": "leɪd"
    }
  },
  {
    "name": "lead",
    "transcription": "lid",
    "translation": ["приводить", "вести"],
    "createdAt": "2022-06-23 09:11:36",
    "type": 3,
    "v2": {
      "name": "led",
      "transcription": "lɛd"
    },
    "v3": {
      "name": "led",
      "transcription": "lɛd"
    }
  },
  {
    "name": "lean",
    "transcription": "lin",
    "translation": ["опираться", "наклоняться"],
    "createdAt": "2022-06-23 09:10:11",
    "type": 1
  },
  {
    "name": "let in",
    "transcription": "lɛt ɪn",
    "translation": ["впустить"],
    "createdAt": "2022-06-23 09:10:52",
    "type": 2
  },
  {
    "name": "let out",
    "transcription": "lɛt aʊt",
    "translation": ["выпускать", "выдохнуть"],
    "createdAt": "2022-06-23 09:09:33",
    "type": 2
  },
  {
    "name": "liability",
    "transcription": "ˌlaɪəˈbɪlɪti",
    "translation": ["обязательство"],
    "createdAt": "2022-06-23 09:14:16",
    "type": 1
  },
  {
    "name": "lie down",
    "transcription": "laɪ daʊn",
    "translation": ["лечь", "прилечь"],
    "createdAt": "2022-06-23 09:12:14",
    "type": 2
  },
  {
    "name": "literally",
    "transcription": "ˈlɪtərəli",
    "translation": ["буквально"],
    "createdAt": "2022-06-23 09:11:19",
    "type": 1
  },
  {
    "name": "look back",
    "transcription": "lʊk bæk",
    "translation": ["оглянуться"],
    "createdAt": "2022-06-23 09:14:53",
    "type": 2
  },
  {
    "name": "look forward to",
    "transcription": "lʊk ˈfɔrwərd tu",
    "translation": ["ждать с нетерпением"],
    "createdAt": "2022-06-23 09:08:27",
    "type": 2
  },
  {
    "name": "look out",
    "transcription": "lʊk aʊt",
    "translation": ["приглядывать", "выглянуть"],
    "createdAt": "2022-06-23 09:12:57",
    "type": 2
  },
  {
    "name": "look over",
    "transcription": "lʊk ˈoʊvər",
    "translation": ["осмотреть"],
    "createdAt": "2022-06-23 09:14:51",
    "type": 2
  },
  {
    "name": "losses",
    "transcription": "ˈlɔsəz",
    "translation": ["убытки"],
    "createdAt": "2022-06-23 09:13:43",
    "type": 1
  },
  {
    "name": "luggage",
    "transcription": "ˈlʌgəʤ",
    "translation": ["багаж"],
    "createdAt": "2022-06-23 09:14:04",
    "type": 1
  },
  {
    "name": "lunatic",
    "transcription": "ˈlunəˌtɪk",
    "translation": ["сумасшедший"],
    "createdAt": "2022-06-23 09:12:25",
    "type": 1
  },
  {
    "name": "mad",
    "transcription": "mæd",
    "translation": ["безумный"],
    "createdAt": "2022-06-23 09:13:17",
    "type": 1
  },
  {
    "name": "main course",
    "transcription": "meɪn kɔrs",
    "translation": ["основное блюдо"],
    "createdAt": "2022-06-23 09:09:08",
    "type": 4
  },
  {
    "name": "make out",
    "transcription": "meɪk aʊt",
    "translation": ["справляться", "разглядеть"],
    "createdAt": "2022-06-23 09:11:23",
    "type": 2
  },
  {
    "name": "make up",
    "transcription": "meɪk ʌp",
    "translation": ["решить", "макияж", "принять (решение)", "выдумать"],
    "createdAt": "2022-06-23 09:09:57",
    "type": 2
  },
  {
    "name": "manufacture",
    "transcription": "ˌmænjəˈfækʧər",
    "translation": ["производство"],
    "createdAt": "2022-06-23 09:07:58",
    "type": 1
  },
  {
    "name": "matter",
    "transcription": "ˈmætər",
    "translation": ["значить", "важно"],
    "createdAt": "2022-06-23 09:13:18",
    "type": 1
  },
  {
    "name": "mean",
    "transcription": "min",
    "translation": ["среднее"],
    "createdAt": "2022-06-23 09:14:00",
    "type": 3,
    "v2": {
      "name": "meant",
      "transcription": "mɛnt"
    },
    "v3": {
      "name": "meant",
      "transcription": "mɛnt"
    }
  },
  {
    "name": "mine",
    "transcription": "maɪn",
    "translation": ["мой"],
    "createdAt": "2022-06-23 09:09:17",
    "type": 1
  },
  {
    "name": "mortgage",
    "transcription": "ˈmɔrgəʤ",
    "translation": ["ипотечный кредит"],
    "createdAt": "2022-06-23 09:10:24",
    "type": 1
  },
  {
    "name": "move in",
    "transcription": "muv ɪn",
    "translation": ["въезжать", "переехать (в)"],
    "createdAt": "2022-06-23 09:07:37",
    "type": 2
  },
  {
    "name": "move on",
    "transcription": "muv ɑn",
    "translation": ["двигаться", "продолжить движение"],
    "createdAt": "2022-06-23 09:13:32",
    "type": 2
  },
  {
    "name": "move out",
    "transcription": "muv aʊt",
    "translation": ["съехать", "переехать (из)"],
    "createdAt": "2022-06-23 09:12:16",
    "type": 2
  },
  {
    "name": "myself",
    "transcription": "ˌmaɪˈsɛlf",
    "translation": ["сам"],
    "createdAt": "2022-06-23 09:08:40",
    "type": 1
  },
  {
    "name": "napkin",
    "transcription": "ˈnæpkɪn",
    "translation": ["салфетка"],
    "createdAt": "2022-06-23 09:11:31",
    "type": 1
  },
  {
    "name": "narrow",
    "transcription": "ˈnɛroʊ",
    "translation": ["узкий"],
    "createdAt": "2022-06-23 09:10:22",
    "type": 1
  },
  {
    "name": "native",
    "transcription": "ˈneɪtɪv",
    "translation": ["родной"],
    "createdAt": "2022-06-23 09:12:11",
    "type": 1
  },
  {
    "name": "neighboring",
    "transcription": "ˈneɪbərɪŋ",
    "translation": ["соседний"],
    "createdAt": "2022-06-23 09:08:42",
    "type": 1
  },
  {
    "name": "nurse",
    "transcription": "nɜrs",
    "translation": ["медсестра"],
    "createdAt": "2022-06-23 09:13:26",
    "type": 1
  },
  {
    "name": "occur",
    "transcription": "əˈkɜr",
    "translation": ["происходить"],
    "createdAt": "2022-06-23 09:07:39",
    "type": 1
  },
  {
    "name": "often",
    "transcription": "ˈɔfən",
    "translation": ["часто"],
    "createdAt": "2022-06-23 09:10:04",
    "type": 1
  },
  {
    "name": "once",
    "transcription": "wʌns",
    "translation": ["однажды", "раз"],
    "createdAt": "2022-06-23 09:07:28",
    "type": 1
  },
  {
    "name": "orderly",
    "transcription": "ˈɔrdərli",
    "translation": ["аккуратный"],
    "createdAt": "2022-06-23 09:14:28",
    "type": 1
  },
  {
    "name": "originally",
    "transcription": "əˈrɪʤənəli",
    "translation": ["первоначально", "оригинально"],
    "createdAt": "2022-06-23 09:10:25",
    "type": 1
  },
  {
    "name": "outside",
    "transcription": "ˈaʊtˈsaɪd",
    "translation": ["вне"],
    "createdAt": "2022-06-23 09:10:05",
    "type": 1
  },
  {
    "name": "over",
    "transcription": "ˈoʊvər",
    "translation": ["над", "выше"],
    "createdAt": "2022-06-23 09:13:51",
    "type": 1
  },
  {
    "name": "ownership",
    "transcription": "ˈoʊnərˌʃɪp",
    "translation": ["собственность"],
    "createdAt": "2022-06-23 09:11:45",
    "type": 1
  },
  {
    "name": "pair",
    "transcription": "pɛr",
    "translation": ["пара"],
    "createdAt": "2022-06-23 09:14:09",
    "type": 1
  },
  {
    "name": "particular",
    "transcription": "pərˈtɪkjələr",
    "translation": ["конкретный"],
    "createdAt": "2022-06-23 09:11:38",
    "type": 1
  },
  {
    "name": "pass",
    "transcription": "pæs",
    "translation": ["проходить", "передавать"],
    "createdAt": "2022-06-23 09:08:46",
    "type": 1,
    "phrases": ["pass out"]
  },
  {
    "name": "pass out",
    "transcription": "pæs aʊt",
    "translation": ["раздавать", "потерять сознание"],
    "createdAt": "2022-06-23 09:08:35",
    "type": 2
  },
  {
    "name": "peninsula",
    "transcription": "pəˈnɪnsələ",
    "translation": ["полуостров"],
    "createdAt": "2022-06-23 09:09:25",
    "type": 1
  },
  {
    "name": "perhaps",
    "transcription": "pərˈhæps",
    "translation": ["возможно"],
    "createdAt": "2022-06-23 09:14:39",
    "type": 1
  },
  {
    "name": "persuasion",
    "transcription": "pərˈsweɪʒən",
    "translation": ["убеждение", "уговор"],
    "createdAt": "2022-06-23 09:13:38",
    "type": 1
  },
  {
    "name": "pick up",
    "transcription": "pɪk ʌp",
    "translation": ["подобрать", "забрать"],
    "createdAt": "2022-06-23 09:14:15",
    "type": 2
  },
  {
    "name": "plumber",
    "transcription": "ˈplʌmər",
    "translation": ["сантехник"],
    "createdAt": "2022-06-23 09:13:46",
    "type": 1
  },
  {
    "name": "point out",
    "transcription": "pɔɪnt aʊt",
    "translation": ["указать"],
    "createdAt": "2022-06-23 09:10:42",
    "type": 2
  },
  {
    "name": "precipitation",
    "transcription": "prɪˌsɪpɪˈteɪʃən",
    "translation": ["осадки"],
    "createdAt": "2022-06-23 09:11:06",
    "type": 1
  },
  {
    "name": "pretend",
    "transcription": "priˈtɛnd",
    "translation": ["притворяться", "претендовать"],
    "createdAt": "2022-06-23 09:07:44",
    "type": 1
  },
  {
    "name": "prevention",
    "transcription": "priˈvɛnʃən",
    "translation": ["профилактика", "предупреждение"],
    "createdAt": "2022-06-23 09:08:34",
    "type": 1
  },
  {
    "name": "probably",
    "transcription": "ˈprɑbəbli",
    "translation": ["вероятно"],
    "createdAt": "2022-06-23 09:10:41",
    "type": 1
  },
  {
    "name": "produce",
    "transcription": "ˈproʊdus",
    "translation": ["производить"],
    "createdAt": "2022-06-23 09:09:00",
    "type": 1
  },
  {
    "name": "prohibit",
    "transcription": "proʊˈhɪbət",
    "translation": ["запрещать", "препятствовать"],
    "createdAt": "2022-06-23 09:13:47",
    "type": 1,
    "phrases": ["prohibition"]
  },
  {
    "name": "prohibition",
    "transcription": "ˌproʊəˈbɪʃən",
    "translation": ["запрет"],
    "createdAt": "2022-06-23 09:08:28",
    "type": 1
  },
  {
    "name": "promotion",
    "transcription": "prəˈmoʊʃən",
    "translation": ["повышение", "продвижение"],
    "createdAt": "2022-06-23 09:10:58",
    "type": 1
  },
  {
    "name": "prove",
    "transcription": "pruv",
    "translation": ["доказывать"],
    "createdAt": "2022-06-23 09:13:08",
    "type": 1
  },
  {
    "name": "pull away",
    "transcription": "pʊl əˈweɪ",
    "translation": ["оттолкнуть"],
    "createdAt": "2022-06-23 09:14:05",
    "type": 2
  },
  {
    "name": "pull off",
    "transcription": "pʊl ɔf",
    "translation": ["стянуть", "оттащить"],
    "createdAt": "2022-06-23 09:08:12",
    "type": 2
  },
  {
    "name": "pull on",
    "transcription": "pʊl ɑn",
    "translation": ["натянуть", "надеть"],
    "createdAt": "2022-06-23 09:14:42",
    "type": 2
  },
  {
    "name": "pull out",
    "transcription": "pʊl aʊt",
    "translation": ["вытаскивать"],
    "createdAt": "2022-06-23 09:11:58",
    "type": 2
  },
  {
    "name": "pull up",
    "transcription": "pʊl ʌp",
    "translation": ["придвинуть"],
    "createdAt": "2022-06-23 09:08:17",
    "type": 2
  },
  {
    "name": "puppy",
    "transcription": "ˈpʌpi",
    "translation": ["щенок"],
    "createdAt": "2022-06-23 09:08:56",
    "type": 1
  },
  {
    "name": "put away",
    "transcription": "pʊt əˈweɪ",
    "translation": ["убрать", "избавиться"],
    "createdAt": "2022-06-23 09:10:50",
    "type": 2
  },
  {
    "name": "put down",
    "transcription": "pʊt daʊn",
    "translation": ["положить (в сторону)"],
    "createdAt": "2022-06-23 09:14:22",
    "type": 2
  },
  {
    "name": "put in",
    "transcription": "pʊt ɪn",
    "translation": ["вставить", "установить"],
    "createdAt": "2022-06-23 09:12:38",
    "type": 2
  },
  {
    "name": "put on",
    "transcription": "pʊt ɑn",
    "translation": ["надеть", "набрать вес"],
    "createdAt": "2022-06-23 09:14:18",
    "type": 2
  },
  {
    "name": "put out",
    "transcription": "pʊt aʊt",
    "translation": ["потушить", "выкладывать"],
    "createdAt": "2022-06-23 09:12:04",
    "type": 2
  },
  {
    "name": "put up",
    "transcription": "pʊt ʌp",
    "translation": ["мириться", "вкладывать"],
    "createdAt": "2022-06-23 09:10:30",
    "type": 2
  },
  {
    "name": "quite",
    "transcription": "kwaɪt",
    "translation": ["довольно", "вполне"],
    "createdAt": "2022-06-23 09:10:55",
    "type": 1
  },
  {
    "name": "raise",
    "transcription": "reɪz",
    "translation": ["поднимать"],
    "createdAt": "2022-06-23 09:14:31",
    "type": 1
  },
  {
    "name": "rare",
    "transcription": "rɛr",
    "translation": ["редкий"],
    "createdAt": "2022-06-23 09:11:40",
    "type": 1,
    "phrases": ["rarely"]
  },
  {
    "name": "rarely",
    "transcription": "ˈrɛrli",
    "translation": ["редко"],
    "createdAt": "2022-06-23 09:07:59",
    "type": 1
  },
  {
    "name": "reach",
    "transcription": "riʧ",
    "translation": ["достигать"],
    "createdAt": "2022-06-23 09:13:45",
    "type": 1
  },
  {
    "name": "realize",
    "transcription": "ˈriəˌlaɪz",
    "translation": ["понимать", "осознавать"],
    "createdAt": "2022-06-23 09:13:35",
    "type": 1
  },
  {
    "name": "recognition",
    "transcription": "ˌrɛkəgˈnɪʃən",
    "translation": ["признание", "распознавание"],
    "createdAt": "2022-06-23 09:10:59",
    "type": 1
  },
  {
    "name": "recognize",
    "transcription": "ˈrɛkəgˌnaɪz",
    "translation": ["распознавать"],
    "createdAt": "2022-06-23 09:08:04",
    "type": 1,
    "phrases": ["recognition"]
  },
  {
    "name": "regrettable",
    "transcription": "rɪˈgrɛtəbəl",
    "translation": ["достойный сожаления"],
    "createdAt": "2022-06-23 09:14:08",
    "type": 1
  },
  {
    "name": "remain",
    "transcription": "rɪˈmeɪn",
    "translation": ["оставаться"],
    "createdAt": "2022-06-23 09:14:03",
    "type": 1
  },
  {
    "name": "represent",
    "transcription": "ˌrɛprəˈzɛnt",
    "translation": ["представлять"],
    "createdAt": "2022-06-23 09:13:27",
    "type": 1
  },
  {
    "name": "reserve",
    "transcription": "rɪˈzɜrv",
    "translation": ["резервировать"],
    "createdAt": "2022-06-23 09:13:39",
    "type": 1
  },
  {
    "name": "respond",
    "transcription": "rɪˈspɑnd",
    "translation": ["реагировать", "отвечать"],
    "createdAt": "2022-06-23 09:09:24",
    "type": 1
  },
  {
    "name": "reunion",
    "transcription": "riˈunjən",
    "translation": ["воссоединение"],
    "createdAt": "2022-06-23 09:10:00",
    "type": 1
  },
  {
    "name": "right away",
    "transcription": "raɪt əˈweɪ",
    "translation": ["немедленно"],
    "createdAt": "2022-06-23 09:11:34",
    "type": 4
  },
  {
    "name": "rise",
    "transcription": "raɪz",
    "translation": ["подниматься"],
    "createdAt": "2022-06-23 09:08:38",
    "type": 3,
    "v2": {
      "name": "rose",
      "transcription": "roʊz"
    },
    "v3": {
      "name": "risen",
      "transcription": "ˈrɪzən"
    }
  },
  {
    "name": "rock",
    "transcription": "rɑk",
    "translation": ["рок", "камень", "скала"],
    "createdAt": "2022-06-23 09:12:15",
    "type": 1
  },
  {
    "name": "roll",
    "transcription": "roʊl",
    "translation": ["катиться", "рулон", "рулет"],
    "createdAt": "2022-06-23 09:13:28",
    "type": 1
  },
  {
    "name": "roof",
    "transcription": "ruf",
    "translation": ["крыша"],
    "createdAt": "2022-06-23 09:14:32",
    "type": 1
  },
  {
    "name": "rude",
    "transcription": "rud",
    "translation": ["грубый"],
    "createdAt": "2022-06-23 09:10:40",
    "type": 1
  },
  {
    "name": "run away",
    "transcription": "rʌn əˈweɪ",
    "translation": ["убегать", "сбегать"],
    "createdAt": "2022-06-23 09:13:40",
    "type": 2
  },
  {
    "name": "run into",
    "transcription": "rʌn ˈɪntu",
    "translation": ["столкнуться", "встретить"],
    "createdAt": "2022-06-23 09:12:24",
    "type": 2
  },
  {
    "name": "run off",
    "transcription": "rʌn ɔf",
    "translation": ["убегать", "сбегать"],
    "createdAt": "2022-06-23 09:09:04",
    "type": 2
  },
  {
    "name": "run out",
    "transcription": "rʌn aʊt",
    "translation": ["закончиться"],
    "createdAt": "2022-06-23 09:09:22",
    "type": 2
  },
  {
    "name": "run over",
    "transcription": "rʌn ˈoʊvər",
    "translation": ["переехать", "переливаться (через край)"],
    "createdAt": "2022-06-23 09:12:55",
    "type": 2
  },
  {
    "name": "rush",
    "transcription": "rʌʃ",
    "translation": ["прилив", "торопиться", "срочный"],
    "createdAt": "2022-06-23 09:07:27",
    "type": 1
  },
  {
    "name": "sack",
    "transcription": "sæk",
    "translation": ["мешок"],
    "createdAt": "2022-06-23 09:10:23",
    "type": 1
  },
  {
    "name": "sadly",
    "transcription": "ˈsædli",
    "translation": ["грустно"],
    "createdAt": "2022-06-23 09:11:50",
    "type": 1
  },
  {
    "name": "safe",
    "transcription": "seɪf",
    "translation": ["безопасный"],
    "createdAt": "2022-06-23 09:13:14",
    "type": 1
  },
  {
    "name": "salary",
    "transcription": "ˈsæləri",
    "translation": ["оклад", "зарплата"],
    "createdAt": "2022-06-23 09:09:38",
    "type": 1
  },
  {
    "name": "same",
    "transcription": "seɪm",
    "translation": ["такой же", "тот же", "то же"],
    "createdAt": "2022-06-23 09:14:33",
    "type": 1
  },
  {
    "name": "sauce",
    "transcription": "sɔs",
    "translation": ["соус"],
    "createdAt": "2022-06-23 09:09:34",
    "type": 1
  },
  {
    "name": "scissors",
    "transcription": "ˈsɪzərz",
    "translation": ["ножницы"],
    "createdAt": "2022-06-23 09:10:02",
    "type": 1
  },
  {
    "name": "scream",
    "transcription": "skrim",
    "translation": ["кричать"],
    "createdAt": "2022-06-23 09:07:29",
    "type": 1
  },
  {
    "name": "seem",
    "transcription": "sim",
    "translation": ["казаться"],
    "createdAt": "2022-06-23 09:10:01",
    "type": 1
  },
  {
    "name": "sensation",
    "transcription": "sɛnˈseɪʃən",
    "translation": ["ощущение"],
    "createdAt": "2022-06-23 09:10:12",
    "type": 1
  },
  {
    "name": "sense",
    "transcription": "sɛns",
    "translation": ["чувство"],
    "createdAt": "2022-06-23 09:14:11",
    "type": 1,
    "phrases": ["sensation"]
  },
  {
    "name": "set off",
    "transcription": "sɛt ɔf",
    "translation": ["отправляться"],
    "createdAt": "2022-06-23 09:10:29",
    "type": 1
  },
  {
    "name": "set up",
    "transcription": "sɛt ʌp",
    "translation": ["установить", "организовать"],
    "createdAt": "2022-06-23 09:12:41",
    "type": 2
  },
  {
    "name": "short-term",
    "transcription": "ʃɔrt-tɜrm",
    "translation": ["краткосрочный"],
    "createdAt": "2022-06-23 09:10:38",
    "type": 1
  },
  {
    "name": "shortened",
    "transcription": "ˈʃɔrtənd",
    "translation": ["укороченный"],
    "createdAt": "2022-06-23 09:13:52",
    "type": 1
  },
  {
    "name": "shout",
    "transcription": "ʃaʊt",
    "translation": ["кричать"],
    "createdAt": "2022-06-23 09:09:44",
    "type": 1
  },
  {
    "name": "show up",
    "transcription": "ʃoʊ ʌp",
    "translation": ["объявиться"],
    "createdAt": "2022-06-23 09:11:13",
    "type": 2
  },
  {
    "name": "shut",
    "transcription": "ʃʌt",
    "translation": ["закрывать", "перекрывать"],
    "createdAt": "2022-06-23 09:12:30",
    "type": 1
  },
  {
    "name": "shut down",
    "transcription": "ʃʌt daʊn",
    "translation": ["закрыть", "прекратить работу"],
    "createdAt": "2022-06-23 09:11:11",
    "type": 2
  },
  {
    "name": "shut up",
    "transcription": "ʃʌt ʌp",
    "translation": ["замолчать"],
    "createdAt": "2022-06-23 09:11:22",
    "type": 2
  },
  {
    "name": "sight",
    "transcription": "saɪt",
    "translation": ["взгляд", "зрение"],
    "createdAt": "2022-06-23 09:12:59",
    "type": 1
  },
  {
    "name": "silence",
    "transcription": "ˈsaɪləns",
    "translation": ["тишина", "молчание"],
    "createdAt": "2022-06-23 09:09:13",
    "type": 1
  },
  {
    "name": "silent",
    "transcription": "ˈsaɪlənt",
    "translation": ["тихий"],
    "createdAt": "2022-06-23 09:09:15",
    "type": 1,
    "phrases": ["silence"]
  },
  {
    "name": "since",
    "transcription": "sɪns",
    "translation": ["с тех пор", "с", "после"],
    "createdAt": "2022-06-23 09:08:00",
    "type": 1
  },
  {
    "name": "sink",
    "transcription": "sɪŋk",
    "translation": ["раковина", "тонуть", "топить"],
    "createdAt": "2022-06-23 09:14:06",
    "type": 3,
    "v2": {
      "name": "sank",
      "transcription": "sæŋk"
    },
    "v3": {
      "name": "sunk",
      "transcription": "sʌŋk"
    }
  },
  {
    "name": "sit back",
    "transcription": "sɪt bæk",
    "translation": ["расслабиться", "сидеть сложа руки"],
    "createdAt": "2022-06-23 09:14:20",
    "type": 2
  },
  {
    "name": "sit down",
    "transcription": "sɪt daʊn",
    "translation": ["сядьте"],
    "createdAt": "2022-06-23 09:10:17",
    "type": 2
  },
  {
    "name": "sit up",
    "transcription": "sɪt ʌp",
    "translation": ["сидеть"],
    "createdAt": "2022-06-23 09:13:31",
    "type": 2
  },
  {
    "name": "skyscraper",
    "transcription": "ˈskaɪˌskreɪpər",
    "translation": ["небоскреб"],
    "createdAt": "2022-06-23 09:09:53",
    "type": 1
  },
  {
    "name": "soap",
    "transcription": "soʊp",
    "translation": ["мыло", "мыльный"],
    "createdAt": "2022-06-23 09:12:48",
    "type": 1
  },
  {
    "name": "society",
    "transcription": "səˈsaɪəti",
    "translation": ["общество", "общественность"],
    "createdAt": "2022-06-23 09:13:24",
    "type": 1
  },
  {
    "name": "soft drinks",
    "transcription": "sɑft drɪŋks",
    "translation": ["безалкогольные напитки"],
    "createdAt": "2022-06-23 09:14:40",
    "type": 1
  },
  {
    "name": "some",
    "transcription": "sʌm",
    "translation": ["некоторые", "немного"],
    "createdAt": "2022-06-23 09:13:22",
    "type": 1
  },
  {
    "name": "spectacles",
    "transcription": "ˈspɛktəkəlz",
    "translation": ["очки"],
    "createdAt": "2022-06-23 09:13:48",
    "type": 1
  },
  {
    "name": "spoil",
    "transcription": "spɔɪl",
    "translation": ["портить"],
    "createdAt": "2022-06-23 09:09:58",
    "type": 1
  },
  {
    "name": "stand by",
    "transcription": "stænd baɪ",
    "translation": ["ожидать", "поддержать", "помочь"],
    "createdAt": "2022-06-23 09:12:58",
    "type": 2
  },
  {
    "name": "stand out",
    "transcription": "stænd aʊt",
    "translation": ["выделяться"],
    "createdAt": "2022-06-23 09:12:10",
    "type": 2
  },
  {
    "name": "stand up",
    "transcription": "stænd ʌp",
    "translation": ["встать"],
    "createdAt": "2022-06-23 09:11:21",
    "type": 2
  },
  {
    "name": "starve",
    "transcription": "stɑrv",
    "translation": ["голодать"],
    "createdAt": "2022-06-23 09:11:37",
    "type": 1
  },
  {
    "name": "stay",
    "transcription": "steɪ",
    "translation": ["оставаться"],
    "createdAt": "2022-06-23 09:08:25",
    "type": 1
  },
  {
    "name": "stick",
    "transcription": "stɪk",
    "translation": ["палка", "держаться", "крепить", "прилипать"],
    "createdAt": "2022-06-23 09:13:12",
    "type": 3,
    "v2": {
      "name": "stuck",
      "transcription": "stʌk"
    },
    "v3": {
      "name": "stuck",
      "transcription": "stʌk"
    },
    "phrases": ["sticky"]
  },
  {
    "name": "sticky",
    "transcription": "ˈstɪki",
    "translation": ["липкий"],
    "createdAt": "2022-06-23 09:13:30",
    "type": 1
  },
  {
    "name": "still",
    "transcription": "stɪl",
    "translation": ["все еще", "по прежнему"],
    "createdAt": "2022-06-23 09:07:53",
    "type": 1
  },
  {
    "name": "straight",
    "transcription": "streɪt",
    "translation": ["прямой", "прямо"],
    "createdAt": "2022-06-23 09:13:56",
    "type": 1
  },
  {
    "name": "stretch",
    "transcription": "strɛʧ",
    "translation": ["растягиваться"],
    "createdAt": "2022-06-23 09:13:23",
    "type": 1
  },
  {
    "name": "stuff",
    "transcription": "stʌf",
    "translation": ["вещь"],
    "createdAt": "2022-06-23 09:13:09",
    "type": 1
  },
  {
    "name": "subway",
    "transcription": "ˈsʌˌbweɪ",
    "translation": ["метро"],
    "createdAt": "2022-06-23 09:07:56",
    "type": 1
  },
  {
    "name": "such",
    "transcription": "sʌʧ",
    "translation": ["такой"],
    "createdAt": "2022-06-23 09:08:41",
    "type": 1
  },
  {
    "name": "suddenly",
    "transcription": "ˈsʌdənli",
    "translation": ["вдруг", "внезапно"],
    "createdAt": "2022-06-23 09:13:55",
    "type": 1
  },
  {
    "name": "sunshine",
    "transcription": "ˈsʌnˌʃaɪn",
    "translation": ["солнечный свет"],
    "createdAt": "2022-06-23 09:14:45",
    "type": 1
  },
  {
    "name": "suppose",
    "transcription": "səˈpoʊz",
    "translation": ["предполагать", "полагать"],
    "createdAt": "2022-06-23 09:08:24",
    "type": 1
  },
  {
    "name": "surprisingly",
    "transcription": "sərˈpraɪzɪŋli",
    "translation": ["удивительно"],
    "createdAt": "2022-06-23 09:14:24",
    "type": 1
  },
  {
    "name": "sympathy",
    "transcription": "ˈsɪmpəθi",
    "translation": ["сочувствие", "симпатия"],
    "createdAt": "2022-06-23 09:10:06",
    "type": 1
  },
  {
    "name": "t-shirt",
    "transcription": "ti-ʃɜrt",
    "translation": ["футболка"],
    "createdAt": "2022-06-23 09:09:46",
    "type": 1
  },
  {
    "name": "take away",
    "transcription": "teɪk əˈweɪ",
    "translation": ["забрать", "еда на вынос", "снять (боль)"],
    "createdAt": "2022-06-23 09:13:13",
    "type": 2
  },
  {
    "name": "take back",
    "transcription": "teɪk bæk",
    "translation": ["взять обратно", "забрать"],
    "createdAt": "2022-06-23 09:08:43",
    "type": 2
  },
  {
    "name": "take in",
    "transcription": "teɪk ɪn",
    "translation": ["принимать", "внимательно слушать"],
    "createdAt": "2022-06-23 09:10:51",
    "type": 2
  },
  {
    "name": "take off",
    "transcription": "teɪk ɔf",
    "translation": ["взлетать", "снимать (вещи)", "сбросить вес"],
    "createdAt": "2022-06-23 09:10:21",
    "type": 2
  },
  {
    "name": "take up",
    "transcription": "teɪk ʌp",
    "translation": ["приняться за что-то новое", "занимать (время, место)"],
    "createdAt": "2022-06-23 09:07:33",
    "type": 2
  },
  {
    "name": "tear",
    "transcription": "tɛr",
    "translation": ["рвать"],
    "createdAt": "2022-06-23 09:11:53",
    "type": 3,
    "v2": {
      "name": "tore",
      "transcription": "tɔr"
    },
    "v3": {
      "name": "torn",
      "transcription": "tɔrn"
    }
  },
  {
    "name": "teeth",
    "transcription": "tiθ",
    "translation": ["зубы"],
    "createdAt": "2022-06-23 09:12:52",
    "type": 1
  },
  {
    "name": "terrify",
    "transcription": "ˈtɛrəˌfaɪ",
    "translation": ["пугать", "запугивать"],
    "createdAt": "2022-06-23 09:10:19",
    "type": 1
  },
  {
    "name": "thankfully",
    "transcription": "ˈθæŋkfəli",
    "translation": ["к счастью"],
    "createdAt": "2022-06-23 09:10:03",
    "type": 1
  },
  {
    "name": "the mediterranean",
    "transcription": "ðə ˌmɛdətəˈreɪniən",
    "translation": ["средиземное море"],
    "createdAt": "2022-06-23 09:09:43",
    "type": 1
  },
  {
    "name": "the movie theater",
    "transcription": "ðə ˈmuvi ˈθiətər",
    "translation": ["кинотеатр"],
    "createdAt": "2022-06-23 09:09:06",
    "type": 4
  },
  {
    "name": "the point",
    "transcription": "ðə pɔɪnt",
    "translation": ["смысл"],
    "createdAt": "2022-06-23 09:14:19",
    "type": 4
  },
  {
    "name": "thick",
    "transcription": "θɪk",
    "translation": ["толстый", "густой"],
    "createdAt": "2022-06-23 09:09:19",
    "type": 1,
    "phrases": ["thickness"]
  },
  {
    "name": "thickness",
    "transcription": "ˈθɪknəs",
    "translation": ["толщина"],
    "createdAt": "2022-06-23 09:09:49",
    "type": 1
  },
  {
    "name": "thin",
    "transcription": "θɪn",
    "translation": ["худой", "тонкий"],
    "createdAt": "2022-06-23 09:12:00",
    "type": 1
  },
  {
    "name": "though",
    "transcription": "ðoʊ",
    "translation": ["хотя"],
    "createdAt": "2022-06-23 09:10:44",
    "type": 1
  },
  {
    "name": "thoughtful",
    "transcription": "ˈθɔtfəl",
    "translation": ["заботливый"],
    "createdAt": "2022-06-23 09:14:52",
    "type": 1
  },
  {
    "name": "throughout",
    "transcription": "θruˈaʊt",
    "translation": ["через", "по всему"],
    "createdAt": "2022-06-23 09:07:42",
    "type": 1
  },
  {
    "name": "throw",
    "transcription": "θroʊ",
    "translation": ["бросать", "выбрасывать"],
    "createdAt": "2022-06-23 09:12:50",
    "type": 3,
    "v2": {
      "name": "threw",
      "transcription": "θru"
    },
    "v3": {
      "name": "thrown",
      "transcription": "θroʊn"
    },
    "phrases": ["throw up"]
  },
  {
    "name": "throw up",
    "transcription": "θroʊ ʌp",
    "translation": ["подкидывать", "подбрасывать"],
    "createdAt": "2022-06-23 09:08:50",
    "type": 2
  },
  {
    "name": "thunder",
    "transcription": "ˈθʌndər",
    "translation": ["гром"],
    "createdAt": "2022-06-23 09:09:41",
    "type": 1
  },
  {
    "name": "thunderstorm",
    "transcription": "ˈθʌndərˌstɔrm",
    "translation": ["гроза"],
    "createdAt": "2022-06-23 09:12:39",
    "type": 1
  },
  {
    "name": "tight",
    "transcription": "taɪt",
    "translation": ["тугой", "плотный"],
    "createdAt": "2022-06-23 09:11:05",
    "type": 1
  },
  {
    "name": "tire",
    "transcription": "ˈtaɪər",
    "translation": ["утомлять", "шина"],
    "createdAt": "2022-06-23 09:08:47",
    "type": 1
  },
  {
    "name": "tooth",
    "transcription": "tuθ",
    "translation": ["зуб"],
    "createdAt": "2022-06-23 09:13:15",
    "type": 1
  },
  {
    "name": "toward",
    "transcription": "təˈwɔrd",
    "translation": ["к"],
    "createdAt": "2022-06-23 09:12:17",
    "type": 1
  },
  {
    "name": "towards",
    "transcription": "təˈwɔrdz",
    "translation": ["по направлению к"],
    "createdAt": "2022-06-23 09:09:26",
    "type": 1
  },
  {
    "name": "towel",
    "transcription": "ˈtaʊəl",
    "translation": ["полотенце"],
    "createdAt": "2022-06-23 09:14:35",
    "type": 1
  },
  {
    "name": "treat",
    "transcription": "trit",
    "translation": ["удовольствие", "относиться", "лечить"],
    "createdAt": "2022-06-23 09:14:46",
    "type": 1
  },
  {
    "name": "tricky",
    "transcription": "ˈtrɪki",
    "translation": ["хитрый"],
    "createdAt": "2022-06-23 09:09:48",
    "type": 1
  },
  {
    "name": "trouble",
    "transcription": "ˈtrʌbəl",
    "translation": ["беда"],
    "createdAt": "2022-06-23 09:11:17",
    "type": 1
  },
  {
    "name": "turkey",
    "transcription": "ˈtɜrki",
    "translation": ["индейка"],
    "createdAt": "2022-06-23 09:07:50",
    "type": 1
  },
  {
    "name": "turn around",
    "transcription": "tɜrn əˈraʊnd",
    "translation": ["повернись", "налаживаться"],
    "createdAt": "2022-06-23 09:13:07",
    "type": 2
  },
  {
    "name": "turn back",
    "transcription": "tɜrn bæk",
    "translation": ["повернуться", "вернуться", "повернуть вспять"],
    "createdAt": "2022-06-23 09:11:07",
    "type": 2
  },
  {
    "name": "turn down",
    "transcription": "tɜrn daʊn",
    "translation": ["отказаться", "отклонить"],
    "createdAt": "2022-06-23 09:08:52",
    "type": 2
  },
  {
    "name": "turn over",
    "transcription": "tɜrn ˈoʊvər",
    "translation": ["перевернуть"],
    "createdAt": "2022-06-23 09:09:51",
    "type": 2
  },
  {
    "name": "turnover",
    "transcription": "ˈtɜrˌnoʊvər",
    "translation": ["товарооборот"],
    "createdAt": "2022-06-23 09:14:23",
    "type": 1
  },
  {
    "name": "unable",
    "transcription": "əˈneɪbəl",
    "translation": ["неспособный"],
    "createdAt": "2022-06-23 09:08:14",
    "type": 1
  },
  {
    "name": "uncomfortable",
    "transcription": "ənˈkʌmfərtəbəl",
    "translation": ["неудобный"],
    "createdAt": "2022-06-23 09:12:19",
    "type": 1
  },
  {
    "name": "understandable",
    "transcription": "ˌʌndərˈstændəbəl",
    "translation": ["понятный"],
    "createdAt": "2022-06-23 09:11:55",
    "type": 1
  },
  {
    "name": "uneducated",
    "transcription": "əˈnɛʤʊˌkeɪtɪd",
    "translation": ["необразованный"],
    "createdAt": "2022-06-23 09:07:31",
    "type": 1
  },
  {
    "name": "unexpectedly",
    "transcription": "ˌʌnɪkˈspɛktɪdli",
    "translation": ["неожиданно"],
    "createdAt": "2022-06-23 09:12:05",
    "type": 1
  },
  {
    "name": "unless",
    "transcription": "ənˈlɛs",
    "translation": ["пока не", "если не"],
    "createdAt": "2022-06-23 09:09:07",
    "type": 1
  },
  {
    "name": "unlike",
    "transcription": "ənˈlaɪk",
    "translation": ["в отличие от"],
    "createdAt": "2022-06-23 09:13:44",
    "type": 1
  },
  {
    "name": "upset",
    "transcription": "əpˈsɛt",
    "translation": ["расстройство"],
    "createdAt": "2022-06-23 09:13:21",
    "type": 1
  },
  {
    "name": "urbane",
    "transcription": "ərˈbeɪn",
    "translation": ["вежливый"],
    "createdAt": "2022-06-23 09:09:36",
    "type": 1
  },
  {
    "name": "urge",
    "transcription": "ɜrʤ",
    "translation": ["побуждать", "убеждать"],
    "createdAt": "2022-06-23 09:11:32",
    "type": 1
  },
  {
    "name": "vacation",
    "transcription": "veɪˈkeɪʃən",
    "translation": ["отпуск", "каникулы"],
    "createdAt": "2022-06-23 09:11:01",
    "type": 1
  },
  {
    "name": "valuable",
    "transcription": "ˈvæljəbəl",
    "translation": ["ценный", "полезный"],
    "createdAt": "2022-06-23 09:11:00",
    "type": 1
  },
  {
    "name": "value",
    "transcription": "ˈvælju",
    "translation": ["ценность", "значение", "стоимость"],
    "createdAt": "2022-06-23 09:12:06",
    "type": 1
  },
  {
    "name": "violent",
    "transcription": "ˈvaɪələnt",
    "translation": ["жестокий", "яростный"],
    "createdAt": "2022-06-23 09:09:03",
    "type": 1
  },
  {
    "name": "vocabulary",
    "transcription": "voʊˈkæbjəˌlɛri",
    "translation": ["запас слов"],
    "createdAt": "2022-06-23 09:11:57",
    "type": 1
  },
  {
    "name": "volume",
    "transcription": "ˈvɑljum",
    "translation": ["объем", "громкость"],
    "createdAt": "2022-06-23 09:12:07",
    "type": 1
  },
  {
    "name": "wage",
    "transcription": "weɪʤ",
    "translation": ["зарплата (почасовая или суточная)"],
    "createdAt": "2022-06-23 09:13:06",
    "type": 1
  },
  {
    "name": "waiter",
    "transcription": "ˈweɪtər",
    "translation": ["официант"],
    "createdAt": "2022-06-23 09:08:06",
    "type": 1
  },
  {
    "name": "walk around",
    "transcription": "wɔk əˈraʊnd",
    "translation": ["прогуливаться"],
    "createdAt": "2022-06-23 09:08:58",
    "type": 2
  },
  {
    "name": "walk away",
    "transcription": "wɔk əˈweɪ",
    "translation": ["уходи"],
    "createdAt": "2022-06-23 09:14:07",
    "type": 2
  },
  {
    "name": "walking",
    "transcription": "ˈwɔkɪŋ",
    "translation": ["ходьба"],
    "createdAt": "2022-06-23 09:09:02",
    "type": 1
  },
  {
    "name": "watch out",
    "transcription": "wɑʧ aʊt",
    "translation": ["осторожно"],
    "createdAt": "2022-06-23 09:08:16",
    "type": 2
  },
  {
    "name": "wealth",
    "transcription": "wɛlθ",
    "translation": ["богатство", "благосостояние"],
    "createdAt": "2022-06-23 09:11:14",
    "type": 1
  },
  {
    "name": "wear",
    "transcription": "wɛr",
    "translation": ["носить"],
    "createdAt": "2022-06-23 09:11:12",
    "type": 3,
    "v2": {
      "name": "wore",
      "transcription": "wɔr"
    },
    "v3": {
      "name": "worn",
      "transcription": "wɔrn"
    }
  },
  {
    "name": "whatever",
    "transcription": "ˌwʌˈtɛvər",
    "translation": ["что бы ни", "любой"],
    "createdAt": "2022-06-23 09:07:36",
    "type": 1
  },
  {
    "name": "whether",
    "transcription": "ˈwɛðər",
    "translation": ["будь то"],
    "createdAt": "2022-06-23 09:09:35",
    "type": 1
  },
  {
    "name": "which",
    "transcription": "wɪʧ",
    "translation": ["какой", "который"],
    "createdAt": "2022-06-23 09:13:36",
    "type": 1
  },
  {
    "name": "whom",
    "transcription": "hum",
    "translation": ["кого", "кому"],
    "createdAt": "2022-06-23 09:08:07",
    "type": 1
  },
  {
    "name": "wink",
    "transcription": "wɪŋk",
    "translation": ["подмигивание"],
    "createdAt": "2022-06-23 09:09:50",
    "type": 1
  },
  {
    "name": "withdraw",
    "transcription": "wɪðˈdrɔ",
    "translation": ["снимать средства со счета"],
    "createdAt": "2022-06-23 09:13:57",
    "type": 1
  },
  {
    "name": "within",
    "transcription": "wɪˈðɪn",
    "translation": ["в пределах"],
    "createdAt": "2022-06-23 09:14:50",
    "type": 1
  },
  {
    "name": "wonder",
    "transcription": "ˈwʌndər",
    "translation": ["удивляться"],
    "createdAt": "2022-06-23 09:08:30",
    "type": 1
  },
  {
    "name": "work out",
    "transcription": "wɜrk aʊt",
    "translation": ["проработать", "разработать"],
    "createdAt": "2022-06-23 09:07:55",
    "type": 2
  },
  {
    "name": "write down",
    "transcription": "raɪt daʊn",
    "translation": ["записывать"],
    "createdAt": "2022-06-23 09:13:11",
    "type": 2
  },
  {
    "name": "your spirits",
    "transcription": "jʊər ˈspɪrɪts",
    "translation": ["ваше настроение"],
    "createdAt": "2022-06-23 09:12:01",
    "type": 4
  },
  {
    "name": "yourself",
    "transcription": "jərˈsɛlf",
    "translation": ["сам", "себя"],
    "createdAt": "2022-06-23 09:07:34",
    "type": 1
  },
  {
    "name": "acceptable",
    "transcription": "ækˈsɛptəbəl",
    "translation": ["приемлемый"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "accidental",
    "transcription": "ˌæksəˈdɛntəl",
    "translation": ["случайный"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "ally",
    "transcription": "ˈælaɪ",
    "translation": ["союзник", "объединяться"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "among",
    "transcription": "əˈmʌŋ",
    "translation": ["среди", "между"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "anniversary",
    "transcription": "ˌænəˈvɜrsəri",
    "translation": ["годовщина", "юбилей"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "belong",
    "transcription": "bɪˈlɔŋ",
    "translation": ["принадлежать", "относиться"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "broad",
    "transcription": "brɔd",
    "translation": ["широкий"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "brutal",
    "transcription": "ˈbrutəl",
    "translation": ["жестокий"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "closely",
    "transcription": "ˈkloʊsli",
    "translation": ["тесно", "близко"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "copycat",
    "transcription": "ˈkɑpiˌkæt",
    "translation": ["подражатель"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "cruel",
    "transcription": "ˈkruəl",
    "translation": ["жестокий"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "curiosity",
    "transcription": "ˌkjʊriˈɑsəti",
    "translation": ["любопытство"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "demand",
    "transcription": "dɪˈmænd",
    "translation": ["требование", "требовать"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "develop",
    "transcription": "dɪˈvɛləp",
    "translation": ["развивать"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "dishonest",
    "transcription": "dɪˈsɑnəst",
    "translation": ["нечестный"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "dispute",
    "transcription": "dɪˈspjut",
    "translation": ["спор"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "envy",
    "transcription": "ˈɛnvi",
    "translation": ["завидовать"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "establish",
    "transcription": "ɪˈstæblɪʃ",
    "translation": ["учреждать"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "furniture",
    "transcription": "ˈfɜrnɪʧər",
    "translation": ["мебель"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "galvanize",
    "transcription": "ˈgælvəˌnaɪz",
    "translation": ["электризовать", "стимулировать", "оживлять"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "greenback",
    "transcription": "ˈgrinˌbæk",
    "translation": ["доллар США"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "guilty",
    "transcription": "ˈgɪlti",
    "translation": ["виновный"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "hedge",
    "transcription": "hɛʤ",
    "translation": ["изгородь", "ограждать"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "highly",
    "transcription": "ˈhaɪli",
    "translation": ["очень", "высоко"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "influence",
    "transcription": "ˈɪnfluəns",
    "translation": ["влияние"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "injury",
    "transcription": "ˈɪnʤəri",
    "translation": ["рана"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "involve",
    "transcription": "ɪnˈvɑlv",
    "translation": ["вовлекать"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "likely",
    "transcription": "ˈlaɪkli",
    "translation": ["вероятно"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "mainly",
    "transcription": "ˈmeɪnli",
    "translation": ["в основном"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "misbehave",
    "transcription": "ˌmɪsbəˈheɪv",
    "translation": ["плохо себя вести"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "obey",
    "transcription": "oʊˈbeɪ",
    "translation": ["подчиниться"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "occuse",
    "transcription": "əˈkjuz",
    "translation": ["занимать"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "overnight",
    "transcription": "ˈoʊvərˈnaɪt",
    "translation": ["с ночевкой", "ночной"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "permit",
    "transcription": "ˈpɜrˌmɪt",
    "translation": ["разрешать"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "property",
    "transcription": "ˈprɑpərti",
    "translation": ["имущество", "собственность"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "proud",
    "transcription": "praʊd",
    "translation": ["гордый"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "reduce",
    "transcription": "rəˈdus",
    "translation": ["уменьшать"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1,
    "phrases": ["reduction"]
  },
  {
    "name": "reduction",
    "transcription": "rəˈdʌkʃən",
    "translation": ["снижение"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "refuse",
    "transcription": "rɪˈfjuz",
    "translation": ["мусор"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "severely",
    "transcription": "səˈvɪrli",
    "translation": ["строго"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "spot",
    "transcription": "spɑt",
    "translation": ["место"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "suspect",
    "transcription": "ˈsʌˌspɛkt",
    "translation": ["подозревать"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "term",
    "transcription": "tɜrm",
    "translation": ["срок"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "till",
    "transcription": "tɪl",
    "translation": ["пока"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "unfair",
    "transcription": "ənˈfɛr",
    "translation": ["несправедливый"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "unnatural",
    "transcription": "ənˈnæʧərəl",
    "translation": ["неестественный"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "warmth",
    "transcription": "wɔrmθ",
    "translation": ["теплота"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  },
  {
    "name": "worthless",
    "transcription": "ˈwɜrθləs",
    "translation": ["бесполезный"],
    "createdAt": "2022-07-17 08:53:19",
    "type": 1
  }
]'::jsonb) AS elem;

-- Вставка английских слов
INSERT INTO dictionary (name, lang, type)
SELECT name, 'en', type FROM temp_en_words;

-- Вставка русских переводов
INSERT INTO dictionary (name, lang, type)
SELECT DISTINCT jsonb_array_elements_text(translations), 'ru', 1
FROM temp_en_words;

-- Вставка переводов (прямые)
INSERT INTO translation (dictionary_id, translation_id)
SELECT e.id, r.id
FROM temp_en_words t
JOIN dictionary e ON e.name = t.name AND e.lang = 'en'
CROSS JOIN jsonb_array_elements_text(t.translations) AS tr
JOIN dictionary r ON r.name = tr AND r.lang = 'ru';

-- Вставка обратных переводов
INSERT INTO translation (dictionary_id, translation_id)
SELECT r.id, e.id
FROM temp_en_words t
JOIN dictionary e ON e.name = t.name AND e.lang = 'en'
CROSS JOIN jsonb_array_elements_text(t.translations) AS tr
JOIN dictionary r ON r.name = tr AND r.lang = 'ru';

-- Вставка предложений
INSERT INTO sentence (text_en, text_ru)
SELECT 
    (s->>'sentence')::VARCHAR,
    (s->>'translation')::VARCHAR
FROM temp_en_words t
CROSS JOIN jsonb_array_elements(t.sentences) AS s;

-- Связывание предложений с записями словаря
INSERT INTO dictionary_sentence (dictionary_id, sentence_id)
SELECT d.id, s.id
FROM temp_en_words t
JOIN dictionary d ON d.name = t.name AND d.lang = 'en'
CROSS JOIN jsonb_array_elements(t.sentences) AS sent
JOIN sentence s ON s.text_en = (sent->>'sentence')::VARCHAR
UNION
SELECT dr.id, s.id
FROM temp_en_words t
CROSS JOIN jsonb_array_elements_text(t.translations) AS tr
JOIN dictionary dr ON dr.name = tr AND dr.lang = 'ru'
CROSS JOIN jsonb_array_elements(t.sentences) AS sent
JOIN sentence s ON s.text_en = (sent->>'sentence')::VARCHAR;

-- Удаление временной таблицы
DROP TABLE temp_en_words;
