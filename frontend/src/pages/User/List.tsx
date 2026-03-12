import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  Card,
  Table,
  Button,
  Input,
  Select,
  Space,
  Tag,
  Avatar,
  Popconfirm,
  message,
} from 'antd';
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { getUserList, deleteUser } from '@/api/user';
import { User, UserRole, UserStatus } from '@/types';
import './index.css';

const { Option } = Select;

const roleMap: Record<UserRole, { text: string; color: string }> = {
  admin: { text: '管理员', color: 'red' },
  manager: { text: '经理', color: 'blue' },
  consultant: { text: '顾问', color: 'green' },
  finance: { text: '财务', color: 'orange' },
  viewer: { text: '访客', color: 'default' },
};

const statusMap: Record<UserStatus, { text: string; color: string }> = {
  active: { text: '正常', color: 'success' },
  inactive: { text: '禁用', color: 'default' },
  suspended: { text: '暂停', color: 'warning' },
};

const UserList = () => {
  const [searchParams, setSearchParams] = useState({
    page: 1,
    page_size: 10,
    keyword: '',
    role: '',
    status: '',
  });

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['users', searchParams],
    queryFn: () => getUserList(searchParams),
  });

  const handleDelete = async (id: string) => {
    try {
      await deleteUser(id);
      message.success('删除成功');
      refetch();
    } catch (error) {
      // 错误已在 request 中处理
    }
  };

  const columns = [
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
      render: (_: string, record: User) => (
        <Space>
          <Avatar icon={<UserOutlined />} src={record.avatar} />
          <div>
            <div>{record.real_name || record.username}</div>
            <div style={{ fontSize: 12, color: '#999' }}>{record.email}</div>
          </div>
        </Space>
      ),
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (role: UserRole) => (
        <Tag color={roleMap[role]?.color}>{roleMap[role]?.text || role}</Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: UserStatus) => (
        <Tag color={statusMap[status]?.color}>{statusMap[status]?.text || status}</Tag>
      ),
    },
    {
      title: '电话',
      dataIndex: 'phone',
      key: 'phone',
    },
    {
      title: '最后登录',
      dataIndex: 'last_login',
      key: 'last_login',
      render: (date: string) => date || '-',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: User) => (
        <Space>
          <Button type="text" icon={<EditOutlined />} size="small">
            编辑
          </Button>
          <Popconfirm
            title="确认删除"
            description="确定要删除该用户吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="text" danger icon={<DeleteOutlined />} size="small">
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const list = data?.data?.list || [];
  const total = data?.data?.total || 0;

  return (
    <div className="user-list-page">
      <Card>
        <div className="page-header">
          <h2>用户管理</h2>
          <Button type="primary" icon={<PlusOutlined />}>
            新增用户
          </Button>
        </div>

        <div className="search-form">
          <Space>
            <Input
              placeholder="搜索用户名/邮箱/姓名"
              prefix={<SearchOutlined />}
              value={searchParams.keyword}
              onChange={(e) =>
                setSearchParams({ ...searchParams, keyword: e.target.value, page: 1 })
              }
              allowClear
              style={{ width: 250 }}
            />
            <Select
              placeholder="角色"
              value={searchParams.role || undefined}
              onChange={(value) =>
                setSearchParams({ ...searchParams, role: value, page: 1 })
              }
              allowClear
              style={{ width: 120 }}
            >
              <Option value="admin">管理员</Option>
              <Option value="manager">经理</Option>
              <Option value="consultant">顾问</Option>
              <Option value="finance">财务</Option>
              <Option value="viewer">访客</Option>
            </Select>
            <Select
              placeholder="状态"
              value={searchParams.status || undefined}
              onChange={(value) =>
                setSearchParams({ ...searchParams, status: value, page: 1 })
              }
              allowClear
              style={{ width: 120 }}
            >
              <Option value="active">正常</Option>
              <Option value="inactive">禁用</Option>
              <Option value="suspended">暂停</Option>
            </Select>
          </Space>
        </div>

        <Table
          columns={columns}
          dataSource={list}
          rowKey="id"
          loading={isLoading}
          pagination={{
            current: searchParams.page,
            pageSize: searchParams.page_size,
            total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条`,
            onChange: (page, pageSize) => {
              setSearchParams({ ...searchParams, page, page_size: pageSize || 10 });
            },
          }}
        />
      </Card>
    </div>
  );
};

export default UserList;
