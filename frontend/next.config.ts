import type { NextConfig } from "next";

// ✅ Bundle analyzer (啟用: ANALYZE=true npm run build)
const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
});

const nextConfig: NextConfig = {
  // ✅ 啟用 React 編譯器 (React 19 新特性)
  reactCompiler: true,

  // ✅ 圖片優化配置
  images: {
    // 允許 base64 data URLs (你的應用使用 base64 圖片)
    dangerouslyAllowSVG: false,
    remotePatterns: [
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '8123',
      },
      // 添加生產環境的 hostname
    ],
    formats: ['image/webp', 'image/avif'],
    // 優化內聯小圖片
    minimumCacheTTL: 60,
  },

  // ✅ 生產環境優化
  compress: true,
  poweredByHeader: false,

  // ✅ 嚴格模式
  reactStrictMode: true,

  // ✅ Standalone 輸出模式,優化 Docker 部署
  output: 'standalone',

  // ✅ 實驗性功能 (Next.js 15)
  experimental: {
    // 減少初始 JS bundle
    optimizePackageImports: ['lucide-react'],
    // 優化 CSS
    optimizeCss: true,
  },

  // ✅ Webpack 優化
  webpack: (config, { isServer }) => {
    // 優化 bundle 大小
    if (!isServer) {
      config.resolve.fallback = {
        ...config.resolve.fallback,
        fs: false,
        net: false,
        tls: false,
      };
    }
    return config;
  },

  // ✅ 跳過靜態生成時的類型檢查和 ESLint (加快 Docker build)
  typescript: {
    // 開發時仍會檢查,但 build 時跳過以加快速度
    ignoreBuildErrors: process.env.SKIP_TYPE_CHECK === 'true',
  },
  eslint: {
    ignoreDuringBuilds: process.env.SKIP_LINT === 'true',
  },
};

export default withBundleAnalyzer(nextConfig);
