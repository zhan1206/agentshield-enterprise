# AgentShield Enterprise

<p align="center">
  <img src="docs/logo.png" alt="AgentShield Enterprise Logo" width="200">
</p>

<p align="center">
  <strong>全球首个生产级、面向AI Agent全场景的开源安全沙箱与全生命周期权限管控系统</strong>
</p>

<p align="center">
  <a href="#核心特性">核心特性</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#架构设计">架构设计</a> •
  <a href="#技术栈">技术栈</a> •
  <a href="#兼容生态">兼容生态</a> •
  <a href="#贡献指南">贡献指南</a>
</p>

<p align="center">
  <a href="https://github.com/zhan1206/agentshield-enterprise/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="License">
  </a>
  <a href="https://www.rust-lang.org/">
    <img src="https://img.shields.io/badge/Rust-1.74+-DEA584?logo=rust" alt="Rust Version">
  </a>
  <a href="https://go.dev">
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go" alt="Go Version">
  </a>
  <a href="https://reactjs.org">
    <img src="https://img.shields.io/badge/React-18+-61DAFB?logo=react" alt="React Version">
  </a>
  <a href="https://www.python.org/">
    <img src="https://img.shields.io/badge/Python-3.10+-3776AB?logo=python" alt="Python Version">
  </a>
  <a href="https://github.com/zhan1206/agentshield-enterprise/stargazers">
    <img src="https://img.shields.io/github/stars/zhan1206/agentshield-enterprise?style=social" alt="Stars">
  </a>
</p>

---

## 项目定位

AgentShield Enterprise 是企业级AI Agent集群的**专属安全防护底座、行为管控中枢、合规审计平台**，填补AI Agent规模化落地的核心安全空白。

### 核心解决痛点

| 痛点 | 解决方案 |
|-----|---------|
| 🔴 Agent执行环境隔离失控 | 自研三级隔离沙箱（进程/容器/微虚拟机），适配不同安全等级 |
| 🔴 权限管控无法适配Agent动态执行 | ABAC+RBAC双维度权限模型，动态授权+全链路穿透校验 |
| 🔴 全链路审计缺失 | 不可篡改链式审计系统，一键溯源+合规报告生成 |
| 🔴 数据泄露风险极高 | 敏感数据自动识别脱敏+DLP防泄漏+水印溯源 |
| 🔴 威胁检测与响应能力缺失 | 规则+AI双模式检测+分级自动化响应闭环 |
| 🔴 落地门槛极高 | 全框架零侵入兼容+20+场景模板+零代码可视化管控 |

---

## 核心特性

### 🛡️ 三级隔离沙箱引擎（Rust）
- **进程级轻量隔离**：低损耗、启动快，适用于可信Agent场景，资源占用极低
- **容器级标准隔离**：平衡安全性与性能，隔离网络/文件系统/进程空间
- **微虚拟机级强隔离**：最高安全等级，完全隔离宿主机环境，适用于不可信Agent
- **环境快照与一键回滚**：执行前自动快照，异常时瞬间恢复到安全状态

### 🔐 ABAC+RBAC 双维度权限管控
- 4级细粒度权限：环境级、资源级、操作级、数据级
- 动态授权引擎：基于Agent角色/任务/执行阶段动态调整权限
- 全链路权限穿透校验：每一次工具调用、API访问均需校验
- 权限生命周期管理：自动发放、自动过期、自动回收

### 📋 不可篡改审计系统
- 链式SHA-256存储，不可篡改、不可删除
- 全链路留痕：决策、调用、文件操作、数据访问全程记录
- 一键溯源：可视化还原Agent完整执行链路
- 合规报告：内置等保2.0、GDPR、SOX规则库，一键生成报告

### 🚨 威胁检测与自动响应
- 规则引擎+AI行为分析双轨检测模式
- 覆盖越权访问、恶意代码、数据窃取、资源耗尽等威胁
- 分级自动化响应：告警→拦截→暂停→终止→隔离→回滚
- 7×24小时无人值守防护

### 🔒 数据安全防护
- 100+敏感数据识别规则（身份证、手机号、密钥等）
- 自动脱敏引擎
- DLP防泄漏（API/文件/代码/聊天多渠道拦截）
- 隐形水印溯源

### 🔌 全生态零侵入兼容
- 全兼容OpenHands、Plandex、LangGraph、CrewAI、AutoGen
- 企业现有Agent无需修改代码，一键接入
- 统一工具网关：MCP协议、OpenAPI、数据库、内网服务全兼容

---

## 架构设计

```
┌─────────────────────────────────────────────────────────┐
│                    前端交互与统一管控控制台                │
│   可视化策略编排、安全大盘、审计中心、告警配置、模板市场  │
├─────────────────────────────────────────────────────────┤
│                    安全合规与全链路审计层                  │
│   不可篡改审计、合规报告、溯源分析、数据安全防护          │
├─────────────────────────────────────────────────────────┤
│                    全生命周期权限管控层                    │
│   双维度权限模型、动态授权、全链路校验、权限生命周期管理  │
├─────────────────────────────────────────────────────────┤
│                    核心沙箱引擎层（Rust）                 │
│   三级隔离沙箱、实时执行管控、威胁检测、自动化响应        │
├─────────────────────────────────────────────────────────┤
│                    全生态兼容适配层                        │
│   多框架适配器、Agent注册中心、统一工具网关、模型适配      │
└─────────────────────────────────────────────────────────┘
```

---

## 快速开始

### 方式1：Docker Compose（推荐）

**前置要求：** [Docker](https://docs.docker.com/get-docker/) + [Docker Compose](https://docs.docker.com/compose/install/)

```bash
# 1. 克隆仓库
git clone https://github.com/zhan1206/agentshield-enterprise.git
cd agentshield-enterprise

# 2. 配置环境变量
cp configs/.env.example configs/.env
# 编辑 configs/.env 填入数据库密码、JWT密钥等

# 3. 启动所有服务（含 MySQL, Redis, Kafka, Jaeger, Elasticsearch）
cd deploy/docker
docker-compose up -d

# 4. 查看日志
docker-compose logs -f server

# 5. 访问
#   管控控制台: http://localhost:3000
#   API:        http://localhost:8080/health
#   Jaeger:     http://localhost:16686
```

### 方式2：源码编译

**前置要求：**
- [Rust 1.74+](https://www.rust-lang.org/tools/install)（沙箱引擎）
- [Go 1.21+](https://go.dev/dl/)（后端服务）
- [Node.js 18+](https://nodejs.org/)（前端控制台）
- [Python 3.10+](https://python.org/)（AI行为分析引擎）
- MySQL 8.0+, Redis 7+

```bash
# 1. 克隆仓库
git clone https://github.com/zhan1206/agentshield-enterprise.git
cd agentshield-enterprise

# 2. 配置环境变量
cp configs/.env.example configs/.env
# 编辑 configs/.env

# 3. 编译沙箱引擎（Rust）
cd sandbox
cargo build --release
cd ..

# 4. 启动后端（Go）
cd backend
go mod tidy
go run ./cmd/server
# 默认监听 :8080

# 5. 启动AI分析引擎（Python，新终端）
cd ai
pip install -r requirements.txt
python -m behavior_analyzer.analyzer
# 默认监听 :8001

# 6. 启动前端（新终端）
cd frontend
npm install
npm run dev
# 默认监听 :3000

# 7. 访问控制台
open http://localhost:3000
```

**测试 API：**
```bash
# 健康检查
curl http://localhost:8080/health

# 创建沙箱
curl -X POST http://localhost:8080/api/v1/sandboxes \
  -H 'Content-Type: application/json' \
  -d '{"agent_id":"agent-001","isolation_level":"container","memory_limit_mb":512,"cpu_quota":2.0}'

# 查看Agent列表
curl http://localhost:8080/api/v1/agents
```

### 方式3：Kubernetes

**前置要求：** kubectl + 集群访问权限

```bash
# 1. 克隆并配置
git clone https://github.com/zhan1206/agentshield-enterprise.git
cd agentshield-enterprise

# 2. 创建 Secret
kubectl create secret generic agentshield-secrets \
  --from-literal=db-password=YOUR_DB_PASSWORD \
  --from-literal=jwt-secret=YOUR_JWT_SECRET \
  --from-literal=encryption-key=YOUR_ENCRYPTION_KEY

# 3. 部署
kubectl apply -f deploy/k8s/

# 4. 查看状态
kubectl get pods -l app=agentshield
kubectl logs -l app=agentshield
```

---

## 技术栈

| 架构模块 | 技术选型 | 选型原因 |
|----------|----------|----------|
| 核心沙箱引擎 | **Rust** | 内存安全、高性能、低损耗，天生适配系统级安全沙箱 |
| 后端管控引擎 | **Go** | 高并发、高可用，完美适配云原生分布式管控 |
| AI行为分析 | **Python** | 兼容AI生态，适配主流Agent框架与ML模型 |
| 前端控制台 | **React + TypeScript + Ant Design** | 企业级可视化管控平台 |
| 核心存储 | **MySQL + Redis** | 权限规则+审计元数据+缓存+实时状态 |
| 审计日志 | **Elasticsearch + 对象存储** | 海量日志检索+冷热分离 |
| 隔离底层 | **Linux Namespace/Cgroups + Containerd + KVM** | 三级隔离技术支撑 |
| 可观测 | **OpenTelemetry + Prometheus + Grafana** | 兼容开源标准 |
| 消息队列 | **Kafka** | 高吞吐审计日志与安全事件分发 |
| 部署 | **Docker + Kubernetes** | 云原生容器化部署 |

---

## 兼容生态

### Agent框架（全兼容）
| 框架 | 接入方式 | 侵入性 |
|------|----------|--------|
| OpenHands | 标准适配器 | 零侵入 |
| Plandex | 标准适配器 | 零侵入 |
| LangGraph | 标准适配器 | 零侵入 |
| CrewAI | 标准适配器 | 零侵入 |
| AutoGen | 标准适配器 | 零侵入 |
| 自研Agent | 通用SDK | 轻量接入 |

### 大模型（全兼容）
OpenAI GPT-4/3.5、Anthropic Claude、豆包、通义千问、DeepSeek、Llama、Qwen、GLM、Gemini 等

### 工具协议
MCP协议工具集、OpenAPI插件、数据库连接、内网服务、云资源API

---

## 项目结构

```
agentshield-enterprise/
├── sandbox/                    # Rust 沙箱引擎（核心心脏）
│   ├── Cargo.toml
│   └── src/
│       ├── lib.rs
│       └── engine/
│           ├── mod.rs
│           ├── process.rs      # 进程级轻量隔离
│           ├── container.rs    # 容器级标准隔离
│           ├── microvm.rs      # 微虚拟机级强隔离
│           ├── monitor.rs      # 实时执行监控
│           └── snapshot.rs     # 环境快照与回滚
├── backend/                    # Go 后端管控引擎
│   ├── cmd/server/             # API 服务入口
│   └── internal/
│       ├── core/
│       │   ├── permission/     # ABAC+RBAC权限引擎
│       │   ├── audit/          # 不可篡改审计链
│       │   ├── threat/         # 威胁检测与响应
│       │   └── agent/          # Agent注册管理
│       ├── adapter/
│       │   ├── framework/      # 多框架适配器
│       │   ├── tool/           # 统一工具网关
│       │   └── data/           # 数据安全防护
│       ├── security/
│       │   ├── auth/           # JWT认证
│       │   └── encryption/     # 加密引擎
│       └── observability/
│           ├── tracing/        # 链路追踪
│           └── metrics/        # 指标采集
├── ai/                         # Python AI行为分析
│   ├── requirements.txt
│   ├── behavior_analyzer/      # 行为分析引擎
│   └── compliance/             # 合规规则引擎
├── frontend/                   # React 前端控制台
│   └── src/
│       ├── pages/              # 9个页面模块
│       ├── services/           # API服务
│       └── store/              # 状态管理
├── deploy/                     # 部署配置
│   ├── docker/                 # Docker + Compose
│   └── k8s/                    # Kubernetes
├── configs/                    # 配置文件
├── templates/                  # 安全防护模板
└── docs/                       # 文档
```

---

## 核心API

### 沙箱管理

```bash
# 创建沙箱
POST /api/v1/sandboxes
{
  "agent_id": "agent-001",
  "isolation_level": "container",
  "memory_limit_mb": 512,
  "cpu_quota": 2.0,
  "network_isolated": true,
  "enable_snapshot": true
}

# 在沙箱中执行命令
POST /api/v1/sandboxes/{id}/execute
{
  "command": "ls -la /app",
  "timeout": 30
}

# 创建环境快照
POST /api/v1/sandboxes/{id}/snapshot

# 回滚到快照
POST /api/v1/sandboxes/{id}/rollback
{
  "snapshot_id": "snap-001"
}
```

### 权限管控

```bash
# 创建权限规则
POST /api/v1/permissions
{
  "name": "研发Agent代码库只读",
  "effect": "allow",
  "level": "resource",
  "agent_roles": ["developer"],
  "actions": ["read"],
  "resources": ["code_repo:*"]
}
```

### 合规审计

```bash
# 获取审计日志
GET /api/v1/audit/logs?agent_id=agent-001&limit=50

# 验证审计链完整性
GET /api/v1/audit/chain/verify

# 生成合规报告
GET /api/v1/audit/compliance/djl_2_0
```

---

## 典型落地场景

### 场景1：企业级研发Agent集群安全防护
研发Agent统一接入AgentShield，代码执行全程沙箱隔离，敏感密钥自动脱敏，全链路审计留痕，违规操作实时拦截。

### 场景2：金融机构智能投研Agent合规管控
微虚拟机级强隔离，严格数据访问权限，满足等保三级、数据安全法合规要求，一键生成合规审计报告。

### 场景3：SaaS化Agent平台多租户隔离
为每个租户提供独立隔离环境，租户之间环境/数据/权限完全隔离，杜绝交叉感染。

---

## 路线图

| 版本 | 时间 | 核心目标 |
|-----|------|---------|
| v0.1 MVP | 3个月 | 容器级沙箱、5+框架适配、基础RBAC+审计 |
| v0.5 Beta | 6个月 | 三级隔离、ABAC+RBAC、威胁检测、数据脱敏 |
| v1.0 正式版 | 12个月 | 企业级生产可用、完整合规库、20+场景模板 |
| v2.0 生态版 | 24个月 | 多租户SaaS、插件市场、行业定制方案 |

---

## 贡献指南

我们欢迎所有形式的贡献！

- 🐛 [报告Bug](https://github.com/zhan1206/agentshield-enterprise/issues)
- 💡 [功能建议](https://github.com/zhan1206/agentshield-enterprise/issues)
- 🔒 贡献安全规则与威胁检测模型
- 🔧 贡献Agent框架适配器
- 📖 改进文档

详见 [CONTRIBUTING.md](CONTRIBUTING.md)

---

## 许可证

本项目采用 [Apache 2.0](LICENSE) 协议开源。核心功能永久开源免费，商用友好，无厂商锁定。

---

<p align="center">
  Made with ❤️ by AgentShield Community
</p>
