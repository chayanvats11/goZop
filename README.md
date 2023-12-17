# goZop Garage Application

## Description
This is goZop Garage Application which has been created using Go language and GoFr Framework (Supports accelerated microservice development)

---

## Setup
0. Clone this repository
1. Setup Go on your system (https://go.dev/doc/install)<br>
2. Setup MySQL on local (https://dev.mysql.com/doc/mysql-getting-started/en/)<br>
    I. After Setting up MySQL ensure you have User Created and you have password for it <br>
    ```
    mysql -u root
    CREATE USER 'newuser'@'localhost' IDENTIFIED BY 'password';
    GRANT ALL PRIVILEGES ON *.* TO 'newuser'@'localhost' WITH GRANT OPTION;
    ```
    II. Login to MySQL and Create Database named "cars" <br>
    - For Workbench:
         https://dev.mysql.com/doc/workbench/en/<br>
    - For Terminal: ```mysql -u newuser -p -h 127.0.0.1 -P 3306``` <br>
    Create database: ```create database cars;```<br>

3. Replace with your MySQL username and password in configs/.env file<br>
    - DB_USER=YOUR_USERNAME
    - DB_PASSWORD=YOUR_PASSWORD

4. To run test go to your project repository and use this commmand
    - ```go test```

5. Run the Application using this command<br>
    - ``` go run main.go```

---

## System Design

![System Design] (system_design.png)

## Database Schema and Flow

![Database Schema] (db_schema_flow.png)