# MLANG

## Cодержание
* [Использование](#Использование)
* [Грамматика языка](#Грамматика-языка)
* [Основные правиля языка](#Основные-правила-языка)
* [Примеры](#Примеры)
* [Структура проекта](#Cтруктура-проекта)

## Использование
Для запуска интерактивного выполните:
```bash
go run main.go
```

Чтобы выполнить заготовленный файл с кодом укажите путь до нужного файла:
```bash
go run main.go program.mlang -o result
```

Здесь ключ -o указывает путь до файла, в которых будет записан результат выполнения программы

## Грамматика языка

## Основные правила языка

MLang поддерживает 3 типа данных с которыми может работать пользователь:
* integer - целые 64битные числа
* boolean - логический тип, в программе обозначается литералами *true* и *false*
* function - функции

## Примеры

## Cтруктура проекта

Интерпретатор языка состоит из 7 основных пакетов:

* **token** - Описание разрешенных токенов в языке
* **lexer** - Производит преобразование исходного кода на mlang в последовательность токенов для последующей обработки парсером
* **ast** - Описание абстрактного синтаксического дерева, задающего структуру выполнения программы
* **parser** - Преобразует последовательность токенов полученных из модуля лексера в абстрактное синтаксическое дерево
* **evaluator** - Производит разбор *аст*, выполняя описаннные в нем вычисления
* **object** - Описание внутренних объектов и типов языка
* **repl** - Собственно интерпретатор