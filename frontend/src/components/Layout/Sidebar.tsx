import React, { useState, useEffect } from 'react';
import { Menu, Layout } from 'antd';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  DashboardOutlined,
  UserOutlined,
  TeamOutlined,
  SolutionOutlined,
  FileTextOutlined,
  DollarOutlined,
  PieChartOutlined,
  BellOutlined,
  FolderOutlined,
  FileProtectOutlined,
  SettingOutlined,
} from '@ant-design/icons';
import { useAuthStore } from '@/stores/auth';
import { useAppStore } from '@/stores/app';

const { Sider } = Layout;

interface MenuItem {
  key: string;
  icon?: React.ReactNode;
  label: string;
  path?: string;
  children?: MenuItem[];
  permissions?: string[];
}

const menuItems: MenuItem[] = [
  {
    key: 'dashboard',
    icon: <DashboardOutlined />,
    label: '仪表盘',
    path: '/',
  },
  {
    key: 'user',
    icon: <UserOutlined />,
    label: '用户权限',
    permissions: ['user:view'],
    children: [
      { key: 'users', label: '用户管理', path: '/users', permissions: ['user:view'] },
      { key: 'roles', label: '角色管理', path: '/roles', permissions: ['role:view'] },
    ],
  },
  {
    key: 'customer',
    icon: <TeamOutlined />,
    label: '客户管理',
    permissions: ['customer:view'],
    children: [
      { key: 'customers', label: '客户列表', path: '/customers', permissions: ['customer:view'] },
    ],
  },
  {
    key: 'service',
    icon: <SolutionOutlined />,
    label: '服务管理',
    permissions: ['service:view'],
    children: [
      { key: 'service-types', label: '服务类型', path: '/service-types', permissions: ['service:type:view'] },
      { key: 'service-orders', label: '服务订单', path: '/service-orders', permissions: ['service:order:view'] },
    ],
  },
  {
    key: 'finance',
    icon: <DollarOutlined />,
    label: '财务管理',
    permissions: ['finance:view'],
    children: [
      { key: 'invoices', label: '开票管理', path: '/invoices', permissions: ['invoice:view'] },
      { key: 'payment-plans', label: '收款计划', path: '/payment-plans', permissions: ['payment:view'] },
      { key: 'receipts', label: '收款记录', path: '/receipts', permissions: ['receipt:view'] },
      { key: 'refunds', label: '退款管理', path: '/refunds', permissions: ['refund:view'] },
    ],
  },
  {
    key: 'contract',
    icon: <FileProtectOutlined />,
    label: '合同管理',
    permissions: ['contract:view'],
    children: [
      { key: 'contracts', label: '合同列表', path: '/contracts', permissions: ['contract:view'] },
    ],
  },
  {
    key: 'document',
    icon: <FolderOutlined />,
    label: '文档管理',
    permissions: ['document:view'],
    children: [
      { key: 'documents', label: '文档列表', path: '/documents', permissions: ['document:view'] },
    ],
  },
  {
    key: 'statistics',
    icon: <PieChartOutlined />,
    label: '统计分析',
    permissions: ['statistics:view'],
    children: [
      { key: 'statistics-customers', label: '客户分析', path: '/statistics/customers', permissions: ['statistics:customer:view'] },
      { key: 'statistics-services', label: '服务分析', path: '/statistics/services', permissions: ['statistics:service:view'] },
      { key: 'statistics-finance', label: '财务分析', path: '/statistics/finance', permissions: ['statistics:finance:view'] },
      { key: 'statistics-performance', label: '绩效分析', path: '/statistics/performance', permissions: ['statistics:performance:view'] },
    ],
  },
  {
    key: 'notice',
    icon: <BellOutlined />,
    label: '通知公告',
    permissions: ['notice:view'],
    children: [
      { key: 'notices', label: '通知列表', path: '/notices', permissions: ['notice:view'] },
    ],
  },
  {
    key: 'setting',
    icon: <SettingOutlined />,
    label: '系统设置',
    permissions: ['setting:view'],
    children: [
      { key: 'settings', label: '系统配置', path: '/settings', permissions: ['setting:config:view'] },
      { key: 'logs', label: '操作日志', path: '/logs', permissions: ['setting:log:view'] },
    ],
  },
];

const Sidebar: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { collapsed } = useAppStore();
  const { hasPermission } = useAuthStore();
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
  const [openKeys, setOpenKeys] = useState<string[]>([]);

  // 根据权限过滤菜单
  const filterMenuByPermission = (items: MenuItem[]): MenuItem[] => {
    return items
      .filter((item) => {
        if (!item.permissions || item.permissions.length === 0) return true;
        return item.permissions.some((p) => hasPermission(p));
      })
      .map((item) => {
        if (item.children) {
          return {
            ...item,
            children: filterMenuByPermission(item.children),
          };
        }
        return item;
      })
      .filter((item) => !item.children || item.children.length > 0);
  };

  const filteredMenuItems = filterMenuByPermission(menuItems);

  // 将菜单项转换为Ant Design Menu组件需要的格式
  const convertMenuItems = (items: MenuItem[]) => {
    return items.map((item) => {
      const menuItem: {
        key: string;
        icon?: React.ReactNode;
        label: React.ReactNode;
        children?: unknown[];
      } = {
        key: item.key,
        icon: item.icon,
        label: item.path ? (
          <span onClick={() => navigate(item.path!)}>{item.label}</span>
        ) : (
          item.label
        ),
      };

      if (item.children && item.children.length > 0) {
        menuItem.children = convertMenuItems(item.children);
      }

      return menuItem;
    });
  };

  // 根据当前路径设置选中的菜单项
  useEffect(() => {
    const findMenuKeyByPath = (items: MenuItem[], path: string): string | null => {
      for (const item of items) {
        if (item.path === path) {
          return item.key;
        }
        if (item.children) {
          const childKey = findMenuKeyByPath(item.children, path);
          if (childKey) {
            return childKey;
          }
        }
      }
      return null;
    };

    const key = findMenuKeyByPath(menuItems, location.pathname);
    if (key) {
      setSelectedKeys([key]);
      
      // 设置展开的父菜单
      const findParentKey = (items: MenuItem[], targetKey: string, parentKey?: string): string | null => {
        for (const item of items) {
          if (item.key === targetKey) {
            return parentKey || null;
          }
          if (item.children) {
            const result = findParentKey(item.children, targetKey, item.key);
            if (result) return result;
          }
        }
        return null;
      };
      
      const parentKey = findParentKey(menuItems, key);
      if (parentKey && !collapsed) {
        setOpenKeys([parentKey]);
      }
    }
  }, [location.pathname, collapsed]);

  const handleMenuClick = ({ key }: { key: string }) => {
    const findMenuItemByKey = (items: MenuItem[], targetKey: string): MenuItem | null => {
      for (const item of items) {
        if (item.key === targetKey) {
          return item;
        }
        if (item.children) {
          const child = findMenuItemByKey(item.children, targetKey);
          if (child) return child;
        }
      }
      return null;
    };

    const menuItem = findMenuItemByKey(menuItems, key);
    if (menuItem?.path) {
      navigate(menuItem.path);
    }
  };

  return (
    <Sider
      trigger={null}
      collapsible
      collapsed={collapsed}
      theme="light"
      style={{
        boxShadow: '2px 0 8px rgba(0,0,0,0.1)',
        zIndex: 10,
      }}
    >
      <div
        style={{
          height: 64,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          borderBottom: '1px solid #f0f0f0',
        }}
      >
        {collapsed ? (
          <div style={{ fontSize: 24, fontWeight: 'bold', color: '#1890ff' }}>C</div>
        ) : (
          <div style={{ fontSize: 18, fontWeight: 'bold', color: '#1890ff' }}>
            咨询管理系统
          </div>
        )}
      </div>
      <Menu
        mode="inline"
        selectedKeys={selectedKeys}
        openKeys={openKeys}
        onOpenChange={setOpenKeys}
        items={convertMenuItems(filteredMenuItems) as unknown[]}
        onClick={handleMenuClick}
        style={{ borderRight: 0 }}
      />
    </Sider>
  );
};

export default Sidebar;
