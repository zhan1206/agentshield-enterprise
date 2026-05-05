use crate::engine::IsolationLevel;
use crate::Result;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;

/// Monitoring data for a sandbox
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MonitorData {
    pub sandbox_id: String,
    pub memory_usage_mb: f64,
    pub cpu_usage_percent: f64,
    pub network_bytes_in: u64,
    pub network_bytes_out: u64,
    pub process_count: u32,
    pub uptime_seconds: u64,
}

/// Sandbox monitor for real-time resource tracking
pub struct SandboxMonitor {
    data: Arc<RwLock<HashMap<String, MonitorData>>>,
}

impl SandboxMonitor {
    /// Create a new monitor
    pub fn new() -> Self {
        Self {
            data: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    /// Register a sandbox for monitoring
    pub async fn register(&self, sandbox_id: &str) {
        let monitor_data = MonitorData {
            sandbox_id: sandbox_id.to_string(),
            memory_usage_mb: 0.0,
            cpu_usage_percent: 0.0,
            network_bytes_in: 0,
            network_bytes_out: 0,
            process_count: 1,
            uptime_seconds: 0,
        };
        self.data.write().await.insert(sandbox_id.to_string(), monitor_data);
    }

    /// Update monitoring data for a sandbox
    pub async fn update(&self, sandbox_id: &str, memory_mb: f64, cpu_pct: f64) {
        let mut data = self.data.write().await;
        if let Some(m) = data.get_mut(sandbox_id) {
            m.memory_usage_mb = memory_mb;
            m.cpu_usage_percent = cpu_pct;
            m.uptime_seconds += 1;
        }
    }

    /// Get monitoring data for a sandbox
    pub async fn get(&self, sandbox_id: &str) -> Option<MonitorData> {
        self.data.read().await.get(sandbox_id).cloned()
    }

    /// Unregister a sandbox from monitoring
    pub async fn unregister(&self, sandbox_id: &str) {
        self.data.write().await.remove(sandbox_id);
    }

    /// Check if any sandbox exceeds resource limits
    pub async fn check_limits(&self, memory_limit_mb: f64, cpu_limit_pct: f64) -> Vec<String> {
        let data = self.data.read().await;
        let mut violations = Vec::new();
        for (id, m) in data.iter() {
            if m.memory_usage_mb > memory_limit_mb || m.cpu_usage_percent > cpu_limit_pct {
                violations.push(id.clone());
            }
        }
        violations
    }
}
