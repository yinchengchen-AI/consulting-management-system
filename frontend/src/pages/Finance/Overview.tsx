import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Card, Row, Col, Statistic, DatePicker, Space } from 'antd';
import {
  ArrowUpOutlined,
  ArrowDownOutlined,
  DollarOutlined,
  WalletOutlined,
  PieChartOutlined,
} from '@ant-design/icons';
import ReactECharts from 'echarts-for-react';
import dayjs from 'dayjs';
import { getFinanceSummary } from '@/api/finance';

const { RangePicker } = DatePicker;

const FinanceOverview = () => {
  const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs]>([
    dayjs().startOf('month'),
    dayjs(),
  ]);

  const { data, isLoading } = useQuery({
    queryKey: ['finance-summary', dateRange],
    queryFn: () =>
      getFinanceSummary({
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
      }),
  });

  const summary = data?.data;

  // 月度趋势图配置
  const trendChartOption = {
    title: {
      text: '月度收支趋势',
      left: 'center',
    },
    tooltip: {
      trigger: 'axis',
    },
    legend: {
      data: ['收入', '支出', '利润'],
      bottom: 0,
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '15%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: ['1月', '2月', '3月', '4月', '5月', '6月'],
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: (value: number) => `¥${(value / 10000).toFixed(0)}万`,
      },
    },
    series: [
      {
        name: '收入',
        type: 'line',
        data: [1200000, 1320000, 1010000, 1340000, 900000, 2300000],
        smooth: true,
        itemStyle: { color: '#52c41a' },
      },
      {
        name: '支出',
        type: 'line',
        data: [800000, 900000, 700000, 850000, 600000, 1200000],
        smooth: true,
        itemStyle: { color: '#ff4d4f' },
      },
      {
        name: '利润',
        type: 'line',
        data: [400000, 420000, 310000, 490000, 300000, 1100000],
        smooth: true,
        itemStyle: { color: '#1890ff' },
      },
    ],
  };

  // 收入构成饼图
  const incomePieOption = {
    title: {
      text: '收入构成',
      left: 'center',
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: ¥{c} ({d}%)',
    },
    legend: {
      orient: 'vertical',
      left: 'left',
    },
    series: [
      {
        name: '收入类型',
        type: 'pie',
        radius: '60%',
        data: [
          { value: 4500000, name: '项目收入', itemStyle: { color: '#1890ff' } },
          { value: 2800000, name: '合同收入', itemStyle: { color: '#52c41a' } },
          { value: 500000, name: '其他收入', itemStyle: { color: '#faad14' } },
        ],
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

  // 支出构成饼图
  const expensePieOption = {
    title: {
      text: '支出构成',
      left: 'center',
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: ¥{c} ({d}%)',
    },
    legend: {
      orient: 'vertical',
      left: 'left',
    },
    series: [
      {
        name: '支出类型',
        type: 'pie',
        radius: '60%',
        data: [
          { value: 3200000, name: '工资', itemStyle: { color: '#ff4d4f' } },
          { value: 1200000, name: '项目成本', itemStyle: { color: '#faad14' } },
          { value: 500000, name: '办公费用', itemStyle: { color: '#722ed1' } },
          { value: 400000, name: '差旅费', itemStyle: { color: '#13c2c2' } },
          { value: 300000, name: '其他', itemStyle: { color: '#bfbfbf' } },
        ],
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

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Space>
          <span>日期范围:</span>
          <RangePicker
            value={dateRange}
            onChange={(dates) => {
              if (dates) {
                setDateRange([dates[0]!, dates[1]!]);
              }
            }}
          />
        </Space>
      </Card>

      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={isLoading}>
            <Statistic
              title="总收入"
              value={summary?.total_income || 0}
              precision={2}
              prefix="¥"
              valueStyle={{ color: '#52c41a' }}
              formatter={(value) => (value as number).toLocaleString()}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={isLoading}>
            <Statistic
              title="总支出"
              value={summary?.total_expense || 0}
              precision={2}
              prefix="¥"
              valueStyle={{ color: '#ff4d4f' }}
              formatter={(value) => (value as number).toLocaleString()}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={isLoading}>
            <Statistic
              title="净利润"
              value={summary?.net_profit || 0}
              precision={2}
              prefix="¥"
              valueStyle={{ color: (summary?.net_profit || 0) >= 0 ? '#52c41a' : '#ff4d4f' }}
              formatter={(value) => (value as number).toLocaleString()}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={isLoading}>
            <Statistic
              title="利润率"
              value={summary?.profit_margin || 0}
              precision={2}
              suffix="%"
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} lg={12}>
          <Card>
            <ReactECharts option={trendChartOption} style={{ height: 350 }} />
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Row gutter={[16, 16]}>
            <Col span={24}>
              <Card>
                <ReactECharts option={incomePieOption} style={{ height: 250 }} />
              </Card>
            </Col>
          </Row>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={12}>
          <Card>
            <ReactECharts option={expensePieOption} style={{ height: 300 }} />
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card title="待处理事项">
            <div style={{ padding: '16px 0' }}>
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  padding: '12px 0',
                  borderBottom: '1px solid #f0f0f0',
                }}
              >
                <span>待确认收入</span>
                <span style={{ color: '#faad14', fontWeight: 600 }}>
                  ¥{(summary?.pending_income || 0).toLocaleString()}
                </span>
              </div>
              <div
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  padding: '12px 0',
                }}
              >
                <span>待支付支出</span>
                <span style={{ color: '#ff4d4f', fontWeight: 600 }}>
                  ¥{(summary?.pending_expense || 0).toLocaleString()}
                </span>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default FinanceOverview;
