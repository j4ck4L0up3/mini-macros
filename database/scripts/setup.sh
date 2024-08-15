#!/bin/bash

hostname="0.0.0.0"
if [$ENV != "development"]; then
  hostname="$POSTGRES_HOST"
fi

psql -v ON_ERROR_STOP=1 --host "$hostname" --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  REVOKE ALL ON SCHEMA public FROM public;
  
  CREATE SCHEMA [IF NOT EXISTS] "$DB_SCHEMA_NAME";
  
  CREATE ROLE administrator WITH LOGIN PASSWORD "$DB_ADMIN_PASSWORD" INHERIT; 
  CREATE ROLE admin WITH NOLOGIN NOINHERIT;
  
  GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA minimacros TO admin;
  GRANT admin TO administrator;
EOSQL
