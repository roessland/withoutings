# Install PostgreSQL on Ubuntu

```
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql.service

edit pg_hba.conf

sudo systemctl restart postgresql

sudo -u postgres psql

```

Append to .pgpass
```
sudo su <withoutings-user>
cd
touch .pgpass
chmod 600 .pgpass
echo "localhost:5432:wot:wotsa:<POSTGRESPASSWORD>" >> .pgpass
echo "localhost:5432:wot:wotrw:<POSTGRESPASSWORD>" >> .pgpass
```

Append to .profile for application user so that you can connect with 
simply "sudo -u <withoutings-user> -i psql".
```
sudo su <withoutings-user>
export PGHOST=localhost
export PGUSER=wotsa
export PGDATABASE=wot
```