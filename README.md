# Getting Started with Docker-Compose

This project was created to use docker-compose to pull Azure edge sql for the database

To start off with using this project:
1. Clone teh repo
2. Use docker-compose in the project directory
#### `sudo docker-compose up -d`
3. After docker is up and running, you can try to run migrations in the migration folder by connecting to Azure studio to connect to the Azure sql edge and in turn create
database by the name `CompanyDB` and then run the queries in each of the migration file on the database table as I was having issues using goose to run migrations
  or while in the top directory of the project run: 
#### `./goose.sh up`
4. After docker is up and running, to set up the application on a unix kernel system run :
#### `./run.sh`
5. 


Once APP is up it runs on:
Open [http://localhost:7070] to view in your browser.

Still in development: 
Check the playground for the documentation and schema to run
