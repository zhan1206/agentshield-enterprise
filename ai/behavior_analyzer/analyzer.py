"""AgentShield Enterprise - AI Behavior Analyzer

Detects anomalous agent behavior using ML models and rule-based analysis.
"""

import time
import logging
from typing import Any
from dataclasses import dataclass, field
from enum import Enum

logger = logging.getLogger(__name__)


class AnomalyLevel(str, Enum):
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"
    CRITICAL = "critical"


@dataclass
class BehaviorEvent:
    agent_id: str
    action: str
    resource: str
    result: str
    timestamp: float = field(default_factory=time.time)
    details: dict[str, Any] = field(default_factory=dict)


@dataclass
class AnomalyResult:
    agent_id: str
    anomaly_score: float
    level: AnomalyLevel
    description: str
    indicators: list[str] = field(default_factory=list)


class BehaviorAnalyzer:
    """Analyzes agent behavior patterns to detect anomalies."""

    def __init__(self, threshold: float = 0.8):
        self.threshold = threshold
        self.history: dict[str, list[BehaviorEvent]] = {}
        self.action_counts: dict[str, dict[str, int]] = {}

    def record_event(self, event: BehaviorEvent) -> None:
        """Record a behavior event for an agent."""
        if event.agent_id not in self.history:
            self.history[event.agent_id] = []
            self.action_counts[event.agent_id] = {}
        self.history[event.agent_id].append(event)
        action_key = f"{event.action}:{event.resource}"
        self.action_counts[event.agent_id][action_key] = \
            self.action_counts[event.agent_id].get(action_key, 0) + 1

        # Keep only last 1000 events per agent
        if len(self.history[event.agent_id]) > 1000:
            self.history[event.agent_id] = self.history[event.agent_id][-1000:]

    def analyze(self, event: BehaviorEvent) -> AnomalyResult | None:
        """Analyze a behavior event and return anomaly result if detected."""
        score = self._compute_anomaly_score(event)
        if score < self.threshold:
            return None

        level = self._score_to_level(score)
        indicators = self._get_indicators(event, score)

        result = AnomalyResult(
            agent_id=event.agent_id,
            anomaly_score=score,
            level=level,
            description=f"Anomalous behavior detected for agent {event.agent_id}",
            indicators=indicators,
        )
        logger.warning(f"Anomaly detected: {result}")
        return result

    def _compute_anomaly_score(self, event: BehaviorEvent) -> float:
        """Compute anomaly score based on behavior patterns."""
        score = 0.0

        # Check for failed operations
        if event.result == "denied":
            score += 0.3
        elif event.result == "failed":
            score += 0.2

        # Check for suspicious actions
        suspicious_actions = ["execute", "delete", "modify", "export", "send"]
        if event.action in suspicious_actions:
            score += 0.2

        # Check for high-frequency actions
        agent_actions = self.action_counts.get(event.agent_id, {})
        action_key = f"{event.action}:{event.resource}"
        count = agent_actions.get(action_key, 0)
        if count > 50:
            score += 0.3
        elif count > 20:
            score += 0.15

        # Check for rapid sequential actions
        events = self.history.get(event.agent_id, [])
        if len(events) >= 5:
            recent = events[-5:]
            time_span = recent[-1].timestamp - recent[0].timestamp
            if time_span < 1.0:  # 5 actions in less than 1 second
                score += 0.4

        return min(score, 1.0)

    def _score_to_level(self, score: float) -> AnomalyLevel:
        """Convert anomaly score to severity level."""
        if score >= 0.9:
            return AnomalyLevel.CRITICAL
        elif score >= 0.7:
            return AnomalyLevel.HIGH
        elif score >= 0.5:
            return AnomalyLevel.MEDIUM
        return AnomalyLevel.LOW

    def _get_indicators(self, event: BehaviorEvent, score: float) -> list[str]:
        """Get human-readable indicators for the anomaly."""
        indicators = []
        if event.result == "denied":
            indicators.append("Permission denied")
        if event.result == "failed":
            indicators.append("Operation failed")
        if event.action in ["execute", "delete"]:
            indicators.append(f"Suspicious action: {event.action}")
        agent_actions = self.action_counts.get(event.agent_id, {})
        action_key = f"{event.action}:{event.resource}"
        count = agent_actions.get(action_key, 0)
        if count > 20:
            indicators.append(f"High frequency action ({count} times)")
        events = self.history.get(event.agent_id, [])
        if len(events) >= 5:
            recent = events[-5:]
            time_span = recent[-1].timestamp - recent[0].timestamp
            if time_span < 1.0:
                indicators.append("Rapid sequential actions")
        return indicators

    def get_agent_summary(self, agent_id: str) -> dict[str, Any]:
        """Get a summary of an agent's behavior."""
        events = self.history.get(agent_id, [])
        actions = self.action_counts.get(agent_id, {})
        return {
            "agent_id": agent_id,
            "total_events": len(events),
            "unique_actions": len(actions),
            "top_actions": sorted(actions.items(), key=lambda x: x[1], reverse=True)[:5],
        }
