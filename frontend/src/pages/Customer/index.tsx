import React, { useState } from 'react';
import { Card, Button, Space, Tag, Popconfirm, message, Drawer, Tabs } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined, PhoneOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { ProTable, ProDescriptions } from '@ant-design/pro-components';
import { getCustomerList, deleteCustomer, getCustomerTags, getIndustries } from '@/api/customer';
import { formatDateTime } from '@/utils';
import type { Customer } from '@/types';

const CustomerList: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [detailVisible, setDetailVisible] = useState(false);
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(null);

  // 获取客户标签
  const { data: tags } = useQuery({
    queryKey: ['customerTags'],
    queryFn: getCustomerTags,
  });

  // 获取行业列表
  const { data: industries } = useQuery({
    queryKey: ['industries'],
    queryFn: getIndustries,
  });

  // 删除客户
  const deleteMutation = useMutation({
    mutationFn: deleteCustomer,
    onSuccess: () => {
      message.success('删除成功');
      queryClient.invalidateQueries({ queryKey: ['customers'] });
    },
  });

  const columns = [
    {
      title: '客户名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '简称',
      dataIndex: 'shortName',
      key: 'shortName',
      render: (text: string) => text || '-',
    },
    {
      title: '行业',
      dataIndex: 'industry',
      key: 'industry',
      render: (text: string) => text || '-',
    },
    {
      title: '联系人',
      dataIndex: 'contactName',
      key: 'contactName',
      render: (text: string) => text || '-',
    },
    {
      title: '联系电话',
      dataIndex: 'contactPhone',
      key: 'contactPhone',
      render: (text: string) => text || '-',
    },
    {
      title: '标签',
      dataIndex: 'tags',
      key: 'tags',
      render: (tags: string[]) => (
        <Space wrap>
          {tags?.map((tag) => (
            <Tag key={tag} color="blue">
              {tag}
            </Tag>
          ))}
        </Space>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) => (
        <Tag color={status === 1 ? 'success' : 'default'}>
          {status === 1 ? '合作中' : '已终止'}
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
      render: (_: unknown, record: Customer) => (
        <Space>
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => {
              setSelectedCustomer(record);
              setDetailVisible(true);
            }}
          >
            详情
          </Button>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => navigate(`/customers/edit/${record.id}`)}
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
    <>
      <Card>
        <ProTable<Customer>
          headerTitle="客户列表"
          rowKey="id"
          search={{
            labelWidth: 120,
          }}
          toolBarRender={() => [
            <Button
              key="add"
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => navigate('/customers/create')}
            >
              新建客户
            </Button>,
          ]}
          request={async (params) => {
            const data = await getCustomerList({
              page: params.current,
              size: params.pageSize,
              keyword: params.name,
              status: params.status,
              industry: params.industry,
              tag: params.tag,
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

      {/* 客户详情抽屉 */}
      <Drawer
        title="客户详情"
        width={600}
        open={detailVisible}
        onClose={() => setDetailVisible(false)}
      >
        {selectedCustomer && (
          <Tabs
            items={[
              {
                key: 'basic',
                label: '基本信息',
                children: (
                  <ProDescriptions column={1}>
                    <ProDescriptions.Item label="客户名称">{selectedCustomer.name}</ProDescriptions.Item>
                    <ProDescriptions.Item label="简称">{selectedCustomer.shortName || '-'}</ProDescriptions.Item>
                    <ProDescriptions.Item label="行业">{selectedCustomer.industry || '-'}</ProDescriptions.Item>
                    <ProDescriptions.Item label="规模">{selectedCustomer.scale || '-'}</ProDescriptions.Item>
                    <ProDescriptions.Item label="网站">{selectedCustomer.website || '-'}</ProDescriptions.Item>
                    <ProDescriptions.Item label="地址">{selectedCustomer.address || '-'}</ProDescriptions.Item>
                    <ProDescriptions.Item label="标签">
                      <Space wrap>
                        {selectedCustomer.tags?.map((tag) => (
                          <Tag key={tag} color="blue">{tag}</Tag>
                        ))}
                      </Space>
                    </ProDescriptions.Item>
                    <ProDescriptions.Item label="备注">{selectedCustomer.remark || '-'}</ProDescriptions.Item>
                  </ProDescriptions>
                ),
              },
              {
                key: 'contact',
                label: '联系信息',
                children: (
                  <ProDescriptions column={1}>
                    <ProDescriptions.Item label="联系人">{selectedCustomer.contactName || '-'}</ProDescriptions.Item>
                    <ProDescriptions.Item label="联系电话">{selectedCustomer.contactPhone || '-'}</ProDescriptions.Item>
                    <ProDescriptions.Item label="联系邮箱">{selectedCustomer.contactEmail || '-'}</ProDescriptions.Item>
                  </ProDescriptions>
                ),
              },
            ]}
          />
        )}
      </Drawer>
    </>
  );
};

export default CustomerList;
