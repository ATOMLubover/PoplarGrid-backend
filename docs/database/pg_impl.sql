-- PostgreSQL Version for the Localization Team Project Database (v2)
-- Reflects new requirements for Moetran integration.

BEGIN;


-- Table: members
-- Stores user information.
CREATE TABLE members (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,
    
    nickname VARCHAR(128) NOT NULL UNIQUE,
    email VARCHAR(128) NOT NULL UNIQUE,
    password_hash VARCHAR(256) NOT NULL,
    moetran_id TEXT NOT NULL UNIQUE, -- Changed from longyi_id, now NOT NULL
    poplar_is_admin BOOLEAN NOT NULL DEFAULT FALSE, -- Renamed from is_admin
    labors INTEGER NOT NULL DEFAULT 0,
    remark TEXT NULL,
    last_active TIMESTAMPTZ(3) NULL
);
CREATE INDEX idx_members_nickname ON members(nickname);
CREATE INDEX idx_members_deleted_at ON members(deleted_at);
CREATE INDEX idx_members_poplar_is_admin ON members(poplar_is_admin);
CREATE INDEX idx_members_labors ON members(labors DESC);


-- Table: teams (NEW)
-- Stores translation team info synchronized from Moetran.
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,

    name VARCHAR(256) NOT NULL UNIQUE,
    moetran_id TEXT NOT NULL UNIQUE
);
CREATE INDEX idx_teams_deleted_at ON teams(deleted_at);

-- Table: worksets
-- Stores workset (series/collection) info synchronized from Moetran.
CREATE TABLE worksets (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,
    
    title TEXT NOT NULL, -- Title can be fetched from Moetran
    moetran_id TEXT NOT NULL UNIQUE -- Added for synchronization
);
CREATE INDEX idx_worksets_deleted_at ON worksets(deleted_at);
CREATE INDEX idx_worksets_title_trgm ON worksets USING GIN (title gin_trgm_ops);
CREATE INDEX idx_worksets_updated_at ON worksets(updated_at DESC);

-- Table: works (NEW)
-- Stores individual work (chapter/article) info synchronized from Moetran.
CREATE TABLE works (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,

    last_edit TIMESTAMPTZ(3) NULL, -- New field for tracking last edit time
    title TEXT NOT NULL UNIQUE,
    moetran_id TEXT NOT NULL UNIQUE,
    description TEXT NULL
);
CREATE INDEX idx_works_deleted_at ON works(deleted_at);
CREATE INDEX idx_works_last_edit ON works(last_edit DESC NULLS LAST);

-- Table: tags
-- Stores system-wide tags for categorization.
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,
    
    name VARCHAR(128) NOT NULL UNIQUE,
    description VARCHAR(256) NULL
);


-- Table: projects (HEAVILY MODIFIED)
-- Core table tracking the progress of a localization project.
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,
    
    legacy_id INTEGER NULL,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE, -- Changed from team_affiliated_to
    workset_id INTEGER NOT NULL REFERENCES worksets(id) ON DELETE CASCADE,
    work_id INTEGER NOT NULL REFERENCES works(id) ON DELETE CASCADE, -- NEW foreign key
    status INTEGER NOT NULL,
    urgency SMALLINT NOT NULL
);
CREATE INDEX idx_projects_legacy_id ON projects(legacy_id);
CREATE INDEX idx_projects_status ON projects(status);


-- Table: member_preferences (MODIFIED)
-- Stores tags that a member is not good at or wants to avoid.
CREATE TABLE member_preferences (
    member_id INTEGER NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    is_resisted BOOLEAN NOT NULL DEFAULT FALSE, -- Renamed from is_prefered, meaning is inverted
    
    UNIQUE (member_id, tag_id)
);


-- Table: project_tags (MODIFIED)
-- Join table between projects and tags. Using composite primary key.
CREATE TABLE work_tags (
    work_id INTEGER NOT NULL REFERENCES work(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    
    PRIMARY KEY (work_id, tag_id) -- More accurately reflects a pure join table
);


-- Table: project_labor_divisions
-- Join table for project assignments.
CREATE TABLE project_labor_divisions (
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    member_id INTEGER NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    labor_role INTEGER NOT NULL,

    UNIQUE(project_id, member_id) -- A member should only have one labor entry per project
);
CREATE INDEX idx_labor_project_id ON project_labor_divisions(project_id);
CREATE INDEX idx_labor_member_id ON project_labor_divisions(member_id);


COMMIT;
