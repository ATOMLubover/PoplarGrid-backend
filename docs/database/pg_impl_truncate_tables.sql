BEGIN;

TRUNCATE TABLE applications RESTART IDENTITY CASCADE;
TRUNCATE TABLE invitations RESTART IDENTITY CASCADE;
TRUNCATE TABLE labors RESTART IDENTITY CASCADE;
TRUNCATE TABLE projects RESTART IDENTITY CASCADE;
TRUNCATE TABLE members RESTART IDENTITY CASCADE;
TRUNCATE TABLE worksets RESTART IDENTITY CASCADE;
TRUNCATE TABLE teams RESTART IDENTITY CASCADE;
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

REFRESH MATERIALIZED VIEW CONCURRENTLY project_stats_mv;

DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT relname AS sequence_name
              FROM pg_class
              WHERE relkind = 'S'
                AND relnamespace = (SELECT oid FROM pg_namespace WHERE nspname = current_schema())
                AND relname LIKE 'workset_project_index_seq_%'
              ORDER BY relname)
    LOOP
        EXECUTE format('DROP SEQUENCE %I', r.sequence_name);
    END LOOP;
END $$;

DROP TRIGGER IF EXISTS trg_create_workset_project_sequence ON worksets;
DROP FUNCTION IF EXISTS create_workset_project_sequence();

DROP TRIGGER IF EXISTS trg_set_project_workset_index ON projects;
DROP FUNCTION IF EXISTS set_project_workset_index();

COMMIT;
