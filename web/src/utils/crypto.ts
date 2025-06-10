import CryptoJS from 'crypto-js';

// 密码加密 - 使用 MD5 作为前端预处理
// 注意：这只是前端的预处理，后端会再次使用 bcrypt 进行加密
export function encryptPassword(password: string): string {
  // 先进行 MD5，再加上固定的盐值，再进行一次 MD5
  const md5Pass = CryptoJS.MD5(password).toString();
  const saltedPass = md5Pass + 'datax-admin';
  return CryptoJS.MD5(saltedPass).toString();
}
