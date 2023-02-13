##  test_colver

Скрипт, который проходит опрос, расположенный по адресу http://185.204.3.165

## Запуск
1. Склонировать репозиторий:
```
git clone git@github.com:Krynegal/test_solver.git
```
2. В папке проекта прописать:
```
go run main.go -w <кол-во параллельных потоков>
```
или
```
go run main.go 
```
В этом случае количество потоков будет равно значению WORKERS_NUM в /configs/.env (1 по умолчанию)