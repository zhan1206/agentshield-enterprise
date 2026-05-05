import React, { useState } from 'react';
import { Table, Tabs, Button, Tag, Modal, Form, Input, Select, InputNumber, Switch, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';

const effectColors: Record<string, string> = { allow: 'green', deny: 'red' };

const mockRules = [
  { key: '1', id: 'perm-001', name: '研发Agent代码库只读', effect: 'allow', level: 'resource', actions: 'read', priority: 10, enabled: true },
  { key: '2', id: 'perm-002', name: '禁止访问生产数据库', effect: 'deny', level: 'environment', actions: '*', priority: 100, enabled: true },
  { key: '3', id: 'perm-003', name: '金融Agent数据导出限制', effect: 'deny', level: 'data', actions: 'export,send', priority: 50, enabled: true },
];

const mockRoles = [
  { key: '1', id: 'role-dev', name: '开发者Agent', desc: '研发场景Agent，代码库读写权限', perms: 8 },
  { key: '2', id: 'role-finance', name: '金融Agent', desc: '金融场景Agent，严格数据访问限制', perms: 5 },
  { key: '3', id: 'role-ops', name: '运维Agent', desc: '运维场景Agent，基础设施操作权限', perms: 12 },
];

const Permissions: React.FC = () => {
  const [modalOpen, setModalOpen] = useState(false);
  const [form] = Form.useForm();

  return (
    <Tabs defaultActiveKey="rules" items={[
      {
        key: 'rules',
        label: '权限规则',
        children: (
          <>
            <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)} style={{ marginBottom: 16 }}>创建权限规则</Button>
            <Table dataSource={mockRules} size="small"
              columns={[
                { title: 'ID', dataIndex: 'id', width: 110 },
                { title: '名称', dataIndex: 'name', width: 200 },
                { title: '效果', dataIndex: 'effect', width: 80, render: (v: string) => <Tag color={effectColors[v]}>{v}</Tag> },
                { title: '级别', dataIndex: 'level', width: 100 },
                { title: '操作', dataIndex: 'actions', width: 120 },
                { title: '优先级', dataIndex: 'priority', width: 80 },
                { title: '启用', dataIndex: 'enabled', width: 80, render: (v: boolean) => <Switch checked={v} size="small" /> },
              ]}
            />
          </>
        ),
      },
      {
        key: 'roles',
        label: '角色定义',
        children: (
          <Table dataSource={mockRoles} size="small"
            columns={[
              { title: 'ID', dataIndex: 'id', width: 120 },
              { title: '名称', dataIndex: 'name', width: 150 },
              { title: '描述', dataIndex: 'desc', width: 300 },
              { title: '基础权限数', dataIndex: 'perms', width: 100 },
              { title: '操作', width: 100, render: () => <Button size="small">编辑</Button> },
            ]}
          />
        ),
      },
    ]} />
  );
};

export default Permissions;
