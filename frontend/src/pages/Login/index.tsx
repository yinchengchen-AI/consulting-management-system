import React, { useState, useEffect } from 'react';
import { Form, Input, Button, Checkbox, Card, message, Spin } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { login, getCurrentUser } from '@/api/auth';
import { useAuthStore } from '@/stores/auth';
import type { LoginParams } from '@/types';

const Login: React.FC = () => {
  const navigate = useNavigate();
  const { setToken, setUserInfo, isAuthenticated } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();

  // 如果已登录，跳转到首页
  useEffect(() => {
    if (isAuthenticated) {
      navigate('/');
    }
  }, [isAuthenticated, navigate]);

  // 处理登录
  const handleLogin = async (values: LoginParams) => {
    setLoading(true);
    try {
      const result = await login(values);
      
      // 保存token
      setToken(result.token, result.refreshToken);
      
      // 获取用户信息
      const userInfo = await getCurrentUser();
      setUserInfo(userInfo);
      
      message.success('登录成功');
      navigate('/');
    } catch (error) {
      // 错误已在请求拦截器中处理
      console.error('登录失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        height: '100vh',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      }}
    >
      <Card
        title={
          <div style={{ textAlign: 'center', fontSize: 24, fontWeight: 'bold' }}>
            咨询公司业务管理系统
          </div>
        }
        style={{ width: 400, boxShadow: '0 4px 12px rgba(0,0,0,0.15)' }}
      >
        <Spin spinning={loading}>
          <Form
            form={form}
            name="login"
            initialValues={{ remember: true }}
            onFinish={handleLogin}
            autoComplete="off"
            size="large"
          >
            <Form.Item
              name="username"
              rules={[
                { required: true, message: '请输入用户名' },
                { min: 3, message: '用户名至少3个字符' },
              ]}
            >
              <Input
                prefix={<UserOutlined />}
                placeholder="用户名"
                autoFocus
              />
            </Form.Item>

            <Form.Item
              name="password"
              rules={[
                { required: true, message: '请输入密码' },
                { min: 6, message: '密码至少6个字符' },
              ]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="密码"
              />
            </Form.Item>

            <Form.Item>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <Form.Item name="remember" valuePropName="checked" noStyle>
                  <Checkbox>记住我</Checkbox>
                </Form.Item>
                <a href="#">忘记密码？</a>
              </div>
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                block
                loading={loading}
              >
                登录
              </Button>
            </Form.Item>
          </Form>
        </Spin>
      </Card>
    </div>
  );
};

export default Login;
