#!/bin/bash

# Параметры
JWT_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ3c1A2RW9SZUFYYlRmWTZBMTU3NEt4SFdPZlZXUTJwNTN3eEtIUjR2N0VFIn0.eyJleHAiOjE3NDY4MDIxMzYsImlhdCI6MTc0Njc2NjEzNiwianRpIjoiYzVlODIwN2YtOWFjNS00NDNiLWEzNzAtMTJkZGE4MjdjODYwIiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmxvY2FsaG9zdC9yZWFsbXMvbGluZ3VhLWNhdC1nbyIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiIxMWM2NWU0MS0yNDk2LTQzYWYtYWM0Yy1kYWE4OThjMjQ2NjQiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJsaW5ndWEtY2F0LWdvLWRldiIsInNpZCI6IjQyNmQxNTZlLWNiZWItNDUzNC1hMjliLTA3YjdkYjg3ZDIzYSIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovL2xpbmd1YS1jYXQtZ28ubG9jYWxob3N0Il0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLWxpbmd1YS1jYXQtZ28iLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsibGluZ3VhLWNhdC1nby1kZXYiOnsicm9sZXMiOlsiVklTSVRPUiJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJkdW1teS11c2VyIiwiZW1haWwiOiJkZXYtdXNlckBtYWlsLmRldiJ9.WudelMYtQexBOALS8DT_m-zAnGJGIRnuAeuEDdwvOCb1CFEIijXs5pHN5kL7FsnomPfti1ItCZCWl-81nzZUuNpf6d_azWoHzfYulJL5AY8jFyFZ44MuAsQAXjK9job9rran0jr84VhfZDQR60POKq7QQIiJsWTwNe06d2BsOO4R1WZwL28l1G9yPVmMUPaaJWxRaME2DpXtNW-ysq42t621QqJe-VCxWx_WAw9ZBj1dQBq-CZ9udI6eZkR-SUPBaOdiTP0_1dlrLQcT4raRMhsVvS-7pE9entTozdSN-88XtOn7o7hdkucKacfkIsW8jo37UYoCTrAmb6UDupVUag"  # Замените на ваш JWT-токен
REQUESTS=100            # Количество итераций (аналог -n в ab)
CONCURRENCY=100          # Количество параллельных потоков (аналог -c в ab)
OUTPUT_FILE="ab_test_results.txt"

# Функция для выполнения цепочки запросов
run_chain() {
    local iteration=$1
    echo "Starting iteration $iteration" >> $OUTPUT_FILE

    # 1. Создать упражнение
    EXERCISE_RESPONSE=$(curl -s -X POST "http://api.lingua-cat-go.localhost/exercise/v1/exercise" \
        -H "Content-Type: application/json" \
        -H "Accept: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d '{"lang":"en","task_amount":1}')

    if [ -z "$EXERCISE_RESPONSE" ]; then
        echo "Iteration $iteration: Failed to create exercise" >> $OUTPUT_FILE
        return 1
    fi

    EXERCISE_ID=$(echo $EXERCISE_RESPONSE | jq -r '.data.exercise.id')
    if [ "$EXERCISE_ID" == "null" ] || [ -z "$EXERCISE_ID" ]; then
        echo "Iteration $iteration: Failed to parse exercise_id" >> $OUTPUT_FILE
        return 1
    fi
    echo "Iteration $iteration: Created exercise with ID $EXERCISE_ID" >> $OUTPUT_FILE

    # 2. Создать задачу
    TASK_RESPONSE=$(curl -s -X POST "http://api.lingua-cat-go.localhost/exercise/v1/exercise/$EXERCISE_ID/task" \
        -H "Content-Type: application/json" \
        -H "Accept: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN")

    if [ -z "$TASK_RESPONSE" ]; then
        echo "Iteration $iteration: Failed to create task" >> $OUTPUT_FILE
        return 1
    fi

    TASK_ID=$(echo $TASK_RESPONSE | jq -r '.data.task.id')
    WORD_CORRECT_ID=$(echo $TASK_RESPONSE | jq -r '.data.task.word_correct.id')
    if [ "$TASK_ID" == "null" ] || [ -z "$TASK_ID" ] || [ "$WORD_CORRECT_ID" == "null" ] || [ -z "$WORD_CORRECT_ID" ]; then
        echo "Iteration $iteration: Failed to parse task_id or word_correct_id" >> $OUTPUT_FILE
        return 1
    fi
    echo "Iteration $iteration: Created task with ID $TASK_ID, word_correct_id $WORD_CORRECT_ID" >> $OUTPUT_FILE

    # 3. Выбрать корректное слово
    SELECT_WORD_RESPONSE=$(curl -s -X POST "http://api.lingua-cat-go.localhost/exercise/v1/exercise/$EXERCISE_ID/task/$TASK_ID/word-selected" \
        -H "Content-Type: application/json" \
        -H "Accept: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d "{\"word_select\":$WORD_CORRECT_ID}")

    if [ -z "$SELECT_WORD_RESPONSE" ]; then
        echo "Iteration $iteration: Failed to select word" >> $OUTPUT_FILE
        return 1
    fi
    echo "Iteration $iteration: Successfully selected word for task $TASK_ID" >> $OUTPUT_FILE
}

# Очистка файла результатов
> $OUTPUT_FILE

# Запуск цепочек в параллельных потоках
for ((i=1; i<=REQUESTS; i++)); do
    # Запускаем цепочку в фоновом режиме, ограничивая количество параллельных процессов
    run_chain $i &

    # Ограничиваем количество параллельных процессов до CONCURRENCY
    if (( $(jobs -r | wc -l) >= CONCURRENCY )); then
        wait -n
    fi
done

# Ожидаем завершения всех фоновых процессов
wait

echo "Test completed. Results written to $OUTPUT_FILE"
