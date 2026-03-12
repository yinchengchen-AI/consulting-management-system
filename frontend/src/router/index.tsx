import React, { Suspense } from 'react';
import { createBrowserRouter, Navigate, Outlet } from 'react-router-dom';
import { Spin } from 'antd';
import { useAuthStore } from '@/stores/auth';
import MainLayout from '@/components/Layout/MainLayout';

// 懒加载页面组件
const Login = React.lazy(() => import('@/pages/Login'));
const Dashboard = React.lazy(() => import('@/pages/Dashboard'));

// 用户管理
const UserList = React.lazy(() => import('@/pages/User/UserList'));
const UserForm = React.lazy(() => import('@/pages/User/UserForm'));
const RoleList = React.lazy(() => import('@/pages/User/RoleList'));
const RoleForm = React.lazy(() => import('@/pages/User/RoleForm'));

// 客户管理
const CustomerList = React.lazy(() => import('@/pages/Customer'));
const CustomerDetail = React.lazy(() => import('@/pages/Customer/Detail'));

// 服务类型管理
const ServiceTypeList = React.lazy(() => import('@/pages/Service/ServiceTypeList'));

// 服务订单管理
const ServiceOrderList = React.lazy(() => import('@/pages/ServiceOrder'));
const ServiceOrderDetail = React.lazy(() => import('@/pages/ServiceOrder/Detail'));
const ServiceOrderForm = React.lazy(() => import('@/pages/ServiceOrder/Form'));

// 开票管理
const InvoiceList = React.lazy(() => import('@/pages/Invoice'));

// 收款管理
const ReceiptList = React.lazy(() => import('@/pages/Receipt'));

// 合同管理
const ContractList = React.lazy(() => import('@/pages/Contract/List'));
const ContractDetail = React.lazy(() => import('@/pages/Contract/Detail'));

// 统计分析
const StatisticsDashboard = React.lazy(() => import('@/pages/Statistics'));

// 404页面
const NotFound = React.lazy(() => import('@/pages/NotFound'));

// 路由守卫组件
const PrivateRoute: React.FC = () => {
  const { isAuthenticated } = useAuthStore();
  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />;
};

// 公共路由守卫
const PublicRoute: React.FC = () => {
  const { isAuthenticated } = useAuthStore();
  return isAuthenticated ? <Navigate to="/" replace /> : <Outlet />;
};

// 加载中组件
const PageLoading: React.FC = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
    <Spin size="large" />
  </div>
);

// 路由配置
export const router = createBrowserRouter([
  {
    path: '/login',
    element: (
      <Suspense fallback={<PageLoading />}>
        <PublicRoute />
      </Suspense>
    ),
    children: [
      {
        index: true,
        element: <Login />,
      },
    ],
  },
  {
    path: '/',
    element: (
      <Suspense fallback={<PageLoading />}>
        <PrivateRoute />
      </Suspense>
    ),
    children: [
      {
        element: <MainLayout />,
        children: [
          {
            index: true,
            element: <Dashboard />,
          },
          // 用户管理
          {
            path: 'users',
            element: <UserList />,
          },
          {
            path: 'users/create',
            element: <UserForm />,
          },
          {
            path: 'users/edit/:id',
            element: <UserForm />,
          },
          // 角色管理
          {
            path: 'roles',
            element: <RoleList />,
          },
          {
            path: 'roles/create',
            element: <RoleForm />,
          },
          {
            path: 'roles/edit/:id',
            element: <RoleForm />,
          },
          // 客户管理
          {
            path: 'customers',
            element: <CustomerList />,
          },
          {
            path: 'customers/:id',
            element: <CustomerDetail />,
          },
          // 服务类型管理
          {
            path: 'service-types',
            element: <ServiceTypeList />,
          },
          // 服务订单管理
          {
            path: 'service-orders',
            element: <ServiceOrderList />,
          },
          {
            path: 'service-orders/create',
            element: <ServiceOrderForm />,
          },
          {
            path: 'service-orders/edit/:id',
            element: <ServiceOrderForm />,
          },
          {
            path: 'service-orders/:id',
            element: <ServiceOrderDetail />,
          },
          // 开票管理
          {
            path: 'invoices',
            element: <InvoiceList />,
          },
          // 收款管理
          {
            path: 'receipts',
            element: <ReceiptList />,
          },
          {
            path: 'payment-plans',
            element: <ReceiptList />,
          },
          // 合同管理
          {
            path: 'contracts',
            element: <ContractList />,
          },
          {
            path: 'contracts/:id',
            element: <ContractDetail />,
          },
          // 统计分析
          {
            path: 'statistics',
            element: <StatisticsDashboard />,
          },
          // 404
          {
            path: '*',
            element: <NotFound />,
          },
        ],
      },
    ],
  },
]);

export default router;
