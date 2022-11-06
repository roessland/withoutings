# Install PostgreSQL on MacOS (ARM)

Copied from my personal notes [[PostgreSQL setup on MacOS]].


## Install

- Install v14 since that is the stable version for Ubuntu
  - https://www.postgresql.org/support/versioning/

```bash
brew install postgresql@14
```

## Start service

- Start the service now and on every boot
  - https://wiki.postgresql.org/wiki/Homebrew
  - Use `run` instead of `start` to only start it once right now.
```bash
brew services start postgresql@15
```

## Connect

Check if psql is already on path.
```bash
psql postgres
```

If not, find the path:
```bash
brew info postgresql@15
ls /opt/homebrew/opt/postgresql@15/bin
echo 'export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"' >> ~/.zshrc
exec zsh
psql postgres
```

## Find data directory

```bash
ps aux | grep postgres
cd /opt/homebrew/var/postgresql@15
```

## Change password for $USER
Make a new password and add it in a password manager.
```bash
psql postgres
\password
# type a long password
```

### Add password to `.pgpass`
See [[pgpass]].
```
cd
touch .pgpass
chmod 600 .pgpass
echo "localhost:5432:*:$USER:<POSTGRESPASSWORD>" >> .pgpass
```

## Remove "trust" authentication in `pg_hba.conf`
```bash
cd /opt/homebrew/var/postgresql@15
vi pg_hba.conf
```
Change it to something like this:
```ini
# "local" is for Unix domain socket connections only
local   all             all                                     scram-sha-256
# IPv4 local connections:
host    all             all             127.0.0.1/32            scram-sha-256
# IPv6 local connections:
host    all             all             ::1/128                 scram-sha-256
```
Then restart the [[brew]] service:
```bash
brew services restart postgresql@15
```