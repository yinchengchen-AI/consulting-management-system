import React, { useEffect } from 'react';
import { Card, Form, Input, Select, Button, message, Space } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { ProForm, ProFormText, ProFormSelect } from '@ant-design/pro-components';
import { getUserById, createUser, updateUser, getAllRoles } from '@/api/user';
import type { UserFormData, Role } from '@/types';

const UserForm: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const isEdit = !!id;
  const [form] = Form.useForm();
  const queryClient = useQueryClient();

  // 获取角色列表
  const { data: roles } = useQuery<Role[]>({
    queryKey: ['roles'],
    queryFn: getAllRoles,
  });

  // 获取用户详情
  const { data: userData, isLoading } = useQuery({
    queryKey: ['user', id],
    queryFn: () => getUserById(Number(id)),
    enabled: isEdit,
  });

  // 设置表单数据
  useEffect(() => {
    if (userData) {
      form.setFieldsValue({
        ...userData,
        roleIds: userData.roles?.map((role) => role.id),
      });
    }
  }, [userData, form]);

  // 创建/更新用户
  const mutation = useMutation({
    mutationFn: (values: UserFormData) => {
      if (isEdit) {
        return updateUser(Number(id), values);
      }
      return createUser(values);
    },
    onSuccess: () => {
      message.success(isEdit ? '更新成功' : '创建成功');
      queryClient.invalidateQueries({ queryKey: ['users'] });
      navigate('/users');
    },
  });

  const handleSubmit = async (values: UserFormData) => {
    mutation.mutate(values);
  };

  return (
    <Card title={isEdit ? '编辑用户' : '新建用户'} loading={isLoading}>
      <ProForm
        form={form}
        onFinish={handleSubmit}
        layout="horizontal"
        labelCol={{ span: 4 }}
        wrapperCol={{ span: 16 }}
        submitter={{
          render: (props) => (
            <Form.Item wrapperCol={{ offset: 4, span: 16 }}>
              <Space>
                <Button type="primary" onClick={() => props.submit?.()}>
                  保存
                </Button>
                <Button onClick={() => navigate('/users')}>取消</Button>
              </Space>
            </Form.Item>
          ),
        }}
      >
        <ProFormText
          name="username"
          label="用户名"
          rules={[
            { required: true, message: '请输入用户名' },
            { min: 3, message: '用户名至少3个字符' },
          ]}
          placeholder="请输入用户名"
        />

        {!isEdit && (
          <ProFormText.Password
            name="password"
            label="密码"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 6, message: '密码至少6个字符' },
            ]}
            placeholder="请输入密码"
          />
        )}

        <ProFormText
          name="realName"
          label="姓名"
          rules={[{ required: true, message: '请输入姓名' }]}
          placeholder="请输入姓名"
        />

        <ProFormText
          name="email"
          label="邮箱"
          rules={[
            { type: 'email', message: '请输入有效的邮箱地址' },
          ]}
          placeholder="请输入邮箱"
        />

        <ProFormText
          name="phone"
          label="手机号"
          placeholder="请输入手机号"
        />

        <ProFormSelect
          name="departmentId"
          label="部门"
          options={[]}
          placeholder="请选择部门"
        />

        <ProFormText
          name="position"
          label="职位"
          placeholder="请输入职位"
        />

        <Form.Item
          name="roleIds"
          label="角色"
          rules={[{ required: true, message: '请选择角色' }]}
        >
          <Select
            mode="multiple"
            placeholder="请选择角色"
            options={roles?.map((role) => ({ label: role.name, value: role.id }))}
          />
        </Form.Item>

        <ProFormSelect
          name="status"
          label="状态"
          initialValue={1}
          options={[
            { label: '启用', value: 1 },
            { label: '禁用', value: 0 },
          ]}
        />
      </ProForm>
    </Card>
  );
};

export default UserForm;
