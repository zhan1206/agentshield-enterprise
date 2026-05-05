import { create } from 'zustand';

interface AppState {
  sandboxList: any[];
  agentList: any[];
  threatAlerts: any[];
  metrics: any;
  fetchSandboxList: () => Promise<void>;
  fetchAgentList: () => Promise<void>;
  setMetrics: (m: any) => void;
  addThreatAlert: (alert: any) => void;
}

const useStore = create<AppState>((set) => ({
  sandboxList: [],
  agentList: [],
  threatAlerts: [],
  metrics: { active_sandboxes: 0, running_agents: 0, security_events: 0, blocked_operations: 0, threat_alerts: 0 },

  fetchSandboxList: async () => {
    try {
      const { data } = await (await import('./services/api')).default.get('/sandboxes');
      set({ sandboxList: data.sandboxes || [] });
    } catch { /* ignore */ }
  },

  fetchAgentList: async () => {
    try {
      const { data } = await (await import('./services/api')).default.get('/agents');
      set({ agentList: data.agents || [] });
    } catch { /* ignore */ }
  },

  setMetrics: (metrics) => set({ metrics }),
  addThreatAlert: (alert) => set((s) => ({ threatAlerts: [...s.threatAlerts, alert] })),
}));

export default useStore;
