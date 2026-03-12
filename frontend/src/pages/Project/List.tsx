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
  Progress,
  Popconfirm,
  message,
} from 'antd';
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { getProjectList, deleteProject } from '@/api/project';
import { Project, ProjectStatus, ProjectType } from '@/types';

const { Option } = Select;

const statusMap: Record<ProjectStatus, { text: string; color: string }> = {
  pending: { text: '待启动', color: 'default' },
  active: { text: '进行中', color: 'processing' },
  paused: { text: '已暂停', color: 'warning' },
  completed: { text: '已完成', color: 'success' },
  cancelled: { text: '已取消', color: 'error' },
};

const typeMap: Record<ProjectType, string> = {
  strategy: '战略咨询',
  management: '管理咨询',
  technology: '技术咨询',
  finance: '财务咨询',
  hr: '人力资源',
  marketing: '市场营销',
  other: '其他',
};

const ProjectList = () => {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useState({
    page: 1,
    page_size: 10,
    keyword: '',
    status: '',
    type: '',
  });

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['projects', searchParams],
    queryFn: () => getProjectList(searchParams),
  });

  const handleDelete = async (id: string) => {
    try {
      await deleteProject(id);
      message.success('删除成功');
      refetch();
    } catch (error) {
      // 错误已在 request 中处理
    }
  };

  const columns = [
    {
      title: '项目编号',
      dataIndex: 'code',
      key: 'code',
      width: 150,
    },
    {
      title: '项目名称',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: Project) => (
        <div>
          <div style={{ fontWeight: 500 }}>{name}</div>
          <div style={{ fontSize: 12, color: '#999' }}>{typeMap[record.type]}</div>
        </div>
      ),
    },
    {
      title: '客户',
      dataIndex: 'customer',
      key: 'customer',
      render: (customer: any) => customer?.name || '-',
    },
    {
      title: '项目经理',
      dataIndex: 'manager',
      key: 'manager',
      render: (manager: any) => manager?.real_name || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: ProjectStatus) => (
        <Tag color={statusMap[status]?.color}>{statusMap[status]?.text}</Tag>
      ),
    },
    {
      title: '进度',
      dataIndex: 'progress',
      key: 'progress',
      render: (progress: number) => (
        <Progress percent={progress} size="small" />
      ),
    },
    {
      title: '预算',
      dataIndex: 'budget',
      key: 'budget',
      render: (budget: number) => `¥${(budget || 0).toLocaleString()}`,
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Project) => (
        <Space>
          <Button
            type="text"
            icon={<EyeOutlined />}
            size="small"
            onClick={() => navigate(`/projects/${record.id}`)}
          >
            详情
          </Button>
          <Button type="text" icon={<EditOutlined />} size="small">
            编辑
          </Button>
          <Popconfirm
            title="确认删除"
            description="确定要删除该项目吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="text" danger icon={<DeleteOutlined />} size="small">
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const list = data?.data?.list || [];
  const total = data?.data?.total || 0;

  return (
    <div>
      <Card>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 24 }}>
          <h2 style={{ margin: 0 }}>项目管理</h2>
          <Button type="primary" icon={<PlusOutlined />}>
            新增项目
          </Button>
        </div>

        <div style={{ marginBottom: 16, padding: 16, background: '#fafafa', borderRadius: 8 }}>
          <Space>
            <Input
              placeholder="搜索项目名称/编号"
              prefix={<SearchOutlined />}
              value={searchParams.keyword}
              onChange={(e) =>
                setSearchParams({ ...searchParams, keyword: e.target.value, page: 1 })
              }
              allowClear
              style={{ width: 250 }}
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
              <Option value="pending">待启动</Option>
              <Option value="active">进行中</Option>
              <Option value="paused">已暂停</Option>
              <Option value="completed">已完成</Option>
              <Option value="cancelled">已取消</Option>
            </Select>
            <Select
              placeholder="类型"
              value={searchParams.type || undefined}
              onChange={(value) =>
                setSearchParams({ ...searchParams, type: value, page: 1 })
              }
              allowClear
              style={{ width: 120 }}
            >
              <Option value="strategy">战略咨询</Option>
              <Option value="management">管理咨询</Option>
              <Option value="technology">技术咨询</Option>
              <Option value="finance">财务咨询</Option>
              <Option value="hr">人力资源</Option>
              <Option value="marketing">市场营销</Option>
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

export default ProjectList;
