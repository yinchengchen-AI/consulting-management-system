import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Statistic, DatePicker, Tabs, Spin } from 'antd';
import {
  UserOutlined,
  ProjectOutlined,
  DollarOutlined,
  FileTextOutlined,
  RiseOutlined,
  FallOutlined,
} from '@ant-design/icons';
import ReactECharts from 'echarts-for-react';
import { getCustomerStats, getServiceStats, getFinanceStats } from '@/api/statistics';
import type { CustomerStats, ServiceStats, FinanceStats } from '@/types';

const { RangePicker } = DatePicker;
const { TabPane } = Tabs;

const StatisticsDashboard: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [customerStats, setCustomerStats] = useState<CustomerStats | null>(null);
  const [serviceStats, setServiceStats] = useState<ServiceStats | null>(null);
  const [financeStats, setFinanceStats] = useState<FinanceStats | null>(null);

  // 获取统计数据
  const fetchStats = async () => {
    setLoading(true);
    try {
      const [customerData, serviceData, financeData] = await Promise.all([
        getCustomerStats(),
        getServiceStats(),
        getFinanceStats(),
      ]);
      setCustomerStats(customerData);
      setServiceStats(serviceData);
      setFinanceStats(financeData);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, []);

  // 客户行业分布图表配置
  const industryChartOption = {
    title: { text: '客户行业分布', left: 'center' },
    tooltip: { trigger: 'item' },
    legend: { bottom: '0%' },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        data: customerStats?.industryDistribution || [],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
        },
      },
    ],
  };

  // 服务类型分布图表配置
  const serviceTypeChartOption = {
    title: { text: '服务类型分布', left: 'center' },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: serviceStats?.typeDistribution?.map((item: any) => item.name) || [],
    },
    yAxis: { type: 'value' },
    series: [
      {
        type: 'bar',
        data: serviceStats?.typeDistribution?.map((item: any) => item.value) || [],
        itemStyle: {
          color: '#1890ff',
        },
      },
    ],
  };

  // 财务趋势图表配置
  const financeTrendChartOption = {
    title: { text: '财务趋势', left: 'center' },
    tooltip: { trigger: 'axis' },
    legend: { data: ['开票金额', '收款金额'], bottom: '0%' },
    xAxis: {
      type: 'category',
      data: financeStats?.monthlyTrend?.map((item: any) => item.month) || [],
    },
    yAxis: { type: 'value' },
    series: [
      {
        name: '开票金额',
        type: 'line',
        data: financeStats?.monthlyTrend?.map((item: any) => item.invoice) || [],
        smooth: true,
      },
      {
        name: '收款金额',
        type: 'line',
        data: financeStats?.monthlyTrend?.map((item: any) => item.receipt) || [],
        smooth: true,
      },
    ],
  };

  return (
    <Spin spinning={loading}>
      <div style={{ marginBottom: 16 }}>
        <RangePicker style={{ width: 300 }} />
      </div>

      <Tabs defaultActiveKey="overview">
        <TabPane tab="数据概览" key="overview">
          {/* 统计卡片 */}
          <Row gutter={16} style={{ marginBottom: 24 }}>
            <Col span={6}>
              <Card>
                <Statistic
                  title="总客户数"
                  value={customerStats?.total || 0}
                  prefix={<UserOutlined />}
                  valueStyle={{ color: '#1890ff' }}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="服务订单数"
                  value={serviceStats?.total || 0}
                  prefix={<ProjectOutlined />}
                  valueStyle={{ color: '#52c41a' }}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="总开票金额"
                  value={financeStats?.totalInvoice || 0}
                  prefix={<FileTextOutlined />}
                  suffix="元"
                  valueStyle={{ color: '#faad14' }}
                  formatter={(value) => `¥${Number(value).toLocaleString()}`}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="总收款金额"
                  value={financeStats?.totalReceipt || 0}
                  prefix={<DollarOutlined />}
                  suffix="元"
                  valueStyle={{ color: '#13c2c2' }}
                  formatter={(value) => `¥${Number(value).toLocaleString()}`}
                />
              </Card>
            </Col>
          </Row>

          <Row gutter={16} style={{ marginBottom: 24 }}>
            <Col span={6}>
              <Card>
                <Statistic
                  title="新增客户(本月)"
                  value={customerStats?.newThisMonth || 0}
                  prefix={<RiseOutlined />}
                  valueStyle={{ color: '#52c41a' }}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="进行中服务"
                  value={serviceStats?.inProgress || 0}
                  prefix={<ProjectOutlined />}
                  valueStyle={{ color: '#1890ff' }}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="回款率"
                  value={financeStats?.collectionRate || 0}
                  suffix="%"
                  prefix={financeStats?.collectionRate >= 80 ? <RiseOutlined /> : <FallOutlined />}
                  valueStyle={{ color: financeStats?.collectionRate >= 80 ? '#52c41a' : '#ff4d4f' }}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="逾期未收款"
                  value={financeStats?.overdueAmount || 0}
                  prefix={<FallOutlined />}
                  suffix="元"
                  valueStyle={{ color: '#ff4d4f' }}
                  formatter={(value) => `¥${Number(value).toLocaleString()}`}
                />
              </Card>
            </Col>
          </Row>

          {/* 图表 */}
          <Row gutter={16}>
            <Col span={12}>
              <Card title="客户行业分布">
                <ReactECharts option={industryChartOption} style={{ height: 300 }} />
              </Card>
            </Col>
            <Col span={12}>
              <Card title="服务类型分布">
                <ReactECharts option={serviceTypeChartOption} style={{ height: 300 }} />
              </Card>
            </Col>
          </Row>

          <Row gutter={16} style={{ marginTop: 16 }}>
            <Col span={24}>
              <Card title="财务趋势">
                <ReactECharts option={financeTrendChartOption} style={{ height: 350 }} />
              </Card>
            </Col>
          </Row>
        </TabPane>

        <TabPane tab="客户分析" key="customer">
          <Row gutter={16}>
            <Col span={12}>
              <Card title="客户状态分布">
                <ReactECharts option={industryChartOption} style={{ height: 350 }} />
              </Card>
            </Col>
            <Col span={12}>
              <Card title="高价值客户排行">
                <Table
                  dataSource={customerStats?.topCustomers || []}
                  columns={[
                    { title: '客户名称', dataIndex: 'name' },
                    { title: '合作金额', dataIndex: 'amount', render: (v: number) => `¥${v?.toLocaleString()}` },
                  ]}
                  pagination={false}
                />
              </Card>
            </Col>
          </Row>
        </TabPane>

        <TabPane tab="财务分析" key="finance">
          <Row gutter={16}>
            <Col span={24}>
              <Card title="月度财务统计">
                <ReactECharts option={financeTrendChartOption} style={{ height: 400 }} />
              </Card>
            </Col>
          </Row>
        </TabPane>
      </Tabs>
    </Spin>
  );
};

// 引入Table用于高价值客户排行
import { Table } from 'antd';

export default StatisticsDashboard;
