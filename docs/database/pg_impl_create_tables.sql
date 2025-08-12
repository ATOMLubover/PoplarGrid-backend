BEGIN;

-- Table: users
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    nickname VARCHAR(128) NOT NULL,
    email VARCHAR(128),
    password_hash VARCHAR(256) NOT NULL,
    qq_number BIGINT DEFAULT NULL,
    moetran_id TEXT,
    moetran_jwt TEXT,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    remark TEXT DEFAULT NULL
);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_qq_number ON users(qq_number);
CREATE UNIQUE INDEX idx_users_nickname_unique_active ON users (nickname) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_users_email_unique_active ON users (email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_users_qq_number_unique_active ON users (qq_number) WHERE deleted_at IS NULL AND qq_number IS NOT NULL;
CREATE UNIQUE INDEX idx_users_moetran_id_unique ON users (moetran_id);

-- Table: teams
CREATE TABLE teams (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    moetran_id TEXT
);

CREATE INDEX idx_teams_deleted_at ON teams(deleted_at);
CREATE INDEX idx_teams_moetran_id ON teams(moetran_id);
CREATE UNIQUE INDEX idx_teams_name_unique_active ON teams (name) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_teams_moetran_id_unique ON teams (moetran_id);


-- Table: members
CREATE TABLE members (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_id BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    is_source_provider BOOLEAN NOT NULL DEFAULT FALSE,
    is_perfector BOOLEAN NOT NULL DEFAULT FALSE,
    is_translator BOOLEAN NOT NULL DEFAULT FALSE,
    is_proofreader BOOLEAN NOT NULL DEFAULT FALSE,
    is_letterer BOOLEAN NOT NULL DEFAULT FALSE,
    is_reviewer BOOLEAN NOT NULL DEFAULT FALSE,
    is_publisher BOOLEAN NOT NULL DEFAULT FALSE
);

ALTER TABLE members DROP CONSTRAINT IF EXISTS unique_user_in_team;
CREATE INDEX idx_members_deleted_at ON members(deleted_at);
CREATE INDEX idx_members_user_id ON members(user_id);
CREATE INDEX idx_members_team_id ON members(team_id);
CREATE UNIQUE INDEX idx_members_user_team_unique ON members (user_id, team_id);


-- Table: worksets
CREATE TABLE worksets (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    team_id BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    project_sequence_name TEXT NOT NULL,
    moetran_id TEXT
);

ALTER TABLE worksets DROP CONSTRAINT IF EXISTS worksets_name_key;
ALTER TABLE worksets DROP CONSTRAINT IF EXISTS worksets_project_sequence_name_key;
CREATE INDEX idx_worksets_deleted_at ON worksets(deleted_at);
CREATE INDEX idx_worksets_team_id ON worksets(team_id);
CREATE INDEX idx_worksets_moetran_id ON worksets(moetran_id);
CREATE UNIQUE INDEX idx_worksets_name_team_unique_active ON worksets (team_id, name) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_worksets_project_sequence_name_unique_active ON worksets (project_sequence_name) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_worksets_moetran_id_unique ON worksets (moetran_id);


-- Table: projects
CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    title VARCHAR(128) NOT NULL UNIQUE,
    description TEXT,
    moetran_id TEXT,
    legacy_id BIGINT,
    workset_id BIGINT NOT NULL REFERENCES worksets(id) ON DELETE CASCADE,
    principal_id BIGINT NOT NULL REFERENCES members(id) ON DELETE RESTRICT,
    workset_index BIGINT NOT NULL,
    translate_status SMALLINT NOT NULL DEFAULT 0,
    proofread_status SMALLINT NOT NULL DEFAULT 0,
    letter_status SMALLINT NOT NULL DEFAULT 0,
    review_status SMALLINT NOT NULL DEFAULT 0,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    allow_auto_join BOOLEAN NOT NULL DEFAULT FALSE,
    is_hidden BOOLEAN NOT NULL DEFAULT FALSE
);

ALTER TABLE projects DROP CONSTRAINT IF EXISTS unique_project_in_workset;
CREATE INDEX idx_projects_deleted_at ON projects(deleted_at);
CREATE INDEX idx_projects_legacy_id ON projects(legacy_id);
CREATE INDEX idx_projects_moetran_id ON projects(moetran_id);
CREATE INDEX idx_projects_principal_id ON projects(principal_id);
CREATE INDEX idx_projects_workset_id ON projects(workset_id);
CREATE UNIQUE INDEX idx_projects_workset_index_unique_active ON projects (workset_id, workset_index) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_projects_moetran_id_unique ON projects (moetran_id);


-- Table: labors
CREATE TABLE labors (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    labor_mask INTEGER NOT NULL DEFAULT 0
);

ALTER TABLE labors DROP CONSTRAINT IF EXISTS unique_member_in_project;
CREATE INDEX idx_labors_deleted_at ON labors(deleted_at);
CREATE INDEX idx_labors_project_id ON labors(project_id);
CREATE INDEX idx_labors_member_id ON labors(member_id);
CREATE UNIQUE INDEX idx_labors_project_member_unique_active ON labors (project_id, member_id) WHERE deleted_at IS NULL;


-- Table: invitations
CREATE TABLE invitations (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    invitor_member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    invitee_member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    target_labor_mask INTEGER NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 0
);

CREATE INDEX idx_invitations_deleted_at ON invitations(deleted_at);
CREATE INDEX idx_invitations_invitor_member_id ON invitations(invitor_member_id);
CREATE INDEX idx_invitations_invitee_member_id ON invitations(invitee_member_id);
CREATE INDEX idx_invitations_project_id ON invitations(project_id);
CREATE UNIQUE INDEX idx_invitations_unique_pending_active ON invitations (invitor_member_id, invitee_member_id, project_id) WHERE deleted_at IS NULL AND status = 0;


-- Table: applications
CREATE TABLE applications (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    applicant_member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    processor_member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE RESTRICT,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    target_labor_mask INTEGER NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 0
);

CREATE INDEX idx_applications_deleted_at ON applications(deleted_at);
CREATE INDEX idx_applications_applicant_member_id ON applications(applicant_member_id);
CREATE INDEX idx_applications_processor_member_id ON applications(processor_member_id);
CREATE INDEX idx_applications_project_id ON applications(project_id);
CREATE UNIQUE INDEX idx_applications_unique_pending_active ON applications (applicant_member_id, project_id) WHERE deleted_at IS NULL AND status = 0;

DROP MATERIALIZED VIEW IF EXISTS project_stats_mv;

CREATE MATERIALIZED VIEW project_stats_mv AS
SELECT
    p.workset_id,
    COUNT(*) AS total,
    COUNT(CASE WHEN p.translate_status = 0 THEN 1 END) AS not_translating,
    COUNT(CASE WHEN p.translate_status = 1 THEN 1 END) AS translate_in_progress,
    COUNT(CASE WHEN p.translate_status = 2 THEN 1 END) AS translate_completed,
    COUNT(CASE WHEN p.proofread_status = 0 THEN 1 END) AS not_prooving,
    COUNT(CASE WHEN p.proofread_status = 1 THEN 1 END) AS proof_in_progress,
    COUNT(CASE WHEN p.proofread_status = 2 THEN 1 END) AS proof_completed,
    COUNT(CASE WHEN p.letter_status = 0 THEN 1 END) AS not_lettering,
    COUNT(CASE WHEN p.letter_status = 1 THEN 1 END) AS letter_in_progress,
    COUNT(CASE WHEN p.letter_status = 2 THEN 1 END) AS letter_completed,
    COUNT(CASE WHEN p.review_status = 0 THEN 1 END) AS not_reviewing,
    COUNT(CASE WHEN p.review_status = 1 THEN 1 END) AS review_in_progress,
    COUNT(CASE WHEN p.review_status = 2 THEN 1 END) AS review_completed,
    COUNT(CASE WHEN p.is_published = FALSE THEN 1 END) AS not_published,
    COUNT(CASE WHEN p.is_published = TRUE THEN 1 END) AS published
FROM
    projects AS p
WHERE
    p.deleted_at IS NULL
GROUP BY
    p.workset_id;

CREATE UNIQUE INDEX IF NOT EXISTS idx_project_stats_mv_workset_id ON project_stats_mv (workset_id);

CREATE OR REPLACE FUNCTION create_workset_project_sequence()
RETURNS TRIGGER AS $$
DECLARE
    v_sequence_name TEXT;
BEGIN
    IF NEW.id IS NULL THEN
        RETURN NULL;
    END IF;
    v_sequence_name := format('workset_project_index_seq_%s', NEW.id);
    NEW.project_sequence_name := v_sequence_name;
    EXECUTE format('CREATE SEQUENCE %I START 1', v_sequence_name);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_workset_project_sequence
BEFORE INSERT ON worksets
FOR EACH ROW
EXECUTE FUNCTION create_workset_project_sequence();

CREATE OR REPLACE FUNCTION set_project_workset_index()
RETURNS TRIGGER AS $$
DECLARE
    v_sequence_name TEXT;
    v_next_val BIGINT;
BEGIN
    SELECT project_sequence_name INTO v_sequence_name
    FROM worksets
    WHERE id = NEW.workset_id
    FOR UPDATE;

    IF v_sequence_name IS NULL OR v_sequence_name = '' THEN
        RETURN NULL;
    END IF;

    EXECUTE 'SELECT nextval($1::regclass)' INTO v_next_val USING v_sequence_name;
    NEW.workset_index := v_next_val;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_set_project_workset_index
BEFORE INSERT ON projects
FOR EACH ROW
EXECUTE FUNCTION set_project_workset_index();

COMMIT;
