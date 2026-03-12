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
  Popconfirm,
  message,
} from 'antd';
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { getCustomerList, deleteCustomer } from '@/api/customer';
import { Customer, CustomerType, CustomerLevel, CustomerStatus } from '@/types';

const { Option } = Select;

const typeMap: Record<CustomerType, string> = {
  enterprise: '企业',
  government: '政府',
  individual: '个人',
  other: '其他',
};

const levelMap: Record<CustomerLevel, { text: string; color: string }> = {
  A: { text: 'A级', color: 'red' },
  B: { text: 'B级', color: 'orange' },
  C: { text: 'C级', color: 'blue' },
  D: { text: 'D级', color: 'default' },
};

const statusMap: Record<CustomerStatus, { text: string; color: string }> = {
  active: { text: '正常', color: 'success' },
  inactive: { text: '停用', color: 'default' },
  potential: { text: '潜在', color: 'warning' },
};

const CustomerList = () => {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useState({
    page: 1,
    page_size: 10,
    keyword: '',
    type: '',
    level: '',
    status: '',
  });

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['customers', searchParams],
    queryFn: () => getCustomerList(searchParams),
  });

  const handleDelete = async (id: string) => {
    try {
      await deleteCustomer(id);
      message.success('删除成功');
      refetch();
    } catch (error) {
      // 错误已在 request 中处理
    }
  };

  const columns = [
    {
      title: '客户名称',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: Customer) => (
        <div>
          <div style={{ fontWeight: 500 }}>{name}</div>
          <div style={{ fontSize: 12, color: '#999' }}>{record.industry}</div>
        </div>
      ),
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: CustomerType) => typeMap[type] || type,
    },
    {
      title: '等级',
      dataIndex: 'level',
      key: 'level',
      render: (level: CustomerLevel) => (
        <Tag color={levelMap[level]?.color}>{levelMap[level]?.text}</Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: CustomerStatus) => (
        <Tag color={statusMap[status]?.color}>{statusMap[status]?.text}</Tag>
      ),
    },
    {
      title: '联系人',
      dataIndex: 'contact_name',
      key: 'contact_name',
      render: (name: string, record: Customer) => (
        <div>
          <div>{name || '-'}</div>
          {record.contact_phone && (
            <div style={{ fontSize: 12, color: '#999' }}>{record.contact_phone}</div>
          )}
        </div>
      ),
    },
    {
      title: '合同数',
      dataIndex: 'contract_count',
      key: 'contract_count',
      render: (count: number) => count || 0,
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Customer) => (
        <Space>
          <Button
            type="text"
            icon={<EyeOutlined />}
            size="small"
            onClick={() => navigate(`/customers/${record.id}`)}
          >
            详情
          </Button>
          <Button type="text" icon={<EditOutlined />} size="small">
            编辑
          </Button>
          <Popconfirm
            title="确认删除"
            description="确定要删除该客户吗？"
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
    <div>
      <Card>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 24 }}>
          <h2 style={{ margin: 0 }}>客户管理</h2>
          <Button type="primary" icon={<PlusOutlined />}>
            新增客户
          </Button>
        </div>

        <div style={{ marginBottom: 16, padding: 16, background: '#fafafa', borderRadius: 8 }}>
          <Space>
            <Input
              placeholder="搜索客户名称/联系人"
              prefix={<SearchOutlined />}
              value={searchParams.keyword}
              onChange={(e) =>
                setSearchParams({ ...searchParams, keyword: e.target.value, page: 1 })
              }
              allowClear
              style={{ width: 250 }}
            />
            <Select
              placeholder="类型"
              value={searchParams.type || undefined}
              onChange={(value) =>
                setSearchParams({ ...searchParams, type: value, page: 1 })
              }
              allowClear
              style={{ width: 120 }}
            >
              <Option value="enterprise">企业</Option>
              <Option value="government">政府</Option>
              <Option value="individual">个人</Option>
              <Option value="other">其他</Option>
            </Select>
            <Select
              placeholder="等级"
              value={searchParams.level || undefined}
              onChange={(value) =>
                setSearchParams({ ...searchParams, level: value, page: 1 })
              }
              allowClear
              style={{ width: 120 }}
            >
              <Option value="A">A级</Option>
              <Option value="B">B级</Option>
              <Option value="C">C级</Option>
              <Option value="D">D级</Option>
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
              <Option value="inactive">停用</Option>
              <Option value="potential">潜在</Option>
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

export default CustomerList;
