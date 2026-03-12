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
import { getContractList, deleteContract } from '@/api/contract';
import { Contract, ContractStatus, ContractType } from '@/types';

const { Option } = Select;

const statusMap: Record<ContractStatus, { text: string; color: string }> = {
  draft: { text: '草稿', color: 'default' },
  pending: { text: '待审批', color: 'warning' },
  active: { text: '生效中', color: 'success' },
  completed: { text: '已完成', color: 'blue' },
  cancelled: { text: '已取消', color: 'error' },
  expired: { text: '已过期', color: 'default' },
};

const typeMap: Record<ContractType, string> = {
  project: '项目合同',
  retainer: '年度框架',
  consulting: '单次咨询',
};

const ContractList = () => {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useState({
    page: 1,
    page_size: 10,
    keyword: '',
    status: '',
    type: '',
  });

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['contracts', searchParams],
    queryFn: () => getContractList(searchParams),
  });

  const handleDelete = async (id: string) => {
    try {
      await deleteContract(id);
      message.success('删除成功');
      refetch();
    } catch (error) {
      // 错误已在 request 中处理
    }
  };

  const columns = [
    {
      title: '合同编号',
      dataIndex: 'code',
      key: 'code',
      width: 150,
    },
    {
      title: '合同名称',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: Contract) => (
        <div>
          <div style={{ fontWeight: 500 }}>{name}</div>
          <div style={{ fontSize: 12, color: '#999' }}>{typeMap[record.type]}</div>
        </div>
      ),
    },
    {
      title: '客户',
      dataIndex: 'customer',
      key: 'customer',
      render: (customer: any) => customer?.name || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: ContractStatus) => (
        <Tag color={statusMap[status]?.color}>{statusMap[status]?.text}</Tag>
      ),
    },
    {
      title: '合同金额',
      dataIndex: 'total_amount',
      key: 'total_amount',
      render: (amount: number) => `¥${(amount || 0).toLocaleString()}`,
    },
    {
      title: '已收款',
      dataIndex: 'paid_amount',
      key: 'paid_amount',
      render: (amount: number) => `¥${(amount || 0).toLocaleString()}`,
    },
    {
      title: '签约日期',
      dataIndex: 'signed_date',
      key: 'signed_date',
      render: (date: string) => (date ? new Date(date).toLocaleDateString() : '-'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Contract) => (
        <Space>
          <Button
            type="text"
            icon={<EyeOutlined />}
            size="small"
            onClick={() => navigate(`/contracts/${record.id}`)}
          >
            详情
          </Button>
          <Button type="text" icon={<EditOutlined />} size="small">
            编辑
          </Button>
          <Popconfirm
            title="确认删除"
            description="确定要删除该合同吗？"
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
          <h2 style={{ margin: 0 }}>合同管理</h2>
          <Button type="primary" icon={<PlusOutlined />}>
            新增合同
          </Button>
        </div>

        <div style={{ marginBottom: 16, padding: 16, background: '#fafafa', borderRadius: 8 }}>
          <Space>
            <Input
              placeholder="搜索合同名称/编号"
              prefix={<SearchOutlined />}
              value={searchParams.keyword}
              onChange={(e) =>
                setSearchParams({ ...searchParams, keyword: e.target.value, page: 1 })
              }
              allowClear
              style={{ width: 250 }}
            />
            <Select
              placeholder="状态"
              value={searchParams.status || undefined}
              onChange={(value) =>
                setSearchParams({ ...searchParams, status: value, page: 1 })
              }
              allowClear
              style={{ width: 120 }}
            >
              <Option value="draft">草稿</Option>
              <Option value="pending">待审批</Option>
              <Option value="active">生效中</Option>
              <Option value="completed">已完成</Option>
              <Option value="cancelled">已取消</Option>
            </Select>
            <Select
              placeholder="类型"
              value={searchParams.type || undefined}
              onChange={(value) =>
                setSearchParams({ ...searchParams, type: value, page: 1 })
              }
              allowClear
              style={{ width: 120 }}
            >
              <Option value="project">项目合同</Option>
              <Option value="retainer">年度框架</Option>
              <Option value="consulting">单次咨询</Option>
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

export default ContractList;
