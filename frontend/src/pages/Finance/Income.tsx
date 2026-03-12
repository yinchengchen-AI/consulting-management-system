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
  DatePicker,
} from 'antd';
import { PlusOutlined, SearchOutlined } from '@ant-design/icons';
import { getIncomeList } from '@/api/finance';
import { Income, IncomeStatus, IncomeType } from '@/types';
import dayjs from 'dayjs';

const { Option } = Select;
const { RangePicker } = DatePicker;

const statusMap: Record<IncomeStatus, { text: string; color: string }> = {
  pending: { text: '待确认', color: 'warning' },
  received: { text: '已收款', color: 'success' },
  invoiced: { text: '已开票', color: 'processing' },
};

const typeMap: Record<IncomeType, string> = {
  contract: '合同收入',
  project: '项目收入',
  other: '其他收入',
};

const IncomeList = () => {
  const [searchParams, setSearchParams] = useState({
    page: 1,
    page_size: 10,
    start_date: '',
    end_date: '',
    customer_id: '',
    status: '',
  });

  const { data, isLoading } = useQuery({
    queryKey: ['incomes', searchParams],
    queryFn: () => getIncomeList(searchParams),
  });

  const columns = [
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: IncomeType) => typeMap[type] || type,
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount: number) => `¥${(amount || 0).toLocaleString()}`,
    },
    {
      title: '客户',
      dataIndex: 'customer',
      key: 'customer',
      render: (customer: any) => customer?.name || '-',
    },
    {
      title: '关联合同',
      dataIndex: 'contract',
      key: 'contract',
      render: (contract: any) => contract?.name || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: IncomeStatus) => (
        <Tag color={statusMap[status]?.color}>{statusMap[status]?.text}</Tag>
      ),
    },
    {
      title: '发票号',
      dataIndex: 'invoice_no',
      key: 'invoice_no',
      render: (no: string) => no || '-',
    },
    {
      title: '开票日期',
      dataIndex: 'invoice_date',
      key: 'invoice_date',
      render: (date: string) => (date ? new Date(date).toLocaleDateString() : '-'),
    },
    {
      title: '收款日期',
      dataIndex: 'received_date',
      key: 'received_date',
      render: (date: string) => (date ? new Date(date).toLocaleDateString() : '-'),
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
          <h2 style={{ margin: 0 }}>收入管理</h2>
          <Button type="primary" icon={<PlusOutlined />}>
            新增收入
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
              placeholder="状态"
              value={searchParams.status || undefined}
              onChange={(value) =>
                setSearchParams({ ...searchParams, status: value, page: 1 })
              }
              allowClear
              style={{ width: 120 }}
            >
              <Option value="pending">待确认</Option>
              <Option value="received">已收款</Option>
              <Option value="invoiced">已开票</Option>
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

export default IncomeList;
