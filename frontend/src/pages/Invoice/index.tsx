import React, { useState, useEffect } from 'react';
import { Table, Card, Button, Input, Tag, Space, message, Popconfirm } from 'antd';
import { PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined, EyeOutlined, CheckOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { getInvoices, deleteInvoice, auditInvoice } from '@/api/finance';
import type { Invoice } from '@/types';

const InvoiceList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<Invoice[]>([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [searchKeyword, setSearchKeyword] = useState('');

  // 获取开票列表
  const fetchData = async (page = 1, pageSize = 10) => {
    setLoading(true);
    try {
      const result = await getInvoices({
        page,
        pageSize,
        keyword: searchKeyword,
      });
      setData(result.list);
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
    fetchData();
  }, []);

  // 处理删除
  const handleDelete = async (id: number) => {
    try {
      await deleteInvoice(id);
      message.success('删除成功');
      fetchData(pagination.current, pagination.pageSize);
    } catch (error) {
      console.error('删除失败:', error);
    }
  };

  // 处理审核
  const handleAudit = async (id: number, approved: boolean) => {
    try {
      await auditInvoice(id, { approved });
      message.success(approved ? '审核通过' : '审核拒绝');
      fetchData(pagination.current, pagination.pageSize);
    } catch (error) {
      console.error('审核失败:', error);
    }
  };

  // 状态标签
  const getStatusTag = (status: number) => {
    const statusMap: Record<number, { color: string; text: string }> = {
      1: { color: 'default', text: '待开票' },
      2: { color: 'processing', text: '已开票' },
      3: { color: 'error', text: '已作废' },
    };
    const { color, text } = statusMap[status] || { color: 'default', text: '未知' };
    return <Tag color={color}>{text}</Tag>;
  };

  // 发票类型标签
  const getTypeTag = (type: number) => {
    return type === 2 ? <Tag color="blue">专票</Tag> : <Tag>普票</Tag>;
  };

  // 表格列
  const columns = [
    {
      title: '发票号码',
      dataIndex: 'invoice_no',
      key: 'invoiceNo',
      render: (no: string) => no || '-',
    },
    {
      title: '客户',
      dataIndex: ['customer', 'name'],
      key: 'customer',
    },
    {
      title: '发票类型',
      dataIndex: 'invoice_type',
      key: 'invoiceType',
      render: getTypeTag,
    },
    {
      title: '开票金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount: number) => `¥${amount?.toLocaleString() || 0}`,
    },
    {
      title: '税率',
      dataIndex: 'tax_rate',
      key: 'taxRate',
      render: (rate: number) => `${rate}%`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: getStatusTag,
    },
    {
      title: '开票日期',
      dataIndex: 'invoice_date',
      key: 'invoiceDate',
      render: (date: string) => date ? new Date(date).toLocaleDateString() : '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: Invoice) => (
        <Space>
          <Button
            type="text"
            icon={<EyeOutlined />}
            onClick={() => navigate(`/invoices/${record.id}`)}
          />
          <Button
            type="text"
            icon={<EditOutlined />}
            onClick={() => navigate(`/invoices/edit/${record.id}`)}
          />
          {record.status === 1 && (
            <>
              <Popconfirm
                title="审核通过"
                onConfirm={() => handleAudit(record.id, true)}
              >
                <Button type="text" icon={<CheckOutlined />} title="审核通过" />
              </Popconfirm>
            </>
          )}
          <Popconfirm
            title="确认删除"
            description="确定要删除该开票记录吗？"
            onConfirm={() => handleDelete(record.id)}
          >
            <Button type="text" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card
      title="开票管理"
      extra={
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => navigate('/invoices/create')}
        >
          申请开票
        </Button>
      }
    >
      <div style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="搜索发票号码或客户"
          allowClear
          enterButton={<><SearchOutlined /> 搜索</>}
          value={searchKeyword}
          onChange={(e) => setSearchKeyword(e.target.value)}
          onSearch={() => fetchData(1, pagination.pageSize)}
          style={{ width: 350 }}
        />
      </div>
      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        pagination={{
          ...pagination,
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条`,
          onChange: (page, pageSize) => fetchData(page, pageSize),
        }}
      />
    </Card>
  );
};

export default InvoiceList;
