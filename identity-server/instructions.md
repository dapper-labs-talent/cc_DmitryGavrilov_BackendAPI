START SERVER FROM DEVELOPER MACHINE, DB IN DOCKER

The current repository contains  a docker-compose file with configured Postgres image
1. We need to start a Postgres container, for that simply use the docker-compose up command, which will pull a Postgres image and it will create and start the container. The docker-compose file is configured to run init.sql (./scripts/postgres/init.sql) file that will create a axiomzen database, with a user the same user. axiomze:axiomze. 
2. the repository contains "identity-server" binary file compiled on mac (intel)
3. please start the above mentioned file with a "migrate" command - e.g. - ./identity-server migrate. This command will create a schema in the remote database
4. now, we are ready to start the identity-server. The default (root command) will start the server, therefore you just need to start the binary file with no commands - e.g. - ./identity-server.   A configuration file (./config/config.ini) contains a bunch of settings. One of them is a port to listen to. By default, it is going to listen:8888
5. A Postman collection, environment are located in postman folder, that contains 4 endpoints to send request



START SERVER AND DB in DOCKER

1. We need to build two containers: one for the server, and another one for the database. Please run the following command from the root: - ./runInDocker/db-server-compose.sh build
2. Now we are ready to start containers. Please use the following command: ./runInDocker/db-server-compose.sh up
3. If no errors happened, the server should be available on :8080 port (could be changed in configuration file. For docker deployment please refer to ./runInDocker/config.ini)
4. A Postman collection, environment are located in postman folder, that contains 4 endpoints to send request
