FROM php:8.1-fpm

RUN docker-php-ext-install pdo_mysql

WORKDIR /app
COPY  . .

RUN curl -sS https://getcomposer.org/installer | php

RUN mv composer.phar /usr/bin/composer

RUN composer install

EXPOSE 8000

