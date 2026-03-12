import React, { useState } from 'react';
import { Card, Button, Space, Tag, Popconfirm, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, KeyOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { ProTable } from '@ant-design/pro-components';
import { getUserList, deleteUser, resetPassword, toggleUserStatus } from '@/api/user';
import { formatDateTime } from '@/utils';
import type { User } from '@/types';

const UserList: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // 删除用户
  const deleteMutation = useMutation({
    mutationFn: deleteUser,
    onSuccess: () => {
      message.success('删除成功');
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });

  // 重置密码
  const resetPasswordMutation = useMutation({
    mutationFn: resetPassword,
    onSuccess: (data) => {
      message.success(`密码已重置，新密码: ${data.password}`);
    },
  });

  // 切换状态
  const toggleStatusMutation = useMutation({
    mutationFn: ({ id, status }: { id: number; status: number }) => toggleUserStatus(id, status),
    onSuccess: () => {
      message.success('状态更新成功');
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });

  const columns = [
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '姓名',
      dataIndex: 'realName',
      key: 'realName',
    },
    {
      title: '部门',
      dataIndex: 'departmentName',
      key: 'departmentName',
      render: (text: string) => text || '-',
    },
    {
      title: '职位',
      dataIndex: 'position',
      key: 'position',
      render: (text: string) => text || '-',
    },
    {
      title: '角色',
      dataIndex: 'roles',
      key: 'roles',
      render: (roles: { name: string }[]) => (
        <Space>
          {roles?.map((role) => (
            <Tag key={role.name} color="blue">
              {role.name}
            </Tag>
          ))}
        </Space>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number, record: User) => (
        <Tag
          color={status === 1 ? 'success' : 'default'}
          style={{ cursor: 'pointer' }}
          onClick={() => toggleStatusMutation.mutate({ id: record.id, status: status === 1 ? 0 : 1 })}
        >
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '最后登录',
      dataIndex: 'lastLoginTime',
      key: 'lastLoginTime',
      render: (text: string) => formatDateTime(text),
    },
    {
      title: '创建时间',
      dataIndex: 'createTime',
      key: 'createTime',
      render: (text: string) => formatDateTime(text),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: User) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => navigate(`/users/edit/${record.id}`)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定重置密码吗？"
            onConfirm={() => resetPasswordMutation.mutate(record.id)}
          >
            <Button type="link" icon={<KeyOutlined />}>
              重置密码
            </Button>
          </Popconfirm>
          <Popconfirm
            title="确定删除吗？"
            onConfirm={() => deleteMutation.mutate(record.id)}
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card>
      <ProTable<User>
        headerTitle="用户列表"
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button
            key="add"
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => navigate('/users/create')}
          >
            新建用户
          </Button>,
        ]}
        request={async (params) => {
          const data = await getUserList({
            page: params.current,
            size: params.pageSize,
            keyword: params.username || params.realName,
            status: params.status,
          });
          return {
            data: data.list,
            total: data.total,
            success: true,
          };
        }}
        columns={columns}
        pagination={{
          showSizeChanger: true,
          showQuickJumper: true,
        }}
      />
    </Card>
  );
};

export default UserList;
