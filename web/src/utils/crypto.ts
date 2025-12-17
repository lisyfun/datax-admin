import CryptoJS from 'crypto-js';

// 密码加密 - 使用 MD5 作为前端预处理
// 注意：这只是前端的预处理，后端会再次使用 bcrypt 进行加密
export function encryptPassword(password: string): string {
  // 先进行 MD5，再加上固定的盐值，再进行一次 MD5
  const md5Pass = CryptoJS.MD5(password).toString();
  const saltedPass = md5Pass + 'datax-admin';
  return CryptoJS.MD5(saltedPass).toString();
}

// 加密配置
interface CryptoConfig {
  secretKey: string;
  enableLog?: boolean;
}

// 加密请求数据结构
interface EncryptedRequest {
  data: string;
  timestamp: number;
  signature: string;
}

// 加密响应数据结构
interface EncryptedResponse {
  data: string;
  timestamp: number;
  signature: string;
}

class CryptoManager {
  private config: CryptoConfig;

  constructor(config: CryptoConfig) {
    this.config = config;
  }

  /**
   * 加密数据 - 使用AES-CBC模式
   */
  encrypt(plaintext: string): string {
    try {
      // 生成随机IV
      const iv = CryptoJS.lib.WordArray.random(16);

      // 使用AES-CBC加密
      const encrypted = CryptoJS.AES.encrypt(plaintext, CryptoJS.enc.Utf8.parse(this.config.secretKey), {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
      });

      // 组合IV和密文
      const combined = iv.concat(encrypted.ciphertext);
      return CryptoJS.enc.Base64.stringify(combined);
    } catch (error) {
      console.error('加密失败:', error);
      throw new Error('数据加密失败');
    }
  }

  /**
   * 解密数据 - 使用AES-CBC模式
   */
  decrypt(ciphertext: string): string {
    try {
      // Base64解码
      const combined = CryptoJS.enc.Base64.parse(ciphertext);

      // 提取IV (前16字节)
      const iv = CryptoJS.lib.WordArray.create(combined.words.slice(0, 4));

      // 提取密文 (剩余部分)
      const encrypted = CryptoJS.lib.WordArray.create(combined.words.slice(4));

      // 解密
      const decrypted = CryptoJS.AES.decrypt(
        { ciphertext: encrypted } as any,
        CryptoJS.enc.Utf8.parse(this.config.secretKey),
        {
          iv: iv,
          mode: CryptoJS.mode.CBC,
          padding: CryptoJS.pad.Pkcs7
        }
      );

      return decrypted.toString(CryptoJS.enc.Utf8);
    } catch (error) {
      console.error('解密失败:', error);
      throw new Error('数据解密失败');
    }
  }

  /**
   * 生成签名
   */
  generateSignature(data: string, timestamp: number): string {
    const message = `${data}${timestamp}`;
    return CryptoJS.HmacSHA256(message, this.config.secretKey).toString();
  }

  /**
   * 验证签名
   */
  verifySignature(data: string, timestamp: number, signature: string): boolean {
    // 检查时间戳是否在有效范围内（5分钟）
    const now = Math.floor(Date.now() / 1000);
    if (Math.abs(now - timestamp) > 300) {
      return false;
    }

    const expectedSignature = this.generateSignature(data, timestamp);
    return expectedSignature === signature;
  }

  /**
   * 加密请求数据
   */
  encryptRequest(data: any): EncryptedRequest {
    try {
      // 序列化数据
      const jsonData = JSON.stringify(data);

      // 加密数据
      const encryptedData = this.encrypt(jsonData);

      // 生成时间戳
      const timestamp = Math.floor(Date.now() / 1000);

      // 生成签名
      const signature = this.generateSignature(encryptedData, timestamp);

      return {
        data: encryptedData,
        timestamp,
        signature,
      };
    } catch (error) {
      console.error('加密请求失败:', error);
      throw new Error('请求数据加密失败');
    }
  }

  /**
   * 解密响应数据
   */
  decryptResponse(encryptedResponse: EncryptedResponse): any {
    try {
      // 验证签名
      if (!this.verifySignature(encryptedResponse.data, encryptedResponse.timestamp, encryptedResponse.signature)) {
        throw new Error('响应签名验证失败');
      }

      // 解密数据
      const decryptedData = this.decrypt(encryptedResponse.data);

      // 解析JSON
      return JSON.parse(decryptedData);
    } catch (error) {
      console.error('解密响应失败:', error);
      throw new Error('响应数据解密失败');
    }
  }
}

// 创建全局加密管理器实例
const cryptoManager = new CryptoManager({
  secretKey: import.meta.env.VITE_CRYPTO_SECRET_KEY || 'your-32-character-secret-key-here',
  enableLog: import.meta.env.DEV,
});

// 导出加密管理器
export { CryptoManager, cryptoManager };

// 导出类型
export type { EncryptedRequest, EncryptedResponse };

// 便捷方法
export const encryptData = (data: any): EncryptedRequest => {
  return cryptoManager.encryptRequest(data);
};

export const decryptData = (encryptedResponse: EncryptedResponse): any => {
  return cryptoManager.decryptResponse(encryptedResponse);
};

// 检查是否启用加密
export const isEncryptionEnabled = (): boolean => {
  return import.meta.env.VITE_ENABLE_ENCRYPTION === 'true';
};
