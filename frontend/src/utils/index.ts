import dayjs from 'dayjs';

// ==================== 日期工具 ====================

export const formatDate = (date: string | Date | number, format = 'YYYY-MM-DD'): string => {
  if (!date) return '-';
  return dayjs(date).format(format);
};

export const formatDateTime = (date: string | Date | number): string => {
  if (!date) return '-';
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
};

export const formatTime = (date: string | Date | number): string => {
  if (!date) return '-';
  return dayjs(date).format('HH:mm:ss');
};

// ==================== 数字工具 ====================

export const formatNumber = (num: number | string, decimals = 0): string => {
  if (num === null || num === undefined) return '-';
  const n = typeof num === 'string' ? parseFloat(num) : num;
  if (isNaN(n)) return '-';
  return n.toLocaleString('zh-CN', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  });
};

export const formatMoney = (amount: number | string, prefix = '¥'): string => {
  if (amount === null || amount === undefined) return '-';
  const n = typeof amount === 'string' ? parseFloat(amount) : amount;
  if (isNaN(n)) return '-';
  return `${prefix}${n.toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })}`;
};

export const formatPercent = (value: number | string, decimals = 2): string => {
  if (value === null || value === undefined) return '-';
  const n = typeof value === 'string' ? parseFloat(value) : value;
  if (isNaN(n)) return '-';
  return `${(n * 100).toFixed(decimals)}%`;
};

// ==================== 文件工具 ====================

export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
};

export const getFileExtension = (filename: string): string => {
  if (!filename) return '';
  const index = filename.lastIndexOf('.');
  return index === -1 ? '' : filename.slice(index + 1).toLowerCase();
};

// ==================== 验证工具 ====================

export const isValidEmail = (email: string): boolean => {
  const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return regex.test(email);
};

export const isValidPhone = (phone: string): boolean => {
  const regex = /^1[3-9]\d{9}$/;
  return regex.test(phone);
};

export const isValidURL = (url: string): boolean => {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
};

// ==================== 数组工具 ====================

export const uniqueArray = <T>(arr: T[]): T[] => {
  return [...new Set(arr)];
};

export const groupBy = <T>(arr: T[], key: keyof T): Record<string, T[]> => {
  return arr.reduce((result, item) => {
    const groupKey = String(item[key]);
    if (!result[groupKey]) {
      result[groupKey] = [];
    }
    result[groupKey].push(item);
    return result;
  }, {} as Record<string, T[]>);
};

// ==================== 对象工具 ====================

export const pick = <T extends Record<string, unknown>, K extends keyof T>(
  obj: T,
  keys: K[]
): Pick<T, K> => {
  const result = {} as Pick<T, K>;
  keys.forEach((key) => {
    if (key in obj) {
      result[key] = obj[key];
    }
  });
  return result;
};

export const omit = <T extends Record<string, unknown>, K extends keyof T>(
  obj: T,
  keys: K[]
): Omit<T, K> => {
  const result = { ...obj };
  keys.forEach((key) => {
    delete result[key];
  });
  return result as Omit<T, K>;
};

// ==================== 树形数据工具 ====================

interface TreeNode {
  id: number;
  parentId?: number;
  children?: TreeNode[];
  [key: string]: unknown;
}

export const arrayToTree = <T extends TreeNode>(arr: T[]): T[] => {
  const map: Record<number, T> = {};
  const result: T[] = [];
  
  arr.forEach((item) => {
    map[item.id] = { ...item, children: [] };
  });
  
  arr.forEach((item) => {
    const node = map[item.id];
    if (item.parentId && map[item.parentId]) {
      map[item.parentId].children = map[item.parentId].children || [];
      map[item.parentId].children!.push(node);
    } else {
      result.push(node);
    }
  });
  
  return result;
};

export const treeToArray = <T extends TreeNode>(tree: T[]): T[] => {
  const result: T[] = [];
  
  const traverse = (nodes: T[]) => {
    nodes.forEach((node) => {
      const { children, ...rest } = node;
      result.push(rest as T);
      if (children && children.length > 0) {
        traverse(children as T[]);
      }
    });
  };
  
  traverse(tree);
  return result;
};

// ==================== 其他工具 ====================

export const debounce = <T extends (...args: unknown[]) => unknown>(
  fn: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let timer: ReturnType<typeof setTimeout> | null = null;
  return (...args: Parameters<T>) => {
    if (timer) clearTimeout(timer);
    timer = setTimeout(() => {
      fn(...args);
    }, delay);
  };
};

export const throttle = <T extends (...args: unknown[]) => unknown>(
  fn: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let lastTime = 0;
  return (...args: Parameters<T>) => {
    const now = Date.now();
    if (now - lastTime >= delay) {
      fn(...args);
      lastTime = now;
    }
  };
};

export const generateId = (): string => {
  return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
};

export const copyToClipboard = async (text: string): Promise<boolean> => {
  try {
    await navigator.clipboard.writeText(text);
    return true;
  } catch {
    return false;
  }
};
