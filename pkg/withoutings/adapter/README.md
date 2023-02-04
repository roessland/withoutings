# adapter

Contains implementations of domain repos tied to specific storage engines.

For example, AccountRepo using PostgreSQL storage.

## Reading
- https://threedots.tech/post/introducing-clean-architecture/
  - If the project grows in size, you may find it helpful to add another level of subdirectories. For example `adapters/account/pg_repo.go`.