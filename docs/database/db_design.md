# 数据库设计（Markdown 文档版）

## **数据库设计文档 (v3)**

### **一、核心实体表**

#### **1.1 用户表 (users)**

**`存储平台用户（汉化组成员）的基本信息`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| **id** | **BIGINT** | **PK, SERIAL** | **主键ID** |
| created_at | TIMESTAMPTZ(6) | NOT NULL | 创建时间 |
| updated_at | TIMESTAMPTZ(6) | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMPTZ(6) | | 软删除时间 |
| nickname | VARCHAR(128) | UNIQUE, NOT NULL | 用户名 |
| email | VARCHAR(128) | UNIQUE, NOT NULL | 邮箱 |
| password_hash | VARCHAR(256) | NOT NULL | 密码哈希值 |
| moetran_id | TEXT | UNIQUE | 尨译系统对应的用户ID |
| moetran_auth | TEXT | | 尨译系统认证信息 |
| poplar_is_admin | BOOLEAN | NOT NULL, DEFAULT `false` | 是否为本平台管理员 |
| remark | TEXT | | 备注 |
| qq_number | BIGINT | | QQ号码 |
| last_active | TIMESTAMPTZ(6) | | 上次活跃时间 |

#### **1.2 汉化组表 (teams)**

**`存储从尨译同步的汉化组信息`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| **id** | **BIGINT** | **PK, SERIAL** | **主键ID** |
| created_at | TIMESTAMPTZ(6) | NOT NULL | 创建时间 |
| updated_at | TIMESTAMPTZ(6) | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMPTZ(6) | | 软删除时间 |
| name | VARCHAR(256) | UNIQUE, NOT NULL | 汉化组名称 |
| moetran_id | TEXT | UNIQUE | 尨译系统对应的团队ID |

#### **1.3 作品集表 (worksets)**

**`存储从尨译同步的作品集信息`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| **id** | **BIGINT** | **PK, SERIAL** | **主键ID** |
| created_at | TIMESTAMPTZ(6) | NOT NULL | 创建时间 |
| updated_at | TIMESTAMPTZ(6) | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMPTZ(6) | | 软删除时间 |
| **team_id** | **BIGINT** | **FK, NOT NULL** | **所属汉化组ID (关联 teams.id)** |
| name | TEXT | UNIQUE, NOT NULL | 作品集名称 |
| moetran_id | TEXT | UNIQUE | 尨译系统对应的作品集ID |
| project_sequence_name | TEXT | UNIQUE, NOT NULL | 作品集内项目序号的序列名称 |

### **二、项目与关联表**

#### **2.1 项目表 (projects)**

**`存储核心的汉化项目进度信息`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| **id** | **BIGINT** | **PK, SERIAL** | **主键ID** |
| created_at | TIMESTAMPTZ(6) | NOT NULL | 创建时间 |
| updated_at | TIMESTAMPTZ(6) | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMPTZ(6) | | 软删除时间 |
| title | TEXT | NOT NULL | 作品标题 |
| description | TEXT | | 作品描述 |
| moetran_id | TEXT | UNIQUE | 尨译系统对应的作品ID |
| legacy_id | BIGINT | INDEX | 历史遗留序号 (兼容旧数据) |
| **workset_id** | **BIGINT** | **FK, NOT NULL** | **所属作品集ID (关联 worksets.id)** |
| **principal_id** | **BIGINT** | **FK, NOT NULL** | **项目负责人ID (关联 users.id)** |
| **workset_index** | **BIGINT** | **UNIQUE (workset_id, deleted_at IS NULL), NOT NULL** | **作品集内项目序号** |
| translate_status | SMALLINT | NOT NULL, DEFAULT `0` | 翻译状态 |
| proof_status | SMALLINT | NOT NULL, DEFAULT `0` | 校对状态 |
| letter_status | SMALLINT | NOT NULL, DEFAULT `0` | 嵌字状态 |
| review_status | SMALLINT | NOT NULL, DEFAULT `0` | 审核状态 |
| is_published | BOOLEAN | NOT NULL, DEFAULT `false` | 是否已发布 |
| allow_auto_join | BOOLEAN | NOT NULL, DEFAULT `false` | 是否允许自动加入 |
| is_hidden | BOOLEAN | NOT NULL, DEFAULT `false` | 是否隐藏 |

#### **2.2 汉化组成员关系表 (team_members)**

**`存储汉化组与成员的关系 (多对多)`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| **id** | **BIGINT** | **PK, SERIAL** | **主键ID** |
| created_at | TIMESTAMPTZ(6) | NOT NULL | 创建时间 |
| updated_at | TIMESTAMPTZ(6) | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMPTZ(6) | | 软删除时间 |
| **user_id** | **BIGINT** | **FK, NOT NULL** | **成员ID (关联 users.id)** |
| **team_id** | **BIGINT** | **FK, NOT NULL** | **汉化组ID (关联 teams.id)** |
| role | BIGINT | NOT NULL, DEFAULT `0` | 成员在组内的角色掩码 |

#### **2.3 项目分工表 (project_labor_divisions)**

**`存储项目与成员的分工关系 (多对多)`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| **id** | **BIGINT** | **PK, SERIAL** | **主键ID** |
| created_at | TIMESTAMPTZ(6) | NOT NULL | 创建时间 |
| updated_at | TIMESTAMPTZ(6) | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMPTZ(6) | | 软删除时间 |
| **project_id** | **BIGINT** | **FK, NOT NULL** | **项目ID (关联 projects.id)** |
| **user_id** | **BIGINT** | **FK, NOT NULL** | **成员ID (关联 users.id)** |
| labor_role | BIGINT | NOT NULL | 职责掩码 (按位存储) |

#### **2.4 项目申请表 (project_applications)**

**`存储用户申请加入项目的记录`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| **id** | **BIGINT** | **PK, SERIAL** | **主键ID** |
| created_at | TIMESTAMPTZ(6) | NOT NULL | 创建时间 |
| updated_at | TIMESTAMPTZ(6) | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMPTZ(6) | | 软删除时间 |
| **applicant_id** | **BIGINT** | **FK, NOT NULL** | **申请人ID (关联 users.id)** |
| **project_id** | **BIGINT** | **FK, NOT NULL** | **项目ID (关联 projects.id)** |
| **principal_id** | **BIGINT** | **FK, NOT NULL** | **处理人/项目负责人ID (关联 users.id)** |
| target_role | BIGINT | NOT NULL, DEFAULT `0` | 申请加入的目标角色 |
| status | SMALLINT | NOT NULL, DEFAULT `0` | 申请状态 (例如：待处理、批准、拒绝) |

#### **2.5 项目邀请表 (project_invitations)**

**`存储项目负责人邀请用户加入项目的记录`**

| 字段名 | 类型 | 约束 | 描述 |
| :---- | :---- | :---- | :---- |
| **id** | **BIGINT** | **PK, SERIAL** | **主键ID** |
| created_at | TIMESTAMPTZ(6) | NOT NULL | 创建时间 |
| updated_at | TIMESTAMPTZ(6) | NOT NULL | 更新时间 |
| deleted_at | TIMESTAMPTZ(6) | | 软删除时间 |
| **inviter_id** | **BIGINT** | **FK, NOT NULL** | **邀请人ID (关联 users.id)** |
| **invitee_id** | **BIGINT** | **FK, NOT NULL** | **被邀请人ID (关联 users.id)** |
| **project_id** | **BIGINT** | **FK, NOT NULL** | **项目ID (关联 projects.id)** |
| target_role | BIGINT | NOT NULL, DEFAULT `0` | 邀请加入的目标角色 |
| status | SMALLINT | NOT NULL, DEFAULT `0` | 邀请状态 (例如：待处理、接受、拒绝) |

### **三、函数与触发器**

#### **3.1 `create_workset_project_sequence()` 函数**

**`用途：在 worksets 表插入新记录前，动态创建用于该作品集下项目序号的序列。`**

* **触发时机：** `BEFORE INSERT ON worksets`
* **功能：**
  * 根据新插入的 `workset.id` 生成唯一的序列名称，格式为 `workset_project_index_seq_` + `workset.id`。
  * 将生成的序列名称存储到 `NEW.project_sequence_name` 字段。
  * 动态执行 SQL 语句 `CREATE SEQUENCE [生成的序列名称] START 1` 来创建新的序列。
  * 返回 `NEW` 以允许插入操作继续。

#### **3.2 `set_project_workset_index()` 函数**

**`用途：在 projects 表插入新记录前，自动为其分配作品集内的项目序号。`**

* **触发时机：** `BEFORE INSERT ON projects`
* **功能：**
  * 根据新插入的 `project.workset_id` 查询 `worksets` 表，获取对应的 `project_sequence_name`。
  * 动态执行 SQL 语句 `SELECT nextval('[序列名称]')`，从获取到的序列中获取下一个值。
  * 将获取到的序列值赋值给 `NEW.workset_index` 字段。
  * 返回 `NEW` 以允许插入操作继续。

### **四、外键关系图谱**

* **`worksets.team_id`** → **`teams.id`**
* **`projects.workset_id`** → **`worksets.id`**
* **`projects.principal_id`** → **`users.id`**
* **`team_members.user_id`** → **`users.id`**
* **`team_members.team_id`** → **`teams.id`**
* **`project_labor_divisions.project_id`** → **`projects.id`**
* **`project_labor_divisions.user_id`** → **`users.id`**
* **`project_applications.applicant_id`** → **`users.id`**
* **`project_applications.project_id`** → **`projects.id`**
* **`project_applications.principal_id`** → **`users.id`**
* **`project_invitations.inviter_id`** → **`users.id`**
* **`project_invitations.invitee_id`** → **`users.id`**
* **`project_invitations.project_id`** → **`projects.id`**
