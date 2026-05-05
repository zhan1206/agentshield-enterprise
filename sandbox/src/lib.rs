//! AgentShield Enterprise - Three-tier Isolation Sandbox Engine
//!
//! This crate provides the core sandbox engine for AgentShield Enterprise,
//! offering three levels of isolation:
//! - Process-level: Lightweight isolation for trusted agents
//! - Container-level: Standard isolation using Linux namespaces/cgroups
//! - MicroVM-level: Strongest isolation using lightweight VMs

pub mod engine;

pub use engine::{SandboxEngine, IsolationLevel, SandboxConfig, SandboxStatus};

/// Error type for sandbox operations
#[derive(Debug, thiserror::Error)]
pub enum SandboxError {
    #[error("Sandbox not found: {0}")]
    NotFound(String),

    #[error("Sandbox already running: {0}")]
    AlreadyRunning(String),

    #[error("Sandbox not running: {0}")]
    NotRunning(String),

    #[error("Isolation level not supported: {0}")]
    UnsupportedIsolation(String),

    #[error("Resource limit exceeded: {0}")]
    ResourceLimitExceeded(String),

    #[error("Snapshot error: {0}")]
    SnapshotError(String),

    #[error("Execution error: {0}")]
    ExecutionError(String),

    #[error("Internal error: {0}")]
    Internal(#[from] anyhow::Error),
}

/// Result type for sandbox operations
pub type Result<T> = std::result::Result<T, SandboxError>;
