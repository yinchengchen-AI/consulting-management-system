import React, { useEffect } from 'react';
import { Layout, Button, Badge, Dropdown, Avatar, Space, Tooltip, theme } from 'antd';
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  BellOutlined,
  UserOutlined,
  LogoutOutlined,
  SettingOutlined,
  FullscreenOutlined,
  FullscreenExitOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useAppStore } from '@/stores/app';
import { useAuthStore } from '@/stores/auth';
import { getUnreadNoticeCount } from '@/api/notice';

const { Header: AntHeader } = Layout;

const Header: React.FC = () => {
  const navigate = useNavigate();
  const { token } = theme.useToken();
  const { collapsed, toggleCollapsed, unreadNoticeCount, setUnreadNoticeCount } = useAppStore();
  const { userInfo, logout } = useAuthStore();
  const [isFullscreen, setIsFullscreen] = React.useState(false);

  // 获取未读通知数量
  useEffect(() => {
    const fetchUnreadCount = async () => {
      try {
        const count = await getUnreadNoticeCount();
        setUnreadNoticeCount(count);
      } catch {
        // 忽略错误
      }
    };
    fetchUnreadCount();
    
    // 每5分钟刷新一次
    const timer = setInterval(fetchUnreadCount, 5 * 60 * 1000);
    return () => clearInterval(timer);
  }, [setUnreadNoticeCount]);

  // 处理全屏切换
  const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen().then(() => {
        setIsFullscreen(true);
      });
    } else {
      document.exitFullscreen().then(() => {
        setIsFullscreen(false);
      });
    }
  };

  // 处理登出
  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  // 用户下拉菜单
  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
      onClick: () => navigate('/profile'),
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '账号设置',
      onClick: () => navigate('/settings'),
    },
    {
      type: 'divider' as const,
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  return (
    <AntHeader
      style={{
        padding: '0 24px',
        background: token.colorBgContainer,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        boxShadow: '0 1px 4px rgba(0,0,0,0.1)',
        zIndex: 9,
      }}
    >
      <Space>
        <Button
          type="text"
          icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          onClick={toggleCollapsed}
          style={{ fontSize: 16 }}
        />
      </Space>

      <Space size={16}>
        {/* 全屏按钮 */}
        <Tooltip title={isFullscreen ? '退出全屏' : '全屏'}>
          <Button
            type="text"
            icon={isFullscreen ? <FullscreenExitOutlined /> : <FullscreenOutlined />}
            onClick={toggleFullscreen}
          />
        </Tooltip>

        {/* 通知按钮 */}
        <Tooltip title="通知公告">
          <Badge count={unreadNoticeCount} size="small">
            <Button
              type="text"
              icon={<BellOutlined />}
              onClick={() => navigate('/notices')}
            />
          </Badge>
        </Tooltip>

        {/* 用户下拉菜单 */}
        <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
          <Space style={{ cursor: 'pointer' }}>
            <Avatar
              src={userInfo?.avatar}
              icon={!userInfo?.avatar && <UserOutlined />}
              size="small"
            />
            <span>{userInfo?.realName || userInfo?.username || '用户'}</span>
          </Space>
        </Dropdown>
      </Space>
    </AntHeader>
  );
};

export default Header;
