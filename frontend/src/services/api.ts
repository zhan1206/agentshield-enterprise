import axios from 'axios';

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/';
    }
    return Promise.reject(err);
  }
);

export const getSandboxList = () => api.get('/sandboxes');
export const createSandbox = (data: any) => api.post('/sandboxes', data);
export const getSandboxDetail = (id: string) => api.get(`/sandboxes/${id}`);
export const deleteSandbox = (id: string) => api.delete(`/sandboxes/${id}`);
export const executeInSandbox = (id: string, data: any) => api.post(`/sandboxes/${id}/execute`, data);
export const getAgentList = () => api.get('/agents');
export const createAgent = (data: any) => api.post('/agents', data);
export const getAgentDetail = (id: string) => api.get(`/agents/${id}`);
export const getPermissionList = () => api.get('/permissions');
export const getAuditLogs = (params?: any) => api.get('/audit/logs', { params });
export const getThreatAlerts = () => api.get('/threats/alerts');
export const getMetrics = () => api.get('/metrics');
export const getDataPolicies = () => api.get('/data-security/policies');

export default api;
