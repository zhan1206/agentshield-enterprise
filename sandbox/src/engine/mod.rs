pub mod process;
pub mod container;
pub mod microvm;
pub mod monitor;
pub mod snapshot;

use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;
use uuid::Uuid;
use chrono::Utc;

/// Isolation level for sandbox environments
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum IsolationLevel {
    Process,
    Container,
    MicroVM,
}

/// Sandbox configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SandboxConfig {
    pub agent_id: String,
    pub isolation_level: IsolationLevel,
    pub memory_limit_mb: u64,
    pub cpu_quota: f64,
    pub network_isolated: bool,
    pub enable_snapshot: bool,
    pub env_vars: HashMap<String, String>,
}

impl Default for SandboxConfig {
    fn default() -> Self {
        Self {
            agent_id: String::new(),
            isolation_level: IsolationLevel::Container,
            memory_limit_mb: 512,
            cpu_quota: 2.0,
            network_isolated: true,
            enable_snapshot: true,
            env_vars: HashMap::new(),
        }
    }
}

/// Sandbox status
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum SandboxStatus {
    Creating,
    Running,
    Stopped,
    Error,
}

/// Information about a running sandbox
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SandboxInfo {
    pub id: String,
    pub agent_id: String,
    pub isolation_level: IsolationLevel,
    pub status: SandboxStatus,
    pub config: SandboxConfig,
    pub created_at: String,
    pub memory_usage_mb: f64,
    pub cpu_usage_percent: f64,
}

/// Execution result from sandbox
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExecutionResult {
    pub exit_code: i32,
    pub stdout: String,
    pub stderr: String,
    pub duration_ms: u64,
}

/// The main sandbox engine managing all sandbox instances
pub struct SandboxEngine {
    sandboxes: Arc<RwLock<HashMap<String, SandboxInfo>>>,
}

impl SandboxEngine {
    /// Create a new sandbox engine
    pub fn new() -> Self {
        Self {
            sandboxes: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    /// Create a new sandbox with the given configuration
    pub async fn create(&self, config: SandboxConfig) -> crate::Result<SandboxInfo> {
        let id = format!("sandbox-{}", Uuid::new_v4().to_string().split('-').next().unwrap());

        let info = SandboxInfo {
            id: id.clone(),
            agent_id: config.agent_id.clone(),
            isolation_level: config.isolation_level.clone(),
            status: SandboxStatus::Running,
            config: config.clone(),
            created_at: Utc::now().to_rfc3339(),
            memory_usage_mb: 0.0,
            cpu_usage_percent: 0.0,
        };

        match config.isolation_level {
            IsolationLevel::Process => {
                process::ProcessSandbox::create(&id, &config).await?;
            }
            IsolationLevel::Container => {
                container::ContainerSandbox::create(&id, &config).await?;
            }
            IsolationLevel::MicroVM => {
                microvm::MicroVMSandbox::create(&id, &config).await?;
            }
        }

        self.sandboxes.write().await.insert(id.clone(), info.clone());
        tracing::info!("Created sandbox {} with {:?} isolation for agent {}", id, config.isolation_level, config.agent_id);
        Ok(info)
    }

    /// Execute a command in a sandbox
    pub async fn execute(&self, sandbox_id: &str, command: &str, timeout_secs: u64) -> crate::Result<ExecutionResult> {
        let sandboxes = self.sandboxes.read().await;
        let info = sandboxes.get(sandbox_id)
            .ok_or_else(|| crate::SandboxError::NotFound(sandbox_id.to_string()))?;

        if info.status != SandboxStatus::Running {
            return Err(crate::SandboxError::NotRunning(sandbox_id.to_string()));
        }

        let isolation = info.isolation_level.clone();
        drop(sandboxes);

        let result = match isolation {
            IsolationLevel::Process => process::ProcessSandbox::execute(sandbox_id, command, timeout_secs).await?,
            IsolationLevel::Container => container::ContainerSandbox::execute(sandbox_id, command, timeout_secs).await?,
            IsolationLevel::MicroVM => microvm::MicroVMSandbox::execute(sandbox_id, command, timeout_secs).await?,
        };

        Ok(result)
    }

    /// Stop a sandbox
    pub async fn stop(&self, sandbox_id: &str) -> crate::Result<()> {
        let mut sandboxes = self.sandboxes.write().await;
        let info = sandboxes.get_mut(sandbox_id)
            .ok_or_else(|| crate::SandboxError::NotFound(sandbox_id.to_string()))?;

        if info.status != SandboxStatus::Running {
            return Err(crate::SandboxError::NotRunning(sandbox_id.to_string()));
        }

        info.status = SandboxStatus::Stopped;
        tracing::info!("Stopped sandbox {}", sandbox_id);
        Ok(())
    }

    /// Delete a sandbox
    pub async fn delete(&self, sandbox_id: &str) -> crate::Result<()> {
        let mut sandboxes = self.sandboxes.write().await;
        sandboxes.remove(sandbox_id)
            .ok_or_else(|| crate::SandboxError::NotFound(sandbox_id.to_string()))?;
        tracing::info!("Deleted sandbox {}", sandbox_id);
        Ok(())
    }

    /// Get sandbox info
    pub async fn get(&self, sandbox_id: &str) -> crate::Result<SandboxInfo> {
        let sandboxes = self.sandboxes.read().await;
        sandboxes.get(sandbox_id)
            .cloned()
            .ok_or_else(|| crate::SandboxError::NotFound(sandbox_id.to_string()))
    }

    /// List all sandboxes
    pub async fn list(&self) -> Vec<SandboxInfo> {
        self.sandboxes.read().await.values().cloned().collect()
    }

    /// Create a snapshot of a sandbox
    pub async fn snapshot(&self, sandbox_id: &str) -> crate::Result<String> {
        let sandboxes = self.sandboxes.read().await;
        let info = sandboxes.get(sandbox_id)
            .ok_or_else(|| crate::SandboxError::NotFound(sandbox_id.to_string()))?;
        if !info.config.enable_snapshot {
            return Err(crate::SandboxError::SnapshotError("Snapshots not enabled for this sandbox".to_string()));
        }
        let snap_id = snapshot::SnapshotManager::create(sandbox_id, &info.isolation_level).await?;
        Ok(snap_id)
    }

    /// Rollback a sandbox to a snapshot
    pub async fn rollback(&self, sandbox_id: &str, snapshot_id: &str) -> crate::Result<()> {
        let sandboxes = self.sandboxes.read().await;
        let _ = sandboxes.get(sandbox_id)
            .ok_or_else(|| crate::SandboxError::NotFound(sandbox_id.to_string()))?;
        snapshot::SnapshotManager::restore(sandbox_id, snapshot_id).await
    }
}
