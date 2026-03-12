import React, { useEffect, useState } from 'react';
import { Card, Form, Tree, Button, message, Space } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { ProForm, ProFormText, ProFormTextArea, ProFormSelect } from '@ant-design/pro-components';
import { getRoleById, createRole, updateRole, getPermissionTree } from '@/api/user';
import type { RoleFormData, Permission } from '@/types';

const RoleForm: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const isEdit = !!id;
  const [form] = Form.useForm();
  const queryClient = useQueryClient();
  const [selectedPermissions, setSelectedPermissions] = useState<number[]>([]);

  // 获取权限树
  const { data: permissionTree, isLoading: treeLoading } = useQuery({
    queryKey: ['permissionTree'],
    queryFn: getPermissionTree,
  });

  // 获取角色详情
  const { data: roleData, isLoading } = useQuery({
    queryKey: ['role', id],
    queryFn: () => getRoleById(Number(id)),
    enabled: isEdit,
  });

  // 设置表单数据和权限
  useEffect(() => {
    if (roleData) {
      form.setFieldsValue({
        name: roleData.name,
        code: roleData.code,
        description: roleData.description,
        status: roleData.status,
      });
      setSelectedPermissions(roleData.permissions?.map((p) => p.id) || []);
    }
  }, [roleData, form]);

  // 创建/更新角色
  const mutation = useMutation({
    mutationFn: (values: RoleFormData) => {
      const data = { ...values, permissionIds: selectedPermissions };
      if (isEdit) {
        return updateRole(Number(id), data);
      }
      return createRole(data);
    },
    onSuccess: () => {
      message.success(isEdit ? '更新成功' : '创建成功');
      queryClient.invalidateQueries({ queryKey: ['roles'] });
      navigate('/roles');
    },
  });

  const handleSubmit = async (values: RoleFormData) => {
    if (selectedPermissions.length === 0) {
      message.warning('请至少选择一个权限');
      return;
    }
    mutation.mutate(values);
  };

  // 转换权限树
  const convertPermissionTree = (permissions: Permission[] | undefined) => {
    if (!permissions) return [];
    return permissions.map((perm) => ({
      title: perm.name,
      key: perm.id,
      children: perm.children ? convertPermissionTree(perm.children) : undefined,
    }));
  };

  return (
    <Card title={isEdit ? '编辑角色' : '新建角色'} loading={isLoading}>
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
                <Button onClick={() => navigate('/roles')}>取消</Button>
              </Space>
            </Form.Item>
          ),
        }}
      >
        <ProFormText
          name="name"
          label="角色名称"
          rules={[
            { required: true, message: '请输入角色名称' },
          ]}
          placeholder="请输入角色名称"
        />

        <ProFormText
          name="code"
          label="角色编码"
          rules={[
            { required: true, message: '请输入角色编码' },
          ]}
          placeholder="请输入角色编码"
        />

        <ProFormTextArea
          name="description"
          label="描述"
          placeholder="请输入描述"
        />

        <ProFormSelect
          name="status"
          label="状态"
          initialValue={1}
          options={[
            { label: '启用', value: 1 },
            { label: '禁用', value: 0 },
          ]}
        />

        <Form.Item
          label="权限"
          required
          rules={[{ required: true, message: '请选择权限' }]}
        >
          <Card loading={treeLoading} style={{ maxHeight: 400, overflow: 'auto' }}>
            <Tree
              checkable
              treeData={convertPermissionTree(permissionTree)}
              checkedKeys={selectedPermissions}
              onCheck={(checkedKeys) => {
                setSelectedPermissions(checkedKeys as number[]);
              }}
            />
          </Card>
        </Form.Item>
      </ProForm>
    </Card>
  );
};

export default RoleForm;
