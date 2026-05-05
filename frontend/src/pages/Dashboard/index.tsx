import React from 'react';
import { Row, Col, Card, Statistic, Table, Tag } from 'antd';
import {
  CloudServerOutlined, RobotOutlined, AlertOutlined,
  StopOutlined, AuditOutlined, DatabaseOutlined, SafetyOutlined
} from '@ant-design/icons';
import { AreaChart, Area, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

const event24h = Array.from({ length: 24 }, (_, i) => ({
  hour: `${i}:00`,
  events: Math.floor(Math.random() * 30) + 5,
  blocked: Math.floor(Math.random() * 10),
}));

const threatDist = [
  { type: '越权访问', count: 45 },
  { type: '恶意代码', count: 23 },
  { type: '数据外泄', count: 18 },
  { type: '资源滥用', count: 31 },
  { type: '权限提升', count: 12 },
  { type: '异常行为', count: 27 },
];

const recentEvents = [
  { key: '1', time: '12:30:15', agent: 'agent-rd-001', type: '越权访问', level: 'high', status: '已拦截' },
  { key: '2', time: '12:28:42', agent: 'agent-finance-003', type: '数据外泄', level: 'critical', status: '已拦截' },
  { key: '3', time: '12:25:18', agent: 'agent-ops-002', type: '资源滥用', level: 'medium', status: '已告警' },
  { key: '4', time: '12:20:33', agent: 'agent-rd-004', type: '恶意代码', level: 'critical', status: '已终止' },
  { key: '5', time: '12:15:07', agent: 'agent-cs-001', type: '异常行为', level: 'low', status: '观察中' },
];

const levelColors: Record<string, string> = { critical: 'red', high: 'orange', medium: 'gold', low: 'blue' };

const Dashboard: React.FC = () => (
  <div>
    <Row gutter={[16, 16]}>
      <Col span={6}><Card><Statistic title="活跃沙箱" value={24} prefix={<CloudServerOutlined />} valueStyle={{ color: '#1677ff' }} /></Card></Col>
      <Col span={6}><Card><Statistic title="运行Agent" value={38} prefix={<RobotOutlined />} valueStyle={{ color: '#52c41a' }} /></Card></Col>
      <Col span={6}><Card><Statistic title="安全事件" value={156} prefix={<AlertOutlined />} valueStyle={{ color: '#fa8c16' }} /></Card></Col>
      <Col span={6}><Card><Statistic title="已拦截操作" value={23} prefix={<StopOutlined />} valueStyle={{ color: '#f5222d' }} /></Card></Col>
    </Row>
    <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
      <Col span={6}><Card><Statistic title="威胁告警" value={8} prefix={<AlertOutlined />} valueStyle={{ color: '#ff4d4f' }} /></Card></Col>
      <Col span={6}><Card><Statistic title="审计条目" value={12453} prefix={<AuditOutlined />} /></Card></Col>
      <Col span={6}><Card><Statistic title="数据扫描" value={3421} prefix={<DatabaseOutlined />} /></Card></Col>
      <Col span={6}><Card><Statistic title="合规评分" value={95.5} suffix="%" prefix={<SafetyOutlined />} valueStyle={{ color: '#52c41a' }} /></Card></Col>
    </Row>
    <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
      <Col span={14}>
        <Card title="安全事件趋势（24小时）">
          <ResponsiveContainer width="100%" height={280}>
            <AreaChart data={event24h}><CartesianGrid strokeDasharray="3 3" /><XAxis dataKey="hour" /><YAxis /><Tooltip /><Area type="monotone" dataKey="events" stroke="#1677ff" fill="#1677ff" fillOpacity={0.2} /><Area type="monotone" dataKey="blocked" stroke="#f5222d" fill="#f5222d" fillOpacity={0.2} /></AreaChart>
          </ResponsiveContainer>
        </Card>
      </Col>
      <Col span={10}>
        <Card title="威胁类型分布">
          <ResponsiveContainer width="100%" height={280}>
            <BarChart data={threatDist}><CartesianGrid strokeDasharray="3 3" /><XAxis dataKey="type" /><YAxis /><Tooltip /><Bar dataKey="count" fill="#1677ff" /></BarChart>
          </ResponsiveContainer>
        </Card>
      </Col>
    </Row>
    <Card title="最近安全事件" style={{ marginTop: 16 }}>
      <Table dataSource={recentEvents} pagination={false} size="small"
        columns={[
          { title: '时间', dataIndex: 'time', width: 100 },
          { title: 'Agent', dataIndex: 'agent', width: 160 },
          { title: '类型', dataIndex: 'type', width: 100 },
          { title: '级别', dataIndex: 'level', width: 80, render: (v: string) => <Tag color={levelColors[v]}>{v.toUpperCase()}</Tag> },
          { title: '状态', dataIndex: 'status', width: 100 },
        ]}
      />
    </Card>
  </div>
);

export default Dashboard;
