import { Suspense } from 'react';
import HomeClient from '@/components/HomeClient';
import LoadingState from '@/components/HomeImageGrid/LoadingState';

export default function Home() {
  return (
    <Suspense fallback={<LoadingState message="Loading application..." />}>
      <HomeClient />
    </Suspense>
  );
}
