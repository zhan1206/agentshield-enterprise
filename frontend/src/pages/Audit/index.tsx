import React, { useState } from 'react';
import { Table, Button, Tag, Space, Input, Select, DatePicker, Modal } from 'antd';

const resultColors: Record<string, string> = { success: 'green', failed: 'red', denied: 'orange' };

const mockLogs = [
  { key: '1', id: 'aud-001', time: '2026-05-05 12:30:15', agent: 'agent-rd-001', action: 'execute', resource: '/app/main.py', result: 'success', hash: 'a1b2c3d4...e5f6' },
  { key: '2', id: 'aud-002', time: '2026-05-05 12:28:42', agent: 'agent-finance-003', action: 'access', resource: '/data/financial.db', result: 'denied', hash: 'b2c3d4e5...f6a1' },
  { key: '3', id: 'aud-003', time: '2026-05-05 12:25:18', agent: 'agent-ops-002', action: 'modify', resource: '/etc/nginx.conf', result: 'success', hash: 'c3d4e5f6...a1b2' },
  { key: '4', id: 'aud-004', time: '2026-05-05 12:20:33', agent: 'agent-rd-004', action: 'delete', resource: '/tmp/cache', result: 'failed', hash: 'd4e5f6a1...b2c3' },
];

const Audit: React.FC = () => {
  const [detailOpen, setDetailOpen] = useState(false);

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Input placeholder="Agent ID" style={{ width: 160 }} />
        <Select defaultValue="all" style={{ width: 120 }} options={[{ value: 'all', label: '全部操作' }, { value: 'execute', label: '执行' }, { value: 'access', label: '访问' }, { value: 'modify', label: '修改' }, { value: 'delete', label: '删除' }]} />
        <DatePicker.RangePicker />
        <Button type="primary">搜索</Button>
      </Space>
      <Table dataSource={mockLogs} size="small"
        columns={[
          { title: 'ID', dataIndex: 'id', width: 100 },
          { title: '时间', dataIndex: 'time', width: 160 },
          { title: 'Agent', dataIndex: 'agent', width: 150 },
          { title: '操作', dataIndex: 'action', width: 80, render: (v: string) => <Tag color="blue">{v}</Tag> },
          { title: '资源', dataIndex: 'resource', width: 180 },
          { title: '结果', dataIndex: 'result', width: 80, render: (v: string) => <Tag color={resultColors[v]}>{v}</Tag> },
          { title: '哈希', dataIndex: 'hash', width: 130, render: (v: string) => <code style={{ fontSize: 12 }}>{v}</code> },
          { title: '操作', width: 100, render: () => <Button size="small" onClick={() => setDetailOpen(true)}>详情</Button> },
        ]}
      />
      <Modal title="审计链详情" open={detailOpen} onCancel={() => setDetailOpen(false)} footer={null} width={600}>
        <p><strong>条目ID:</strong> aud-001</p>
        <p><strong>时间:</strong> 2026-05-05 12:30:15</p>
        <p><strong>Agent:</strong> agent-rd-001</p>
        <p><strong>操作:</strong> execute</p>
        <p><strong>资源:</strong> /app/main.py</p>
        <p><strong>结果:</strong> <Tag color="green">success</Tag></p>
        <p><strong>前一条哈希:</strong> <code>prev_hash_value...</code></p>
        <p><strong>当前哈希:</strong> <code>a1b2c3d4e5f6...</code></p>
        <p><strong>链完整性:</strong> <Tag color="green">✅ 验证通过</Tag></p>
      </Modal>
    </div>
  );
};

export default Audit;
