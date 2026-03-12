import React, { useState, useEffect } from 'react';
import { Card, Descriptions, Tag, Timeline, Button, message, Spin, Divider } from 'antd';
import { ArrowLeftOutlined, EditOutlined, FileTextOutlined } from '@ant-design/icons';
import { useParams, useNavigate } from 'react-router-dom';
import { getServiceOrder, getCommunications } from '@/api/service';
import type { ServiceOrder, Communication } from '@/types';

const ServiceOrderDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<ServiceOrder | null>(null);
  const [communications, setCommunications] = useState<Communication[]>([]);

  // 获取服务订单详情
  const fetchDetail = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const result = await getServiceOrder(Number(id));
      setData(result);
      // 获取沟通纪要
      const commResult = await getCommunications(Number(id));
      setCommunications(commResult);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDetail();
  }, [id]);

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

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!data) {
    return <Card>服务订单不存在</Card>;
  }

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/service-orders')}>
          返回列表
        </Button>
        <Button
          type="primary"
          icon={<EditOutlined />}
          style={{ marginLeft: 8 }}
          onClick={() => navigate(`/service-orders/edit/${id}`)}
        >
          编辑
        </Button>
      </div>

      <Card title="基本信息" style={{ marginBottom: 16 }}>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="服务编号">{data.code}</Descriptions.Item>
          <Descriptions.Item label="服务名称">{data.name}</Descriptions.Item>
          <Descriptions.Item label="客户">{data.customer?.name}</Descriptions.Item>
          <Descriptions.Item label="服务类型">{data.service_type?.name}</Descriptions.Item>
          <Descriptions.Item label="服务金额">¥{data.amount?.toLocaleString()}</Descriptions.Item>
          <Descriptions.Item label="状态">{getStatusTag(data.status)}</Descriptions.Item>
          <Descriptions.Item label="进度">{data.progress}%</Descriptions.Item>
          <Descriptions.Item label="负责人">
            {data.participants?.join(', ') || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="开始日期">
            {data.start_date ? new Date(data.start_date).toLocaleDateString() : '-'}
          </Descriptions.Item>
          <Descriptions.Item label="结束日期">
            {data.end_date ? new Date(data.end_date).toLocaleDateString() : '-'}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="服务描述" style={{ marginBottom: 16 }}>
        <div style={{ whiteSpace: 'pre-wrap' }}>{data.description || '暂无描述'}</div>
      </Card>

      <Card title="沟通纪要" extra={<Button icon={<FileTextOutlined />}>添加纪要</Button>}>
        {communications.length > 0 ? (
          <Timeline
            items={communications.map((comm) => ({
              children: (
                <div>
                  <div style={{ color: '#666', fontSize: 12 }}>
                    {new Date(comm.created_at).toLocaleString()} - {comm.user?.real_name}
                  </div>
                  <div style={{ marginTop: 4 }}>{comm.content}</div>
                </div>
              ),
            }))}
          />
        ) : (
          <div style={{ textAlign: 'center', color: '#999', padding: '20px' }}>暂无沟通纪要</div>
        )}
      </Card>
    </div>
  );
};

export default ServiceOrderDetail;
