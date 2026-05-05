import React, { useState } from 'react';
import { Table, Button, Tag, Space, Modal, Form, Input, Select, InputNumber, Switch, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';

const levelColors: Record<string, string> = { process: 'blue', container: 'green', microvm: 'red' };
const statusColors: Record<string, string> = { Running: 'green', Stopped: 'default', Error: 'red' };

const mockData = [
  { key: '1', id: 'sandbox-a1b2', agent: 'agent-rd-001', level: 'container', status: 'Running', memory: '256MB', cpu: '1.5', created: '2026-05-05 10:00' },
  { key: '2', id: 'sandbox-c3d4', agent: 'agent-finance-003', level: 'microvm', status: 'Running', memory: '512MB', cpu: '2.0', created: '2026-05-05 09:30' },
  { key: '3', id: 'sandbox-e5f6', agent: 'agent-ops-002', level: 'process', status: 'Stopped', memory: '128MB', cpu: '0.5', created: '2026-05-05 08:00' },
];

const Sandbox: React.FC = () => {
  const [modalOpen, setModalOpen] = useState(false);
  const [form] = Form.useForm();

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>创建沙箱</Button>
        <Select defaultValue="all" style={{ width: 150 }} options={[{ value: 'all', label: '全部隔离级别' }, { value: 'process', label: '进程级' }, { value: 'container', label: '容器级' }, { value: 'microvm', label: '微虚拟机级' }]} />
        <Input.Search placeholder="搜索Agent ID" style={{ width: 200 }} />
      </Space>
      <Table dataSource={mockData} size="small"
        columns={[
          { title: 'ID', dataIndex: 'id', width: 140 },
          { title: 'Agent', dataIndex: 'agent', width: 160 },
          { title: '隔离级别', dataIndex: 'level', width: 100, render: (v: string) => <Tag color={levelColors[v]}>{v}</Tag> },
          { title: '状态', dataIndex: 'status', width: 80, render: (v: string) => <Tag color={statusColors[v]}>{v}</Tag> },
          { title: '内存', dataIndex: 'memory', width: 80 },
          { title: 'CPU', dataIndex: 'cpu', width: 80 },
          { title: '创建时间', dataIndex: 'created', width: 150 },
          { title: '操作', width: 240, render: () => <Space size="small"><Button size="small">执行</Button><Button size="small">快照</Button><Button size="small" danger>停止</Button><Button size="small" danger type="primary">删除</Button></Space> },
        ]}
      />
      <Modal title="创建沙箱" open={modalOpen} onCancel={() => setModalOpen(false)} onOk={() => { message.success('沙箱创建成功'); setModalOpen(false); }}>
        <Form form={form} layout="vertical">
          <Form.Item name="agent_id" label="Agent ID" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="isolation_level" label="隔离级别" rules={[{ required: true }]}><Select options={[{ value: 'process', label: '进程级（轻量）' }, { value: 'container', label: '容器级（标准）' }, { value: 'microvm', label: '微虚拟机级（强隔离）' }]} /></Form.Item>
          <Form.Item name="memory_limit_mb" label="内存限制 (MB)" initialValue={512}><InputNumber min={64} max={8192} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="cpu_quota" label="CPU配额 (核)" initialValue={2}><InputNumber min={0.5} max={16} step={0.5} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="network_isolated" label="网络隔离" valuePropName="checked" initialValue={true}><Switch /></Form.Item>
          <Form.Item name="enable_snapshot" label="启用快照" valuePropName="checked" initialValue={true}><Switch /></Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Sandbox;
