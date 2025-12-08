Traveler local SQLite database

This folder holds the local development SQLite database assets.

What is included
- schema.sql: The canonical SQL schema for initializing the local DB.
- .gitignore: Ensures generated .db files are not committed.

Recommended usage
1) Create or recreate the database from schema.sql:
   make db-init

2) Remove the generated database file:
   make db-clean

Notes
- By default, the database file will be created at db/traveler.db.
- The application now initializes the SQLite database automatically on startup using the schema at db/schema.sql. If the DB file doesn't exist, it will be created, and the schema (including the specials seeds) will be applied idempotently.
- You can still manage the DB manually with the make targets above during development.
