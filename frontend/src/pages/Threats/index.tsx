import React from 'react';
import { Row, Col, Card, Statistic, Table, Button, Tag, Space } from 'antd';
import { AlertOutlined } from '@ant-design/icons';

const alertCards = [
  { title: '严重', count: 2, color: '#f5222d' },
  { title: '高危', count: 5, color: '#fa8c16' },
  { title: '中危', count: 12, color: '#fadb14' },
  { title: '低危', count: 8, color: '#1677ff' },
];

const levelColors: Record<string, string> = { critical: 'red', high: 'orange', medium: 'gold', low: 'blue' };

const mockAlerts = [
  { key: '1', time: '12:30:15', agent: 'agent-rd-004', type: '恶意代码执行', level: 'critical', desc: '检测到rm -rf命令', status: 'active' },
  { key: '2', time: '12:28:42', agent: 'agent-finance-003', type: '数据外泄', level: 'high', desc: '尝试发送敏感数据到外部', status: 'active' },
  { key: '3', time: '12:25:18', agent: 'agent-ops-002', type: '资源滥用', level: 'medium', desc: 'CPU使用率超过90%', status: 'responded' },
  { key: '4', time: '12:20:33', agent: 'agent-rd-001', type: '越权访问', level: 'low', desc: '访问未授权目录', status: 'resolved' },
];

const Threats: React.FC = () => (
  <div>
    <Row gutter={[16, 16]}>
      {alertCards.map(c => (
        <Col span={6} key={c.title}>
          <Card><Statistic title={c.title} value={c.count} prefix={<AlertOutlined />} valueStyle={{ color: c.color }} /></Card>
        </Col>
      ))}
    </Row>
    <Card title="威胁告警" style={{ marginTop: 16 }}>
      <Table dataSource={mockAlerts} size="small"
        columns={[
          { title: '时间', dataIndex: 'time', width: 80 },
          { title: 'Agent', dataIndex: 'agent', width: 150 },
          { title: '威胁类型', dataIndex: 'type', width: 120 },
          { title: '级别', dataIndex: 'level', width: 80, render: (v: string) => <Tag color={levelColors[v]}>{v}</Tag> },
          { title: '描述', dataIndex: 'desc', width: 220 },
          { title: '状态', dataIndex: 'status', width: 90, render: (v: string) => <Tag color={v === 'active' ? 'red' : v === 'responded' ? 'orange' : 'green'}>{v === 'active' ? '活跃' : v === 'responded' ? '已响应' : '已解决'}</Tag> },
          { title: '响应', width: 280, render: () => <Space size="small"><Button size="small" danger>拦截</Button><Button size="small">暂停</Button><Button size="small">终止</Button><Button size="small">隔离</Button><Button size="small">回滚</Button></Space> },
        ]}
      />
    </Card>
  </div>
);

export default Threats;
