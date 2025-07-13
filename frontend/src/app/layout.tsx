import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import '@/styles/globals.css';

import { JobDatasetProvider } from '@/components/JobDatasetContext';
import LeftSidebar from '@/components/HomeLeftSidebar/index';
import RightSidebar from '@/components/HomeRightSidebar/index';


const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Image Verify Viewer',
  description: 'A simple image viewer for verification',
};

export default function RootLayout({ children }: {
  children: React.ReactNode;
}) {
  
  return (
    <html lang="en">
      <body className={inter.className}>
        <JobDatasetProvider>
          <main className="flex h-screen">
            <LeftSidebar />
            <div className="flex-1">{children}</div>
            <RightSidebar />
          </main>
        </JobDatasetProvider>
      </body>
    </html>
  );
}