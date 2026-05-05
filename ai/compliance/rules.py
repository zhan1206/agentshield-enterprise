"""AgentShield Enterprise - Compliance Rules Engine

Evaluates compliance against regulatory frameworks (等保2.0, GDPR, SOX, 数据安全法).
"""

from typing import Any
from dataclasses import dataclass, field
from enum import Enum


class ComplianceStandard(str, Enum):
    DJL_2_0 = "djl_2_0"
    GDPR = "gdpr"
    SOX = "sox"
    DATA_SECURITY_LAW = "data_security_law"


class CheckStatus(str, Enum):
    PASS = "pass"
    WARNING = "warning"
    FAIL = "fail"


@dataclass
class ComplianceRule:
    id: str
    name: str
    standard: ComplianceStandard
    category: str
    required: bool = True
    description: str = ""
    check_function: str = ""


@dataclass
class ComplianceCheckResult:
    rule_id: str
    rule_name: str
    status: CheckStatus
    score: float
    details: str = ""


class ComplianceEngine:
    """Evaluates compliance against multiple regulatory frameworks."""

    def __init__(self):
        self.rules: dict[str, ComplianceRule] = {}
        self._load_builtin_rules()

    def _load_builtin_rules(self) -> None:
        """Load built-in compliance rules."""
        builtin = [
            # 等保2.0
            ComplianceRule("djl2-audit-1", "审计日志完整性", ComplianceStandard.DJL_2_0, "安全审计", True, "审计日志不可篡改", "verify_audit_chain"),
            ComplianceRule("djl2-audit-2", "审计日志保留期限", ComplianceStandard.DJL_2_0, "安全审计", True, "审计日志保留≥180天", "check_audit_retention"),
            ComplianceRule("djl2-access-1", "访问控制策略", ComplianceStandard.DJL_2_0, "访问控制", True, "实施ABAC+RBAC访问控制", "check_access_control"),
            ComplianceRule("djl2-data-1", "数据加密存储", ComplianceStandard.DJL_2_0, "数据安全", True, "敏感数据加密存储", "check_data_encryption"),
            ComplianceRule("djl2-data-2", "数据分类分级", ComplianceStandard.DJL_2_0, "数据安全", True, "实施数据分类分级管理", "check_data_classification"),
            # GDPR
            ComplianceRule("gdpr-1", "数据处理合法性", ComplianceStandard.GDPR, "数据处理", True, "数据处理需有合法基础", "check_lawful_basis"),
            ComplianceRule("gdpr-2", "数据主体权利", ComplianceStandard.GDPR, "数据权利", True, "保障数据主体访问/删除权利", "check_data_subject_rights"),
            ComplianceRule("gdpr-3", "数据泄露通知", ComplianceStandard.GDPR, "数据泄露", True, "72小时内通知监管机构", "check_breach_notification"),
            # SOX
            ComplianceRule("sox-1", "内部控制有效性", ComplianceStandard.SOX, "内部控制", True, "财务相关系统需有效内部控制", "check_internal_controls"),
            ComplianceRule("sox-2", "审计追踪", ComplianceStandard.SOX, "审计追踪", True, "关键操作需完整审计追踪", "check_audit_trail"),
            # 数据安全法
            ComplianceRule("dsl-1", "数据安全管理制度", ComplianceStandard.DATA_SECURITY_LAW, "管理制度", True, "建立数据安全管理制度", "check_data_security_mgmt"),
            ComplianceRule("dsl-2", "数据风险评估", ComplianceStandard.DATA_SECURITY_LAW, "风险评估", True, "定期开展数据风险评估", "check_risk_assessment"),
        ]
        for rule in builtin:
            self.rules[rule.id] = rule

    def evaluate(self, standard: ComplianceStandard, context: dict[str, Any]) -> list[ComplianceCheckResult]:
        """Evaluate compliance for a given standard."""
        results = []
        for rule in self.rules.values():
            if rule.standard != standard:
                continue
            result = self._evaluate_rule(rule, context)
            results.append(result)
        return results

    def _evaluate_rule(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        """Evaluate a single compliance rule against context."""
        check_fn = getattr(self, f"_check_{rule.check_function}", None)
        if check_fn:
            return check_fn(rule, context)

        # Default: check context for matching key
        key = rule.check_function
        if key in context:
            value = context[key]
            if value is True or (isinstance(value, (int, float)) and value >= 80):
                return ComplianceCheckResult(rule.id, rule.name, CheckStatus.PASS, 100.0)
            elif isinstance(value, (int, float)) and value >= 60:
                return ComplianceCheckResult(rule.id, rule.name, CheckStatus.WARNING, value)
            else:
                return ComplianceCheckResult(rule.id, rule.name, CheckStatus.FAIL, value if isinstance(value, (int, float)) else 0.0)

        return ComplianceCheckResult(rule.id, rule.name, CheckStatus.WARNING, 50.0, "Not evaluated")

    def _check_verify_audit_chain(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        chain_valid = context.get("audit_chain_valid", True)
        status = CheckStatus.PASS if chain_valid else CheckStatus.FAIL
        score = 100.0 if chain_valid else 0.0
        return ComplianceCheckResult(rule.id, rule.name, status, score, f"Audit chain valid: {chain_valid}")

    def _check_audit_retention(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        retention_days = context.get("audit_retention_days", 180)
        status = CheckStatus.PASS if retention_days >= 180 else CheckStatus.FAIL
        score = min(100.0, retention_days / 180 * 100)
        return ComplianceCheckResult(rule.id, rule.name, status, score, f"Retention: {retention_days} days")

    def _check_access_control(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        has_abac = context.get("has_abac", True)
        has_rbac = context.get("has_rbac", True)
        score = (50.0 if has_abac else 0) + (50.0 if has_rbac else 0)
        status = CheckStatus.PASS if score >= 80 else CheckStatus.WARNING if score >= 50 else CheckStatus.FAIL
        return ComplianceCheckResult(rule.id, rule.name, status, score)

    def _check_data_encryption(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        encrypted = context.get("data_encrypted", True)
        status = CheckStatus.PASS if encrypted else CheckStatus.FAIL
        return ComplianceCheckResult(rule.id, rule.name, status, 100.0 if encrypted else 0.0)

    def _check_data_classification(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        classified = context.get("data_classified", True)
        status = CheckStatus.PASS if classified else CheckStatus.WARNING
        return ComplianceCheckResult(rule.id, rule.name, status, 90.0 if classified else 40.0)

    def _check_lawful_basis(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        has_basis = context.get("lawful_basis", True)
        status = CheckStatus.PASS if has_basis else CheckStatus.FAIL
        return ComplianceCheckResult(rule.id, rule.name, status, 100.0 if has_basis else 0.0)

    def _check_data_subject_rights(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        rights_enabled = context.get("data_subject_rights_enabled", True)
        status = CheckStatus.PASS if rights_enabled else CheckStatus.FAIL
        return ComplianceCheckResult(rule.id, rule.name, status, 100.0 if rights_enabled else 0.0)

    def _check_breach_notification(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        notification_setup = context.get("breach_notification_setup", True)
        status = CheckStatus.PASS if notification_setup else CheckStatus.WARNING
        return ComplianceCheckResult(rule.id, rule.name, status, 100.0 if notification_setup else 50.0)

    def _check_internal_controls(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        controls = context.get("internal_controls_score", 85)
        status = CheckStatus.PASS if controls >= 80 else CheckStatus.WARNING if controls >= 60 else CheckStatus.FAIL
        return ComplianceCheckResult(rule.id, rule.name, status, float(controls))

    def _check_audit_trail(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        trail_complete = context.get("audit_trail_complete", True)
        status = CheckStatus.PASS if trail_complete else CheckStatus.FAIL
        return ComplianceCheckResult(rule.id, rule.name, status, 100.0 if trail_complete else 30.0)

    def _check_data_security_mgmt(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        mgmt_exists = context.get("data_security_management", True)
        status = CheckStatus.PASS if mgmt_exists else CheckStatus.FAIL
        return ComplianceCheckResult(rule.id, rule.name, status, 100.0 if mgmt_exists else 0.0)

    def _check_risk_assessment(self, rule: ComplianceRule, context: dict[str, Any]) -> ComplianceCheckResult:
        last_assessment = context.get("last_risk_assessment_days", 30)
        status = CheckStatus.PASS if last_assessment <= 90 else CheckStatus.WARNING if last_assessment <= 180 else CheckStatus.FAIL
        score = max(0, 100 - (last_assessment - 90) * 0.5)
        return ComplianceCheckResult(rule.id, rule.name, status, score, f"Last assessment: {last_assessment} days ago")

    def get_rules(self, standard: ComplianceStandard | None = None) -> list[ComplianceRule]:
        """Get compliance rules, optionally filtered by standard."""
        if standard:
            return [r for r in self.rules.values() if r.standard == standard]
        return list(self.rules.values())
