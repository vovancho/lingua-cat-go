#!/usr/bin/env bash
set -euo pipefail

###############
#  ПАРАМЕТРЫ  #
###############
# Значения по умолчанию
REQUESTS=100
CONCURRENCY=10
KEYCLOAK_URL="http://keycloak.localhost/realms/lingua-cat-go/protocol/openid-connect/token"
CLIENT_ID="lingua-cat-go-dev"
CLIENT_SECRET="GatPbS9gsEfplvCpiNitwBdmIRc0QqyQ"
USERNAME="dummy-user"
PASSWORD="password"
API_BASE="http://api.lingua-cat-go.localhost/exercise/v1"

#########################
#  РАЗБОР ОПЦИЙ SHELL   #
#########################
print_usage() {
  cat <<EOF
Usage: $0 [--requests=<N>] [--concurrency=<M>]

  --requests      Количество итераций (аналог -n в ab). Default: $REQUESTS
  --concurrency   Количество параллельных потоков (аналог -c в ab). Default: $CONCURRENCY
EOF
  exit 1
}

# разбор long options
for arg in "$@"; do
  case $arg in
    --requests=*)
      REQUESTS="${arg#*=}"
      shift
      ;;
    --concurrency=*)
      CONCURRENCY="${arg#*=}"
      shift
      ;;
    --help|-h)
      print_usage
      ;;
    *)
      echo "Unknown option: $arg"
      print_usage
      ;;
  esac
done

###################################
#  ПОЛУЧАЕМ JWT-ТOKEN ОДИН РАЗ    #
###################################
echo "Получаем access_token из Keycloak..."
TOKEN_RESPONSE="$(curl -s -X POST "$KEYCLOAK_URL" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password&scope=openid&client_id=$CLIENT_ID&client_secret=$CLIENT_SECRET&username=$USERNAME&password=$PASSWORD")"

JWT_TOKEN=$(jq -r '.access_token' <<<"$TOKEN_RESPONSE")

if [[ -z "$JWT_TOKEN" || "$JWT_TOKEN" == "null" ]]; then
  echo "Ошибка: не удалось получить access_token." >&2
  echo "Ответ Keycloak: $TOKEN_RESPONSE" >&2
  exit 1
fi
echo "Получен access_token."

#############################
#  ФУНКЦИЯ RUN_CHAIN       #
#############################
run_chain() {
  local iteration=$1
  echo "Starting iteration $iteration" >> "$OUTPUT_FILE"

  # 1. Создать упражнение
  local CRE=$(curl -s -X POST "$API_BASE/exercise" \
    -H "Content-Type: application/json" \
    -H "Accept: application/json" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -d '{"lang":"en","task_amount":1}')
  local EXERCISE_ID=$(jq -r '.data.exercise.id // empty' <<<"$CRE")
  if [[ -z "$EXERCISE_ID" ]]; then
    echo "Iteration $iteration: Failed to create exercise" >> "$OUTPUT_FILE"
    return 1
  fi
  echo "Iteration $iteration: Created exercise ID $EXERCISE_ID" >> "$OUTPUT_FILE"

  # 2. Создать задачу
  local CTR=$(curl -s -X POST "$API_BASE/exercise/$EXERCISE_ID/task" \
    -H "Accept: application/json" \
    -H "Authorization: Bearer $JWT_TOKEN")
  local TASK_ID=$(jq -r '.data.task.id // empty' <<<"$CTR")
  local WORD_CORRECT_ID=$(jq -r '.data.task.word_correct.id // empty' <<<"$CTR")
  if [[ -z "$TASK_ID" || -z "$WORD_CORRECT_ID" ]]; then
    echo "Iteration $iteration: Failed to create task or parse IDs" >> "$OUTPUT_FILE"
    return 1
  fi
  echo "Iteration $iteration: Task $TASK_ID, correct word $WORD_CORRECT_ID" >> "$OUTPUT_FILE"

  # 3. Выбрать слово
  local SLR=$(curl -s -X POST "$API_BASE/exercise/$EXERCISE_ID/task/$TASK_ID/word-selected" \
    -H "Content-Type: application/json" \
    -H "Accept: application/json" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -d "{\"word_select\":$WORD_CORRECT_ID}")
  if [[ -z "$SLR" ]]; then
    echo "Iteration $iteration: Failed to select word" >> "$OUTPUT_FILE"
    return 1
  fi
  echo "Iteration $iteration: Selected word for task $TASK_ID" >> "$OUTPUT_FILE"
}

#############################
#   ЗАПУСК В ПАРАЛЛЕЛИ      #
#############################
OUTPUT_FILE="ab_test_results.txt"
> "$OUTPUT_FILE"

echo "Запуск $REQUESTS итераций с concurrency=$CONCURRENCY..."
for ((i = 1; i <= REQUESTS; i++)); do
  run_chain "$i" &
  # ограничиваем число фоновых задач
  if (( $(jobs -rp | wc -l) >= CONCURRENCY )); then
    wait -n
  fi
done
wait

echo "Тест завершён. Результаты в $OUTPUT_FILE."
