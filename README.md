Скрипт для миграции -   
migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' up  
Скрипт для развертывания Postgres в Docker -   
run —name tg-bot-messages -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=bot-msg -d postgres
