import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { Card, Descriptions, Tag, Button, Spin, Empty, Progress } from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { getProjectDetail } from '@/api/project';
import { ProjectStatus, ProjectType } from '@/types';

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

const ProjectDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data, isLoading } = useQuery({
    queryKey: ['project', id],
    queryFn: () => getProjectDetail(id!),
    enabled: !!id,
  });

  const project = data?.data;

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: 50 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!project) {
    return <Empty description="项目不存在" />;
  }

  return (
    <div>
      <Card>
        <div style={{ marginBottom: 24 }}>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/projects')}>
            返回列表
          </Button>
        </div>

        <Descriptions title="基本信息" bordered column={2}>
          <Descriptions.Item label="项目编号">{project.code}</Descriptions.Item>
          <Descriptions.Item label="项目名称">{project.name}</Descriptions.Item>
          <Descriptions.Item label="项目类型">
            {typeMap[project.type] || project.type}
          </Descriptions.Item>
          <Descriptions.Item label="项目状态">
            <Tag color={statusMap[project.status]?.color}>
              {statusMap[project.status]?.text}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="客户">{project.customer?.name || '-'}</Descriptions.Item>
          <Descriptions.Item label="关联合同">{project.contract?.name || '-'}</Descriptions.Item>
          <Descriptions.Item label="项目经理">
            {project.manager?.real_name || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="优先级">P{project.priority}</Descriptions.Item>
        </Descriptions>

        <Descriptions title="时间信息" bordered column={2} style={{ marginTop: 24 }}>
          <Descriptions.Item label="开始日期">
            {project.start_date ? new Date(project.start_date).toLocaleDateString() : '-'}
          </Descriptions.Item>
          <Descriptions.Item label="结束日期">
            {project.end_date ? new Date(project.end_date).toLocaleDateString() : '-'}
          </Descriptions.Item>
        </Descriptions>

        <Descriptions title="财务信息" bordered column={2} style={{ marginTop: 24 }}>
          <Descriptions.Item label="预算金额">
            ¥{(project.budget || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="实际成本">
            ¥{(project.actual_cost || 0).toLocaleString()}
          </Descriptions.Item>
        </Descriptions>

        <Descriptions title="项目进度" bordered column={1} style={{ marginTop: 24 }}>
          <Descriptions.Item label="当前进度">
            <Progress percent={project.progress} />
          </Descriptions.Item>
        </Descriptions>

        {project.description && (
          <Descriptions title="项目描述" bordered column={1} style={{ marginTop: 24 }}>
            <Descriptions.Item label="描述">{project.description}</Descriptions.Item>
          </Descriptions>
        )}

        {project.deliverables && (
          <Descriptions title="交付物" bordered column={1} style={{ marginTop: 24 }}>
            <Descriptions.Item label="交付内容">{project.deliverables}</Descriptions.Item>
          </Descriptions>
        )}
      </Card>
    </div>
  );
};

export default ProjectDetail;
