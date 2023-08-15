import { Card, Title } from '@tremor/react';
import Search from './search';

export const dynamic = 'force-dynamic';

export default async function IndexPage({
  searchParams
}: {
  searchParams: { q: string };
}) {
  const search = searchParams.q ?? '';
  return (
    <main className="p-4 md:p-10 mx-auto max-w-7xl">
      <Title>Dashboard</Title>

      <Search />
      <Card className="mt-6">
        Some statistic...
      </Card>
    </main>
  );
}
