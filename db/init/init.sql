DROP DATABASE IF EXISTS `isucon`;
CREATE DATABASE `isucon`;

DROP USER IF EXISTS 'isucon'@'localhost';
CREATE USER 'isucon'@'localhost' IDENTIFIED BY 'isucon';
GRANT ALL PRIVILEGES ON `isucon`.* TO 'isucon'@'localhost';

DROP USER IF EXISTS 'isucon'@'%';
CREATE USER 'isucon'@'%' IDENTIFIED BY 'isucon';
GRANT ALL PRIVILEGES ON `isucon`.* TO 'isucon'@'%';