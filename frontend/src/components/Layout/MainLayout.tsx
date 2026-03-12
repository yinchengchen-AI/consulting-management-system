import React, { Suspense } from 'react';
import { Layout, Spin, Breadcrumb } from 'antd';
import { Outlet, useLocation, Link } from 'react-router-dom';
import Sidebar from './Sidebar';
import Header from './Header';

const { Content } = Layout;

// 路由映射表，用于生成面包屑
const routeMap: Record<string, { name: string; parent?: string }> = {
  '/': { name: '仪表盘' },
  '/users': { name: '用户管理', parent: '/user' },
  '/users/create': { name: '创建用户', parent: '/users' },
  '/users/edit': { name: '编辑用户', parent: '/users' },
  '/roles': { name: '角色管理', parent: '/user' },
  '/roles/create': { name: '创建角色', parent: '/roles' },
  '/roles/edit': { name: '编辑角色', parent: '/roles' },
  '/customers': { name: '客户列表', parent: '/customer' },
  '/customers/create': { name: '创建客户', parent: '/customers' },
  '/customers/edit': { name: '编辑客户', parent: '/customers' },
  '/customers/detail': { name: '客户详情', parent: '/customers' },
  '/service-types': { name: '服务类型', parent: '/service' },
  '/service-types/create': { name: '创建服务类型', parent: '/service-types' },
  '/service-types/edit': { name: '编辑服务类型', parent: '/service-types' },
  '/service-orders': { name: '服务订单', parent: '/service' },
  '/service-orders/create': { name: '创建服务订单', parent: '/service-orders' },
  '/service-orders/edit': { name: '编辑服务订单', parent: '/service-orders' },
  '/service-orders/detail': { name: '服务订单详情', parent: '/service-orders' },
  '/invoices': { name: '开票管理', parent: '/finance' },
  '/invoices/create': { name: '创建开票', parent: '/invoices' },
  '/invoices/edit': { name: '编辑开票', parent: '/invoices' },
  '/payment-plans': { name: '收款计划', parent: '/finance' },
  '/receipts': { name: '收款记录', parent: '/finance' },
  '/refunds': { name: '退款管理', parent: '/finance' },
  '/contracts': { name: '合同列表', parent: '/contract' },
  '/contracts/create': { name: '创建合同', parent: '/contracts' },
  '/contracts/edit': { name: '编辑合同', parent: '/contracts' },
  '/contracts/detail': { name: '合同详情', parent: '/contracts' },
  '/documents': { name: '文档列表', parent: '/document' },
  '/statistics/customers': { name: '客户分析', parent: '/statistics' },
  '/statistics/services': { name: '服务分析', parent: '/statistics' },
  '/statistics/finance': { name: '财务分析', parent: '/statistics' },
  '/statistics/performance': { name: '绩效分析', parent: '/statistics' },
  '/notices': { name: '通知列表', parent: '/notice' },
  '/notices/detail': { name: '通知详情', parent: '/notices' },
  '/settings': { name: '系统配置', parent: '/setting' },
  '/logs': { name: '操作日志', parent: '/setting' },
};

const MainLayout: React.FC = () => {
  const location = useLocation();

  // 生成面包屑
  const generateBreadcrumbs = () => {
    const breadcrumbs: { title: React.ReactNode }[] = [];
    const path = location.pathname;
    
    // 处理动态路由参数
    const basePath = path.replace(/\/\d+$/, '').replace(/\/edit\/\d+$/, '/edit');
    
    const buildBreadcrumb = (currentPath: string) => {
      const route = routeMap[currentPath];
      if (route) {
        if (route.parent) {
          buildBreadcrumb(route.parent);
        }
        
        // 如果是当前页面，不添加链接
        if (currentPath === basePath || currentPath === path) {
          breadcrumbs.push({ title: route.name });
        } else {
          breadcrumbs.push({
            title: <Link to={currentPath}>{route.name}</Link>,
          });
        }
      }
    };

    buildBreadcrumb(basePath);
    return breadcrumbs;
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sidebar />
      <Layout>
        <Header />
        <Content style={{ margin: '16px', overflow: 'initial' }}>
          <Breadcrumb
            items={generateBreadcrumbs()}
            style={{ marginBottom: 16 }}
          />
          <div
            style={{
              padding: 24,
              background: '#fff',
              borderRadius: 4,
              minHeight: 'calc(100vh - 180px)',
            }}
          >
            <Suspense
              fallback={
                <div style={{ textAlign: 'center', padding: '50px 0' }}>
                  <Spin size="large" />
                </div>
              }
            >
              <Outlet />
            </Suspense>
          </div>
        </Content>
      </Layout>
    </Layout>
  );
};

export default MainLayout;
