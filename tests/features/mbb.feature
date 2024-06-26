# language: ru
Функционал: Защита от брутфорс-атак
  Как система авторизации
  Я хочу проверить, является ли попытка авторизации легитимной или это попытка брутфорс-атаки

  Сценарий: Успешная авторизация
    Допустим, у нас есть IP-адрес "192.168.2.1", логин "user2" и пароль "password2"
    Когда я вызываю метод Check с IP-адресом "192.168.2.1", логином "user2" и паролем "password2"
    Тогда метод Check возвращает "true"

  Сценарий: Несколько успешных авторизаций
    Допустим, у нас есть IP-адрес "192.168.3.1", логин "user3" и пароль "password3"
    Когда я вызываю метод Check с IP-адресом "192.168.3.1", логином "user3" и паролем "password3" 2 раза за 12 секунд
    Тогда метод Check возвращает "true"

  Сценарий: Брутфорс-атака
    Допустим, у нас есть IP-адрес "192.168.4.1", логин "user4" и пароль "password4"
    Когда я вызываю метод Check с IP-адресом "192.168.4.1", логином "user4" и паролем "password4" 2 раза за 5 секунд
    Тогда метод Check возвращает "false"

  Сценарий: IP-адрес в черном списке
    Допустим, у нас есть IP-адрес "192.168.5.1", логин "user5" и пароль "password5"
    Когда я вызываю функцию Deny с подсетью "192.168.5.0/24"
    И я вызываю метод Check с IP-адресом "192.168.5.1", логином "user5" и паролем "password5"
    Тогда метод Check возвращает "false"

  Сценарий: IP-адрес в белом списке
    Допустим, у нас есть IP-адрес "192.168.6.1", логин "user6" и пароль "password6"
    Когда я вызываю функцию Allow с подсетью "192.168.6.0/24"
    И я вызываю метод Check с IP-адресом "192.168.6.1", логином "user6" и паролем "password6" 22 раза за 5 секунд
    Тогда метод Check возвращает "true"
