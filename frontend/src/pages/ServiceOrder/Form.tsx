import React, { useEffect, useState } from 'react';
import { Form, Input, Select, DatePicker, InputNumber, Button, Card, message, Spin } from 'antd';
import { ArrowLeftOutlined, SaveOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import { createServiceOrder, updateServiceOrder, getServiceOrder } from '@/api/service';
import { getCustomers } from '@/api/customer';
import { getServiceTypes } from '@/api/service';
import type { Customer, ServiceType } from '@/types';
import dayjs from 'dayjs';

const { TextArea } = Input;
const { RangePicker } = DatePicker;

const ServiceOrderForm: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const isEdit = !!id;
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [saveLoading, setSaveLoading] = useState(false);
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [serviceTypes, setServiceTypes] = useState<ServiceType[]>([]);

  // 获取客户列表
  const fetchCustomers = async () => {
    try {
      const result = await getCustomers({ pageSize: 1000 });
      setCustomers(result.list);
    } catch (error) {
      console.error('获取客户列表失败:', error);
    }
  };

  // 获取服务类型列表
  const fetchServiceTypes = async () => {
    try {
      const result = await getServiceTypes({ pageSize: 1000 });
      setServiceTypes(result.list);
    } catch (error) {
      console.error('获取服务类型列表失败:', error);
    }
  };

  // 获取服务订单详情
  const fetchDetail = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const result = await getServiceOrder(Number(id));
      form.setFieldsValue({
        ...result,
        dateRange: result.start_date && result.end_date 
          ? [dayjs(result.start_date), dayjs(result.end_date)] 
          : undefined,
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCustomers();
    fetchServiceTypes();
    if (isEdit) {
      fetchDetail();
    }
  }, [id]);

  // 提交表单
  const handleSubmit = async (values: any) => {
    setSaveLoading(true);
    try {
      const data = {
        ...values,
        start_date: values.dateRange?.[0]?.format('YYYY-MM-DD'),
        end_date: values.dateRange?.[1]?.format('YYYY-MM-DD'),
      };
      delete data.dateRange;

      if (isEdit) {
        await updateServiceOrder(Number(id), data);
        message.success('更新成功');
      } else {
        await createServiceOrder(data);
        message.success('创建成功');
      }
      navigate('/service-orders');
    } catch (error) {
      console.error('保存失败:', error);
    } finally {
      setSaveLoading(false);
    }
  };

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/service-orders')}>
          返回列表
        </Button>
      </div>

      <Spin spinning={loading}>
        <Card title={isEdit ? '编辑服务订单' : '新建服务订单'}>
          <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
            style={{ maxWidth: 800 }}
          >
            <Form.Item
              name="code"
              label="服务编号"
              rules={[{ required: true, message: '请输入服务编号' }]}
            >
              <Input placeholder="请输入服务编号" disabled={isEdit} />
            </Form.Item>

            <Form.Item
              name="name"
              label="服务名称"
              rules={[{ required: true, message: '请输入服务名称' }]}
            >
              <Input placeholder="请输入服务名称" />
            </Form.Item>

            <Form.Item
              name="customer_id"
              label="客户"
              rules={[{ required: true, message: '请选择客户' }]}
            >
              <Select
                placeholder="请选择客户"
                showSearch
                optionFilterProp="children"
                options={customers.map(c => ({ label: c.name, value: c.id }))}
              />
            </Form.Item>

            <Form.Item
              name="service_type_id"
              label="服务类型"
              rules={[{ required: true, message: '请选择服务类型' }]}
            >
              <Select
                placeholder="请选择服务类型"
                options={serviceTypes.map(s => ({ label: s.name, value: s.id }))}
              />
            </Form.Item>

            <Form.Item
              name="amount"
              label="服务金额"
              rules={[{ required: true, message: '请输入服务金额' }]}
            >
              <InputNumber
                style={{ width: '100%' }}
                prefix="¥"
                placeholder="请输入服务金额"
                min={0}
                precision={2}
              />
            </Form.Item>

            <Form.Item name="dateRange" label="服务周期">
              <RangePicker style={{ width: '100%' }} />
            </Form.Item>

            <Form.Item name="status" label="状态" initialValue={1}>
              <Select
                options={[
                  { label: '待启动', value: 1 },
                  { label: '进行中', value: 2 },
                  { label: '已完成', value: 3 },
                  { label: '已暂停', value: 4 },
                ]}
              />
            </Form.Item>

            <Form.Item name="progress" label="进度(%)" initialValue={0}>
              <InputNumber style={{ width: '100%' }} min={0} max={100} />
            </Form.Item>

            <Form.Item name="description" label="服务描述">
              <TextArea rows={4} placeholder="请输入服务描述" />
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                icon={<SaveOutlined />}
                loading={saveLoading}
              >
                保存
              </Button>
              <Button style={{ marginLeft: 8 }} onClick={() => navigate('/service-orders')}>
                取消
              </Button>
            </Form.Item>
          </Form>
        </Card>
      </Spin>
    </div>
  );
};

export default ServiceOrderForm;
