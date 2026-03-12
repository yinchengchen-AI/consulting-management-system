import React, { useEffect, useRef } from 'react';
import { Row, Col, Card, Statistic, List, Tag, Badge, Spin, Empty } from 'antd';
import {
  TeamOutlined,
  SolutionOutlined,
  DollarOutlined,
  FileTextOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined,
  BellOutlined,
  CheckCircleOutlined,
} from '@ant-design/icons';
import { useQuery } from '@tanstack/react-query';
import * as echarts from 'echarts';
import { getDashboardStats, getTodoList } from '@/api/statistics';
import { formatMoney, formatNumber } from '@/utils';
import type { DashboardStats, TodoItem } from '@/types';

const Dashboard: React.FC = () => {
  const customerChartRef = useRef<HTMLDivElement>(null);
  const serviceChartRef = useRef<HTMLDivElement>(null);
  const financeChartRef = useRef<HTMLDivElement>(null);

  // 获取仪表盘统计数据
  const { data: stats, isLoading: statsLoading } = useQuery<DashboardStats>({
    queryKey: ['dashboardStats'],
    queryFn: getDashboardStats,
  });

  // 获取待办事项
  const { data: todos, isLoading: todosLoading } = useQuery<TodoItem[]>({
    queryKey: ['todos'],
    queryFn: getTodoList,
  });

  // 初始化图表
  useEffect(() => {
    if (!customerChartRef.current || !serviceChartRef.current || !financeChartRef.current) return;

    // 客户统计图表
    const customerChart = echarts.init(customerChartRef.current);
    customerChart.setOption({
      title: { text: '客户增长趋势', left: 'center', textStyle: { fontSize: 14 } },
      tooltip: { trigger: 'axis' },
      xAxis: {
        type: 'category',
        data: ['1月', '2月', '3月', '4月', '5月', '6月'],
      },
      yAxis: { type: 'value' },
      series: [
        {
          name: '新增客户',
          type: 'line',
          data: [12, 19, 15, 22, 28, 35],
          smooth: true,
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(24, 144, 255, 0.3)' },
              { offset: 1, color: 'rgba(24, 144, 255, 0.05)' },
            ]),
          },
          itemStyle: { color: '#1890ff' },
        },
      ],
    });

    // 服务统计图表
    const serviceChart = echarts.init(serviceChartRef.current);
    serviceChart.setOption({
      title: { text: '服务类型分布', left: 'center', textStyle: { fontSize: 14 } },
      tooltip: { trigger: 'item' },
      legend: { bottom: 0 },
      series: [
        {
          name: '服务类型',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2,
          },
          label: { show: false },
          emphasis: {
            label: { show: true, fontSize: 16, fontWeight: 'bold' },
          },
          data: [
            { value: 35, name: '战略咨询' },
            { value: 28, name: '管理咨询' },
            { value: 22, name: '技术咨询' },
            { value: 15, name: '财务咨询' },
          ],
        },
      ],
    });

    // 财务统计图表
    const financeChart = echarts.init(financeChartRef.current);
    financeChart.setOption({
      title: { text: '月度收支统计', left: 'center', textStyle: { fontSize: 14 } },
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      legend: { bottom: 0 },
      xAxis: {
        type: 'category',
        data: ['1月', '2月', '3月', '4月', '5月', '6月'],
      },
      yAxis: { type: 'value' },
      series: [
        {
          name: '收入',
          type: 'bar',
          data: [120000, 150000, 180000, 165000, 210000, 240000],
          itemStyle: { color: '#52c41a' },
        },
        {
          name: '开票',
          type: 'bar',
          data: [100000, 130000, 160000, 140000, 180000, 200000],
          itemStyle: { color: '#1890ff' },
        },
      ],
    });

    // 响应式处理
    const handleResize = () => {
      customerChart.resize();
      serviceChart.resize();
      financeChart.resize();
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      customerChart.dispose();
      serviceChart.dispose();
      financeChart.dispose();
    };
  }, []);

  // 获取优先级标签颜色
  const getPriorityColor = (priority: string) => {
    const colors: Record<string, string> = {
      low: 'default',
      normal: 'processing',
      high: 'warning',
      urgent: 'error',
    };
    return colors[priority] || 'default';
  };

  // 获取优先级标签文本
  const getPriorityText = (priority: string) => {
    const texts: Record<string, string> = {
      low: '低',
      normal: '普通',
      high: '高',
      urgent: '紧急',
    };
    return texts[priority] || priority;
  };

  return (
    <div>
      {/* 数据概览卡片 */}
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={statsLoading}>
            <Statistic
              title="客户总数"
              value={stats?.customerCount || 0}
              prefix={<TeamOutlined />}
              suffix={
                <span style={{ fontSize: 14, marginLeft: 8 }}>
                  {stats?.customerGrowth > 0 ? (
                    <span style={{ color: '#52c41a' }}>
                      <ArrowUpOutlined /> {stats?.customerGrowth}%
                    </span>
                  ) : (
                    <span style={{ color: '#ff4d4f' }}>
                      <ArrowDownOutlined /> {Math.abs(stats?.customerGrowth || 0)}%
                    </span>
                  )}
                </span>
              }
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={statsLoading}>
            <Statistic
              title="服务订单"
              value={stats?.orderCount || 0}
              prefix={<SolutionOutlined />}
              suffix={
                <span style={{ fontSize: 14, marginLeft: 8 }}>
                  {stats?.orderGrowth > 0 ? (
                    <span style={{ color: '#52c41a' }}>
                      <ArrowUpOutlined /> {stats?.orderGrowth}%
                    </span>
                  ) : (
                    <span style={{ color: '#ff4d4f' }}>
                      <ArrowDownOutlined /> {Math.abs(stats?.orderGrowth || 0)}%
                    </span>
                  )}
                </span>
              }
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={statsLoading}>
            <Statistic
              title="本月收入"
              value={formatMoney(stats?.revenue || 0, '')}
              prefix={<DollarOutlined />}
              suffix={
                <span style={{ fontSize: 14, marginLeft: 8 }}>
                  {stats?.revenueGrowth > 0 ? (
                    <span style={{ color: '#52c41a' }}>
                      <ArrowUpOutlined /> {stats?.revenueGrowth}%
                    </span>
                  ) : (
                    <span style={{ color: '#ff4d4f' }}>
                      <ArrowDownOutlined /> {Math.abs(stats?.revenueGrowth || 0)}%
                    </span>
                  )}
                </span>
              }
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={statsLoading}>
            <Statistic
              title="本月开票"
              value={formatMoney(stats?.invoiceAmount || 0, '')}
              prefix={<FileTextOutlined />}
              suffix={
                <span style={{ fontSize: 14, marginLeft: 8 }}>
                  {stats?.invoiceGrowth > 0 ? (
                    <span style={{ color: '#52c41a' }}>
                      <ArrowUpOutlined /> {stats?.invoiceGrowth}%
                    </span>
                  ) : (
                    <span style={{ color: '#ff4d4f' }}>
                      <ArrowDownOutlined /> {Math.abs(stats?.invoiceGrowth || 0)}%
                    </span>
                  )}
                </span>
              }
            />
          </Card>
        </Col>
      </Row>

      {/* 统计图表 */}
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} lg={8}>
          <Card>
            <div ref={customerChartRef} style={{ height: 300 }} />
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card>
            <div ref={serviceChartRef} style={{ height: 300 }} />
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card>
            <div ref={financeChartRef} style={{ height: 300 }} />
          </Card>
        </Col>
      </Row>

      {/* 待办事项 */}
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} lg={12}>
          <Card
            title={
              <span>
                <BellOutlined style={{ marginRight: 8 }} />
                待办事项
                {stats?.pendingTasks ? (
                  <Badge count={stats.pendingTasks} style={{ marginLeft: 8 }} />
                ) : null}
              </span>
            }
          >
            <Spin spinning={todosLoading}>
              {todos && todos.length > 0 ? (
                <List
                  dataSource={todos.slice(0, 5)}
                  renderItem={(item) => (
                    <List.Item
                      actions={[
                        <Tag color={getPriorityColor(item.priority)}>
                          {getPriorityText(item.priority)}
                        </Tag>,
                      ]}
                    >
                      <List.Item.Meta
                        title={item.title}
                        description={item.deadline && `截止: ${item.deadline}`}
                      />
                    </List.Item>
                  )}
                />
              ) : (
                <Empty description="暂无待办事项" />
              )}
            </Spin>
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card
            title={
              <span>
                <CheckCircleOutlined style={{ marginRight: 8 }} />
                快捷入口
              </span>
            }
          >
            <Row gutter={[16, 16]}>
              <Col span={8}>
                <Card.Grid style={{ width: '100%', textAlign: 'center' }}>
                  <TeamOutlined style={{ fontSize: 32, color: '#1890ff' }} />
                  <div style={{ marginTop: 8 }}>新增客户</div>
                </Card.Grid>
              </Col>
              <Col span={8}>
                <Card.Grid style={{ width: '100%', textAlign: 'center' }}>
                  <SolutionOutlined style={{ fontSize: 32, color: '#52c41a' }} />
                  <div style={{ marginTop: 8 }}>新增订单</div>
                </Card.Grid>
              </Col>
              <Col span={8}>
                <Card.Grid style={{ width: '100%', textAlign: 'center' }}>
                  <FileTextOutlined style={{ fontSize: 32, color: '#faad14' }} />
                  <div style={{ marginTop: 8 }}>开具发票</div>
                </Card.Grid>
              </Col>
              <Col span={8}>
                <Card.Grid style={{ width: '100%', textAlign: 'center' }}>
                  <DollarOutlined style={{ fontSize: 32, color: '#722ed1' }} />
                  <div style={{ marginTop: 8 }}>登记收款</div>
                </Card.Grid>
              </Col>
              <Col span={8}>
                <Card.Grid style={{ width: '100%', textAlign: 'center' }}>
                  <FileTextOutlined style={{ fontSize: 32, color: '#eb2f96' }} />
                  <div style={{ marginTop: 8 }}>新建合同</div>
                </Card.Grid>
              </Col>
              <Col span={8}>
                <Card.Grid style={{ width: '100%', textAlign: 'center' }}>
                  <BellOutlined style={{ fontSize: 32, color: '#13c2c2' }} />
                  <div style={{ marginTop: 8 }}>发布通知</div>
                </Card.Grid>
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
