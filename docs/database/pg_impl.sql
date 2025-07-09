-- PostgreSQL Version for the Localization Team Project Database
-- Corrected version to match the Go models design (v3).
-- Key changes: 'works' table removed, 'projects' restored, and 'project_tags' created.

BEGIN;


-- Table: members
-- Stores user information. Matches the Member struct.
CREATE TABLE members (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,
    
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    nickname VARCHAR(128) NOT NULL,
    email VARCHAR(128) NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    moetran_id TEXT NOT NULL,
    poplar_is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    labors INTEGER NOT NULL DEFAULT 0,
    remark TEXT NULL,
    qq_number VARCHAR(64) NULL,
    last_active TIMESTAMPTZ(3) NULL,

    CONSTRAINT unique_members_nickname UNIQUE (nickname),
    CONSTRAINT unique_members_email UNIQUE (email),
    CONSTRAINT unique_members_moetran_id UNIQUE (moetran_id)
);
-- Indexes from original schema, which are good practice.
CREATE INDEX idx_members_deleted_at ON members(deleted_at);


-- Table: teams
-- Stores translation team info. Matches the Team struct.
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,

    name VARCHAR(256) NOT NULL,
    moetran_id TEXT NOT NULL,

    CONSTRAINT unique_teams_name UNIQUE (name),
    CONSTRAINT unique_teams_moetran_id UNIQUE (moetran_id)
);
CREATE INDEX idx_teams_deleted_at ON teams(deleted_at);


-- Table: worksets
-- Stores workset (series/collection) info. Matches the Workset struct.
CREATE TABLE worksets (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTamptz(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,
    
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    moetran_id TEXT NOT NULL,

    CONSTRAINT unique_worksets_moetran_id UNIQUE (moetran_id),
    CONSTRAINT unique_worksets_name UNIQUE (name)
);
CREATE INDEX idx_worksets_deleted_at ON worksets(deleted_at);


-- Table: tags
-- Stores system-wide tags for categorization. Matches the Tag struct.
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,
    
    name VARCHAR(128) NOT NULL,
    description VARCHAR(256) NULL,

    CONSTRAINT unique_tags_name UNIQUE (name)
);
CREATE INDEX idx_tags_deleted_at ON tags(deleted_at);


-- Table: projects (CORRECTED)
-- Core table tracking the progress of a localization project. Now matches the Project struct.
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,
    
    -- Restored fields from the Go model
    title TEXT NOT NULL,
    moetran_id TEXT NOT NULL,

    -- Fields that were already correct
    legacy_id INTEGER NULL,
    workset_id INTEGER NOT NULL REFERENCES worksets(id) ON DELETE CASCADE,
    status INTEGER NOT NULL,
    urgency SMALLINT NOT NULL,

    -- Add unique constraints as defined in the Go model's gorm tags
    -- CONSTRAINT unique_projects_title_legacy_id_workset_id UNIQUE (moetran_id, workset_id),
    CONSTRAINT unique_projects_moetran_id UNIQUE (moetran_id)
);
CREATE INDEX idx_projects_deleted_at ON projects(deleted_at);
CREATE INDEX idx_projects_legacy_id ON projects(legacy_id);
CREATE INDEX idx_projects_title ON projects(title);


-- Table: member_preferences
-- Stores member's tag preferences. Matches the MemberPreference struct.
CREATE TABLE member_preferences (
    member_id INTEGER NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    is_resisted BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- A member can only have one preference entry per tag.
    PRIMARY KEY (member_id, tag_id)
);


-- Table: project_tags (CORRECTED)
-- Join table between projects and tags. Replaces 'work_tags'. Matches the ProjectTag struct.
CREATE TABLE project_tags (
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    
    -- A tag can only be applied to a project once.
    PRIMARY KEY (project_id, tag_id)
);


-- Table: project_labor_divisions
-- Join table for project assignments. Matches the ProjectLaborDivision struct.
CREATE TABLE project_labor_divisions (
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    member_id INTEGER NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    labor_role INTEGER NOT NULL,

    -- A member should only have one labor entry per project
    -- Using PRIMARY KEY is also an option here, but UNIQUE works perfectly.
    CONSTRAINT unique_pld_project_member UNIQUE (project_id, member_id)
);
-- Indexes are helpful for querying assignments by project or by member.
CREATE INDEX idx_pld_project_id ON project_labor_divisions(project_id);
CREATE INDEX idx_pld_member_id ON project_labor_divisions(member_id);


COMMIT;
