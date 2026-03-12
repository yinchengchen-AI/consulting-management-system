import React from 'react';
import { Card, Button, Space, Tag, Popconfirm, message, Tree } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { ProTable } from '@ant-design/pro-components';
import { getRoleList, deleteRole, getPermissionTree } from '@/api/user';
import { formatDateTime } from '@/utils';
import type { Role } from '@/types';

const RoleList: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // 获取权限树
  const { data: permissionTree } = useQuery({
    queryKey: ['permissionTree'],
    queryFn: getPermissionTree,
  });

  // 删除角色
  const deleteMutation = useMutation({
    mutationFn: deleteRole,
    onSuccess: () => {
      message.success('删除成功');
      queryClient.invalidateQueries({ queryKey: ['roles'] });
    },
  });

  const columns = [
    {
      title: '角色名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '角色编码',
      dataIndex: 'code',
      key: 'code',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      render: (text: string) => text || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) => (
        <Tag color={status === 1 ? 'success' : 'default'}>
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      ),
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
      render: (_: unknown, record: Role) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => navigate(`/roles/edit/${record.id}`)}
          >
            编辑
          </Button>
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
      <ProTable<Role>
        headerTitle="角色列表"
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button
            key="add"
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => navigate('/roles/create')}
          >
            新建角色
          </Button>,
        ]}
        request={async (params) => {
          const data = await getRoleList({
            page: params.current,
            size: params.pageSize,
            keyword: params.name,
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
        expandable={{
          expandedRowRender: (record) => (
            <div style={{ padding: '0 24px' }}>
              <h4>权限列表</h4>
              {record.permissions && record.permissions.length > 0 ? (
                <Space wrap>
                  {record.permissions.map((perm) => (
                    <Tag key={perm.id} color="blue">
                      {perm.name}
                    </Tag>
                  ))}
                </Space>
              ) : (
                <span>暂无权限</span>
              )}
            </div>
          ),
        }}
      />
    </Card>
  );
};

export default RoleList;
