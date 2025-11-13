-- ============================================================================
-- PostgreSQL Setup Script for Bank Assistant Orchestrator
-- ============================================================================
--
-- This script creates the necessary PostgreSQL database and user
-- Run this script as a PostgreSQL superuser (usually 'postgres')
--
-- Usage:
--   psql -U postgres -f setup_db.sql
--
-- ============================================================================

-- Create database
CREATE DATABASE bankassistant
    ENCODING 'UTF8'
    LC_COLLATE 'en_US.UTF-8'
    LC_CTYPE 'en_US.UTF-8'
    TEMPLATE template0;

-- Create user
CREATE USER bankassistant WITH PASSWORD 'your_secure_password_here';

-- Grant privileges to user
GRANT CONNECT ON DATABASE bankassistant TO bankassistant;
GRANT CREATE ON DATABASE bankassistant TO bankassistant;

-- Connect to the new database
\c bankassistant

-- Create schema
CREATE SCHEMA IF NOT EXISTS public;

-- Grant schema privileges
GRANT USAGE ON SCHEMA public TO bankassistant;
GRANT CREATE ON SCHEMA public TO bankassistant;

-- Set default privileges
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO bankassistant;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO bankassistant;

-- ============================================================================
-- Tables are created automatically by the Go application
-- ============================================================================
--
-- The following tables will be auto-created by the application:
--   - users (user_id, phone, pin, name, email, address, date_of_birth, ssn)
--   - sessions (token, user_id, expires_at)
--   - accounts (account_id, user_id, type, balance)
--   - transactions (account_id, from_account, to_account, amount, date)
--
-- You can also manually create them using the schema definitions in database.go
