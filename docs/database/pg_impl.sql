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

--------------
-- 针对实现 workset 内递增 index 的实现
--------------

-- 假设你已经创建了一个新的 workset，它的 ID 是 123
CREATE SEQUENCE projects_seq_workset_123 START 1;

-- 用于动态获取并设置项目的 workset_project_index 的触发器函数
CREATE OR REPLACE FUNCTION set_project_workset_index()
RETURNS TRIGGER AS $$
DECLARE
    -- 声明一个变量来存储从 worksets 表查到的序列名称
    v_seq_name TEXT;
    -- 声明一个变量来存储从序列中获取的下一个值
    v_next_val BIGINT;
BEGIN
    -- 1. 根据当前新插入项目（NEW）的 workset_id，去 worksets 表查找到对应的序列名称
    SELECT project_sequence_name INTO v_seq_name
    FROM worksets
    WHERE id = NEW.workset_id;

    -- 如果没有找到序列名称，说明 workset 数据有问题，抛出异常阻止插入
    IF v_seq_name IS NULL THEN
        RAISE EXCEPTION 'Sequence name not found for workset_id %', NEW.workset_id;
    END IF;

    -- 2. 动态执行 SQL：从找到的序列中获取下一个值
    -- 注意：这里必须使用 EXECUTE，因为序列名称是变量，不能直接写在 nextval() 里
    -- nextval('sequence_name') 是获取下一个值的函数
    EXECUTE 'SELECT nextval(''' || v_seq_name || ''')' INTO v_next_val;

    -- 3. 将获取到的下一个序列值，赋值给新插入项目记录的 workset_project_number 字段
    NEW.workset_index := v_next_val;

    -- 4. 返回 NEW，表示允许插入操作继续，并使用 NEW 中修改后的值
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_projects_set_workset_number
BEFORE INSERT ON projects          -- 在向 projects 表插入数据之前
FOR EACH ROW                       -- 对每一行数据都执行
EXECUTE FUNCTION set_project_workset_index(); -- 调用我们的函数

------------
-- 针对 projects 进行统计的物化视图的实现
------------

-- 假设你的 projects 表和 WorksetId 字段都已存在并包含数据。

-- 如果你已经创建了 project_stats_mv，需要先删除它才能重新创建，
-- 特别是如果字段名发生了变化。
-- DROP MATERIALIZED VIEW IF EXISTS project_stats_mv;

CREATE MATERIALIZED VIEW project_stats_mv AS
SELECT
    p.workset_id, -- WorksetId 作为分组和查询的维度
    COUNT(*) AS total,

    -- 翻译相关统计
    COUNT(CASE WHEN p.translate_status = 0 THEN 1 END) AS not_translating,
    COUNT(CASE WHEN p.translate_status = 1 THEN 1 END) AS translating, -- 对应 TranslateInProgress
    COUNT(CASE WHEN p.translate_status = 2 THEN 1 END) AS translated,  -- 对应 TranslateCompleted

    -- 校对相关统计
    COUNT(CASE WHEN p.proof_status = 0 THEN 1 END) AS not_prooving,
    COUNT(CASE WHEN p.proof_status = 1 THEN 1 END) AS prooving,    -- 对应 ProofInProgress
    COUNT(CASE WHEN p.proof_status = 2 THEN 1 END) AS prooved,     -- 对应 ProofCompleted

    -- 排版相关统计 (原嵌字)
    COUNT(CASE WHEN p.letter_status = 0 THEN 1 END) AS not_lettering,
    COUNT(CASE WHEN p.letter_status = 1 THEN 1 END) AS lettering,   -- 对应 LetterInProgress
    COUNT(CASE WHEN p.letter_status = 2 THEN 1 END) AS letterred,   -- 对应 LetterCompleted

    -- 审核相关统计
    COUNT(CASE WHEN p.review_status = 0 THEN 1 END) AS not_reviewing,
    COUNT(CASE WHEN p.review_status = 1 THEN 1 END) AS reviewing,   -- 对应 ReviewInProgress
    COUNT(CASE WHEN p.review_status = 2 THEN 1 END) AS reviewed,    -- 对应 ReviewCompleted

    -- 发布相关统计
    COUNT(CASE WHEN p.is_published = FALSE THEN 1 END) AS not_published,
    COUNT(CASE WHEN p.is_published = TRUE THEN 1 END) AS published
FROM
    projects AS p
WHERE
    p.deleted_at IS NULL -- 只统计未被软删除的项目
GROUP BY
    p.workset_id; -- 按 WorksetId 分组

-- 推荐：为物化视图创建唯一索引，以支持 CONCURRENTLY 刷新和加速查询
-- WorksetId 是每行唯一的标识
CREATE UNIQUE INDEX IF NOT EXISTS idx_project_stats_mv_workset_id ON project_stats_mv (workset_id);
