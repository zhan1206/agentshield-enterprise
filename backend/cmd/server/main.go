package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zhan1206/agentshield-enterprise/backend/internal/adapter/data"
	"github.com/zhan1206/agentshield-enterprise/backend/internal/adapter/framework"
	"github.com/zhan1206/agentshield-enterprise/backend/internal/adapter/tool"
	"github.com/zhan1206/agentshield-enterprise/backend/internal/core/agent"
	"github.com/zhan1206/agentshield-enterprise/backend/internal/core/audit"
	"github.com/zhan1206/agentshield-enterprise/backend/internal/core/permission"
	"github.com/zhan1206/agentshield-enterprise/backend/internal/core/threat"
	"github.com/zhan1206/agentshield-enterprise/backend/internal/observability/metrics"
	"github.com/zhan1206/agentshield-enterprise/backend/internal/observability/tracing"
	"github.com/zhan1206/agentshield-enterprise/backend/internal/security/auth"
	"github.com/zhan1206/agentshield-enterprise/backend/internal/security/encryption"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize observability
	tracing.InitTracing("agentshield-server", "http://localhost:14268/api/traces")
	defer tracing.Shutdown()

	metricsCollector := metrics.NewCollector()

	// Initialize security
	jwtSecret := getEnv("JWT_SECRET", "agentshield-default-secret-change-me")
	authManager := auth.NewManager(jwtSecret, 24*time.Hour)

	encKey := getEnv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
	encEngine, err := encryption.NewEngine(encKey)
	if err != nil {
		log.Fatalf("Failed to initialize encryption engine: %v", err)
	}

	// Initialize core modules
	permEngine := permission.NewEngine()
	permLifecycle := permission.NewLifecycle(permEngine)
	auditChain := audit.NewChain()
	complianceEngine := audit.NewComplianceEngine()
	threatDetector := threat.NewDetector()
	threatResponder := threat.NewResponder(threatDetector)
	agentRegistry := agent.NewRegistry()

	// Initialize adapters
	fwAdapter := framework.NewAdapter(framework.AdapterConfig{}, nil)
	toolGateway := tool.NewGateway(nil, nil)
	dataGuard := data.NewGuard(encEngine)

	_ = permLifecycle
	_ = complianceEngine
	_ = threatResponder
	_ = fwAdapter
	_ = toolGateway
	_ = dataGuard
	_ = metricsCollector

	// Setup Gin router
	r := gin.Default()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(metrics.GinMiddleware())
	r.Use(tracing.GinMiddleware())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"version":   "0.1.0",
			"timestamp": time.Now().Unix(),
		})
	})

	// API v1 group
	v1 := r.Group("/api/v1")
	v1.Use(auth.AuthMiddleware(authManager))

	// Sandbox routes
	sandboxes := v1.Group("/sandboxes")
	{
		sandboxes.GET("", listSandboxes)
		sandboxes.POST("", createSandbox)
		sandboxes.GET("/:id", getSandbox)
		sandboxes.DELETE("/:id", deleteSandbox)
		sandboxes.POST("/:id/execute", executeInSandbox)
		sandboxes.POST("/:id/snapshot", createSnapshot)
		sandboxes.POST("/:id/rollback", rollbackSandbox)
		sandboxes.POST("/:id/stop", stopSandbox)
	}

	// Agent routes
	agents := v1.Group("/agents")
	{
		agents.GET("", listAgents)
		agents.POST("", registerAgent)
		agents.GET("/:id", getAgent)
		agents.PUT("/:id", updateAgent)
		agents.DELETE("/:id", deleteAgent)
		agents.GET("/:id/status", getAgentStatus)
	}

	// Permission routes
	perms := v1.Group("/permissions")
	{
		perms.GET("", listPermissions)
		perms.POST("", createPermission)
		perms.GET("/:id", getPermission)
		perms.PUT("/:id", updatePermission)
		perms.DELETE("/:id", deletePermission)
		perms.POST("/check", checkPermission)
	}

	// Audit routes
	audits := v1.Group("/audit")
	{
		audits.GET("/logs", listAuditLogs)
		audits.GET("/chain/verify", verifyAuditChain)
		audits.GET("/compliance/:standard", getComplianceReport)
	}

	// Threat routes
	threats := v1.Group("/threats")
	{
		threats.GET("/alerts", listThreatAlerts)
		threats.POST("/alerts/:id/respond", respondToThreat)
		threats.POST("/alerts/:id/resolve", resolveThreat)
		threats.GET("/rules", listThreatRules)
		threats.POST("/rules", createThreatRule)
	}

	// Data security routes
	dataSec := v1.Group("/data-security")
	{
		dataSec.POST("/scan", scanContent)
		dataSec.GET("/policies", listDataPolicies)
		dataSec.POST("/policies", createDataPolicy)
		dataSec.GET("/rules", listDetectionRules)
	}

	// Metrics routes
	v1.GET("/metrics", getMetrics)

	// Start server
	port := getEnv("SERVER_PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("🛡️ AgentShield Enterprise server starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// --- Handler stubs ---

func listSandboxes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"sandboxes": []interface{}{}, "total": 0})
}

func createSandbox(c *gin.Context) {
	var req struct {
		AgentID         string  `json:"agent_id" binding:"required"`
		IsolationLevel  string  `json:"isolation_level" binding:"required"`
		MemoryLimitMB   int     `json:"memory_limit_mb"`
		CPUQuota        float64 `json:"cpu_quota"`
		NetworkIsolated bool    `json:"network_isolated"`
		EnableSnapshot  bool    `json:"enable_snapshot"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":              fmt.Sprintf("sandbox-%d", time.Now().UnixNano()),
		"agent_id":        req.AgentID,
		"isolation_level": req.IsolationLevel,
		"status":          "running",
	})
}

func getSandbox(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "status": "running"})
}

func deleteSandbox(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func executeInSandbox(c *gin.Context) {
	var req struct {
		Command string `json:"command" binding:"required"`
		Timeout int    `json:"timeout"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"output": "executed", "exit_code": 0})
}

func createSnapshot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"snapshot_id": fmt.Sprintf("snap-%d", time.Now().UnixNano())})
}

func rollbackSandbox(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "rolled back"})
}

func stopSandbox(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "stopped"})
}

func listAgents(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"agents": []interface{}{}, "total": 0})
}

func registerAgent(c *gin.Context) {
	var req struct {
		Name          string `json:"name" binding:"required"`
		Framework     string `json:"framework" binding:"required"`
		SecurityLevel string `json:"security_level"`
		TeamID        string `json:"team_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":       fmt.Sprintf("agent-%d", time.Now().UnixNano()),
		"name":     req.Name,
		"framework": req.Framework,
		"status":   "active",
	})
}

func getAgent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "status": "active"})
}

func updateAgent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func deleteAgent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func getAgentStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "status": "active", "heartbeat": time.Now().Unix()})
}

func listPermissions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"permissions": []interface{}{}, "total": 0})
}

func createPermission(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"id": fmt.Sprintf("perm-%d", time.Now().UnixNano())})
}

func getPermission(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id")})
}

func updatePermission(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func deletePermission(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func checkPermission(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"allowed": true})
}

func listAuditLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"logs": []interface{}{}, "total": 0})
}

func verifyAuditChain(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"valid": true, "blocks_checked": 0})
}

func getComplianceReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"standard": c.Param("standard"), "score": 95.5})
}

func listThreatAlerts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"alerts": []interface{}{}, "total": 0})
}

func respondToThreat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "response executed"})
}

func resolveThreat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "resolved"})
}

func listThreatRules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"rules": []interface{}{}, "total": 0})
}

func createThreatRule(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"id": fmt.Sprintf("rule-%d", time.Now().UnixNano())})
}

func scanContent(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"findings": []interface{}{}, "masked_content": req.Content})
}

func listDataPolicies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"policies": []interface{}{}, "total": 0})
}

func createDataPolicy(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"id": fmt.Sprintf("policy-%d", time.Now().UnixNano())})
}

func listDetectionRules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"rules": []interface{}{}, "total": 0})
}

func getMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"active_sandboxes":  0,
		"running_agents":    0,
		"security_events":   0,
		"blocked_operations": 0,
		"threat_alerts":     0,
	})
}
