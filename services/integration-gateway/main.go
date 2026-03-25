package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"enterprise-hcm-lms/pkg/resilience" // Assuming your module name
)

// Define the incoming Webhook structure from HCM/CRM
type WebhookRequest struct {
	Source    string                 `json:"source" binding:"required"` // e.g., "Workday"
	EventID   string                 `json:"event_id" binding:"required"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload" binding:"required"`
}

func main() {
	r := gin.Default()

	// Initialize the Circuit Breaker from your pkg/resilience
	// This protects your internal Kafka/Bus from being overwhelmed
	cb := &resilience.CircuitBreaker{Threshold: 5}

	// 1. Security Middleware (The 'Executive Advisor' requirement)
	r.Use(APIKeyAuth())

	// 2. The Ingestion Endpoint
	r.POST("/v1/ingest", func(c *gin.Context) {
		var req WebhookRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Schema Contract: " + err.Error()})
			return
		}

		// 3. Resilience Check: Execute handoff within the Circuit Breaker
		err := cb.Execute(func() error {
			// In a real PoC, this is where you'd push to Confluent Kafka
			log.Printf("[Nervous-System] Accepted %s event from %s", req.EventID, req.Source)
			return nil 
		})

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"status": "accepted",
			"ingested_at": time.Now().UTC(),
		})
	})

	log.Println("Integration Gateway running on :8080")
	r.Run(":8080")
}

// APIKeyAuth is a simple guardrail for the PoC
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Integration-Key")
		if key != "alf-secret-poc-key" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Missing or Invalid Key"})
			return
		}
		c.Next()
	}
}
