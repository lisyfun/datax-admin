declare module 'sm-crypto' {
  export const sm4: {
    encrypt: (data: string, key: string) => string;
    decrypt: (data: string, key: string) => string;
  };
}
