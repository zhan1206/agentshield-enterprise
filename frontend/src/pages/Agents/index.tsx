import React, { useState } from 'react';
import { Table, Button, Tag, Modal, Form, Input, Select, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';

const secColors: Record<string, string> = { trusted: 'green', standard: 'blue', restricted: 'orange', untrusted: 'red' };

const mockAgents = [
  { key: '1', id: 'agent-rd-001', name: '研发助手Alpha', framework: 'OpenHands', security: 'standard', status: 'active', team: '研发团队', heartbeat: '2026-05-05 12:00' },
  { key: '2', id: 'agent-finance-003', name: '金融分析Agent', framework: 'LangGraph', security: 'restricted', status: 'active', team: '金融团队', heartbeat: '2026-05-05 11:55' },
  { key: '3', id: 'agent-ops-002', name: '运维巡检Agent', framework: 'CrewAI', security: 'trusted', status: 'idle', team: '运维团队', heartbeat: '2026-05-05 11:30' },
];

const Agents: React.FC = () => {
  const [modalOpen, setModalOpen] = useState(false);
  const [form] = Form.useForm();

  return (
    <div>
      <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)} style={{ marginBottom: 16 }}>注册Agent</Button>
      <Table dataSource={mockAgents} size="small"
        columns={[
          { title: 'ID', dataIndex: 'id', width: 150 },
          { title: '名称', dataIndex: 'name', width: 160 },
          { title: '框架', dataIndex: 'framework', width: 100, render: (v: string) => <Tag>{v}</Tag> },
          { title: '安全等级', dataIndex: 'security', width: 100, render: (v: string) => <Tag color={secColors[v]}>{v}</Tag> },
          { title: '状态', dataIndex: 'status', width: 80, render: (v: string) => <Tag color={v === 'active' ? 'green' : 'default'}>{v}</Tag> },
          { title: '团队', dataIndex: 'team', width: 120 },
          { title: '最近心跳', dataIndex: 'heartbeat', width: 150 },
          { title: '操作', width: 120, render: () => <Space size="small"><Button size="small">详情</Button><Button size="small" danger>注销</Button></Space> },
        ]}
      />
      <Modal title="注册Agent" open={modalOpen} onCancel={() => setModalOpen(false)} onOk={() => { message.success('Agent注册成功'); setModalOpen(false); }}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="framework" label="框架" rules={[{ required: true }]}><Select options={[{ value: 'OpenHands', label: 'OpenHands' }, { value: 'Plandex', label: 'Plandex' }, { value: 'LangGraph', label: 'LangGraph' }, { value: 'CrewAI', label: 'CrewAI' }, { value: 'AutoGen', label: 'AutoGen' }, { value: 'Custom', label: '自定义' }]} /></Form.Item>
          <Form.Item name="security_level" label="安全等级" initialValue="standard"><Select options={[{ value: 'trusted', label: '可信' }, { value: 'standard', label: '标准' }, { value: 'restricted', label: '受限' }, { value: 'untrusted', label: '不可信' }]} /></Form.Item>
          <Form.Item name="team_id" label="团队ID"><Input /></Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

// Need Space import
import { Space } from 'antd';

export default Agents;
