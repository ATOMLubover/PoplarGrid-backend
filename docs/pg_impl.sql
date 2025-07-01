-- PostgreSQL Version for the Localization Team Project Database

-- 最佳实践：将所有对象创建包裹在一个事务中
-- 如果任何语句失败，所有操作都将回滚
BEGIN;

-- 表定义
-- 注意：
-- 1. SERIAL/BIGSERIAL 用于自增主键
-- 2. TIMESTAMP WITH TIME ZONE (TIMESTAMPTZ) 是推荐的时间类型
-- 3. INT UNSIGNED -> INTEGER, TINYINT -> SMALLINT
-- 4. 'deleted_at' 的软删除逻辑在应用层处理，数据库层面只记录时间

-- 成员表 (members)
CREATE TABLE members (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,
    
    nickname VARCHAR(128) NOT NULL UNIQUE,
    email VARCHAR(128) NOT NULL UNIQUE,
    password_hash VARCHAR(256) NOT NULL,
    longyi_id VARCHAR(50) NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    labors SMALLINT NOT NULL DEFAULT 0,
    remark TEXT NULL,
    last_active TIMESTAMPTZ(3) NULL
);

-- 标签表 (tags)
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,
    
    name VARCHAR(128) NOT NULL UNIQUE,
    description VARCHAR(256) NULL
);

-- 成员偏好表 (member_preferences)
CREATE TABLE member_preferences (
    id SERIAL PRIMARY KEY,
    member_id INTEGER NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    is_prefered BOOLEAN NOT NULL,
    
    -- 确保一个成员对一个标签只有一条偏好记录
    UNIQUE (member_id, tag_id)
);


-- 作品集表 (worksets)
CREATE TABLE worksets (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,
    
    title VARCHAR(255) NOT NULL
);

-- 项目表 (projects)
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ(3) NULL,
    
    team_affiliated_to VARCHAR(100) NOT NULL,
    legacy_id INTEGER NULL,
    workset_id INTEGER NOT NULL REFERENCES worksets(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    status INTEGER NOT NULL,
    urgency SMALLINT NOT NULL
);

-- 为常用查询字段创建索引
CREATE INDEX idx_projects_legacy_id ON projects(legacy_id);
CREATE INDEX idx_projects_status ON projects(status);

-- 项目标签表 (project_tags)
CREATE TABLE project_tags (
    id SERIAL PRIMARY KEY,
    -- created_at/updated_at 在关联表中通常不是必需的，但如果需要审计则可以保留
    
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    
    -- 确保一个项目和一个标签的关联是唯一的
    UNIQUE (project_id, tag_id)
);

-- 为外键字段创建索引以优化 JOIN 查询
CREATE INDEX idx_project_tags_project_id ON project_tags(project_id);
CREATE INDEX idx_project_tags_tag_id ON project_tags(tag_id);

-- 项目分工表 (project_labor_divisions)
CREATE TABLE project_labor_divisions (
    id SERIAL PRIMARY KEY,
    
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    member_id INTEGER NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    labor_role INTEGER NOT NULL,

    -- 一个成员在一个项目中的分工角色组合应该是唯一的
    UNIQUE(project_id, member_id)
);

-- 为外键字段创建索引以优化 JOIN 查询
CREATE INDEX idx_labor_project_id ON project_labor_divisions(project_id);
CREATE INDEX idx_labor_member_id ON project_labor_divisions(member_id);

-- 提交事务
COMMIT;
