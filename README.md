# map_runners

## Idea
Веб сайт на котором отображается карта мира. Каждая страна имеет два цвета: серый и зеленый. Серый если пользователь еще не бегал на ее территории, зеленый если есть хотя бы одна тренировка больше 5 км. Изначально пользователь должен авторизоваться в сервисе Strava. После этого сервис на backend скачивает все тренировки пользователя и складывает их в PostgreSQL. По мере появляния тренировок в БД, пользователь видит на карте обновления.

## ТЗ

Цель: Разработать веб-сайт с интерактивной картой мира, отображающей страны, в которых пользователь совершал тренировки, полученные из сервиса Strava.

Функциональные требования:
1. Авторизация пользователя через сервис Strava.
2. Автоматическое получение данных о тренировках пользователя из Strava и сохранение их в базе данных PostgreSQL. Обновление данных должно происходить раз в день.
3. Отображение интерактивной карты мира с возможностью масштабирования и перемещения.
4. Окрашивание стран в зеленый цвет, если пользователь совершил хотя бы одну тренировку более 5 км на их территории. В противном случае, страна остается серой.
5. При наведении или клике на страну должна отображаться дополнительная информация: количество тренировок, общая дистанция, даты первой и последней тренировки.
6. Для целей тестирования и разработки необходимо предусмотреть возможность выбора периода, за который отображаются тренировки, но в продакшн-версии всегда должны выбираться все тренировки.

Нефункциональные требования:
1. Минималистичный дизайн веб-сайта.
2. Использование современных веб-технологий и фреймворков для frontend и backend разработки.
3. Обеспечение безопасности передачи и хранения данных пользователя.
4. Оптимизация производительности и скорости загрузки страниц.

Технические требования:
1. Frontend: HTML, CSS, JavaScript, фреймворк (React, Angular, Vue.js) на выбор разработчика.
2. Backend: Go.
3. База данных: PostgreSQL.
4. Интеграция с API Strava для получения данных о тренировках пользователя.
5. Развертывание веб-сайта на выбранной платформе (AWS, Heroku, Digital Ocean и т.д.).

Список задач:
1. Настройка окружения разработки и инициализация проекта.
2. Разработка базовой структуры frontend и backend частей приложения.
3. Интеграция с API Strava и реализация механизма авторизации пользователя.
4. Разработка функционала получения и сохранения данных о тренировках пользователя в базу данных PostgreSQL.
5. Реализация интерактивной карты мира с возможностью масштабирования и перемещения.
6. Реализация логики окрашивания стран в зависимости от наличия тренировок пользователя.
7. Разработка функционала отображения дополнительной информации о тренировках при наведении или клике на страну.
8. Оптимизация производительности и скорости загрузки страниц.
9. Тестирование и отладка приложения.
10. Развертывание веб-сайта на выбранной платформе.