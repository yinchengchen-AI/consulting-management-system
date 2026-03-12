import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { Card, Descriptions, Tag, Button, Spin, Empty } from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { getCustomerDetail } from '@/api/customer';
import { CustomerLevel, CustomerStatus, CustomerType } from '@/types';

const typeMap: Record<CustomerType, string> = {
  enterprise: '企业',
  government: '政府',
  individual: '个人',
  other: '其他',
};

const levelMap: Record<CustomerLevel, { text: string; color: string }> = {
  A: { text: 'A级-战略客户', color: 'red' },
  B: { text: 'B级-重要客户', color: 'orange' },
  C: { text: 'C级-普通客户', color: 'blue' },
  D: { text: 'D级-潜在客户', color: 'default' },
};

const statusMap: Record<CustomerStatus, { text: string; color: string }> = {
  active: { text: '正常', color: 'success' },
  inactive: { text: '停用', color: 'default' },
  potential: { text: '潜在', color: 'warning' },
};

const CustomerDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data, isLoading } = useQuery({
    queryKey: ['customer', id],
    queryFn: () => getCustomerDetail(id!),
    enabled: !!id,
  });

  const customer = data?.data;

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: 50 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!customer) {
    return <Empty description="客户不存在" />;
  }

  return (
    <div>
      <Card>
        <div style={{ marginBottom: 24 }}>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/customers')}>
            返回列表
          </Button>
        </div>

        <Descriptions title="基本信息" bordered column={2}>
          <Descriptions.Item label="客户名称">{customer.name}</Descriptions.Item>
          <Descriptions.Item label="客户类型">
            {typeMap[customer.type] || customer.type}
          </Descriptions.Item>
          <Descriptions.Item label="客户等级">
            <Tag color={levelMap[customer.level]?.color}>
              {levelMap[customer.level]?.text}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={statusMap[customer.status]?.color}>
              {statusMap[customer.status]?.text}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="所属行业">{customer.industry || '-'}</Descriptions.Item>
          <Descriptions.Item label="企业规模">{customer.scale || '-'}</Descriptions.Item>
          <Descriptions.Item label="公司网站">{customer.website || '-'}</Descriptions.Item>
          <Descriptions.Item label="公司地址">{customer.address || '-'}</Descriptions.Item>
        </Descriptions>

        <Descriptions title="联系人信息" bordered column={2} style={{ marginTop: 24 }}>
          <Descriptions.Item label="联系人姓名">{customer.contact_name || '-'}</Descriptions.Item>
          <Descriptions.Item label="联系人职位">{customer.contact_title || '-'}</Descriptions.Item>
          <Descriptions.Item label="联系电话">{customer.contact_phone || '-'}</Descriptions.Item>
          <Descriptions.Item label="电子邮箱">{customer.contact_email || '-'}</Descriptions.Item>
        </Descriptions>

        <Descriptions title="业务信息" bordered column={2} style={{ marginTop: 24 }}>
          <Descriptions.Item label="销售负责人">
            {customer.sales_owner?.real_name || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="合同数量">{customer.contract_count || 0}</Descriptions.Item>
          <Descriptions.Item label="项目数量">{customer.project_count || 0}</Descriptions.Item>
          <Descriptions.Item label="累计金额">
            ¥{(customer.total_amount || 0).toLocaleString()}
          </Descriptions.Item>
        </Descriptions>

        {customer.description && (
          <Descriptions title="备注" bordered column={1} style={{ marginTop: 24 }}>
            <Descriptions.Item label="描述">{customer.description}</Descriptions.Item>
          </Descriptions>
        )}
      </Card>
    </div>
  );
};

export default CustomerDetail;
