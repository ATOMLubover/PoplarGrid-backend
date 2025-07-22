/*
 Navicat Premium Data Transfer

 Topic: poplar         : 白杨表格数据库设计 PostgreSQL 实现版

 Source Server         : local_postgres_root
 Source Server Type    : PostgreSQL
 Source Server Version : 170005 (170005)
 Source Host           : localhost:5432
 Source Catalog        : poplar
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 170005 (170005)
 File Encoding         : 65001

 Date: 22/07/2025 23:14:47
*/


-- ----------------------------
-- Sequence structure for project_applications_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."project_applications_id_seq";
CREATE SEQUENCE "public"."project_applications_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for project_invitations_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."project_invitations_id_seq";
CREATE SEQUENCE "public"."project_invitations_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for project_labor_divisions_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."project_labor_divisions_id_seq";
CREATE SEQUENCE "public"."project_labor_divisions_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for projects_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."projects_id_seq";
CREATE SEQUENCE "public"."projects_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for team_members_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."team_members_id_seq";
CREATE SEQUENCE "public"."team_members_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for teams_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."teams_id_seq";
CREATE SEQUENCE "public"."teams_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for users_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."users_id_seq";
CREATE SEQUENCE "public"."users_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for worksets_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."worksets_id_seq";
CREATE SEQUENCE "public"."worksets_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for project_applications
-- ----------------------------
DROP TABLE IF EXISTS "public"."project_applications";
CREATE TABLE "public"."project_applications" (
  "id" int8 NOT NULL DEFAULT nextval('project_applications_id_seq'::regclass),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6),
  "applicant_id" int8 NOT NULL,
  "project_id" int8 NOT NULL,
  "principal_id" int8 NOT NULL,
  "target_role" int8 NOT NULL DEFAULT 0,
  "status" int2 NOT NULL DEFAULT 0
)
;

-- ----------------------------
-- Table structure for project_invitations
-- ----------------------------
DROP TABLE IF EXISTS "public"."project_invitations";
CREATE TABLE "public"."project_invitations" (
  "id" int8 NOT NULL DEFAULT nextval('project_invitations_id_seq'::regclass),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6),
  "inviter_id" int8 NOT NULL,
  "invitee_id" int8 NOT NULL,
  "project_id" int8 NOT NULL,
  "target_role" int8 NOT NULL DEFAULT 0,
  "status" int2 NOT NULL DEFAULT 0
)
;

-- ----------------------------
-- Table structure for project_labor_divisions
-- ----------------------------
DROP TABLE IF EXISTS "public"."project_labor_divisions";
CREATE TABLE "public"."project_labor_divisions" (
  "id" int8 NOT NULL DEFAULT nextval('project_labor_divisions_id_seq'::regclass),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6),
  "project_id" int8 NOT NULL,
  "user_id" int8 NOT NULL,
  "labor_role" int8 NOT NULL
)
;

-- ----------------------------
-- Table structure for projects
-- ----------------------------
DROP TABLE IF EXISTS "public"."projects";
CREATE TABLE "public"."projects" (
  "id" int8 NOT NULL DEFAULT nextval('projects_id_seq'::regclass),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6),
  "title" text COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "moetran_id" text COLLATE "pg_catalog"."default",
  "legacy_id" int8,
  "workset_id" int8 NOT NULL,
  "principal_id" int8 NOT NULL,
  "workset_index" int8 NOT NULL,
  "translate_status" int2 NOT NULL DEFAULT 0,
  "proof_status" int2 NOT NULL DEFAULT 0,
  "letter_status" int2 NOT NULL DEFAULT 0,
  "review_status" int2 NOT NULL DEFAULT 0,
  "is_published" bool NOT NULL DEFAULT false,
  "allow_auto_join" bool NOT NULL DEFAULT false,
  "is_hidden" bool NOT NULL DEFAULT false
)
;

-- ----------------------------
-- Table structure for team_members
-- ----------------------------
DROP TABLE IF EXISTS "public"."team_members";
CREATE TABLE "public"."team_members" (
  "id" int8 NOT NULL DEFAULT nextval('team_members_id_seq'::regclass),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6),
  "user_id" int8 NOT NULL,
  "team_id" int8 NOT NULL,
  "role" int8 NOT NULL DEFAULT 0
)
;

-- ----------------------------
-- Table structure for teams
-- ----------------------------
DROP TABLE IF EXISTS "public"."teams";
CREATE TABLE "public"."teams" (
  "id" int8 NOT NULL DEFAULT nextval('teams_id_seq'::regclass),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6),
  "name" varchar(256) COLLATE "pg_catalog"."default" NOT NULL,
  "moetran_id" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "id" int8 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6),
  "nickname" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "email" varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
  "password_hash" varchar(256) COLLATE "pg_catalog"."default" NOT NULL,
  "moetran_id" text COLLATE "pg_catalog"."default",
  "moetran_auth" text COLLATE "pg_catalog"."default",
  "poplar_is_admin" bool NOT NULL DEFAULT false,
  "remark" text COLLATE "pg_catalog"."default",
  "qq_number" int8,
  "last_active" timestamptz(6)
)
;

-- ----------------------------
-- Table structure for worksets
-- ----------------------------
DROP TABLE IF EXISTS "public"."worksets";
CREATE TABLE "public"."worksets" (
  "id" int8 NOT NULL DEFAULT nextval('worksets_id_seq'::regclass),
  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
  "deleted_at" timestamptz(6),
  "team_id" int8 NOT NULL,
  "name" text COLLATE "pg_catalog"."default" NOT NULL,
  "moetran_id" text COLLATE "pg_catalog"."default",
  "project_sequence_name" text COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Function structure for create_workset_project_sequence
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."create_workset_project_sequence"();
CREATE OR REPLACE FUNCTION "public"."create_workset_project_sequence"()
  RETURNS "pg_catalog"."trigger" AS $BODY$
DECLARE
    -- 声明一个变量来存储生成的序列名称
    v_sequence_name TEXT;
BEGIN
    -- 检查 NEW.id 是否已存在。对于 BIGSERIAL，ID 在 BEFORE INSERT 触发器执行时就已经生成了。
    IF NEW.id IS NULL THEN
        RAISE EXCEPTION 'Workset ID (NEW.id) cannot be NULL for sequence creation.';
    END IF;

    -- 根据 workset 的 ID 生成唯一的序列名称
    -- 格式与 Go 代码中的 PROJ_IDX_SEQ_PREFIX_FMT 保持一致
    v_sequence_name := format('workset_project_index_seq_%s', NEW.id);

    -- 将生成的序列名称赋值给新行的 project_sequence_name 字段
    -- 这一步必须在创建序列之前，因为 CREATE SEQUENCE 语句中会用到这个名称
    NEW.project_sequence_name := v_sequence_name;

    -- 动态执行 SQL 创建新的 PostgreSQL 序列
    -- START 1 表示序列从 1 开始
    EXECUTE format('CREATE SEQUENCE %I START 1', v_sequence_name);

    -- 返回 NEW，表示允许插入操作继续，并使用修改后的 NEW 记录
    RETURN NEW;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;

-- ----------------------------
-- Function structure for set_project_workset_index
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."set_project_workset_index"();
CREATE OR REPLACE FUNCTION "public"."set_project_workset_index"()
  RETURNS "pg_catalog"."trigger" AS $BODY$
DECLARE
    v_sequence_name TEXT;
    v_next_val BIGINT;
BEGIN
    RAISE NOTICE '--- 触发器 trg_set_project_workset_index 已激活 ---';
    RAISE NOTICE '即将插入的 project.workset_id: %', NEW.workset_id;
    RAISE NOTICE '即将插入的 project.workset_index (初始值): %', NEW.workset_index; -- 观察GORM传入的初始值

    -- 1. 根据 NEW.workset_id 查询 worksets 表获取 project_sequence_name
    BEGIN
        SELECT project_sequence_name INTO v_sequence_name
        FROM worksets
        WHERE id = NEW.workset_id
        FOR UPDATE; -- 锁定 workset 行，确保序列名称的原子性
    EXCEPTION
        WHEN NO_DATA_FOUND THEN
            RAISE EXCEPTION '未找到 workset_id % 对应的 workset 记录。', NEW.workset_id;
        WHEN OTHERS THEN
            RAISE EXCEPTION '查询 workset_id % 对应的序列名时发生未知错误: %', NEW.workset_id, SQLERRM;
    END;

    RAISE NOTICE '从 worksets 表中找到的序列名称: %', v_sequence_name;

    IF v_sequence_name IS NULL OR v_sequence_name = '' THEN
        RAISE EXCEPTION 'workset_id % 对应的 project_sequence_name 为空或无效。', NEW.workset_id;
    END IF;

    -- 2. 动态执行 SQL，从获取到的序列中获取下一个值
    -- 使用 $1::regclass 明确将传入的字符串视为关系（如序列）的名称
    BEGIN
        EXECUTE 'SELECT nextval($1::regclass)' INTO v_next_val USING v_sequence_name;
    EXCEPTION
        WHEN OTHERS THEN
            -- 捕获更具体的错误信息，并重新抛出
            RAISE EXCEPTION '从序列 % 获取 nextval 失败: % (SQLSTATE %)', v_sequence_name, SQLERRM, SQLSTATE;
    END;

    RAISE NOTICE '从序列 % 获取到的下一个值: %', v_sequence_name, v_next_val;

    -- 3. 将获取到的下一个序列值赋值给新插入项目记录的 workset_index 字段
    NEW.workset_index := v_next_val;
    RAISE NOTICE '已将 NEW.workset_index 设置为: %', NEW.workset_index;

    RAISE NOTICE '--- 触发器 trg_set_project_workset_index 执行完毕 ---';
    RETURN NEW;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."project_applications_id_seq"
OWNED BY "public"."project_applications"."id";
SELECT setval('"public"."project_applications_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."project_invitations_id_seq"
OWNED BY "public"."project_invitations"."id";
SELECT setval('"public"."project_invitations_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."project_labor_divisions_id_seq"
OWNED BY "public"."project_labor_divisions"."id";
SELECT setval('"public"."project_labor_divisions_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."projects_id_seq"
OWNED BY "public"."projects"."id";
SELECT setval('"public"."projects_id_seq"', 1358, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."team_members_id_seq"
OWNED BY "public"."team_members"."id";
SELECT setval('"public"."team_members_id_seq"', 297, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."teams_id_seq"
OWNED BY "public"."teams"."id";
SELECT setval('"public"."teams_id_seq"', 1, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."users_id_seq"
OWNED BY "public"."users"."id";
SELECT setval('"public"."users_id_seq"', 297, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_1"', 6, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_10"', 113, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_11"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_12"', 115, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_13"', 6, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_14"', 2, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_2"', 780, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_3"', 176, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_4"', 102, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_5"', 3, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_6"', 3, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_7"', 5, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_8"', 10, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."workset_project_index_seq_9"', 37, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."worksets_id_seq"
OWNED BY "public"."worksets"."id";
SELECT setval('"public"."worksets_id_seq"', 14, true);

-- ----------------------------
-- Indexes structure for table project_applications
-- ----------------------------
CREATE INDEX "idx_project_applications_applicant_id" ON "public"."project_applications" USING btree (
  "applicant_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_project_applications_deleted_at" ON "public"."project_applications" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_project_applications_principal_id" ON "public"."project_applications" USING btree (
  "principal_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_project_applications_project_id" ON "public"."project_applications" USING btree (
  "project_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table project_applications
-- ----------------------------
ALTER TABLE "public"."project_applications" ADD CONSTRAINT "project_applications_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table project_invitations
-- ----------------------------
CREATE INDEX "idx_project_invitations_deleted_at" ON "public"."project_invitations" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_project_invitations_invitee_id" ON "public"."project_invitations" USING btree (
  "invitee_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_project_invitations_inviter_id" ON "public"."project_invitations" USING btree (
  "inviter_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_project_invitations_project_id" ON "public"."project_invitations" USING btree (
  "project_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table project_invitations
-- ----------------------------
ALTER TABLE "public"."project_invitations" ADD CONSTRAINT "project_invitations_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table project_labor_divisions
-- ----------------------------
CREATE INDEX "idx_project_labor_divisions_deleted_at" ON "public"."project_labor_divisions" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table project_labor_divisions
-- ----------------------------
ALTER TABLE "public"."project_labor_divisions" ADD CONSTRAINT "unique_user_in_project" UNIQUE ("project_id", "user_id");

-- ----------------------------
-- Primary Key structure for table project_labor_divisions
-- ----------------------------
ALTER TABLE "public"."project_labor_divisions" ADD CONSTRAINT "project_labor_divisions_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table projects
-- ----------------------------
CREATE INDEX "idx_projects_deleted_at" ON "public"."projects" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_projects_legacy_id" ON "public"."projects" USING btree (
  "legacy_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_projects_principal_id" ON "public"."projects" USING btree (
  "principal_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_projects_title" ON "public"."projects" USING btree (
  "title" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "idx_projects_workset_id_workset_index" ON "public"."projects" USING btree (
  "workset_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "workset_index" "pg_catalog"."int8_ops" ASC NULLS LAST
) WHERE deleted_at IS NULL;

-- ----------------------------
-- Triggers structure for table projects
-- ----------------------------
CREATE TRIGGER "trg_set_project_workset_index" BEFORE INSERT ON "public"."projects"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_project_workset_index"();

-- ----------------------------
-- Uniques structure for table projects
-- ----------------------------
ALTER TABLE "public"."projects" ADD CONSTRAINT "projects_moetran_id_key" UNIQUE ("moetran_id");
ALTER TABLE "public"."projects" ADD CONSTRAINT "unique_project_in_workset" UNIQUE ("workset_id", "workset_index");

-- ----------------------------
-- Primary Key structure for table projects
-- ----------------------------
ALTER TABLE "public"."projects" ADD CONSTRAINT "projects_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table team_members
-- ----------------------------
CREATE INDEX "idx_team_members_deleted_at" ON "public"."team_members" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_team_members_team_id" ON "public"."team_members" USING btree (
  "team_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_team_members_user_id" ON "public"."team_members" USING btree (
  "user_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table team_members
-- ----------------------------
ALTER TABLE "public"."team_members" ADD CONSTRAINT "unique_member_in_team" UNIQUE ("user_id", "team_id");

-- ----------------------------
-- Primary Key structure for table team_members
-- ----------------------------
ALTER TABLE "public"."team_members" ADD CONSTRAINT "team_members_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table teams
-- ----------------------------
CREATE INDEX "idx_teams_deleted_at" ON "public"."teams" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table teams
-- ----------------------------
ALTER TABLE "public"."teams" ADD CONSTRAINT "teams_name_key" UNIQUE ("name");
ALTER TABLE "public"."teams" ADD CONSTRAINT "teams_moetran_id_key" UNIQUE ("moetran_id");

-- ----------------------------
-- Primary Key structure for table teams
-- ----------------------------
ALTER TABLE "public"."teams" ADD CONSTRAINT "teams_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table users
-- ----------------------------
CREATE INDEX "idx_users_deleted_at" ON "public"."users" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_nickname_key" UNIQUE ("nickname");
ALTER TABLE "public"."users" ADD CONSTRAINT "users_email_key" UNIQUE ("email");
ALTER TABLE "public"."users" ADD CONSTRAINT "users_moetran_id_key" UNIQUE ("moetran_id");

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table worksets
-- ----------------------------
CREATE INDEX "idx_worksets_deleted_at" ON "public"."worksets" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Triggers structure for table worksets
-- ----------------------------
CREATE TRIGGER "trg_create_workset_project_sequence" BEFORE INSERT ON "public"."worksets"
FOR EACH ROW
EXECUTE PROCEDURE "public"."create_workset_project_sequence"();

-- ----------------------------
-- Uniques structure for table worksets
-- ----------------------------
ALTER TABLE "public"."worksets" ADD CONSTRAINT "worksets_name_key" UNIQUE ("name");
ALTER TABLE "public"."worksets" ADD CONSTRAINT "worksets_moetran_id_key" UNIQUE ("moetran_id");
ALTER TABLE "public"."worksets" ADD CONSTRAINT "worksets_project_sequence_name_key" UNIQUE ("project_sequence_name");

-- ----------------------------
-- Primary Key structure for table worksets
-- ----------------------------
ALTER TABLE "public"."worksets" ADD CONSTRAINT "worksets_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table project_applications
-- ----------------------------
ALTER TABLE "public"."project_applications" ADD CONSTRAINT "project_applications_applicant_id_fkey" FOREIGN KEY ("applicant_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."project_applications" ADD CONSTRAINT "project_applications_principal_id_fkey" FOREIGN KEY ("principal_id") REFERENCES "public"."users" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION;
ALTER TABLE "public"."project_applications" ADD CONSTRAINT "project_applications_project_id_fkey" FOREIGN KEY ("project_id") REFERENCES "public"."projects" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table project_invitations
-- ----------------------------
ALTER TABLE "public"."project_invitations" ADD CONSTRAINT "project_invitations_invitee_id_fkey" FOREIGN KEY ("invitee_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."project_invitations" ADD CONSTRAINT "project_invitations_inviter_id_fkey" FOREIGN KEY ("inviter_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."project_invitations" ADD CONSTRAINT "project_invitations_project_id_fkey" FOREIGN KEY ("project_id") REFERENCES "public"."projects" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table project_labor_divisions
-- ----------------------------
ALTER TABLE "public"."project_labor_divisions" ADD CONSTRAINT "project_labor_divisions_project_id_fkey" FOREIGN KEY ("project_id") REFERENCES "public"."projects" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."project_labor_divisions" ADD CONSTRAINT "project_labor_divisions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table projects
-- ----------------------------
ALTER TABLE "public"."projects" ADD CONSTRAINT "projects_principal_id_fkey" FOREIGN KEY ("principal_id") REFERENCES "public"."users" ("id") ON DELETE RESTRICT ON UPDATE NO ACTION;
ALTER TABLE "public"."projects" ADD CONSTRAINT "projects_workset_id_fkey" FOREIGN KEY ("workset_id") REFERENCES "public"."worksets" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table team_members
-- ----------------------------
ALTER TABLE "public"."team_members" ADD CONSTRAINT "team_members_team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "public"."teams" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."team_members" ADD CONSTRAINT "team_members_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table worksets
-- ----------------------------
ALTER TABLE "public"."worksets" ADD CONSTRAINT "worksets_team_id_fkey" FOREIGN KEY ("team_id") REFERENCES "public"."teams" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
