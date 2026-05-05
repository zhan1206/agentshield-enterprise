import React from 'react';
import { Tabs, Form, Input, Select, InputNumber, Switch, Button, Statistic, Card, Row, Col } from 'antd';

const Settings: React.FC = () => (
  <Tabs defaultActiveKey="general" items={[
    {
      key: 'general',
      label: '通用设置',
      children: (
        <Form layout="vertical" style={{ maxWidth: 600 }}>
          <Form.Item label="系统名称" initialValue="AgentShield Enterprise"><Input /></Form.Item>
          <Form.Item label="API端口" initialValue={8080}><InputNumber style={{ width: '100%' }} /></Form.Item>
          <Form.Item label="运行模式"><Select defaultValue="cluster" options={[{ value: 'standalone', label: '单机' }, { value: 'cluster', label: '集群' }]} /></Form.Item>
          <Form.Item label="默认隔离级别"><Select defaultValue="container" options={[{ value: 'process', label: '进程级' }, { value: 'container', label: '容器级' }, { value: 'microvm', label: '微虚拟机级' }]} /></Form.Item>
          <Form.Item label="审计保留天数" initialValue={180}><InputNumber style={{ width: '100%' }} /></Form.Item>
          <Form.Item label="最大沙箱数" initialValue={100}><InputNumber style={{ width: '100%' }} /></Form.Item>
          <Form.Item><Button type="primary">保存</Button></Form.Item>
        </Form>
      ),
    },
    {
      key: 'security',
      label: '安全设置',
      children: (
        <Form layout="vertical" style={{ maxWidth: 600 }}>
          <Form.Item label="JWT密钥"><Input.Password defaultValue="••••••••" /></Form.Item>
          <Form.Item label="Token过期时间（小时）" initialValue={24}><InputNumber style={{ width: '100%' }} /></Form.Item>
          <Form.Item label="加密算法"><Select defaultValue="aes-256-gcm" options={[{ value: 'aes-256-gcm', label: 'AES-256-GCM' }, { value: 'chacha20-poly1305', label: 'ChaCha20-Poly1305' }]} /></Form.Item>
          <Form.Item label="内容安全扫描" valuePropName="checked" initialValue={true}><Switch /></Form.Item>
          <Form.Item label="自动威胁响应" valuePropName="checked" initialValue={true}><Switch /></Form.Item>
          <Form.Item><Button type="primary">保存</Button></Form.Item>
        </Form>
      ),
    },
    {
      key: 'compliance',
      label: '合规设置',
      children: (
        <Form layout="vertical" style={{ maxWidth: 600 }}>
          <Form.Item label="默认合规标准"><Select defaultValue="djl_2_0" options={[{ value: 'djl_2_0', label: '等保2.0' }, { value: 'gdpr', label: 'GDPR' }, { value: 'sox', label: 'SOX' }, { value: 'data_security_law', label: '数据安全法' }]} /></Form.Item>
          <Form.Item label="自动合规检查" valuePropName="checked" initialValue={true}><Switch /></Form.Item>
          <Form.Item label="报告自动生成" valuePropName="checked" initialValue={false}><Switch /></Form.Item>
          <Form.Item><Button type="primary">保存</Button></Form.Item>
        </Form>
      ),
    },
    {
      key: 'observability',
      label: '可观测性',
      children: (
        <Form layout="vertical" style={{ maxWidth: 600 }}>
          <Form.Item label="链路追踪"><Select defaultValue="jaeger" options={[{ value: 'jaeger', label: 'Jaeger' }, { value: 'zipkin', label: 'Zipkin' }, { value: 'none', label: '关闭' }]} /></Form.Item>
          <Form.Item label="追踪端点" initialValue="http://localhost:14268/api/traces"><Input /></Form.Item>
          <Form.Item label="采样率 (%)" initialValue={10}><InputNumber min={0} max={100} style={{ width: '100%' }} /></Form.Item>
          <Form.Item label="审计日志" valuePropName="checked" initialValue={true}><Switch /></Form.Item>
          <Form.Item><Button type="primary">保存</Button></Form.Item>
        </Form>
      ),
    },
    {
      key: 'about',
      label: '关于',
      children: (
        <Row gutter={[16, 16]}>
          <Col span={8}><Card><Statistic title="版本" value="v0.1.0 MVP" /></Card></Col>
          <Col span={8}><Card><Statistic title="开源协议" value="Apache 2.0" /></Card></Col>
          <Col span={8}><Card><Statistic title="技术栈" value="Go/Rust/React/Python" /></Card></Col>
        </Row>
      ),
    },
  ]} />
);

export default Settings;
