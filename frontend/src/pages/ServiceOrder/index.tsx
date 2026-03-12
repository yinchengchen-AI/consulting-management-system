import React, { useState, useEffect } from 'react';
import { Table, Card, Button, Input, Tag, Space, Dropdown, message, Popconfirm } from 'antd';
import { PlusOutlined, SearchOutlined, MoreOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { getServiceOrders, deleteServiceOrder } from '@/api/service';
import type { ServiceOrder } from '@/types';

const ServiceOrderList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<ServiceOrder[]>([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [searchKeyword, setSearchKeyword] = useState('');

  // 获取服务订单列表
  const fetchData = async (page = 1, pageSize = 10) => {
    setLoading(true);
    try {
      const result = await getServiceOrders({
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
      await deleteServiceOrder(id);
      message.success('删除成功');
      fetchData(pagination.current, pagination.pageSize);
    } catch (error) {
      console.error('删除失败:', error);
    }
  };

  // 状态标签
  const getStatusTag = (status: number) => {
    const statusMap: Record<number, { color: string; text: string }> = {
      1: { color: 'default', text: '待启动' },
      2: { color: 'processing', text: '进行中' },
      3: { color: 'success', text: '已完成' },
      4: { color: 'warning', text: '已暂停' },
    };
    const { color, text } = statusMap[status] || { color: 'default', text: '未知' };
    return <Tag color={color}>{text}</Tag>;
  };

  // 表格列
  const columns = [
    {
      title: '服务编号',
      dataIndex: 'code',
      key: 'code',
      width: 120,
    },
    {
      title: '服务名称',
      dataIndex: 'name',
      key: 'name',
      ellipsis: true,
    },
    {
      title: '客户',
      dataIndex: ['customer', 'name'],
      key: 'customer',
      ellipsis: true,
    },
    {
      title: '服务类型',
      dataIndex: ['service_type', 'name'],
      key: 'serviceType',
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount: number) => `¥${amount?.toLocaleString() || 0}`,
    },
    {
      title: '进度',
      dataIndex: 'progress',
      key: 'progress',
      render: (progress: number) => `${progress || 0}%`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: getStatusTag,
    },
    {
      title: '开始日期',
      dataIndex: 'start_date',
      key: 'startDate',
      render: (date: string) => date ? new Date(date).toLocaleDateString() : '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: ServiceOrder) => (
        <Space>
          <Button
            type="text"
            icon={<EyeOutlined />}
            onClick={() => navigate(`/service-orders/${record.id}`)}
          />
          <Button
            type="text"
            icon={<EditOutlined />}
            onClick={() => navigate(`/service-orders/edit/${record.id}`)}
          />
          <Dropdown
            menu={{
              items: [
                {
                  key: 'delete',
                  label: (
                    <Popconfirm
                      title="确认删除"
                      description="确定要删除该服务订单吗？"
                      onConfirm={() => handleDelete(record.id)}
                      okText="确定"
                      cancelText="取消"
                    >
                      <span style={{ color: '#ff4d4f' }}>
                        <DeleteOutlined /> 删除
                      </span>
                    </Popconfirm>
                  ),
                },
              ],
            }}
          >
            <Button type="text" icon={<MoreOutlined />} />
          </Dropdown>
        </Space>
      ),
    },
  ];

  return (
    <Card
      title="服务订单管理"
      extra={
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => navigate('/service-orders/create')}
        >
          新建服务订单
        </Button>
      }
    >
      <div style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="搜索服务编号、名称或客户"
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

export default ServiceOrderList;
