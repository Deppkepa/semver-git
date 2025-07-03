# semver-git

## Инструкция по установки: 

1. Установка проекта через сабмодуль:
```
git submodule add https://github.com/Deppkepa/semver-git.git scripts
git submodule update --init
```
2. Собрать файлы:
```
cd scripts/describe && go build -o describe
cd scripts/version_check && go build -o version_check
```
Потом запускаем нужный файл вот так:
```
./scripts/describe/describe -h
./scripts/version_check/version_check -h
```
