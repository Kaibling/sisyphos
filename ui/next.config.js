/** @type {import('next').NextConfig} */
const backendURL = process.env.NEXT_PUBLIC_BACKEND_URL;

const nextConfig = {
  images: {
    domains: ['avatars.githubusercontent.com', 'avatar.vercel.sh']
  },
  eslint: {
    // Warning: This allows production builds to successfully complete even if
    // your project has ESLint errors.
    ignoreDuringBuilds: true,
  },
  typescript: {
    // !! WARN !!
    // Dangerously allow production builds to successfully complete even if
    // your project has type errors.
    // !! WARN !!
    ignoreBuildErrors: true,
  },
  async rewrites() {
    return [
      {
        source: '/app/:path*',
        destination: `${backendURL}/:path*`,
      },
    ]
  },
};

module.exports = nextConfig;
