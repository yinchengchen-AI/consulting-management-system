import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  Card,
  Table,
  Button,
  Select,
  Space,
  Tag,
  DatePicker,
} from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { getExpenseList } from '@/api/finance';
import { Expense, ExpenseStatus, ExpenseType } from '@/types';

const { Option } = Select;
const { RangePicker } = DatePicker;

const statusMap: Record<ExpenseStatus, { text: string; color: string }> = {
  pending: { text: '待支付', color: 'warning' },
  paid: { text: '已支付', color: 'success' },
  cancelled: { text: '已取消', color: 'default' },
};

const typeMap: Record<ExpenseType, string> = {
  salary: '工资',
  bonus: '奖金',
  office: '办公费用',
  travel: '差旅费',
  marketing: '市场费用',
  project: '项目成本',
  tax: '税费',
  other: '其他',
};

const ExpenseList = () => {
  const [searchParams, setSearchParams] = useState({
    page: 1,
    page_size: 10,
    start_date: '',
    end_date: '',
    project_id: '',
    type: '',
    status: '',
  });

  const { data, isLoading } = useQuery({
    queryKey: ['expenses', searchParams],
    queryFn: () => getExpenseList(searchParams),
  });

  const columns = [
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: ExpenseType) => typeMap[type] || type,
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount: number) => `¥${(amount || 0).toLocaleString()}`,
    },
    {
      title: '项目',
      dataIndex: 'project',
      key: 'project',
      render: (project: any) => project?.name || '-',
    },
    {
      title: '收款方',
      dataIndex: 'payee',
      key: 'payee',
      render: (payee: string) => payee || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: ExpenseStatus) => (
        <Tag color={statusMap[status]?.color}>{statusMap[status]?.text}</Tag>
      ),
    },
    {
      title: '支付日期',
      dataIndex: 'paid_date',
      key: 'paid_date',
      render: (date: string) => (date ? new Date(date).toLocaleDateString() : '-'),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      render: (desc: string) => desc || '-',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
  ];

  const list = data?.data?.list || [];
  const total = data?.data?.total || 0;

  return (
    <div>
      <Card>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 24 }}>
          <h2 style={{ margin: 0 }}>支出管理</h2>
          <Button type="primary" icon={<PlusOutlined />}>
            新增支出
          </Button>
        </div>

        <div style={{ marginBottom: 16, padding: 16, background: '#fafafa', borderRadius: 8 }}>
          <Space>
            <RangePicker
              onChange={(dates) => {
                if (dates) {
                  setSearchParams({
                    ...searchParams,
                    start_date: dates[0]?.format('YYYY-MM-DD') || '',
                    end_date: dates[1]?.format('YYYY-MM-DD') || '',
                    page: 1,
                  });
                }
              }}
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
              <Option value="salary">工资</Option>
              <Option value="bonus">奖金</Option>
              <Option value="office">办公费用</Option>
              <Option value="travel">差旅费</Option>
              <Option value="marketing">市场费用</Option>
              <Option value="project">项目成本</Option>
              <Option value="tax">税费</Option>
              <Option value="other">其他</Option>
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
              <Option value="pending">待支付</Option>
              <Option value="paid">已支付</Option>
              <Option value="cancelled">已取消</Option>
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

export default ExpenseList;
