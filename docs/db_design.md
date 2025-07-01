# 数据库设计文档

## 一、人员基础信息表 (members)

**`存储人员基本信息`**

| 字段名        | 类型         | 约束             | 描述         |
| ------------- | ------------ | ---------------- | ------------ |
| id            | SERIAL       | PK               | 主键ID       |
| created_at    | TIMESTAMP    |                  | 创建时间     |
| updated_at    | TIMESTAMP    |                  | 更新时间     |
| deleted_at    | TIMESTAMP    |                  | 删除时间     |
| nickname      | VARCHAR(128) | UNIQUE, NOT NULL | 用户名       |
| email         | VARCHAR(128) | UNIQUE, NOT NULL | 邮箱         |
| password_hash | VARCHAR(256) | NOT NULL         | 密码哈希值   |
| longyi_id     | VARCHAR(50)  |                  | 龙译ID       |
| is_admin      | BOOLEAN      |                  | 管理员身份   |
| labors        | SMALLINT     |                  | 职责掩码     |
| remark        | TEXT         |                  | 备注         |
| last_active   | TIMESTAMP    |                  | 上次活跃时间 |

---

## 二、标签表 (tags)

**`存储系统标签`**

| 字段名      | 类型         | 约束             | 描述     |
| ----------- | ------------ | ---------------- | -------- |
| id          | SERIAL       | PK               | 主键ID   |
| created_at  | TIMESTAMP    |                  | 创建时间 |
| updated_at  | TIMESTAMP    |                  | 更新时间 |
| deleted_at  | TIMESTAMP    |                  | 删除时间 |
| name        | VARCHAR(128) | UNIQUE, NOT NULL | 标签名   |
| description | VARCHAR(256) |                  | 标签描述 |

---

## 三、成员偏好表 (member_preferences)

**`存储成员标签偏好`**

| 字段名      | 类型    | 约束         | 描述     |
| ----------- | ------- | ------------ | -------- |
| member_id   | INT     | FK, NOT NULL | 成员ID   |
| tag_id      | INT     | FK, NOT NULL | 标签ID   |
| is_prefered | BOOLEAN | NOT NULL     | 是否偏好 |

---

## 四、作品集表 (worksets)

**`存储作品集信息`**

| 字段名     | 类型         | 约束 | 描述       |
| ---------- | ------------ | ---- | ---------- |
| id         | SERIAL       | PK   | 主键ID     |
| created_at | TIMESTAMP    |      | 创建时间   |
| updated_at | TIMESTAMP    |      | 更新时间   |
| deleted_at | TIMESTAMP    |      | 删除时间   |
| title      | VARCHAR(255) |      | 作品集标题 |

---

## 五、项目表 (projects)

**`存储项目信息`**

| 字段名             | 类型         | 约束     | 描述         |
| ------------------ | ------------ | -------- | ------------ |
| id                 | SERIAL       | PK       | 主键ID       |
| created_at         | TIMESTAMP    |          | 创建时间     |
| updated_at         | TIMESTAMP    |          | 更新时间     |
| deleted_at         | TIMESTAMP    |          | 删除时间     |
| team_affiliated_to | VARCHAR(100) | NOT NULL | 所属团队     |
| legacy_id          | INT          | INDEX    | 历史遗留序号 |
| workset_id         | INT          | FK       | 作品集ID     |
| title              | VARCHAR(255) |          | 项目标题     |
| description        | TEXT         |          | 项目描述     |
| status             | INT          |          | 状态掩码     |
| urgency          | SMALLINT      |          | 紧急度       |

### 项目状态掩码定义（值）

| 状态常量                        | 值         | 描述                               |
| ------------------------------- | ---------- | ---------------------------------- |
| PROJ_STATUS_UNSET_MASK          | X          | 未知状态，当 status = 0 时在此状态 |
| PROJ_STATUS_ON_TRANSLATING_MASK | 2          | 翻译中                             |
| PROJ_STATUS_TRANSLATED_MASK     | 4          | 翻译完成                           |
| PROJ_STATUS_ON_PROOF_MASK       | 8          | 校对中                             |
| PROJ_STATUS_PROVED_MASK         | 16         | 校对完成                           |
| PROJ_STATUS_ON_LETTERING_MASK   | 32         | 嵌字中                             |
| PROJ_STATUS_LETTERED_MASK       | 64         | 嵌字完成                           |
| PROJ_STATUS_ON_REVIEWING_MASK   | 128        | 审核中                             |
| PROJ_STATUS_REVIEWED_MASK       | 256        | 审核完成                           |
| PROJ_STATUS_PUBLISHED_MASK      | 512        | 已发布                             |
| PROJ_STATUS_CANCELED            | 2147483648 | 中止                               |

---

## 六、项目标签表 (project_tags)

**`项目与标签关联`**

| 字段名     | 类型      | 约束      | 描述     |
| ---------- | --------- | --------- | -------- |
| id         | SERIAL    | PK        | 主键ID   |
| created_at | TIMESTAMP |           | 创建时间 |
| updated_at | TIMESTAMP |           | 更新时间 |
| deleted_at | TIMESTAMP |           | 删除时间 |
| project_id | INT       | FK, INDEX | 项目ID   |
| tag_id     | INT       | FK, INDEX | 标签ID   |

---

## 七、项目分工表 (project_labor_divisions)

**`项目人员分工`**

| 字段名     | 类型      | 约束 | 描述     |
| ---------- | --------- | ---- | -------- |
| id         | SERIAL    | PK   | 主键ID   |
| created_at | TIMESTAMP |      | 创建时间 |
| updated_at | TIMESTAMP |      | 更新时间 |
| deleted_at | TIMESTAMP |      | 删除时间 |
| project_id | INT       | FK   | 项目ID   |
| member_id  | INT       | FK   | 成员ID   |
| labor_role | INT  |      | 职责掩码 |

### 职责掩码定义（值）

| 职责常量                 | 值  | 描述   |
| ------------------------ | --- | ------ |
| LABOR_CREATOR_SHIFT      | 1   | 创建者 |
| LABOR_PRICINPAL_SHIFT    | 2   | 负责人 |
| LABOR_SRC_PROV_SHIFT     | 4   | 图源   |
| LABOR_CLEANER_SHIFT      | 8   | 修图   |
| LABOR_GRAPHIC_PROC_SHIFT | 16  | 美工   |
| LABOR_TRANSLATOR_SHIFT   | 32  | 翻译   |
| LABOR_PROOF_SHIFT        | 64  | 校对   |
| LABOR_LETTERER_SHIFT     | 128 | 嵌字   |
| LABOR_REVIEWER_SHIFT     | 256 | 审核   |

---

## 外键关系

1. `member_preferences.member_id` → `members.id`
2. `member_preferences.tag_id` → `tags.id`
3. `projects.workset_id` → `worksets.id`
4. `project_tags.project_id` → `projects.id`
5. `project_tags.tag_id` → `tags.id`
6. `project_labor_divisions.project_id` → `projects.id`
7. `project_labor_divisions.member_id` → `members.id`

## 索引优化

1. `members.last_active` (DESC)
2. `projects.legacy_id`
3. `project_tags.project_id`
4. `project_tags.tag_id`

## 分区建议

```sql
ALTER TABLE projects PARTITION BY LIST (status_category) (
    PARTITION pending VALUES IN (1),  -- 仅PROJ_STATUS_UNSET_MASK
    PARTITION finished VALUES IN (3), -- PROJ_STATUS_PUBLISHED_MASK 或 PROJ_STATUS_CANCELED
    PARTITION active VALUES IN (2)    -- 其他状态
);
```
