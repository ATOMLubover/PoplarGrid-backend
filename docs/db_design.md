# **数据库设计文档 (v2)**

## **一、核心实体表**

### **1.1 成员表 (members)**

**`存储平台用户（汉化组成员）的基本信息`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| id | SERIAL | PK | 主键ID |
| created_at | TIMESTAMP |  | 创建时间 |
| updated_at | TIMESTAMP |  | 更新时间 |
| deleted_at | TIMESTAMP |  | 软删除时间 |
| nickname | VARCHAR(128) | UNIQUE, NOT NULL | 用户名 |
| email | VARCHAR(128) | UNIQUE, NOT NULL | 邮箱 |
| password_hash | VARCHAR(256) | NOT NULL | 密码哈希值 |
| **moetran_id** | **TEXT** | **UNIQUE, NOT NULL** | **尨译系统对应的用户ID** |
| poplar_is_admin | BOOLEAN |  | 是否为本平台管理员 |
| labors | SMALLINT |  | 职责掩码 (按位存储) |
| remark | TEXT |  | 备注 |
| last_active | TIMESTAMP |  | 上次活跃时间 |

### **1.2 汉化组表 (teams)**

**`存储从尨译同步的汉化组信息`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| id | SERIAL | PK | 主键ID |
| created_at | TIMESTAMP |  | 创建时间 |
| updated_at | TIMESTAMP |  | 更新时间 |
| deleted_at | TIMESTAMP |  | 软删除时间 |
| name | VARCHAR(256) | UNIQUE, NOT NULL | 汉化组名称 |
| **moetran_id** | **TEXT** | **UNIQUE, NOT NULL** | **尨译系统对应的团队ID** |

### **1.3 作品集表 (worksets)**

**`存储从尨译同步的作品集信息`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| id | SERIAL | PK | 主键ID |
| created_at | TIMESTAMP |  | 创建时间 |
| updated_at | TIMESTAMP |  | 更新时间 |
| deleted_at | TIMESTAMP |  | 软删除时间 |
| title | TEXT |  | 作品集标题 |
| **moetran_id** | **TEXT** | **UNIQUE, NOT NULL** | **尨译系统对应的作品集ID** |

### **1.4 作品表 (works)**

**`存储从尨译同步的具体作品信息`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| id | SERIAL | PK | 主键ID |
| created_at | TIMESTAMP |  | 创建时间 |
| updated_at | TIMESTAMP |  | 更新时间 |
| deleted_at | TIMESTAMP |  | 软删除时间 |
| title | TEXT | UNIQUE, NOT NULL | 作品标题 |
| **moetran_id** | **TEXT** | **UNIQUE, NOT NULL** | **尨译系统对应的作品ID** |
| description | TEXT |  | 作品描述 |

### **1.5 标签表 (tags)**

**`存储用于项目分类的系统标签`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| id | SERIAL | PK | 主键ID |
| created_at | TIMESTAMP |  | 创建时间 |
| updated_at | TIMESTAMP |  | 更新时间 |
| deleted_at | TIMESTAMP |  | 软删除时间 |
| name | VARCHAR(128) | UNIQUE, NOT NULL | 标签名 |
| description | VARCHAR(256) |  | 标签描述 |

## **二、项目与关联表**

### **2.1 项目表 (projects)**

**`存储核心的汉化项目进度信息`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| id | SERIAL | PK | 主键ID |
| created_at | TIMESTAMP |  | 创建时间 |
| updated_at | TIMESTAMP |  | 更新时间 |
| deleted_at | TIMESTAMP |  | 软删除时间 |
| legacy_id | INTEGER | INDEX | 历史遗留序号 (兼容旧数据) |
| **team_id** | INTEGER | **FK, NOT NULL** | **所属汉化组ID** |
| **workset_id** | INTEGER | **FK, NOT NULL** | **所属作品集ID** |
| **work_id** | INTEGER | **FK, NOT NULL** | **所属作品ID** |
| status | INTEGER | INDEX | 状态掩码 (按位存储) |
| urgency | SMALLINT |  | 紧急度 |

### **2.2 项目分工表 (project_labor_divisions)**

**`存储项目与成员的分工关系 (多对多)`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| id | SERIAL | PK | 主键ID |
| project_id | INTEGER | FK, NOT NULL | 项目ID |
| member_id | INTEGER | FK, NOT NULL | 成员ID |
| labor_role | INTEGER |  | 职责掩码 (按位存储) |

### **2.3 项目标签表 (project_tags)**

**`存储项目与标签的关联关系 (多对多)`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| project_id | INTEGER | PK, FK, NOT NULL | 项目ID |
| tag_id | INTEGER | PK, FK, NOT NULL | 标签ID |

### **2.4 成员偏好表 (member_preferences)**

**`存储成员不擅长的标签 (用于分配参考)`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| id | SERIAL | PK | 主键ID |
| member_id | INTEGER | FK, NOT NULL | 成员ID |
| tag_id | INTEGER | FK, NOT NULL | 标签ID |
| **is_resisted** | BOOLEAN |  | **是否不擅长/抵触该标签** |

## **三、外键关系图谱**

* projects.team_id → teams.id  
* projects.workset_id → worksets.id  
* projects.work_id → works.id  
* project_labor_divisions.project_id → projects.id  
* project_labor_divisions.member_id → members.id  
* project_tags.project_id → projects.id  
* project_tags.tag_id → tags.id  
* member_preferences.member_id → members.id  
* member_preferences.tag_id → tags.id
