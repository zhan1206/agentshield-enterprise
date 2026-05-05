import React, { useState } from 'react';
import { Routes, Route, Navigate, useNavigate, useLocation } from 'react-router-dom';
import { Layout, Menu } from 'antd';
import {
  DashboardOutlined, CloudServerOutlined, RobotOutlined,
  SafetyOutlined, AuditOutlined, AlertOutlined,
  DatabaseOutlined, AppstoreOutlined, SettingOutlined
} from '@ant-design/icons';
import Dashboard from './pages/Dashboard';
import Sandbox from './pages/Sandbox';
import Agents from './pages/Agents';
import Permissions from './pages/Permissions';
import Audit from './pages/Audit';
import Threats from './pages/Threats';
import DataSecurity from './pages/DataSecurity';
import Templates from './pages/Templates';
import Settings from './pages/Settings';

const { Header, Sider, Content } = Layout;

const menuItems = [
  { key: '/', icon: <DashboardOutlined />, label: '安全大盘' },
  { key: '/sandbox', icon: <CloudServerOutlined />, label: '沙箱管理' },
  { key: '/agents', icon: <RobotOutlined />, label: 'Agent管理' },
  { key: '/permissions', icon: <SafetyOutlined />, label: '权限管控' },
  { key: '/audit', icon: <AuditOutlined />, label: '审计中心' },
  { key: '/threats', icon: <AlertOutlined />, label: '威胁中心' },
  { key: '/data-security', icon: <DatabaseOutlined />, label: '数据安全' },
  { key: '/templates', icon: <AppstoreOutlined />, label: '模板市场' },
  { key: '/settings', icon: <SettingOutlined />, label: '系统设置' },
];

const App: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', padding: '0 24px' }}>
        <div style={{ color: 'white', fontSize: 18, fontWeight: 'bold', marginRight: 40 }}>
          🛡️ AgentShield Enterprise
        </div>
        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[location.pathname]}
          items={menuItems.slice(0, 4)}
          onClick={({ key }) => navigate(key)}
          style={{ flex: 1 }}
        />
      </Header>
      <Layout>
        <Sider
          collapsible
          collapsed={collapsed}
          onCollapse={setCollapsed}
          theme="light"
          width={200}
        >
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
            style={{ height: '100%', borderRight: 0 }}
          />
        </Sider>
        <Layout style={{ padding: '16px' }}>
          <Content style={{ background: '#fff', padding: 24, borderRadius: 8, minHeight: 280 }}>
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/sandbox" element={<Sandbox />} />
              <Route path="/agents" element={<Agents />} />
              <Route path="/permissions" element={<Permissions />} />
              <Route path="/audit" element={<Audit />} />
              <Route path="/threats" element={<Threats />} />
              <Route path="/data-security" element={<DataSecurity />} />
              <Route path="/templates" element={<Templates />} />
              <Route path="/settings" element={<Settings />} />
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default App;
