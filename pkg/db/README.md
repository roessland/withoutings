# Database

## Known issues
If the primary keys are manually modified, they can conflict with the next auto-incremented value, causing inserts to fail with `ERROR: duplicate key value violates unique constraint \"mytable_pkey\`. This won't be fixed or handled, just don't modify the primary keys.