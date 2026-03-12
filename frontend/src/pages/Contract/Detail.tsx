import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { Card, Descriptions, Tag, Button, Spin, Empty, Progress } from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { getContractDetail } from '@/api/contract';
import { ContractStatus, ContractType, PaymentTerms } from '@/types';

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
  retainer: '年度框架合同',
  consulting: '单次咨询合同',
};

const paymentTermsMap: Record<PaymentTerms, string> = {
  prepay: '预付',
  milestone: '里程碑付款',
  monthly: '月付',
  postpay: '后付',
};

const ContractDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data, isLoading } = useQuery({
    queryKey: ['contract', id],
    queryFn: () => getContractDetail(id!),
    enabled: !!id,
  });

  const contract = data?.data;

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: 50 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!contract) {
    return <Empty description="合同不存在" />;
  }

  const paymentProgress = contract.total_amount > 0
    ? Math.round((contract.paid_amount / contract.total_amount) * 100)
    : 0;

  return (
    <div>
      <Card>
        <div style={{ marginBottom: 24 }}>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/contracts')}>
            返回列表
          </Button>
        </div>

        <Descriptions title="基本信息" bordered column={2}>
          <Descriptions.Item label="合同编号">{contract.code}</Descriptions.Item>
          <Descriptions.Item label="合同名称">{contract.name}</Descriptions.Item>
          <Descriptions.Item label="合同类型">
            {typeMap[contract.type] || contract.type}
          </Descriptions.Item>
          <Descriptions.Item label="合同状态">
            <Tag color={statusMap[contract.status]?.color}>
              {statusMap[contract.status]?.text}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="客户">{contract.customer?.name || '-'}</Descriptions.Item>
          <Descriptions.Item label="销售负责人">
            {contract.sales_owner?.real_name || '-'}
          </Descriptions.Item>
        </Descriptions>

        <Descriptions title="金额信息" bordered column={2} style={{ marginTop: 24 }}>
          <Descriptions.Item label="合同金额">
            ¥{(contract.amount || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="税率">{contract.tax_rate}%</Descriptions.Item>
          <Descriptions.Item label="税额">
            ¥{(contract.tax_amount || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="总金额">
            <strong>¥{(contract.total_amount || 0).toLocaleString()}</strong>
          </Descriptions.Item>
        </Descriptions>

        <Descriptions title="付款信息" bordered column={2} style={{ marginTop: 24 }}>
          <Descriptions.Item label="付款条款">
            {paymentTermsMap[contract.payment_terms] || contract.payment_terms}
          </Descriptions.Item>
          <Descriptions.Item label="签约人">{contract.signed_by || '-'}</Descriptions.Item>
          <Descriptions.Item label="签约日期">
            {contract.signed_date ? new Date(contract.signed_date).toLocaleDateString() : '-'}
          </Descriptions.Item>
          <Descriptions.Item label="生效日期">
            {contract.start_date ? new Date(contract.start_date).toLocaleDateString() : '-'}
          </Descriptions.Item>
          <Descriptions.Item label="到期日期">
            {contract.end_date ? new Date(contract.end_date).toLocaleDateString() : '-'}
          </Descriptions.Item>
        </Descriptions>

        <Descriptions title="收款进度" bordered column={1} style={{ marginTop: 24 }}>
          <Descriptions.Item label="进度">
            <Progress percent={paymentProgress} />
            <div style={{ marginTop: 8 }}>
              已收款: ¥{(contract.paid_amount || 0).toLocaleString()} / 
              待收款: ¥{(contract.remaining_amount || 0).toLocaleString()}
            </div>
          </Descriptions.Item>
        </Descriptions>

        {contract.description && (
          <Descriptions title="合同描述" bordered column={1} style={{ marginTop: 24 }}>
            <Descriptions.Item label="描述">{contract.description}</Descriptions.Item>
          </Descriptions>
        )}

        {contract.terms && (
          <Descriptions title="合同条款" bordered column={1} style={{ marginTop: 24 }}>
            <Descriptions.Item label="条款">{contract.terms}</Descriptions.Item>
          </Descriptions>
        )}
      </Card>
    </div>
  );
};

export default ContractDetail;
