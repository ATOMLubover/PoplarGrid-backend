-----

好的，这是一份根据您最初提出的接口设计方案，整理而成的简明 API 文档总结：

-----

## API 接口设计概览

本 API 设计遵循 RESTful 原则，并全面采用 **查询参数 (Query Parameters)** 进行数据筛选、排序和分页，以确保良好的扩展性和兼容性。

-----

### 1\. 标签接口 (TAG)

**获取标签列表**

  * `GET /tags`
      * **查询参数:**
          * `sort`: `newest` (默认，按 ID 降序) | `popularity` (按关联数降序)
      * **响应示例:**
        ```json
        [
          {
            "id": 1,
            "name": "简体中文",
            "description": "中文翻译"
          }
        ]
        ```

-----

### 2\. 汉化组接口 (TEAM)

**获取汉化组列表**

  * `GET /teams`
      * **查询参数:**
          * `sort`: `oldest` (默认，按 ID 升序) | `newest` (按 ID 降序)
      * **响应示例:**
        ```json
        [
          {
            "id": 101,
            "name": "萌翻组",
            "moetran_id": "MT001",
            "created_at": "2023-01-01"
          }
        ]
        ```

-----

### 3\. 作品集接口 (WORKSET)

**获取作品集列表**

  * `GET /worksets`
      * **查询参数:**
          * `team_id`: **必填** (指定汉化组 ID)
          * `sort`: `oldest` (默认，按 ID 升序) | `newest` (按 ID 降序)
      * **响应示例:**
        ```json
        [
          {
            "id": 201,
            "workset_name": "夏季企划",
            "team_id": 101
          }
        ]
        ```

-----

### 4\. 项目接口 (PROJECT)

#### 4.1. 获取项目列表

  * `GET /projects`
      * **查询参数:**
          * `status`: `all` | `completed` | `pending`
          * `page`: 页码 (默认 `1`)
          * `size`: 每页数量 (默认 `20`)
          * `sort`: `newest` (默认，按 ID 降序) | `oldest` (按 ID 升序)
          * `updated_order`: `true` (按更新时间降序，仅当 `status=pending` 时生效)
          * `tag_id`: 按标签 ID 过滤 (可重复使用，如 `tag_id=3&tag_id=7`)
          * `member_id`: 按参与成员 ID 过滤
      * **响应示例:**
        ```json
        {
          "total": 100,
          "data": [
            {
              "id": 301,
              "title": "作品标题",
              "urgency": 3
            }
          ]
        }
        ```

#### 4.2. 获取项目详情

  * `GET /projects/{project_id}`
      * **响应示例:**
        ```json
        {
          "id": 301,
          "title": "作品标题",
          "urgency": 3,
          "status": "pending",
          "tags": [1, 3]
        }
        ```

-----

### 5\. 成员接口 (MEMBER)

#### 5.1. 获取成员列表

  * `GET /members`
      * **查询参数:**
          * `page`: 页码 (默认 `1`)
          * `size`: 每页数量 (默认 `20`)
          * `sort`: `newest` (默认，按 ID 降序) | `oldest` (按 ID 升序)
      * **响应示例:**
        ```json
        {
          "total": 50,
          "data": [
            {
              "id": 401,
              "nickname": "成员A",
              "labors": ["翻译"]
            }
          ]
        }
        ```

#### 5.2. 获取成员详情

  * `GET /members/{member_id}`
      * **响应示例:**
        ```json
        {
          "id": 401,
          "nickname": "成员A",
          "email": "test@example.com",
          "qq_number": "123456"
        }
        ```

-----

### 关键优化点 (Kev Features)

  * **100% 查询参数实现**: 所有筛选条件通过标准化参数传递，新增参数不影响接口兼容性。
  * **出色的扩展性**: 例如，可通过重复 `tag_id` 参数按多个标签筛选项目。
  * **RESTful 兼容性**: 资源路径保持纯净，遵循 HTTP 方法语义。
  * **文档友好性**: 查询参数可清晰定义可选值和组合关系。
  * **安全控制**: 对敏感数据（如成员详情中的 QQ 号）实施访问控制。
  * **强制分页**: 列表接口默认强制分页，避免全量加载。

-----

这份文档是否清晰地总结了您提出的接口设计？