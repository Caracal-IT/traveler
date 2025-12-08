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
- The actual application wiring to use SQLite is out of scope of this change; this only provides the local DB assets and make targets.
