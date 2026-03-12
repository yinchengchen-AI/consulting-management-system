import React from 'react';
import { Card, Button, Space, Tag, Popconfirm, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { ProTable } from '@ant-design/pro-components';
import { getServiceTypeList, deleteServiceType } from '@/api/service';
import { formatMoney, formatDateTime } from '@/utils';
import type { ServiceType } from '@/types';

const ServiceTypeList: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const deleteMutation = useMutation({
    mutationFn: deleteServiceType,
    onSuccess: () => {
      message.success('删除成功');
      queryClient.invalidateQueries({ queryKey: ['serviceTypes'] });
    },
  });

  const columns = [
    {
      title: '服务名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '服务编码',
      dataIndex: 'code',
      key: 'code',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      render: (text: string) => text || '-',
      ellipsis: true,
    },
    {
      title: '价格',
      dataIndex: 'price',
      key: 'price',
      render: (price: number, record: ServiceType) => 
        price ? formatMoney(price) + (record.unit ? `/${record.unit}` : '') : '-',
    },
    {
      title: '时长(天)',
      dataIndex: 'duration',
      key: 'duration',
      render: (text: number) => text || '-',
    },
    {
      title: '排序',
      dataIndex: 'sort',
      key: 'sort',
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
      render: (_: unknown, record: ServiceType) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => navigate(`/service-types/edit/${record.id}`)}
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
      <ProTable<ServiceType>
        headerTitle="服务类型"
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button
            key="add"
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => navigate('/service-types/create')}
          >
            新建服务类型
          </Button>,
        ]}
        request={async (params) => {
          const data = await getServiceTypeList({
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
      />
    </Card>
  );
};

export default ServiceTypeList;
