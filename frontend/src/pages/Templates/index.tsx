import React from 'react';
import { Row, Col, Card, Tag, Button, Input, Select } from 'antd';

const mockTemplates = [
  { key: '1', name: '研发Agent安全沙箱', desc: '适用于研发团队Agent，代码执行沙箱隔离+敏感密钥脱敏', category: '研发', level: 'standard', downloads: 1234 },
  { key: '2', name: '金融合规管控', desc: '满足等保三级+数据安全法要求，微虚拟机强隔离', category: '金融', level: 'restricted', downloads: 892 },
  { key: '3', name: '政务数据保护', desc: '满足政务数据安全规范，完整审计链+合规报告', category: '政务', level: 'restricted', downloads: 567 },
  { key: '4', name: '运维Agent沙箱', desc: '基础设施操作Agent，容器级隔离+操作审计', category: '运维', level: 'standard', downloads: 2103 },
  { key: '5', name: '客服DLP防护', desc: '客服Agent数据防泄漏，自动敏感信息脱敏', category: '客服', level: 'standard', downloads: 756 },
  { key: '6', name: '多租户隔离方案', desc: 'SaaS平台多租户完全隔离，独立沙箱+权限空间', category: '运维', level: 'trusted', downloads: 445 },
  { key: '7', name: '代码审查安全', desc: '代码审查Agent，只读沙箱+代码水印溯源', category: '研发', level: 'standard', downloads: 321 },
  { key: '8', name: 'API网关防护', desc: 'API调用Agent安全网关，流量监控+异常检测', category: '运维', level: 'trusted', downloads: 678 },
];

const levelColors: Record<string, string> = { trusted: 'green', standard: 'blue', restricted: 'orange' };

const Templates: React.FC = () => (
  <div>
    <div style={{ display: 'flex', gap: 12, marginBottom: 16 }}>
      <Input.Search placeholder="搜索模板" style={{ width: 250 }} />
      <Select defaultValue="all" style={{ width: 150 }} options={[
        { value: 'all', label: '全部分类' }, { value: '研发', label: '研发' }, { value: '金融', label: '金融' },
        { value: '政务', label: '政务' }, { value: '运维', label: '运维' }, { value: '客服', label: '客服' },
      ]} />
    </div>
    <Row gutter={[16, 16]}>
      {mockTemplates.map(t => (
        <Col span={6} key={t.key}>
          <Card title={t.name} size="small" extra={<Tag color={levelColors[t.level]}>{t.level}</Tag>}>
            <p style={{ color: '#666', fontSize: 13, minHeight: 48 }}>{t.desc}</p>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 8 }}>
              <Tag>{t.category}</Tag>
              <span style={{ color: '#999', fontSize: 12 }}>下载 {t.downloads}</span>
            </div>
            <Button type="primary" block style={{ marginTop: 8 }}>使用模板</Button>
          </Card>
        </Col>
      ))}
    </Row>
  </div>
);

export default Templates;
