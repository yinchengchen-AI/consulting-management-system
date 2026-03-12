import React, { useState, useEffect } from 'react';
import { Table, Card, Button, Input, Tag, Space, message, Popconfirm, Tabs } from 'antd';
import { PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { getReceipts, getPaymentPlans, deleteReceipt, deletePaymentPlan } from '@/api/finance';
import type { Receipt, PaymentPlan } from '@/types';

const { TabPane } = Tabs;

const ReceiptList: React.FC = () => {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState('receipts');
  const [loading, setLoading] = useState(false);
  const [receipts, setReceipts] = useState<Receipt[]>([]);
  const [plans, setPlans] = useState<PaymentPlan[]>([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [searchKeyword, setSearchKeyword] = useState('');

  // 获取收款记录
  const fetchReceipts = async (page = 1, pageSize = 10) => {
    setLoading(true);
    try {
      const result = await getReceipts({
        page,
        pageSize,
        keyword: searchKeyword,
      });
      setReceipts(result.list);
      setPagination({
        current: page,
        pageSize,
        total: result.total,
      });
    } finally {
      setLoading(false);
    }
  };

  // 获取收款计划
  const fetchPlans = async (page = 1, pageSize = 10) => {
    setLoading(true);
    try {
      const result = await getPaymentPlans({
        page,
        pageSize,
        keyword: searchKeyword,
      });
      setPlans(result.list);
      setPagination({
        current: page,
        pageSize,
        total: result.total,
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (activeTab === 'receipts') {
      fetchReceipts();
    } else {
      fetchPlans();
    }
  }, [activeTab]);

  // 处理删除收款记录
  const handleDeleteReceipt = async (id: number) => {
    try {
      await deleteReceipt(id);
      message.success('删除成功');
      fetchReceipts(pagination.current, pagination.pageSize);
    } catch (error) {
      console.error('删除失败:', error);
    }
  };

  // 处理删除收款计划
  const handleDeletePlan = async (id: number) => {
    try {
      await deletePaymentPlan(id);
      message.success('删除成功');
      fetchPlans(pagination.current, pagination.pageSize);
    } catch (error) {
      console.error('删除失败:', error);
    }
  };

  // 状态标签
  const getPlanStatusTag = (status: number) => {
    const statusMap: Record<number, { color: string; text: string }> = {
      1: { color: 'default', text: '待收款' },
      2: { color: 'processing', text: '部分收款' },
      3: { color: 'success', text: '已收完' },
    };
    const { color, text } = statusMap[status] || { color: 'default', text: '未知' };
    return <Tag color={color}>{text}</Tag>;
  };

  // 收款记录表格列
  const receiptColumns = [
    {
      title: '客户',
      dataIndex: ['plan', 'customer', 'name'],
      key: 'customer',
    },
    {
      title: '收款金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount: number) => `¥${amount?.toLocaleString() || 0}`,
    },
    {
      title: '收款日期',
      dataIndex: 'received_date',
      key: 'receivedDate',
      render: (date: string) => date ? new Date(date).toLocaleDateString() : '-',
    },
    {
      title: '收款方式',
      dataIndex: 'payment_method',
      key: 'paymentMethod',
      render: (method: number) => {
        const methods = ['', '银行转账', '现金', '支票', '其他'];
        return methods[method] || '其他';
      },
    },
    {
      title: '备注',
      dataIndex: 'remark',
      key: 'remark',
      ellipsis: true,
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Receipt) => (
        <Space>
          <Button
            type="text"
            icon={<EyeOutlined />}
            onClick={() => navigate(`/receipts/${record.id}`)}
          />
          <Button
            type="text"
            icon={<EditOutlined />}
            onClick={() => navigate(`/receipts/edit/${record.id}`)}
          />
          <Popconfirm
            title="确认删除"
            onConfirm={() => handleDeleteReceipt(record.id)}
          >
            <Button type="text" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  // 收款计划表格列
  const planColumns = [
    {
      title: '客户',
      dataIndex: ['customer', 'name'],
      key: 'customer',
    },
    {
      title: '计划金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount: number) => `¥${amount?.toLocaleString() || 0}`,
    },
    {
      title: '计划日期',
      dataIndex: 'planned_date',
      key: 'plannedDate',
      render: (date: string) => date ? new Date(date).toLocaleDateString() : '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: getPlanStatusTag,
    },
    {
      title: '备注',
      dataIndex: 'remark',
      key: 'remark',
      ellipsis: true,
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: PaymentPlan) => (
        <Space>
          <Button
            type="text"
            icon={<EyeOutlined />}
            onClick={() => navigate(`/payment-plans/${record.id}`)}
          />
          <Button
            type="text"
            icon={<EditOutlined />}
            onClick={() => navigate(`/payment-plans/edit/${record.id}`)}
          />
          <Button
            type="primary"
            size="small"
            onClick={() => navigate(`/receipts/create?planId=${record.id}`)}
          >
            收款
          </Button>
          <Popconfirm
            title="确认删除"
            onConfirm={() => handleDeletePlan(record.id)}
          >
            <Button type="text" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card title="收款管理">
      <Tabs activeKey={activeTab} onChange={setActiveTab}>
        <TabPane tab="收款记录" key="receipts">
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
            <Input.Search
              placeholder="搜索客户"
              allowClear
              enterButton={<><SearchOutlined /> 搜索</>}
              value={searchKeyword}
              onChange={(e) => setSearchKeyword(e.target.value)}
              onSearch={() => fetchReceipts(1, pagination.pageSize)}
              style={{ width: 350 }}
            />
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => navigate('/receipts/create')}
            >
              录入收款
            </Button>
          </div>
          <Table
            columns={receiptColumns}
            dataSource={receipts}
            rowKey="id"
            loading={loading}
            pagination={{
              ...pagination,
              showSizeChanger: true,
              showTotal: (total) => `共 ${total} 条`,
              onChange: (page, pageSize) => fetchReceipts(page, pageSize),
            }}
          />
        </TabPane>
        <TabPane tab="收款计划" key="plans">
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
            <Input.Search
              placeholder="搜索客户"
              allowClear
              enterButton={<><SearchOutlined /> 搜索</>}
              value={searchKeyword}
              onChange={(e) => setSearchKeyword(e.target.value)}
              onSearch={() => fetchPlans(1, pagination.pageSize)}
              style={{ width: 350 }}
            />
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => navigate('/payment-plans/create')}
            >
              创建计划
            </Button>
          </div>
          <Table
            columns={planColumns}
            dataSource={plans}
            rowKey="id"
            loading={loading}
            pagination={{
              ...pagination,
              showSizeChanger: true,
              showTotal: (total) => `共 ${total} 条`,
              onChange: (page, pageSize) => fetchPlans(page, pageSize),
            }}
          />
        </TabPane>
      </Tabs>
    </Card>
  );
};

export default ReceiptList;
