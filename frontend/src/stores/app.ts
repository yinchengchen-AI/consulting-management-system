import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AppState {
  // 侧边栏折叠状态
  collapsed: boolean;
  // 当前主题
  theme: 'light' | 'dark';
  // 当前语言
  language: 'zh-CN' | 'en-US';
  // 标签页列表
  tabs: { key: string; label: string; path: string; closable?: boolean }[];
  // 当前激活的标签页
  activeTabKey: string;
  // 未读通知数量
  unreadNoticeCount: number;
  
  // 方法
  toggleCollapsed: () => void;
  setCollapsed: (collapsed: boolean) => void;
  setTheme: (theme: 'light' | 'dark') => void;
  setLanguage: (language: 'zh-CN' | 'en-US') => void;
  addTab: (tab: { key: string; label: string; path: string; closable?: boolean }) => void;
  removeTab: (key: string) => void;
  setActiveTabKey: (key: string) => void;
  setUnreadNoticeCount: (count: number) => void;
}

export const useAppStore = create<AppState>()(
  persist(
    (set, get) => ({
      // 初始状态
      collapsed: false,
      theme: 'light',
      language: 'zh-CN',
      tabs: [{ key: 'dashboard', label: '仪表盘', path: '/', closable: false }],
      activeTabKey: 'dashboard',
      unreadNoticeCount: 0,
      
      // 切换侧边栏折叠状态
      toggleCollapsed: () => {
        set((state) => ({ collapsed: !state.collapsed }));
      },
      
      // 设置侧边栏折叠状态
      setCollapsed: (collapsed: boolean) => {
        set({ collapsed });
      },
      
      // 设置主题
      setTheme: (theme: 'light' | 'dark') => {
        set({ theme });
      },
      
      // 设置语言
      setLanguage: (language: 'zh-CN' | 'en-US') => {
        set({ language });
      },
      
      // 添加标签页
      addTab: (tab) => {
        const { tabs } = get();
        const exists = tabs.find((t) => t.key === tab.key);
        if (!exists) {
          set({ tabs: [...tabs, tab] });
        }
        set({ activeTabKey: tab.key });
      },
      
      // 移除标签页
      removeTab: (key: string) => {
        const { tabs, activeTabKey } = get();
        const newTabs = tabs.filter((t) => t.key !== key);
        
        // 如果关闭的是当前激活的标签页，需要切换到其他标签页
        if (activeTabKey === key && newTabs.length > 0) {
          const lastTab = newTabs[newTabs.length - 1];
          set({ activeTabKey: lastTab.key });
        }
        
        set({ tabs: newTabs });
      },
      
      // 设置当前激活的标签页
      setActiveTabKey: (key: string) => {
        set({ activeTabKey: key });
      },
      
      // 设置未读通知数量
      setUnreadNoticeCount: (count: number) => {
        set({ unreadNoticeCount: count });
      },
    }),
    {
      name: 'app-storage',
      partialize: (state) => ({ 
        collapsed: state.collapsed,
        theme: state.theme,
        language: state.language,
      }),
    }
  )
);
