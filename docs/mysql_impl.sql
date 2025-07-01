-- 启用严格模式
SET sql_mode = 'STRICT_ALL_TABLES';

-- 成员表
CREATE TABLE members (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    
    nickname VARCHAR(128) NOT NULL,
    email VARCHAR(128) NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    longyi_id VARCHAR(50) NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    labors SMALLINT UNSIGNED NOT NULL DEFAULT 0,
    remark TEXT NULL,
    last_active DATETIME(3) NULL,
    
    UNIQUE INDEX idx_members_nickname (nickname),
    UNIQUE INDEX idx_members_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 标签表
CREATE TABLE tags (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    
    name VARCHAR(128) NOT NULL,
    description VARCHAR(256) NULL,
    
    UNIQUE INDEX idx_tags_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 成员偏好表
CREATE TABLE member_preferences (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id INT UNSIGNED NOT NULL,
    tag_id INT UNSIGNED NOT NULL,
    is_prefered BOOLEAN NOT NULL,
    
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    UNIQUE INDEX idx_member_tag_preference (member_id, tag_id) -- 确保偏好唯一性
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 作品集表
CREATE TABLE worksets (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    
    title VARCHAR(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 项目表
CREATE TABLE projects (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    
    team_affiliated_to VARCHAR(100) NOT NULL,
    legacy_id INT UNSIGNED,
    workset_id INT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    status INT UNSIGNED NOT NULL,
    urgency TINYINT NOT NULL,
    
    INDEX idx_projects_legacy_id (legacy_id),
    INDEX idx_projects_status (status), 
    FOREIGN KEY (workset_id) REFERENCES worksets(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 项目标签表
CREATE TABLE project_tags (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    
    project_id INT UNSIGNED NOT NULL,
    tag_id INT UNSIGNED NOT NULL,
    
    INDEX idx_project_tags_project (project_id),
    INDEX idx_project_tags_tag (tag_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 项目分工表
CREATE TABLE project_labor_divisions (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    
    project_id INT UNSIGNED NOT NULL,
    member_id INT UNSIGNED NOT NULL,
    labor_role SMALLINT UNSIGNED NOT NULL,
    
    INDEX idx_labor_project (project_id),
    INDEX idx_labor_member (member_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;