# redirectsrv

This application listens TCP :5000. 
It was hardcoded. Please, make sure the port is free.

`release/nginx` - This directory contains all the server definitions.

`release/backup.sql` - Scheme of DB

`release/config.yml` - File is needed for application configuration.
 The file contains default value for quick start. 
 
 `release/redirectsrv` - FastCGI App.
 
 Simple example: `localhost/TestApp?username=Anil&ran=0.97584378943&pageURL=http://xyz.com`
