import React, { useState } from 'react';
import { Row, Col, Card, Statistic, Tabs, Table, Input, Button, Tag, Switch } from 'antd';
import { DatabaseOutlined, ScanOutlined, SafetyOutlined } from '@ant-design/icons';
import { SearchOutlined } from '@ant-design/icons';

const mockRules = [
  { key: '1', id: 'dr-001', name: '中国身份证号', category: '身份证', sensitivity: 'secret', enabled: true },
  { key: '2', id: 'dr-002', name: '手机号码', category: '手机号', sensitivity: 'confidential', enabled: true },
  { key: '3', id: 'dr-003', name: '银行卡号', category: '银行卡', sensitivity: 'secret', enabled: true },
  { key: '4', id: 'dr-004', name: 'API密钥', category: '密钥', sensitivity: 'secret', enabled: true },
  { key: '5', id: 'dr-005', name: '邮箱地址', category: '邮箱', sensitivity: 'confidential', enabled: true },
];

const sensColors: Record<string, string> = { public: 'green', internal: 'blue', confidential: 'orange', secret: 'red' };

const DataSecurity: React.FC = () => {
  const [scanInput, setScanInput] = useState('');

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col span={8}><Card><Statistic title="今日扫描" value={342} prefix={<ScanOutlined />} /></Card></Col>
        <Col span={8}><Card><Statistic title="发现风险" value={18} prefix={<DatabaseOutlined />} valueStyle={{ color: '#fa8c16' }} /></Card></Col>
        <Col span={8}><Card><Statistic title="已脱敏" value={156} prefix={<SafetyOutlined />} valueStyle={{ color: '#52c41a' }} /></Card></Col>
      </Row>
      <Tabs defaultActiveKey="scan" style={{ marginTop: 16 }} items={[
        {
          key: 'scan',
          label: '内容扫描',
          children: (
            <div>
              <Input.TextArea rows={6} value={scanInput} onChange={e => setScanInput(e.target.value)} placeholder="输入要扫描的内容..." />
              <Button type="primary" icon={<SearchOutlined />} style={{ marginTop: 8 }}>扫描</Button>
              <Card title="扫描结果" style={{ marginTop: 16 }}>
                <p style={{ color: '#999' }}>输入内容后点击扫描按钮查看结果</p>
              </Card>
            </div>
          ),
        },
        {
          key: 'rules',
          label: '检测规则',
          children: (
            <Table dataSource={mockRules} size="small"
              columns={[
                { title: 'ID', dataIndex: 'id', width: 80 },
                { title: '名称', dataIndex: 'name', width: 150 },
                { title: '分类', dataIndex: 'category', width: 100 },
                { title: '敏感级别', dataIndex: 'sensitivity', width: 100, render: (v: string) => <Tag color={sensColors[v]}>{v}</Tag> },
                { title: '启用', dataIndex: 'enabled', width: 80, render: (v: boolean) => <Switch checked={v} size="small" /> },
              ]}
            />
          ),
        },
        {
          key: 'policies',
          label: '安全策略',
          children: <Table dataSource={[]} size="small" columns={[{ title: '策略名称' }, { title: '类型' }, { title: '状态' }, { title: '操作' }]} />,
        },
      ]} />
    </div>
  );
};

export default DataSecurity;
