import { Suspense } from 'react';
import { BrowserRouter } from 'react-router-dom';
import { Spin } from 'antd';
import AppRouter from './router';
import './App.css';

// 全局加载组件
const GlobalLoading = () => (
  <div className="global-loading">
    <Spin size="large" tip="加载中..." />
  </div>
);

function App() {
  return (
    <BrowserRouter>
      <Suspense fallback={<GlobalLoading />}>
        <AppRouter />
      </Suspense>
    </BrowserRouter>
  );
}

export default App;
